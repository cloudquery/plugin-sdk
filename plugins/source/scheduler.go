package source

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

const (
	minTableConcurrency    = 1
	minResourceConcurrency = 100
)

const periodicMetricLoggerInterval = 30 * time.Second

func (p *Plugin) logTablesMetrics(tables schema.Tables, client schema.ClientMeta) {
	clientName := client.ID()
	for _, table := range tables {
		metrics := p.metrics.TableClient[table.Name][clientName]
		p.logger.Info().Str("table", table.Name).Str("client", clientName).Uint64("resources", metrics.Resources).Uint64("errors", metrics.Errors).Msg("table sync finished")
		p.logTablesMetrics(table.Relations, client)
	}
}

func (p *Plugin) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item any) *schema.Resource {
	var validationErr *schema.ValidationError
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
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("table", table.Name)
				sentry.CurrentHub().CaptureMessage(stack)
			})
		}
	}()
	if table.PreResourceResolver != nil {
		if err := table.PreResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Err(err).Msg("pre resource resolver failed")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
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
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
		}
	}
	atomic.AddUint64(&tableMetrics.Resources, 1)
	return resource
}

func (p *Plugin) resolveColumn(ctx context.Context, logger zerolog.Logger, tableMetrics *TableClientMetrics, client schema.ClientMeta, resource *schema.Resource, c schema.Column) {
	var validationErr *schema.ValidationError
	columnStartTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Str("column", c.Name).Interface("error", err).TimeDiff("duration", time.Now(), columnStartTime).Str("stack", stack).Msg("column resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("table", resource.Table.Name)
				scope.SetTag("column", c.Name)
				sentry.CurrentHub().CaptureMessage(stack)
			})
		}
	}()

	if c.Resolver != nil {
		if err := c.Resolver(ctx, client, resource, c); err != nil {
			logger.Error().Err(err).Msg("column resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", resource.Table.Name)
					scope.SetTag("column", c.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
		}
	} else {
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.GetItem(), p.caser.ToPascal(c.Name), funk.WithAllowZero())
		if v != nil {
			err := resource.Set(c.Name, v)
			if err != nil {
				logger.Error().Err(err).Msg("column resolver finished with error")
				atomic.AddUint64(&tableMetrics.Errors, 1)
				if errors.As(err, &validationErr) {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("table", resource.Table.Name)
						scope.SetTag("column", c.Name)
						sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
					})
				}
			}
		}
	}
}

func (p *Plugin) periodicMetricLogger(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(periodicMetricLoggerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.logger.Info().
				Uint64("total_resources", p.metrics.TotalResourcesAtomic()).
				Uint64("total_errors", p.metrics.TotalErrorsAtomic()).
				Uint64("total_panics", p.metrics.TotalPanicsAtomic()).
				Msg("Sync in progress")
		}
	}
}

// unparam's suggestion to remove the second parameter is not good advice here.
// nolint:unparam
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
