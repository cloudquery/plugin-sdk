package source

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/cloudquery/plugin-sdk/helpers"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/getsentry/sentry-go"
	"golang.org/x/sync/semaphore"
)

func (p *Plugin) syncDfs(ctx context.Context, spec specs.Source, client schema.ClientMeta, tables schema.Tables, resolvedResources chan<- *schema.Resource) {
	// This is very similar to the concurrent web crawler problem with some minor changes.
	// We are using DFS to make sure memory usage is capped at O(h) where h is the height of the tree.
	tableConcurrency := max(spec.Concurrency/minResourceConcurrency, minTableConcurrency)
	resourceConcurrency := tableConcurrency * minResourceConcurrency

	p.tableSems = make([]*semaphore.Weighted, p.maxDepth)
	for i := uint64(0); i < p.maxDepth; i++ {
		p.tableSems[i] = semaphore.NewWeighted(int64(tableConcurrency))
		// reduce table concurrency logarithmically for every depth level
		tableConcurrency = max(tableConcurrency/2, minTableConcurrency)
	}
	p.resourceSem = semaphore.NewWeighted(int64(resourceConcurrency))

	// we have this because plugins can return sometimes clients in a random way which will cause
	// differences between this run and the next one.
	preInitialisedClients := make([][]schema.ClientMeta, len(tables))
	for i, table := range tables {
		clients := []schema.ClientMeta{client}
		if table.Multiplex != nil {
			clients = table.Multiplex(client)
		}
		preInitialisedClients[i] = clients
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		p.metrics.initWithClients(table, clients)
	}

	// We start a goroutine that logs the metrics periodically.
	// It needs its own waitgroup
	var logWg sync.WaitGroup
	logWg.Add(1)

	logCtx, logCancel := context.WithCancel(ctx)
	go p.periodicMetricLogger(logCtx, &logWg)

	var wg sync.WaitGroup
	for i, table := range tables {
		table := table
		clients := preInitialisedClients[i]
		for _, client := range clients {
			client := client
			if err := p.tableSems[0].Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				// gracefully shut down the logger goroutine
				logCancel()
				logWg.Wait()
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer p.tableSems[0].Release(1)
				// not checking for error here as nothing much todo.
				// the error is logged and this happens when context is cancelled
				p.resolveTableDfs(ctx, table, client, nil, resolvedResources, 1)
			}()
		}
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()

	// gracefully shut down the logger goroutine
	logCancel()
	logWg.Wait()
}

func (p *Plugin) resolveTableDfs(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resolvedResources chan<- *schema.Resource, depth int) {
	var validationErr *schema.ValidationError
	clientName := client.ID()
	logger := p.logger.With().Str("table", table.Name).Str("client", clientName).Logger()

	if parent == nil { // Log only for root tables, otherwise we spam too much.
		logger.Info().Msg("top level table resolver started")
	}
	tableMetrics := p.metrics.TableClient[table.Name][clientName]

	res := make(chan any)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(stack)
				})
				logger.Error().Interface("error", err).Str("stack", stack).Msg("table resolver finished with panic")
				atomic.AddUint64(&tableMetrics.Panics, 1)
			}
			close(res)
		}()
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			logger.Error().Err(err).Msg("table resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
			return
		}
	}()

	for r := range res {
		p.resolveResourcesDfs(ctx, table, client, parent, r, resolvedResources, depth)
	}

	// we don't need any waitgroups here because we are waiting for the channel to close
	if parent == nil { // Log only for root tables and relations only after resolving is done, otherwise we spam per object instead of per table.
		logger.Info().Uint64("resources", tableMetrics.Resources).Uint64("errors", tableMetrics.Errors).Msg("table sync finished")
		p.logTablesMetrics(table.Relations, client)
	}
}

func (p *Plugin) resolveResourcesDfs(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources any, resolvedResources chan<- *schema.Resource, depth int) {
	resourcesSlice := helpers.InterfaceSlice(resources)
	if len(resourcesSlice) == 0 {
		return
	}
	resourcesChan := make(chan *schema.Resource, len(resourcesSlice))
	go func() {
		defer close(resourcesChan)
		var wg sync.WaitGroup
		sentValidationErrors := sync.Map{}
		for i := range resourcesSlice {
			i := i
			if err := p.resourceSem.Acquire(ctx, 1); err != nil {
				p.logger.Warn().Err(err).Msg("failed to acquire semaphore. context cancelled")
				wg.Wait()
				// we have to continue emptying the channel to exit gracefully
				return
			}
			wg.Add(1)
			go func() {
				defer p.resourceSem.Release(1)
				defer wg.Done()
				//nolint:all
				resolvedResource := p.resolveResource(ctx, table, client, parent, resourcesSlice[i])
				if resolvedResource == nil {
					return
				}

				if err := resolvedResource.CalculateCQID(p.spec.DeterministicCQID); err != nil {
					tableMetrics := p.metrics.TableClient[table.Name][client.ID()]
					p.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with primary key calculation error")
					if _, found := sentValidationErrors.LoadOrStore(table.Name, struct{}{}); !found {
						// send resource validation errors to Sentry only once per table,
						// to avoid sending too many duplicate messages
						sentry.WithScope(func(scope *sentry.Scope) {
							scope.SetTag("table", table.Name)
							sentry.CurrentHub().CaptureMessage(err.Error())
						})
					}
					atomic.AddUint64(&tableMetrics.Errors, 1)
					return
				}
				if err := resolvedResource.Validate(); err != nil {
					tableMetrics := p.metrics.TableClient[table.Name][client.ID()]
					p.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation error")
					if _, found := sentValidationErrors.LoadOrStore(table.Name, struct{}{}); !found {
						// send resource validation errors to Sentry only once per table,
						// to avoid sending too many duplicate messages
						sentry.WithScope(func(scope *sentry.Scope) {
							scope.SetTag("table", table.Name)
							sentry.CurrentHub().CaptureMessage(err.Error())
						})
					}
					atomic.AddUint64(&tableMetrics.Errors, 1)
					return
				}
				resourcesChan <- resolvedResource
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
			if err := p.tableSems[depth].Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer p.tableSems[depth].Release(1)
				p.resolveTableDfs(ctx, relation, client, resource, resolvedResources, depth+1)
			}()
		}
	}
	wg.Wait()
}
