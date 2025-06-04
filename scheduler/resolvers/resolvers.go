package resolvers

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

func resolveColumn(ctx context.Context, logger zerolog.Logger, tableMetrics *metrics.TableClientMetrics, client schema.ClientMeta, resource *schema.Resource, column schema.Column, c *caser.Caser) {
	columnStartTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Str("column", column.Name).Interface("error", err).TimeDiff("duration", time.Now(), columnStartTime).Str("stack", stack).Msg("column resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
		}
	}()

	if column.Resolver != nil {
		if err := column.Resolver(ctx, client, resource, column); err != nil {
			logger.Error().Err(err).Msg("column resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
		}
	} else {
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.GetItem(), c.ToPascal(column.Name), funk.WithAllowZero())
		if v != nil {
			err := resource.Set(column.Name, v)
			if err != nil {
				logger.Error().Err(err).Msg("column resolver finished with error")
				atomic.AddUint64(&tableMetrics.Errors, 1)
			}
		}
	}
}

func ResolveSingleResource(ctx context.Context, logger zerolog.Logger, m *metrics.Metrics, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item any, c *caser.Caser) *schema.Resource {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	resource := schema.NewResourceData(table, parent, item)
	objectStartTime := time.Now()
	clientID := client.ID()
	tableMetrics := m.TableClient[table.Name][clientID]
	tableLogger := logger.With().Str("table", table.Name).Str("client", clientID).Logger()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			tableLogger.Error().Interface("error", err).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("resource resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
		}
	}()
	if table.PreResourceResolver != nil {
		if err := table.PreResourceResolver(ctx, client, resource); err != nil {
			tableLogger.Error().Err(err).Msg("pre resource resolver failed")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			return nil
		}
	}

	for _, column := range table.Columns {
		resolveColumn(ctx, tableLogger, tableMetrics, client, resource, column, c)
	}

	if table.PostResourceResolver != nil {
		if err := table.PostResourceResolver(ctx, client, resource); err != nil {
			tableLogger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
		}
	}

	m.OtelResourcesAdd(ctx, 1, tableMetrics)
	atomic.AddUint64(&tableMetrics.Resources, 1)
	return resource
}
