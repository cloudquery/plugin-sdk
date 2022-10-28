package plugins

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/caser"
	"github.com/cloudquery/plugin-sdk/helpers"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/getsentry/sentry-go"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/semaphore"
)

const (
	minTableConcurrency    = 1
	minResourceConcurrency = 100
)

func (p *SourcePlugin) syncDfs(ctx context.Context, spec specs.Source, client schema.ClientMeta, tables schema.Tables, resolvedResources chan<- *schema.Resource) {
	// current DFS supports only parallization for top level tables and resources.
	// it is possible to extend support for multiple levels but this require benchmarking to find a good fit on how to split
	// gorourtines for each level efficiently.
	// This is very similar to the concurrent web crawler problem with some minor changes.
	// We are using DFS to make sure memory usage is capped at O(h) where h is the height of the tree/subchilds.
	tableConcurrency := spec.Concurrency / minResourceConcurrency
	if tableConcurrency < minTableConcurrency {
		tableConcurrency = minTableConcurrency
	}
	resourceConcurrency := tableConcurrency * minResourceConcurrency

	p.tableSem = semaphore.NewWeighted(int64(tableConcurrency))
	p.resourceSem = semaphore.NewWeighted(int64(resourceConcurrency))

	var wg sync.WaitGroup
	for _, table := range tables {
		table := table
		clients := []schema.ClientMeta{client}
		if table.Multiplex != nil {
			clients = table.Multiplex(client)
		}
		p.stats.initWithClients(table, clients)
		for _, client := range clients {
			client := client
			if err := p.tableSem.Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer p.tableSem.Release(1)
				// not checking for error here as nothing much todo.
				// the error is logged and this happens when context is cancelled
				p.resolveTableDfs(ctx, table, client, nil, resolvedResources)
			}()
		}
	}
	wg.Wait()
}

func (p *SourcePlugin) resolveTableDfs(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resolvedResources chan<- *schema.Resource) {
	clientName := client.Name()
	logger := p.logger.With().Str("table", table.Name).Str("client", clientName).Logger()
	logger.Info().Msg("table resolver started")

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
				atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Panics, 1)
			}
			close(res)
		}()
		logger.Debug().Msg("table resolver started")
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			logger.Error().Err(err).Msg("table resolver finished with error")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
			return
		}
		logger.Debug().Msg("table resolver finished successfully")
	}()

	for r := range res {
		p.resolveResourcesDfs(ctx, table, client, parent, r, resolvedResources)
	}

	// we don't need any waitgroups here because we are waiting for the channel to close
	logger.Info().Msg("fetch table finished")
}

func (p *SourcePlugin) resolveResourcesDfs(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources interface{}, resolvedResources chan<- *schema.Resource) {
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

	for resource := range resourcesChan {
		resolvedResources <- resource
		if resource.Table.Relations == nil {
			continue
		}
		for _, relation := range resource.Table.Relations {
			p.resolveTableDfs(ctx, relation, client, resource, resolvedResources)
		}
	}
}

func (p *SourcePlugin) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item interface{}) *schema.Resource {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	resource := schema.NewResourceData(table, parent, item)
	objectStartTime := time.Now()
	csr := caser.New()
	clientName := client.Name()
	logger := p.logger.With().Str("table", table.Name).Str("client", clientName).Logger()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Interface("error", err).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("object resolver finished with panic")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Panics, 1)
		}
	}()
	if table.PreResourceResolver != nil {
		if err := table.PreResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Err(err).Msg("pre resource resolver failed")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
			return nil
		}
	}

	for _, c := range table.Columns {
		if c.Resolver != nil {
			if err := c.Resolver(ctx, client, resource, c); err != nil {
				logger.Error().Err(err).Msg("column resolver finished with error")
				atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
			}
		} else {
			// base use case: try to get column with CamelCase name
			v := funk.Get(resource.GetItem(), csr.ToPascal(c.Name), funk.WithAllowZero())
			if v != nil {
				_ = resource.Set(c.Name, v)
			}
		}
	}

	if table.PostResourceResolver != nil {
		if err := table.PostResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
		}
	}
	atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Resources, 1)
	return resource
}
