package plugins

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/helpers"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/semaphore"
)

const (
	minTableConcurrency    = 1
	minResourceConcurrency = 100
)

func (p *SourcePlugin) syncDfs(ctx context.Context, spec specs.Source, client schema.ClientMeta, tables schema.Tables, resolvedResources chan<- *schema.Resource) {
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

	for _, table := range tables {
		if table.Parent != nil {
			// skip descendent tables here - they are handled in initWithClients
			continue
		}
		clients := []schema.ClientMeta{client}
		if table.Multiplex != nil {
			clients = table.Multiplex(client)
		}
		// we do this here to avoid locks so we initial the metrics structure once in the main goroutines
		// and then we can just read from it in the other goroutines concurrently given we are not writing to it.
		p.metrics.initWithClients(table, clients)
	}

	var wg sync.WaitGroup
	for _, table := range tables {
		if table.Parent != nil {
			// skip descendent tables here - they will be handled by the recursive
			// depth-first-search later.
			continue
		}
		table := table
		clients := []schema.ClientMeta{client}
		if table.Multiplex != nil {
			clients = table.Multiplex(client)
		}
		for _, client := range clients {
			client := client
			if err := p.tableSems[0].Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer p.tableSems[0].Release(1)
				// not checking for error here as nothing much todo.
				// the error is logged and this happens when context is cancelled
				p.resolveTableDfs(ctx, tables, table, client, nil, resolvedResources, 0)
			}()
		}
	}
	wg.Wait()
}

func (p *SourcePlugin) resolveTableDfs(ctx context.Context, allIncludedTables schema.Tables, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resolvedResources chan<- *schema.Resource, depth int) {
	clientName := client.ID()
	logger := p.logger.With().Str("table", table.Name).Str("client", clientName).Logger()

	if parent == nil { // Log only for root tables, otherwise we spam too much.
		logger.Info().Msg("table resolver started")
	}

	tableMetrics := p.metrics.TableClient[table.Name][clientName]
	res := make(chan interface{})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(stack)
				})
				p.logger.Error().Interface("error", err).Str("stack", stack).Msg("table resolver finished with panic")
				atomic.AddUint64(&tableMetrics.Panics, 1)
			}
		}()
		defer close(res)
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			logger.Error().Err(err).Msg("table resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			return
		}
	}()

	for r := range res {
		p.resolveResourcesDfs(ctx, allIncludedTables, table, client, parent, r, resolvedResources, depth)
	}

	// we don't need any waitgroups here because we are waiting for the channel to close
	if parent == nil { // Log only for root tables, otherwise we spam too much.
		logger.Info().Uint64("resources", tableMetrics.Resources).Uint64("errors", tableMetrics.Errors).Msg("fetch table finished")
	}
}

func (p *SourcePlugin) resolveResourcesDfs(ctx context.Context, allIncludedTables schema.Tables, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources interface{}, resolvedResources chan<- *schema.Resource, depth int) {
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
			if allIncludedTables.GetTopLevel(relation.Name) == nil {
				// this indicates that child table is skipped by user config,
				// so we should not sync it
				continue
			}
			if err := p.tableSems[depth].Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer p.tableSems[depth].Release(1)
				p.resolveTableDfs(ctx, allIncludedTables, relation, client, resource, resolvedResources, depth+1)
			}()
		}
	}
	wg.Wait()
}

func (p *SourcePlugin) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item interface{}) *schema.Resource {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	resource := schema.NewResourceData(table, parent, item)
	objectStartTime := time.Now()
	clientID := client.ID()
	tableMetrics := p.metrics.TableClient[table.Name][clientID]
	logger := p.logger.With().Str("table", table.Name).Str("client", clientID).Logger()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Interface("error", err).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("resource resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
		}
	}()
	if table.PreResourceResolver != nil {
		if err := table.PreResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Err(err).Msg("pre resource resolver failed")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			return nil
		}
	}

	for _, c := range table.Columns {
		p.resolveColumn(ctx, logger, tableMetrics, client, resource, c)
	}

	if table.PostResourceResolver != nil {
		if err := table.PostResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
		}
	}
	atomic.AddUint64(&tableMetrics.Resources, 1)
	return resource
}

func (p *SourcePlugin) resolveColumn(ctx context.Context, logger zerolog.Logger, tableMetrics *TableClientMetrics, client schema.ClientMeta, resource *schema.Resource, c schema.Column) {
	columnStartTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Str("column", c.Name).Interface("error", err).TimeDiff("duration", time.Now(), columnStartTime).Str("stack", stack).Msg("column resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
		}
	}()

	if c.Resolver != nil {
		if err := c.Resolver(ctx, client, resource, c); err != nil {
			logger.Error().Err(err).Msg("column resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
		}
	} else {
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.GetItem(), p.caser.ToPascal(c.Name), funk.WithAllowZero())
		if v != nil {
			_ = resource.Set(c.Name, v)
		}
	}
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
