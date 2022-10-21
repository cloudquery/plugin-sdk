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
	"github.com/getsentry/sentry-go"
	"github.com/thoas/go-funk"
)

func (p *SourcePlugin) resolveTable(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resolvedResources chan<- *schema.Resource) error {
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
		p.logger.Debug().Msg("table resolver started")
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			p.logger.Error().Err(err).Msg("table resolver finished with error")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
			return
		}
		p.logger.Debug().Msg("table resolver finished successfully")
	}()

	resolvedObjects := make(chan *schema.Resource)
	go func() {
		defer close(resolvedObjects)
		for r := range res {
			p.resolveResources(ctx, table, client, parent, r, resolvedObjects)
		}
	}()

	for resolvedObject := range resolvedObjects {
		resolvedResources <- resolvedObject
		for _, rel := range table.Relations {
			if err := p.resolveTable(ctx, rel, client, resolvedObject, resolvedResources); err != nil {
				break
			}
		}
	}
	// we don't need any waitgroups here because we are waiting for the channel to close
	p.logger.Info().Msg("fetch table finished")
	return nil
}

func (p *SourcePlugin) resolveResources(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources interface{}, resolvedResources chan<- *schema.Resource) {
	resourcesSlice := helpers.InterfaceSlice(resources)
	if len(resourcesSlice) == 0 {
		return
	}
	wg := &sync.WaitGroup{}
	for i := range resourcesSlice {
		i := i
		if err := p.resourceSem.Acquire(ctx, 1); err != nil {
			p.logger.Error().Err(err).Msg("failed to acquire semaphore. context cancelled")
			// we have to continue emptying the channel to exit gracefully
			break
		}
		wg.Add(1)
		go func() {
			defer p.resourceSem.Release(1)
			defer wg.Done()
			//nolint:all
			resolvedResource := p.resolveResource(ctx, table, client, parent, resourcesSlice[i])
			if resolvedResource != nil {
				resolvedResources <- resolvedResource
			}
		}()
	}

	wg.Wait()
}

func (p *SourcePlugin) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item interface{}) *schema.Resource {
	resource := schema.NewResourceData(table, parent, item)
	objectStartTime := time.Now()
	csr := caser.New()
	clientName := client.Name()
	logger := p.logger.With().Str("table", table.Name).Str("client", clientName).Logger()
	logger.Info().Msg("object resolver started")
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			// sentry.WithScope(func(scope *sentry.Scope) {
			// 	scope.SetTag("table", table.Name)
			// 	sentry.CurrentHub().CaptureMessage(stack)
			// })
			p.logger.Error().Interface("error", err).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("object resolver finished with panic")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Panics, 1)
		}
		
	}()
	if table.PreResourceResolver != nil {
		logger.Trace().Msg("pre resource resolver started")
		if err := table.PreResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Err(err).Msg("pre resource resolver failed")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
			return nil
		}
		logger.Trace().Msg("pre resource resolver finished successfully")
	}

	for _, c := range table.Columns {
		cl := logger.With().Str("column", c.Name).Logger()
		if c.Resolver != nil {
			cl.Trace().Msg("column resolver custom started")
			if err := c.Resolver(ctx, client, resource, c); err != nil {
				cl.Error().Err(err).Msg("column resolver finished with error")
				atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
			}
			cl.Trace().Msg("column resolver finished successfully")
		} else {
			cl.Trace().Msg("column resolver default started")
			// base use case: try to get column with CamelCase name
			v := funk.Get(resource.GetItem(), csr.ToPascal(c.Name), funk.WithAllowZero())
			if v != nil {
				resource.Set(c.Name, v)
				cl.Trace().Msg("column resolver default finished successfully")
			} else {
				cl.Trace().Msg("column resolver default finished successfully with nil")
			}
		}
	}

	if table.PostResourceResolver != nil {
		logger.Trace().Msg("post resource resolver started")
		if err := table.PostResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
			atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Errors, 1)
		} else {
			logger.Trace().Msg("post resource resolver finished successfully")
		}
	}
	atomic.AddUint64(&p.stats.TableClient[table.Name][clientName].Resources, 1)
	return resource
}
