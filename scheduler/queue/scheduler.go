package queue

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const DefaultWorkerCount = 1000

// WorkUnit is an atomic unit of work that the scheduler syncs.
//
// It is one table resolver (same as all other scheduler strategies).
//
// But if it is a non-top-level table, it is bound to a single parent resource.
type WorkUnit struct {
	Table  *schema.Table
	Client schema.ClientMeta
	Parent *schema.Resource
}

type Scheduler struct {
	workerCount       int
	logger            zerolog.Logger
	caser             *caser.Caser
	deterministicCQID bool
	metrics           *metrics.Metrics
	invocationID      string
	seed              int64

	// storage holds the pluggable queue+resource backend. Required.
	storage storage.Storage
	// codec serializes *schema.Resource blobs for the backend. Required.
	codec *Codec
}

type Option func(*Scheduler)

func WithWorkerCount(workerCount int) Option {
	return func(d *Scheduler) {
		d.workerCount = workerCount
	}
}

func WithCaser(c *caser.Caser) Option {
	return func(d *Scheduler) {
		d.caser = c
	}
}

func WithDeterministicCQID(deterministicCQID bool) Option {
	return func(d *Scheduler) {
		d.deterministicCQID = deterministicCQID
	}
}

func WithInvocationID(invocationID string) Option {
	return func(d *Scheduler) {
		d.invocationID = invocationID
	}
}

// WithStorage sets the Storage backend that holds work and parent resources.
// Required — the scheduler will no-op (with a logged error) if unset.
func WithStorage(s storage.Storage) Option {
	return func(d *Scheduler) {
		d.storage = s
	}
}

// WithCodec sets the Codec used to serialize resources into Storage. Required.
func WithCodec(c *Codec) Option {
	return func(d *Scheduler) {
		d.codec = c
	}
}

func NewShuffleQueueScheduler(logger zerolog.Logger, m *metrics.Metrics, seed int64, opts ...Option) *Scheduler {
	scheduler := &Scheduler{
		logger:       logger,
		metrics:      m,
		workerCount:  DefaultWorkerCount,
		caser:        caser.New(),
		invocationID: uuid.New().String(),
		seed:         seed,
	}

	for _, opt := range opts {
		opt(scheduler)
	}

	return scheduler
}

// workerLookups holds in-process-only references to Tables and Clients so
// workers can reconstitute a (*schema.Table, schema.ClientMeta, *schema.Resource)
// from a SerializedWorkUnit.
type workerLookups struct {
	tables  map[string]*schema.Table
	clients map[string]schema.ClientMeta
}

func (d *Scheduler) Sync(ctx context.Context, tableClients []WorkUnit, resolvedResources chan<- *schema.Resource, msgChan chan<- message.SyncMessage) {
	if len(tableClients) == 0 {
		return
	}
	if d.storage == nil {
		d.logger.Error().Msg("queue scheduler started with nil Storage")
		return
	}
	if d.codec == nil {
		d.logger.Error().Msg("queue scheduler started with nil Codec")
		return
	}

	// Maintain in-memory lookup tables so workers can rehydrate.
	lookups := &workerLookups{
		tables:  make(map[string]*schema.Table),
		clients: make(map[string]schema.ClientMeta),
	}
	// Walk the full table tree so relation tables (which can appear as
	// ParentID on future WorkUnits) are resolvable by name.
	walkTables := func(t *schema.Table) {
		lookups.tables[t.Name] = t
	}
	for _, wu := range tableClients {
		lookups.clients[wu.Client.ID()] = wu.Client
	}
	walk := func(tables []*schema.Table) {
		var do func([]*schema.Table)
		do = func(ts []*schema.Table) {
			for _, t := range ts {
				walkTables(t)
				do(t.Relations)
			}
		}
		do(tables)
	}
	// Collect root tables (dedup) and walk them.
	rootTables := make([]*schema.Table, 0, len(tableClients))
	seenRoot := make(map[string]bool)
	for _, wu := range tableClients {
		if !seenRoot[wu.Table.Name] {
			seenRoot[wu.Table.Name] = true
			rootTables = append(rootTables, wu.Table)
		}
	}
	walk(rootTables)

	// Seed: push the initial (root-level) WorkUnits. ParentID is empty —
	// these have no parent resource in the KV.
	seed := make([]storage.SerializedWorkUnit, 0, len(tableClients))
	for _, wu := range tableClients {
		seed = append(seed, storage.SerializedWorkUnit{
			TableName: wu.Table.Name,
			ClientID:  wu.Client.ID(),
			// ParentID: "" — top-level
		})
	}
	if err := d.storage.PushWorkBatch(ctx, seed); err != nil {
		d.logger.Error().Err(err).Msg("failed to seed work queue")
		return
	}

	jobs := make(chan *storage.SerializedWorkUnit)
	activeWorkSignal := newActiveWorkSignal()

	workerPool, _ := errgroup.WithContext(ctx)
	for w := 0; w < d.workerCount; w++ {
		workerPool.Go(func() error {
			newWorker(
				jobs,
				d.storage,
				d.codec,
				lookups,
				resolvedResources,
				d.logger,
				d.caser,
				d.invocationID,
				d.deterministicCQID,
				d.metrics,
				msgChan,
			).work(ctx, activeWorkSignal)
			return nil
		})
	}

	// Work distribution — pulls from Storage, signals idle.
	go func() {
		defer close(jobs)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				item, err := d.storage.PopWork(ctx)
				if err != nil {
					d.logger.Error().Err(err).Msg("queue backend pop error; aborting sync")
					return
				}
				if item != nil {
					jobs <- item
					continue
				}
				if activeWorkSignal.IsIdle() {
					return
				}
				activeWorkSignal.Wait()
			}
		}
	}()

	_ = workerPool.Wait()
}
