package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ErrNotFound is what a Refresher returns for a resource that no longer exists at
// the source. Refresh turns it into a delete message.
var ErrNotFound = errors.New("resource not found")

// Refresher re-fetches single resources by their primary keys. It is implemented
// by source plugins that support event-based syncing.
type Refresher interface {
	// Refresh returns the client scoped to the resource and the resource's current
	// items, or ErrNotFound if it no longer exists.
	Refresh(ctx context.Context, table *schema.Table, keys arrow.RecordBatch) (schema.ClientMeta, []any, error)
}

// Refresh re-fetches the given resources and emits the result to res: an insert
// for every item the Refresher returns (relations included) and a delete message
// for resources that are gone or explicitly deleted.
//
// Unlike Sync it emits no migrate messages and never multiplexes: the caller
// already knows which table and which client each resource belongs to.
//
// A resource that fails does not stop the others. Every failure is reported as a
// SyncError message and all of them are returned joined once every resource has
// been processed.
func (s *Scheduler) Refresh(ctx context.Context, tables schema.Tables, r Refresher, resources []plugin.RefreshResource, res chan<- message.SyncMessage, opts ...SyncOption) error {
	if len(resources) == 0 {
		return nil
	}

	ctx, span := otel.Tracer(metrics.ResourceName).Start(ctx,
		"refresh",
		trace.WithAttributes(attribute.Key("sync.invocation.id").String(s.invocationID)),
	)
	defer span.End()

	if maxDepth(tables) > s.maxDepth {
		return fmt.Errorf("max depth exceeded, max depth is %d", s.maxDepth)
	}

	// Only the sync-scoped options are meaningful here; the shard, if any, was
	// already applied when the changes were produced.
	settings := &syncClient{scheduler: s}
	for _, opt := range opts {
		opt(settings)
	}

	// One slot per resource plus one for the context's cause.
	errs := make([]error, len(resources), len(resources)+1)
	resolved := make(chan *schema.Resource)
	go func() {
		defer close(resolved)
		var wg sync.WaitGroup
		for i, resource := range resources {
			// Bound the fan-out with the same semaphore top-level tables use.
			if err := s.tableSems[0].Acquire(ctx, 1); err != nil {
				errs[i] = fmt.Errorf("failed to refresh resource of table %q: %w", resource.Table, err)
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer s.tableSems[0].Release(1)
				errs[i] = s.refreshResource(ctx, tables, r, resource, settings.deterministicCQID, res, resolved)
			}()
		}
		wg.Wait()
	}()

	b := s.batchSettings.getBatcher(ctx, res, s.logger)
	defer b.close() // wait for all resources to be processed
	for resource := range resolved {
		b.process(resource)
	}

	return errors.Join(append(errs, context.Cause(ctx))...)
}

// refreshResource handles one resource: it either emits a delete message or runs
// the freshly fetched items through the regular resolve pipeline.
func (s *Scheduler) refreshResource(ctx context.Context, tables schema.Tables, r Refresher, resource plugin.RefreshResource, deterministicCQID bool, res chan<- message.SyncMessage, resolved chan<- *schema.Resource) error {
	table := tables.Get(resource.Table)
	if table == nil {
		return fmt.Errorf("table %q not found", resource.Table)
	}

	if resource.Op == plugin.OpDelete {
		return sendDeleteRecord(ctx, table, resource.Keys, res)
	}

	client, items, err := r.Refresh(ctx, table, resource.Keys)
	switch {
	case errors.Is(err, ErrNotFound):
		return sendDeleteRecord(ctx, table, resource.Keys, res)
	case err != nil:
		err = fmt.Errorf("failed to refresh resource %s of table %q: %w", schema.KeysString(resource.Keys), table.Name, err)
		select {
		case res <- &message.SyncError{TableName: table.Name, Error: err.Error()}:
		case <-ctx.Done():
		}
		return err
	}

	if len(items) == 0 {
		return nil
	}
	if client == nil {
		return fmt.Errorf("refresher returned no client for resource %s of table %q", schema.KeysString(resource.Keys), table.Name)
	}

	// A per-call copy of the table whose resolver replays the items we already
	// have. Everything below it - column resolvers, PostResourceResolver,
	// relations, CQID calculation - runs exactly as it does during a full sync.
	refreshTable := table.Copy(nil)
	refreshTable.Multiplex = nil
	refreshTable.Resolver = func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, out chan<- any) error {
		for _, item := range items {
			select {
			case out <- item:
			case <-ctx.Done():
				return context.Cause(ctx)
			}
		}
		return nil
	}

	// Metrics are per resource: the scheduler-wide structure is written to without
	// locks, so sharing it across the concurrent refreshes would race.
	sc := &syncClient{
		metrics:           metrics.NewMetrics(),
		tables:            schema.Tables{refreshTable},
		client:            client,
		scheduler:         s,
		deterministicCQID: deterministicCQID,
		logger:            s.logger,
		invocationID:      s.invocationID,
		msgChan:           res,
	}
	sc.metrics.InitWithClients(refreshTable, []schema.ClientMeta{client})
	sc.resolveTableDfs(ctx, refreshTable, client, nil, resolved, 1)

	return nil
}

func sendDeleteRecord(ctx context.Context, table *schema.Table, keys arrow.RecordBatch, res chan<- message.SyncMessage) error {
	msg, err := message.NewDeleteRecordByKeys(table, keys)
	if err != nil {
		return fmt.Errorf("failed to build delete message for table %q: %w", table.Name, err)
	}
	select {
	case res <- msg:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}
