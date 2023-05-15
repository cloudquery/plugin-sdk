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
}

func New(opts ...Option) destination.BatchingWriterFunc {
	return func(writer destination.ManagedWriter) destination.BatchingWriter {
		w := &Batching{
			underlyingWriter: writer,
			logger:           zerolog.Logger{},
			metrics:          make(map[string]*destination.Metrics),
			metricsLock:      &sync.RWMutex{},
			workers:          make(map[string]*worker),
			workersLock:      &sync.Mutex{},

			dedupPK:        true,
			batchSize:      10000,
			batchSizeBytes: 5 * 1024 * 1024, // 5 MiB
			batchTimeout:   5 * time.Second,
		}
		for _, opt := range opts {
			opt(w)
		}
		return w
	}
}

func (w *Batching) Metrics() destination.Metrics {
	metrics := destination.Metrics{}
	w.metricsLock.RLock()
	for _, m := range w.metrics {
		metrics.Errors += m.Errors
		metrics.Writes += m.Writes
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
	metrics := &destination.Metrics{}

	w.workersLock.Lock()
	for _, table := range tables {
		table := table
		if w.workers[table.Name] == nil {
			ch := make(chan arrow.Record)
			flush := make(chan chan bool)
			wg := &sync.WaitGroup{}
			w.workers[table.Name] = &worker{
				count:      1,
				ch:         ch,
				flush:      flush,
				wg:         wg,
				sourceSpec: sourceSpec,
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
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
			delete(w.workers, tableName)
		}
	}
	w.workersLock.Unlock()
	return nil
}
