package writers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/internal/pk"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type Writer interface {
	Write(ctx context.Context, res <-chan plugin.Message) error
}

const (
	defaultBatchTimeoutSeconds = 20
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

type BatchWriterClient interface {
	WriteTableBatch(ctx context.Context, table *schema.Table, resources []arrow.Record) error
}

type BatchWriter struct {
	tables      schema.Tables
	client      BatchWriterClient
	workers     map[string]*worker
	workersLock *sync.Mutex

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSize      int
	batchSizeBytes int
}

type Option func(*BatchWriter)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *BatchWriter) {
		p.logger = logger
	}
}

func WithBatchTimeout(timeout time.Duration) Option {
	return func(p *BatchWriter) {
		p.batchTimeout = timeout
	}
}

func WithBatchSize(size int) Option {
	return func(p *BatchWriter) {
		p.batchSize = size
	}
}

func WithBatchSizeBytes(size int) Option {
	return func(p *BatchWriter) {
		p.batchSizeBytes = size
	}
}

type worker struct {
	count int
	wg    *sync.WaitGroup
	ch    chan arrow.Record
	flush chan chan bool
}

func NewBatchWriter(tables schema.Tables, client BatchWriterClient, opts ...Option) (*BatchWriter, error) {
	c := &BatchWriter{
		tables:         tables,
		client:         client,
		workers:        make(map[string]*worker),
		workersLock:    &sync.Mutex{},
		logger:         zerolog.Nop(),
		batchTimeout:   defaultBatchTimeoutSeconds * time.Second,
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (w *BatchWriter) worker(ctx context.Context, table *schema.Table, ch <-chan arrow.Record, flush <-chan chan bool) {
	sizeBytes := int64(0)
	resources := make([]arrow.Record, 0)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					w.flush(ctx, table, resources)
				}
				return
			}
			if uint64(len(resources)) == 1000 || sizeBytes+util.TotalRecordSize(r) > int64(1000) {
				w.flush(ctx, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			resources = append(resources, r)
			sizeBytes += util.TotalRecordSize(r)
		case <-time.After(w.batchTimeout):
			if len(resources) > 0 {
				w.flush(ctx, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				w.flush(ctx, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (w *BatchWriter) flush(ctx context.Context, table *schema.Table, resources []arrow.Record) {
	resources = w.removeDuplicatesByPK(table, resources)
	start := time.Now()
	batchSize := len(resources)
	if err := w.client.WriteTableBatch(ctx, table, resources); err != nil {
		w.logger.Err(err).Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
	} else {
		w.logger.Info().Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
	}
}

func (*BatchWriter) removeDuplicatesByPK(table *schema.Table, resources []arrow.Record) []arrow.Record {
	pkIndices := table.PrimaryKeysIndexes()
	// special case where there's no PK at all
	if len(pkIndices) == 0 {
		return resources
	}

	pks := make(map[string]struct{}, len(resources))
	res := make([]arrow.Record, 0, len(resources))
	for _, r := range resources {
		if r.NumRows() > 1 {
			panic(fmt.Sprintf("record with more than 1 row: %d", r.NumRows()))
		}
		key := pk.String(r)
		_, ok := pks[key]
		if !ok {
			pks[key] = struct{}{}
			res = append(res, r)
			continue
		}
		// duplicate, release
		r.Release()
	}

	return res
}

func (w *BatchWriter) Write(ctx context.Context, res <-chan arrow.Record) error {
	workers := make(map[string]*worker, len(w.tables))

	w.workersLock.Lock()
	for _, table := range w.tables {
		table := table
		if w.workers[table.Name] == nil {
			ch := make(chan arrow.Record)
			flush := make(chan chan bool)
			wg := &sync.WaitGroup{}
			w.workers[table.Name] = &worker{
				count: 1,
				ch:    ch,
				flush: flush,
				wg:    wg,
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.worker(ctx, table, ch, flush)
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
