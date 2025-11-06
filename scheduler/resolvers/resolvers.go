package resolvers

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

func resolveColumn(ctx context.Context, logger zerolog.Logger, m *metrics.Metrics, selector metrics.Selector, client schema.ClientMeta, resource *schema.Resource, column schema.Column, c *caser.Caser) {
	columnStartTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Str("column", column.Name).Interface("error", err).TimeDiff("duration", time.Now(), columnStartTime).Str("stack", stack).Msg("column resolver finished with panic")
			m.AddPanics(ctx, 1, selector)
		}
	}()

	if column.Resolver != nil {
		if err := column.Resolver(ctx, client, resource, column); err != nil {
			logger.Error().Err(err).Msg("column resolver finished with error")
			m.AddErrors(ctx, 1, selector)
		}
	} else {
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.GetItem(), c.ToPascal(column.Name), funk.WithAllowZero())
		if v != nil {
			err := resource.Set(column.Name, v)
			if err != nil {
				logger.Error().Err(err).Msg("column resolver finished with error")
				m.AddErrors(ctx, 1, selector)
			}
		}
	}
}

func ResolveResourcesChunk(ctx context.Context, logger zerolog.Logger, m *metrics.Metrics, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, chunk []any, c *caser.Caser) []*schema.Resource {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	resources := make([]*schema.Resource, len(chunk))
	for i, item := range chunk {
		resources[i] = schema.NewResourceData(table, parent, item)
	}
	objectStartTime := time.Now()

	clientID := client.ID()
	tableLogger := logger.With().Str("table", table.Name).Str("client", clientID).Logger()

	selector := m.NewSelector(clientID, table.Name)

	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			tableLogger.Error().Interface("error", err).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("resource resolver finished with panic")
			m.AddPanics(ctx, 1, selector)
		}
	}()

	if table.PreResourceChunkResolver != nil {
		if err := table.PreResourceChunkResolver.RowsResolver(ctx, client, resources); err != nil {
			tableLogger.Error().Stack().Err(err).Msg("pre resource chunk resolver finished with error")
			m.AddErrors(ctx, 1, selector)
			return nil
		}
	}

	if table.PreResourceResolver != nil {
		for _, resource := range resources {
			if err := table.PreResourceResolver(ctx, client, resource); err != nil {
				tableLogger.Error().Err(err).Msg("pre resource resolver failed")
				m.AddErrors(ctx, 1, selector)
				return nil
			}
		}
	}
	for _, resource := range resources {
		for _, column := range table.Columns {
			resolveColumn(ctx, tableLogger, m, selector, client, resource, column, c)
		}
	}

	if table.PostResourceResolver != nil {
		for _, resource := range resources {
			if err := table.PostResourceResolver(ctx, client, resource); err != nil {
				tableLogger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
				m.AddErrors(ctx, 1, selector)
			}
		}
	}

	m.AddResources(ctx, int64(len(resources)), selector)
	return resources
}
