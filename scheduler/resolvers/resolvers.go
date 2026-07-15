package resolvers

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

func resolveColumn(ctx context.Context, logger zerolog.Logger, m *metrics.Metrics, selector metrics.Selector, client schema.ClientMeta, resource *schema.Resource, column schema.Column, c *caser.Caser, classifier schema.ErrorClassifier) {
	columnStartTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Str("column", column.Name).Interface("error", err).TimeDiff("duration", time.Now(), columnStartTime).Str("stack", stack).Msg("column resolver finished with panic")
			m.AddPanics(ctx, 1, selector)
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("table", resource.Table.Name)
				scope.SetTag("column", column.Name)
				sentry.CurrentHub().CaptureMessage(stack)
			})
		}
	}()

	handleErr := func(err error) {
		event := schema.ErrorEvent{Table: resource.Table, Client: client, Phase: schema.ErrorPhaseColumnResolver, Column: &column}
		if classifier.Suppress(ctx, err, event) {
			logger.Debug().Str("column", column.Name).Err(err).Msg("column resolver finished with suppressed error")
			return
		}
		logger.Error().Err(err).Msg("column resolver finished with error")
		m.AddErrors(ctx, 1, selector)
	}

	if column.Resolver != nil {
		if err := column.Resolver(ctx, client, resource, column); err != nil {
			handleErr(err)
		}
	} else {
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.GetItem(), c.ToPascal(column.Name), funk.WithAllowZero())
		if v != nil {
			if err := resource.Set(column.Name, v); err != nil {
				handleErr(err)
			}
		}
	}
}

func ResolveResourcesChunk(ctx context.Context, logger zerolog.Logger, m *metrics.Metrics, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, chunk []any, c *caser.Caser, classifier schema.ErrorClassifier) []*schema.Resource {
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
			event := schema.ErrorEvent{Table: table, Client: client, Phase: schema.ErrorPhasePreResourceChunkResolver}
			if classifier.Suppress(ctx, err, event) {
				tableLogger.Debug().Err(err).Msg("pre resource chunk resolver finished with suppressed error")
			} else {
				tableLogger.Error().Stack().Err(err).Msg("pre resource chunk resolver finished with error")
				m.AddErrors(ctx, 1, selector)
			}
			return nil
		}
	}

	if table.PreResourceResolver != nil {
		filtered := resources[:0]
		for _, resource := range resources {
			if err := table.PreResourceResolver(ctx, client, resource); err != nil {
				event := schema.ErrorEvent{Table: table, Client: client, Phase: schema.ErrorPhasePreResourceResolver}
				suppress := classifier.Suppress(ctx, err, event)
				switch {
				case suppress && ctx.Err() != nil:
					tableLogger.Debug().Err(err).Msg("pre resource resolver failed, context cancelled (suppressed)")
					return nil
				case suppress:
					tableLogger.Debug().Err(err).Msg("pre resource resolver failed (suppressed)")
					continue
				case ctx.Err() != nil:
					tableLogger.Error().Err(err).Msg("pre resource resolver failed, context cancelled")
					m.AddErrors(ctx, 1, selector)
					return nil
				default:
					tableLogger.Error().Err(err).Msg("pre resource resolver failed")
					m.AddErrors(ctx, 1, selector)
					continue
				}
			}
			filtered = append(filtered, resource)
		}
		resources = filtered
	}
	for _, resource := range resources {
		for _, column := range table.Columns {
			resolveColumn(ctx, tableLogger, m, selector, client, resource, column, c, classifier)
		}
	}

	if table.PostResourceResolver != nil {
		for _, resource := range resources {
			if err := table.PostResourceResolver(ctx, client, resource); err != nil {
				event := schema.ErrorEvent{Table: table, Client: client, Phase: schema.ErrorPhasePostResourceResolver}
				if classifier.Suppress(ctx, err, event) {
					tableLogger.Debug().Err(err).Msg("post resource resolver finished with suppressed error")
				} else {
					tableLogger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
					m.AddErrors(ctx, 1, selector)
				}
			}
		}
	}

	m.AddResources(ctx, int64(len(resources)), selector)
	return resources
}
