package queue

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/helpers"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/resolvers"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type worker struct {
	jobs              <-chan *WorkUnit
	queue             *ConcurrentRandomQueue[WorkUnit]
	resolvedResources chan<- *schema.Resource

	logger            zerolog.Logger
	caser             *caser.Caser
	invocationID      string
	deterministicCQID bool
	metrics           *metrics.Metrics
}

func (w *worker) work(ctx context.Context, activeWorkSignal *activeWorkSignal) {
	for j := range w.jobs {
		activeWorkSignal.Add()

		w.resolveTable(ctx, j.Table, j.Client, j.Parent)

		activeWorkSignal.Done()
	}
}

func newWorker(
	jobs <-chan *WorkUnit,
	queue *ConcurrentRandomQueue[WorkUnit],
	resolvedResources chan<- *schema.Resource,

	logger zerolog.Logger,
	c *caser.Caser,
	invocationID string,
	deterministicCQID bool,
	m *metrics.Metrics,
) *worker {
	return &worker{
		jobs:              jobs,
		queue:             queue,
		resolvedResources: resolvedResources,
		logger:            logger,
		caser:             c,
		deterministicCQID: deterministicCQID,
		invocationID:      invocationID,
		metrics:           m,
	}
}

func (w *worker) resolveTable(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource) {
	clientName := client.ID()
	ctx, span := otel.Tracer(metrics.OtelName).Start(ctx,
		"sync.table."+table.Name,
		trace.WithAttributes(
			attribute.Key("sync.client.id").String(clientName),
			attribute.Key("sync.invocation.id").String(w.invocationID),
		),
	)
	defer span.End()
	logger := w.logger.With().Str("table", table.Name).Str("client", clientName).Logger()
	startTime := time.Now()
	if parent == nil { // Log only for root tables, otherwise we spam too much.
		logger.Info().Msg("top level table resolver started")
	}
	tableMetrics := w.metrics.TableClient[table.Name][clientName]
	defer func() {
		span.AddEvent("sync.finish.stats", trace.WithAttributes(
			attribute.Key("sync.resources").Int64(int64(atomic.LoadUint64(&tableMetrics.Resources))),
			attribute.Key("sync.errors").Int64(int64(atomic.LoadUint64(&tableMetrics.Errors))),
			attribute.Key("sync.panics").Int64(int64(atomic.LoadUint64(&tableMetrics.Panics))),
		))
	}()
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
		w.resolveResource(ctx, table, client, parent, r)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	tableMetrics.Duration.Store(&duration)
	tableMetrics.OtelEndTime(ctx, endTime)
	if parent == nil {
		logger.Info().Uint64("resources", tableMetrics.Resources).Uint64("errors", tableMetrics.Errors).Dur("duration_ms", duration).Msg("table sync finished")
	}
}

func (w *worker) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources any) {
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
			wg.Add(1)
			go func() {
				defer wg.Done()
				resolvedResource := resolvers.ResolveSingleResource(ctx, w.logger, w.metrics, table, client, parent, resourcesSlice[i], w.caser)
				if resolvedResource == nil {
					return
				}

				if err := resolvedResource.CalculateCQID(w.deterministicCQID); err != nil {
					tableMetrics := w.metrics.TableClient[table.Name][client.ID()]
					w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with primary key calculation error")
					atomic.AddUint64(&tableMetrics.Errors, 1)
					return
				}
				if err := resolvedResource.Validate(); err != nil {
					switch err.(type) {
					case *schema.PKError:
						tableMetrics := w.metrics.TableClient[table.Name][client.ID()]
						w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation error")
						atomic.AddUint64(&tableMetrics.Errors, 1)
						return
					case *schema.PKComponentError:
						w.logger.Warn().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation warning")
					}
				}
				select {
				case resourcesChan <- resolvedResource:
				case <-ctx.Done():
				}
			}()
		}
		wg.Wait()
	}()

	for resource := range resourcesChan {
		resource := resource
		w.resolvedResources <- resource
		for _, r := range resource.Table.Relations {
			relation := r
			w.queue.Push(WorkUnit{
				Table:  relation,
				Client: client,
				Parent: resource,
			})
		}
	}
}
