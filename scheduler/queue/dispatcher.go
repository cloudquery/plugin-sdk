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
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const DefaultWorkerCount = 1000

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
		metrics.LogTablesMetrics(w.logger, w.metrics, table.Relations, client)
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
					tableMetrics := w.metrics.TableClient[table.Name][client.ID()]
					w.logger.Error().Err(err).Str("table", table.Name).Str("client", client.ID()).Msg("resource resolver finished with validation error")
					atomic.AddUint64(&tableMetrics.Errors, 1)
					return
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
			newJob := &TableClientPair{
				Table:  relation,
				Client: client,
				Parent: resource,
			}
			w.newJobs <- newJob
		}
	}
}

type worker struct {
	jobs              <-chan *TableClientPair
	newJobs           chan<- *TableClientPair
	resolvedResources chan<- *schema.Resource

	logger            zerolog.Logger
	caser             *caser.Caser
	invocationID      string
	deterministicCQID bool
	metrics           *metrics.Metrics
	onResolveStart    func()
	onResolveDone     func()
}

func newWorker(
	jobs <-chan *TableClientPair,
	newJobs chan<- *TableClientPair,
	resolvedResources chan<- *schema.Resource,

	logger zerolog.Logger,
	c *caser.Caser,
	invocationID string,
	deterministicCQID bool,
	m *metrics.Metrics,
	onResolveStart func(),
	onResolveDone func(),
) *worker {
	return &worker{
		jobs:              jobs,
		newJobs:           newJobs,
		resolvedResources: resolvedResources,
		logger:            logger,
		caser:             c,
		deterministicCQID: deterministicCQID,
		invocationID:      invocationID,
		metrics:           m,
		onResolveStart:    onResolveStart,
		onResolveDone:     onResolveDone,
	}
}

func (w *worker) work(ctx context.Context) {
	for j := range w.jobs {
		w.onResolveStart()
		client := j.Client
		table := j.Table
		parent := j.Parent

		w.resolveTable(ctx, table, client, parent)
		w.onResolveDone()
	}
}

type Dispatcher struct {
	workerCount       int
	logger            zerolog.Logger
	caser             *caser.Caser
	deterministicCQID bool
	metrics           *metrics.Metrics
	invocationID      string
}

type Option func(*Dispatcher)

func WithWorkerCount(workerCount int) Option {
	return func(d *Dispatcher) {
		d.workerCount = workerCount
	}
}

func WithCaser(c *caser.Caser) Option {
	return func(d *Dispatcher) {
		d.caser = c
	}
}

func WithDeterministicCQID(deterministicCQID bool) Option {
	return func(d *Dispatcher) {
		d.deterministicCQID = deterministicCQID
	}
}

func WithInvocationID(invocationID string) Option {
	return func(d *Dispatcher) {
		d.invocationID = invocationID
	}
}

func NewDispatcher(logger zerolog.Logger, m *metrics.Metrics, opts ...Option) *Dispatcher {
	dispatcher := &Dispatcher{
		logger:       logger,
		metrics:      m,
		workerCount:  DefaultWorkerCount,
		caser:        caser.New(),
		invocationID: uuid.New().String(),
	}

	for _, opt := range opts {
		opt(dispatcher)
	}

	return dispatcher
}

func (d *Dispatcher) Dispatch(ctx context.Context, tableClients []TableClientPair, resolvedResources chan<- *schema.Resource) {
	if len(tableClients) == 0 {
		return
	}
	queue := NewConcurrentQueue(tableClients)

	jobs := make(chan *TableClientPair)
	newJobs := make(chan *TableClientPair)
	activeWorkers := &atomic.Uint32{}
	workStarted := &atomic.Bool{}

	onResolveStart := func() {
		workStarted.Store(true)
		activeWorkers.Add(1)
	}
	onResolveDone := func() {
		activeWorkers.Add(^uint32(0))
	}

	wg := sync.WaitGroup{}
	for w := 0; w < d.workerCount; w++ {
		worker := newWorker(
			jobs,
			newJobs,
			resolvedResources,
			d.logger,
			d.caser,
			d.invocationID,
			d.deterministicCQID,
			d.metrics,
			onResolveStart,
			onResolveDone,
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.work(ctx)
		}()
	}

	go func() {
		for {
			if ctx.Err() != nil {
				break
			}
			item := queue.Pop()
			if item == nil {
				if !workStarted.Load() || activeWorkers.Load() != 0 {
					continue
				}
				break
			}
			jobs <- item
		}
		close(jobs)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case result := <-newJobs:
				queue.Push(*result)
			}
		}
	}()

	wg.Wait()
}
