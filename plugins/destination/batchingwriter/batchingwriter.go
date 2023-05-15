package batchingwriter

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/rs/zerolog"
)

type Batching struct {
	underlyingWriter destination.ManagedWriter
	logger           zerolog.Logger

	metrics     map[string]*destination.Metrics
	metricsLock *sync.RWMutex

	workers     map[string]*worker
	workersLock *sync.Mutex

	dedupPK        bool
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration

	// underlyingOCW is set if the given underlyingWriter implements OpenCloseWriter
	underlyingOCW OpenCloseWriter
}

// OpenCloseWriter is an optional interface that can be implemented by a Client which already implements destination.ManagedWriter.
type OpenCloseWriter interface {
	OpenTable(ctx context.Context, sourceSpec specs.Source, table *schema.Table) error
	CloseTable(ctx context.Context, sourceSpec specs.Source, table *schema.Table) error
	destination.ManagedWriter
}

func New(opts ...Option) destination.BatchingWriterFuncFunc {
	return func(spec *specs.Destination) destination.BatchingWriterFunc {
		w := &Batching{
			logger:      zerolog.Logger{},
			metrics:     make(map[string]*destination.Metrics),
			metricsLock: &sync.RWMutex{},
			workers:     make(map[string]*worker),
			workersLock: &sync.Mutex{},

			dedupPK:        true,
			batchSize:      int64(spec.BatchSize),
			batchSizeBytes: int64(spec.BatchSizeBytes),
			batchTimeout:   5 * time.Second,
		}
		for _, opt := range opts {
			opt(w)
		}
		// Set the spec batch sizes again, in case the options changed them. Used in testing
		spec.BatchSize = int(w.batchSize)
		spec.BatchSizeBytes = int(w.batchSizeBytes)
		spec.SetDefaults(
			10000,
			5*1024*1024, // 5 MiB
		)
		w.batchSize, w.batchSizeBytes = int64(spec.BatchSize), int64(spec.BatchSizeBytes)

		return func(writer destination.ManagedWriter) destination.BatchingWriter {
			w.underlyingWriter = writer

			if ocw, ok := writer.(OpenCloseWriter); ok {
				w.underlyingOCW = ocw
			}
			return w
		}
	}
}

func (w *Batching) Metrics() destination.Metrics {
	metrics := destination.Metrics{}
	w.metricsLock.RLock()
	for _, m := range w.metrics {
		metrics.Errors += atomic.LoadUint64(&m.Errors)
		metrics.Writes += atomic.LoadUint64(&m.Writes)
	}
	w.metricsLock.RUnlock()
	return metrics
}

type worker struct {
	count      int
	wg         *sync.WaitGroup
	ch         chan arrow.Record
	flush      chan chan bool
	sourceSpec specs.Source
	openError  bool
}

