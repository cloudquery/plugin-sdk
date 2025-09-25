package queue

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/helpers"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/resolvers"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
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
	// message channel for sending SyncError messages
	msgChan chan<- message.SyncMessage
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
	msgChan chan<- message.SyncMessage,
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
		msgChan:           msgChan,
	}
}

func (w *worker) resolveTable(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource) {
	clientName := client.ID()
	ctx, span := otel.Tracer(metrics.ResourceName).Start(ctx,
		"sync.table."+table.Name,
		trace.WithAttributes(
			attribute.Key("sync.client.id").String(clientName),
			attribute.Key("sync.invocation.id").String(w.invocationID),
		),
	)
	defer span.End()
	logger := w.logger.With().Str("table", table.Name).Str("client", clientName).Logger()
	ctx = logger.WithContext(ctx)
	startTime := time.Now()
	if parent == nil { // Log only for root tables, otherwise we spam too much.
		logger.Info().Msg("top level table resolver started")
	}

	selector := w.metrics.NewSelector(clientName, table.Name)
	defer func() {
		span.AddEvent("sync.finish.stats", trace.WithAttributes(
			attribute.Key("sync.resources").Int64(int64(w.metrics.GetResources(selector))),
			attribute.Key("sync.errors").Int64(int64(w.metrics.GetErrors(selector))),
			attribute.Key("sync.panics").Int64(int64(w.metrics.GetPanics(selector))),
		))
	}()
	w.metrics.StartTime(startTime, selector)

	res := make(chan any)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
				logger.Error().Interface("error", err).Str("stack", stack).Msg("table resolver finished with panic")
				w.metrics.AddPanics(ctx, 1, selector)
			}
			close(res)
		}()
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			logger.Error().Err(err).Msg("table resolver finished with error")
			w.metrics.AddErrors(ctx, 1, selector)
			// Send SyncError message
			syncErrorMsg := &message.SyncError{
				TableName: table.Name,
				Error:     err.Error(),
			}
			w.msgChan <- syncErrorMsg
			return
		}
	}()

	for r := range res {
		w.resolveResource(ctx, table, client, parent, r)
	}

	endTime := time.Now()
	w.metrics.EndTime(ctx, endTime, selector)
	if parent == nil {
		logger.Info().Uint64("resources", w.metrics.GetResources(selector)).Uint64("errors", w.metrics.GetErrors(selector)).Dur("duration_ms", w.metrics.GetDuration(selector)).Msg("table sync finished")
	}
}

func (w *worker) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, resources any) {
	resourcesSlice := helpers.InterfaceSlice(resources)
	if len(resourcesSlice) == 0 {
		return
	}

	selector := w.metrics.NewSelector(client.ID(), table.Name)
	resourcesChan := make(chan *schema.Resource, len(resourcesSlice))
	go func() {
		defer close(resourcesChan)
		var wg sync.WaitGroup
		chunks := [][]any{resourcesSlice}
		if table.PreResourceChunkResolver != nil {
			chunks = lo.Chunk(resourcesSlice, table.PreResourceChunkResolver.ChunkSize)
		}
		for i := range chunks {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resolvedResources := resolvers.ResolveResourcesChunk(ctx, w.logger, w.metrics, table, client, parent, chunks[i], w.caser)
				for _, resolvedResource := range resolvedResources {
					if err := resolvedResource.CalculateCQID(w.deterministicCQID); err != nil {
						w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with primary key calculation error")
						w.metrics.AddErrors(ctx, 1, selector)
						return
					}
					if err := resolvedResource.StoreCQClientID(client.ID()); err != nil {
						w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("failed to store _cq_client_id")
					}
					if err := resolvedResource.Validate(); err != nil {
						switch err.(type) {
						case *schema.PKError:
							w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation error")
							w.metrics.AddErrors(ctx, 1, selector)
							return
						case *schema.PKComponentError:
							w.logger.Warn().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation warning")
						}
					}
					select {
					case resourcesChan <- resolvedResource:
					case <-ctx.Done():
					}
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
