package queue

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/helpers"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/resolvers"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type worker struct {
	jobs              <-chan *storage.SerializedWorkUnit
	store             storage.Storage
	codec             *Codec
	lookups           *workerLookups
	resolvedResources chan<- *schema.Resource

	logger            zerolog.Logger
	caser             *caser.Caser
	invocationID      string
	deterministicCQID bool
	metrics           *metrics.Metrics
	msgChan           chan<- message.SyncMessage
}

func newWorker(
	jobs <-chan *storage.SerializedWorkUnit,
	store storage.Storage,
	codec *Codec,
	lookups *workerLookups,
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
		store:             store,
		codec:             codec,
		lookups:           lookups,
		resolvedResources: resolvedResources,
		logger:            logger,
		caser:             c,
		deterministicCQID: deterministicCQID,
		invocationID:      invocationID,
		metrics:           m,
		msgChan:           msgChan,
	}
}

func (w *worker) work(ctx context.Context, activeWorkSignal *activeWorkSignal) {
	for j := range w.jobs {
		activeWorkSignal.Add()
		w.runJob(ctx, j)
		activeWorkSignal.Done()
	}
}

// runJob processes a single SerializedWorkUnit. Guarantees: on any return
// path (success, error, panic), if j.ParentID != "" AND the unit did not
// transfer its pin to a stored intermediate, exactly one
// DecResourceRefcount call is made.
func (w *worker) runJob(ctx context.Context, j *storage.SerializedWorkUnit) {
	pinTransferred := false
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error().Interface("panic", r).Str("table", j.TableName).Msg("worker panic")
		}
		if j.ParentID != "" && !pinTransferred {
			if err := w.store.DecResourceRefcount(ctx, j.ParentID); err != nil {
				w.logger.Error().Err(err).Str("parent_id", j.ParentID).Msg("failed to dec parent refcount")
			}
		}
	}()

	table, ok := w.lookups.tables[j.TableName]
	if !ok {
		w.logger.Error().Str("table", j.TableName).Msg("unknown table in work unit")
		return
	}
	client, ok := w.lookups.clients[j.ClientID]
	if !ok {
		w.logger.Error().Str("client", j.ClientID).Msg("unknown client in work unit")
		return
	}

	var parent *schema.Resource
	if j.ParentID != "" {
		blob, err := w.store.GetResource(ctx, j.ParentID)
		if err != nil {
			w.logger.Error().Err(err).Str("parent_id", j.ParentID).Msg("failed to load parent resource")
			return
		}
		parent, _, err = w.codec.DecodeResource(blob)
		if err != nil {
			w.logger.Error().Err(err).Str("parent_id", j.ParentID).Msg("failed to decode parent resource")
			return
		}
	}

	transferred := w.resolveTable(ctx, table, client, parent, j.ParentID)
	pinTransferred = transferred
}

// resolveTable resolves a single table+client+parent unit. Returns true if
// the WorkUnit's pin on j.ParentID was transferred to one or more stored
// intermediate resources (so the caller must NOT decrement).
func (w *worker) resolveTable(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, parentID string) (pinTransferred bool) {
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
	if parent == nil {
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
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(stack)
				})
			}
			close(res)
		}()
		if err := table.Resolver(ctx, client, parent, res); err != nil {
			logger.Error().Err(err).Msg("table resolver finished with error")
			w.metrics.AddErrors(ctx, 1, selector)
			w.msgChan <- &message.SyncError{TableName: table.Name, Error: err.Error()}
			return
		}
	}()

	for r := range res {
		if w.resolveResource(ctx, table, client, parent, parentID, r) {
			pinTransferred = true
		}
	}

	endTime := time.Now()
	w.metrics.EndTime(ctx, endTime, selector)
	if parent == nil {
		logger.Info().Uint64("resources", w.metrics.GetResources(selector)).Uint64("errors", w.metrics.GetErrors(selector)).Dur("duration_ms", w.metrics.GetDuration(selector)).Msg("table sync finished")
	}
	return pinTransferred
}

// resolveResource processes one chunk of items returned by a resolver.
// Returns true if at least one stored intermediate resource was created
// with ParentID == parentID, which means the caller's pin must NOT be
// released (it's been transferred to the new intermediate).
func (w *worker) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, parentID string, resources any) (pinTransferred bool) {
	resourcesSlice := helpers.InterfaceSlice(resources)
	if len(resourcesSlice) == 0 {
		return false
	}

	selector := w.metrics.NewSelector(client.ID(), table.Name)
	chunks := [][]any{resourcesSlice}
	if table.PreResourceChunkResolver != nil {
		chunks = lo.Chunk(resourcesSlice, table.PreResourceChunkResolver.ChunkSize)
	}

	for i := range chunks {
		resolved := resolvers.ResolveResourcesChunk(ctx, w.logger, w.metrics, table, client, parent, chunks[i], w.caser)
		for _, r := range resolved {
			if err := r.CalculateCQID(w.deterministicCQID); err != nil {
				w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with primary key calculation error")
				w.metrics.AddErrors(ctx, 1, selector)
				continue
			}
			if err := r.StoreCQClientID(client.ID()); err != nil {
				w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("failed to store _cq_client_id")
			}
			if err := r.Validate(); err != nil {
				switch err.(type) {
				case *schema.PKError:
					w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation error")
					w.metrics.AddErrors(ctx, 1, selector)
					continue
				case *schema.PKComponentError:
					w.logger.Warn().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation warning")
				}
			}

			// Emit to destination pipeline.
			select {
			case w.resolvedResources <- r:
			case <-ctx.Done():
				return pinTransferred
			}

			// If this resource has children, store it and push WorkUnits.
			if len(r.Table.Relations) > 0 {
				newID := uuid.NewString()
				blob, err := w.codec.EncodeResource(r, parentID)
				if err != nil {
					w.logger.Error().Err(err).Str("table", r.Table.Name).Msg("failed to encode resource")
					w.metrics.AddErrors(ctx, 1, selector)
					continue
				}
				if err := w.store.PutResource(ctx, newID, blob, len(r.Table.Relations)); err != nil {
					w.logger.Error().Err(err).Str("table", r.Table.Name).Msg("failed to persist resource")
					w.metrics.AddErrors(ctx, 1, selector)
					continue
				}
				wus := make([]storage.SerializedWorkUnit, 0, len(r.Table.Relations))
				for _, rel := range r.Table.Relations {
					wus = append(wus, storage.SerializedWorkUnit{
						TableName: rel.Name,
						ClientID:  client.ID(),
						ParentID:  newID,
					})
				}
				if err := w.store.PushWorkBatch(ctx, wus); err != nil {
					w.logger.Error().Err(err).Msg("failed to push child work units")
					w.metrics.AddErrors(ctx, 1, selector)
					continue
				}

				// Pin transfer: since this intermediate references parentID,
				// the WorkUnit's pin must NOT be decremented on completion —
				// the intermediate now owns that pin and will release it when
				// its own refcount drains.
				if parentID != "" && !pinTransferred {
					pinTransferred = true
				}
			}
		}
	}
	return pinTransferred
}