func (*Batching) dummyWork(_ context.Context, ch <-chan arrow.Record, flush <-chan chan bool) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case done := <-flush:
			done <- true
		}
	}
}
func (w *Batching) work(ctx context.Context, sourceSpec specs.Source, syncTime time.Time, metrics *destination.Metrics, table *schema.Table, ch <-chan arrow.Record, flush <-chan chan bool) {
	sizeBytes := int64(0)
	resources := make([]arrow.Record, 0)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					w.flush(ctx, sourceSpec, syncTime, metrics, table, resources)
				}
				return
			}
			if int64(len(resources)) == w.batchSize || sizeBytes+util.TotalRecordSize(r) > w.batchSizeBytes {
				w.flush(ctx, sourceSpec, syncTime, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			resources = append(resources, r)
			sizeBytes += util.TotalRecordSize(r)
		case <-time.After(w.batchTimeout):
			if len(resources) > 0 {
				w.flush(ctx, sourceSpec, syncTime, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				w.flush(ctx, sourceSpec, syncTime, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (w *Batching) flush(ctx context.Context, sourceSpec specs.Source, syncTime time.Time, metrics *destination.Metrics, table *schema.Table, resources []arrow.Record) {
	if w.dedupPK {
		resources = destination.RemoveDuplicatesByPK(table, resources)
	}
	start := time.Now()
	batchSize := len(resources)

	if err := w.underlyingWriter.WriteTableBatch(ctx, sourceSpec, table, syncTime, resources); err != nil {
		w.logger.Err(err).Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
		// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
		atomic.AddUint64(&metrics.Errors, uint64(batchSize))
	} else {
		w.logger.Info().Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
		atomic.AddUint64(&metrics.Writes, uint64(batchSize))
	}
}

func (w *Batching) Write(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan arrow.Record) error {
	workers := make(map[string]*worker, len(tables))

	w.workersLock.Lock()
	for _, table := range tables {
		table := table
		metrics := &destination.Metrics{}
		w.metricsLock.Lock()
		w.metrics[table.Name] = metrics
		w.metricsLock.Unlock()

		if w.workers[table.Name] == nil {
			ch := make(chan arrow.Record)
			flush := make(chan chan bool)
			wg := &sync.WaitGroup{}
			var errored bool
			if w.underlyingOCW != nil {
				if err := w.underlyingOCW.OpenTable(ctx, sourceSpec, table); err != nil {
					w.logger.Err(err).Str("table", table.Name).Msg("OpenTable failed")
					// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
					atomic.AddUint64(&metrics.Errors, 1)
					errored = true
				}
			}
			w.workers[table.Name] = &worker{
				count:      1,
				ch:         ch,
				flush:      flush,
				wg:         wg,
				sourceSpec: sourceSpec,
				openError:  errored,
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				if errored {
					w.dummyWork(ctx, ch, flush)
					return
				}
				w.work(ctx, sourceSpec, syncTime, metrics, table, ch, flush)
			}()
		} else {
			w.workers[table.Name].count++
		}
		// we save this locally because we don't want to access the map after that so we can
		// keep the workersLock for as short as possible
		workers[table.Name] = w.workers[table.Name]
	}
	w.workersLock.Unlock()

	for r := range res {
		tableName, ok := r.Schema().Metadata().GetValue(schema.MetadataTableName)
		if !ok {
			return fmt.Errorf("missing table name in record metadata")
		}
		if _, ok := workers[tableName]; !ok {
			return fmt.Errorf("table %s not found in destination", tableName)
		}
		workers[tableName].ch <- r
	}

	// flush and wait for all workers to finish flush before finish and calling delete stale
	// This is because destinations can be longed lived and called from multiple sources
	flushChannels := make(map[string]chan bool, len(workers))
	for tableName, w := range workers {
		flushCh := make(chan bool)
		flushChannels[tableName] = flushCh
		w.flush <- flushCh
	}
	for tableName := range flushChannels {
		<-flushChannels[tableName]
	}

	w.workersLock.Lock()
	for tableName := range workers {
		w.workers[tableName].count--
		if w.workers[tableName].count == 0 {
			close(w.workers[tableName].ch)
			w.workers[tableName].wg.Wait()

			if w.underlyingOCW != nil && !w.workers[tableName].openError {
				if err := w.underlyingOCW.CloseTable(ctx, sourceSpec, tables.Get(tableName)); err != nil {
					w.logger.Err(err).Str("table", tableName).Msg("CloseTable failed")

					w.metricsLock.RLock()
					metrics := w.metrics[tableName]
					w.metricsLock.RUnlock()

					atomic.AddUint64(&metrics.Errors, 1)
				}
			}

			delete(w.workers, tableName)
		}
	}
	w.workersLock.Unlock()
	return nil
}

// BatchSize returns the current batch size, used in testing
func (w *Batching) BatchSize() (batchSize int64, batchSizeBytes int64) {
	return w.batchSize, w.batchSizeBytes
}
