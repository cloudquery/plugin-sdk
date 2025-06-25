package queue

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
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

func (d *Scheduler) Sync(ctx context.Context, tableClients []WorkUnit, resolvedResources chan<- *schema.Resource, msgChan chan<- message.SyncMessage) {
	if len(tableClients) == 0 {
		return
	}
	queue := NewConcurrentRandomQueue[WorkUnit](d.seed, len(tableClients))
	for _, tc := range tableClients {
		queue.Push(tc)
	}

	jobs := make(chan *WorkUnit)
	activeWorkSignal := newActiveWorkSignal()

	// Worker pool
	workerPool, _ := errgroup.WithContext(ctx)
	for w := 0; w < d.workerCount; w++ {
		workerPool.Go(func() error {
			newWorker(
				jobs,
				queue,
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

	// Work distribution
	go func() {
		defer close(jobs)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				item := queue.Pop()

				// There is work to do
				if item != nil {
					jobs <- item
					continue
				}

				// Queue is empty and no active work, done!
				if activeWorkSignal.IsIdle() {
					return
				}

				// Queue is empty and there is active work, wait for changes
				activeWorkSignal.Wait()
			}
		}
	}()

	_ = workerPool.Wait()
}
