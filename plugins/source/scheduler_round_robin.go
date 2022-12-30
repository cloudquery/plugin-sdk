package source

import (
	"context"
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

func (p *Plugin) syncRoundRobin(ctx context.Context, spec specs.Source, client schema.ClientMeta, tables schema.Tables, resolvedResources chan<- *schema.Resource) {
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
	type tableClient struct {
		table  *schema.Table
		client schema.ClientMeta
	}

	// interleave table-clients so that we get:
	// table1-client1, table2-client1, table3-client, table1-client2, table2-client2, table3-client2, ...
	tableClients := make([]tableClient, 0)
	c := 0
	for {
		addedNew := false
		for i, table := range tables {
			if c < len(preInitialisedClients[i]) {
				tableClients = append(tableClients, tableClient{table: table, client: preInitialisedClients[i][c]})
				addedNew = true
			}
		}
		c++
		if !addedNew {
			break
		}
	}

	for _, tc := range tableClients {
		table := tc.table
		cl := tc.client
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
			p.resolveTableRoundRobin(ctx, table, cl, nil, resolvedResources, 1)
		}()
	}

	// Wait for all the worker goroutines to finish
	wg.Wait()

	// gracefully shut down the logger goroutine
	logCancel()
	logWg.Wait()
}

func (p *Plugin) resolveTableRoundRobin(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resolvedResources chan<- *schema.Resource, depth int) {
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
			return
		}
	}()

	for r := range res {
		p.resolveResourcesRoundRobin(ctx, table, client, parent, r, resolvedResources, depth)
	}

	// we don't need any waitgroups here because we are waiting for the channel to close
	if parent == nil { // Log only for root tables and relations only after resolving is done, otherwise we spam per object instead of per table.
		logger.Info().Uint64("resources", tableMetrics.Resources).Uint64("errors", tableMetrics.Errors).Msg("table sync finished")
		p.logTablesMetrics(table.Relations, client)
	}
}

func (p *Plugin) resolveResourcesRoundRobin(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources any, resolvedResources chan<- *schema.Resource, depth int) {
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
				p.resolveTableRoundRobin(ctx, relation, client, resource, resolvedResources, depth+1)
			}()
		}
	}
	wg.Wait()
}
