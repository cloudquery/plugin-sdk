package scheduler

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/helpers"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/resolvers"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/semaphore"
)

func (s *syncClient) syncDfs(ctx context.Context, resolvedResources chan<- *schema.Resource) {
	// we have this because plugins can return sometimes clients in a random way which will cause
	// differences between this run and the next one.
	preInitialisedClients := make([][]schema.ClientMeta, len(s.tables))
	for i, table := range s.tables {
		clients := []schema.ClientMeta{s.client}
		if table.Multiplex != nil {
			clients = table.Multiplex(s.client)
		}
		// Detect duplicate clients while multiplexing
		seenClients := make(map[string]bool)
		for _, c := range clients {
			if _, ok := seenClients[c.ID()]; !ok {
				seenClients[c.ID()] = true
			} else {
				s.logger.Warn().Str("client", c.ID()).Str("table", table.Name).Msg("multiplex returned duplicate client")
			}
		}
		preInitialisedClients[i] = clients
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		s.metrics.InitWithClients(table, clients, s.invocationID)
	}

	tableClients := make([]tableClient, 0)
	for i, table := range s.tables {
		for _, client := range preInitialisedClients[i] {
			tableClients = append(tableClients, tableClient{table: table, client: client})
		}
	}
	tableClients = shardTableClients(tableClients, s.shard)

	var wg sync.WaitGroup
	for _, tc := range tableClients {
		table := tc.table
		cl := tc.client
		if err := s.scheduler.tableSems[0].Acquire(ctx, 1); err != nil {
			// This means context was cancelled
			wg.Wait()
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer s.scheduler.tableSems[0].Release(1)
			// not checking for error here as nothing much to do.
			// the error is logged and this happens when context is cancelled
			// Round Robin currently uses the DFS algorithm to resolve the tables, but this
			// may change in the future.
			s.resolveTableDfs(ctx, table, cl, nil, resolvedResources, 1)
		}()
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()
}

func (s *syncClient) resolveTableDfs(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resolvedResources chan<- *schema.Resource, depth int) {
	clientName := client.ID()
	ctx, span := otel.Tracer(metrics.OtelName).Start(ctx,
		"sync.table."+table.Name,
		trace.WithAttributes(
			attribute.Key("sync.client.id").String(clientName),
			attribute.Key("sync.invocation.id").String(s.invocationID),
		),
	)
	defer span.End()
	logger := s.logger.With().Str("table", table.Name).Str("client", clientName).Logger()

	startTime := time.Now()
	if parent == nil { // Log only for root tables, otherwise we spam too much.
		logger.Info().Msg("top level table resolver started")
	}
	tableMetrics := s.metrics.TableClient[table.Name][clientName]
	tableMetrics.OtelStartTime(ctx, startTime)

	res := make(chan any)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
				logger.Error().Interface("error", err).Str("stack", stack).Msg("table resolver finished with panic")
				tableMetrics.OtelPanicsAdd(ctx, 1)
				atomic.AddUint64(&tableMetrics.Panics, 1)
			}
			close(res)
		}()
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			logger.Error().Err(err).Msg("table resolver finished with error")
			tableMetrics.OtelErrorsAdd(ctx, 1)
			atomic.AddUint64(&tableMetrics.Errors, 1)
			return
		}
	}()

	for r := range res {
		s.resolveResourcesDfs(ctx, table, client, parent, r, resolvedResources, depth)
	}

	// we don't need any waitgroups here because we are waiting for the channel to close
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	tableMetrics.Duration.Store(&duration)
	tableMetrics.OtelEndTime(ctx, endTime)
	if parent == nil { // Log only for root tables and relations only after resolving is done, otherwise we spam per object instead of per table.
		logger.Info().Uint64("resources", tableMetrics.Resources).Uint64("errors", tableMetrics.Errors).Dur("duration_ms", duration).Msg("table sync finished")
		s.logTablesMetrics(table.Relations, client)
	}
}

func (s *syncClient) resolveResourcesDfs(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources any, resolvedResources chan<- *schema.Resource, depth int) {
	resourcesSlice := helpers.InterfaceSlice(resources)
	if len(resourcesSlice) == 0 {
		return
	}
	resourcesChan := make(chan *schema.Resource, len(resourcesSlice))
	go func() {
		defer close(resourcesChan)
		var wg sync.WaitGroup
		for i := range resourcesSlice {
			i := i
			resourceConcurrencyKey := table.Name + "-" + client.ID() + "-" + "resource"
			resourceSemVal, _ := s.scheduler.singleTableConcurrency.LoadOrStore(resourceConcurrencyKey, semaphore.NewWeighted(s.scheduler.singleResourceMaxConcurrency))
			resourceSem := resourceSemVal.(*semaphore.Weighted)
			if err := resourceSem.Acquire(ctx, 1); err != nil {
				s.logger.Warn().Err(err).Msg("failed to acquire semaphore. context cancelled")
				// This means context was cancelled
				wg.Wait()
				// we have to continue emptying the channel to exit gracefully
				return
			}

			// Once Resource semaphore is acquired we can acquire the global semaphore
			if err := s.scheduler.resourceSem.Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				resourceSem.Release(1)
				wg.Wait()
				// we have to continue emptying the channel to exit gracefully
				return
			}
			wg.Add(1)
			go func() {
				defer resourceSem.Release(1)
				defer s.scheduler.resourceSem.Release(1)
				defer wg.Done()
				//nolint:all
				resolvedResource := resolvers.ResolveSingleResource(ctx, s.logger, s.metrics, table, client, parent, resourcesSlice[i], s.scheduler.caser)
				if resolvedResource == nil {
					return
				}

				if err := resolvedResource.CalculateCQID(s.deterministicCQID); err != nil {
					tableMetrics := s.metrics.TableClient[table.Name][client.ID()]
					s.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with primary key calculation error")
					atomic.AddUint64(&tableMetrics.Errors, 1)
					return
				}
				if err := resolvedResource.Validate(); err != nil {
					tableMetrics := s.metrics.TableClient[table.Name][client.ID()]
					s.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation error")
					atomic.AddUint64(&tableMetrics.Errors, 1)
					return
				}
				select {
				case resourcesChan <- resolvedResource:
				case <-ctx.Done():
				}
			}()
		}
		wg.Wait()
	}()

	var wg sync.WaitGroup
	for resource := range resourcesChan {
		resource := resource
		resolvedResources <- resource
		for _, relation := range resource.Table.Relations {
			relation := relation
			tableConcurrencyKey := table.Name + "-" + client.ID()
			// Acquire the semaphore for the table
			tableSemVal, _ := s.scheduler.singleTableConcurrency.LoadOrStore(tableConcurrencyKey, semaphore.NewWeighted(s.scheduler.singleNestedTableMaxConcurrency))
			tableSem := tableSemVal.(*semaphore.Weighted)
			if err := tableSem.Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				return
			}
			// Once table semaphore is acquired we can acquire the global semaphore
			if err := s.scheduler.tableSems[depth].Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				tableSem.Release(1)
				wg.Wait()
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer tableSem.Release(1)
				defer s.scheduler.tableSems[depth].Release(1)
				s.resolveTableDfs(ctx, relation, client, resource, resolvedResources, depth+1)
			}()
		}
	}
	wg.Wait()
}
