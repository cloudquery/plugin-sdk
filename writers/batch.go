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
	"golang.org/x/sync/semaphore"
)

type Writer interface {
	Write(ctx context.Context, writeOptions plugin.WriteOptions, res <-chan plugin.Message) error
}

const (
	defaultBatchTimeoutSeconds = 20
	defaultMaxWorkers          = int64(10000)
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

type BatchWriterClient interface {
	MigrateTables(context.Context, []*plugin.MessageMigrateTable) error
	WriteTableBatch(ctx context.Context, name string, upsert bool, msgs []*plugin.MessageInsert) error
	DeleteStale(context.Context, []*plugin.MessageDeleteStale) error
}

type BatchWriter struct {
	client               BatchWriterClient
	semaphore            *semaphore.Weighted
	workers              map[string]*worker
	workersLock          *sync.RWMutex
	workersWaitGroup     *sync.WaitGroup
	migrateTableMessages []*plugin.MessageMigrateTable
	deleteStaleMessages  []*plugin.MessageDeleteStale

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

func WithMaxWorkers(n int64) Option {
	return func(p *BatchWriter) {
		p.semaphore = semaphore.NewWeighted(n)
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
	ch    chan *plugin.MessageInsert
	flush chan chan bool
}

func NewBatchWriter(client BatchWriterClient, opts ...Option) (*BatchWriter, error) {
	c := &BatchWriter{
		client:           client,
		workers:          make(map[string]*worker),
		workersLock:      &sync.RWMutex{},
		workersWaitGroup: &sync.WaitGroup{},
		logger:           zerolog.Nop(),
		batchTimeout:     defaultBatchTimeoutSeconds * time.Second,
		batchSize:        defaultBatchSize,
		batchSizeBytes:   defaultBatchSizeBytes,
		semaphore:        semaphore.NewWeighted(defaultMaxWorkers),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.migrateTableMessages = make([]*plugin.MessageMigrateTable, 0, c.batchSize)
	c.deleteStaleMessages = make([]*plugin.MessageDeleteStale, 0, c.batchSize)
	return c, nil
}

func (w *BatchWriter) Flush(ctx context.Context) error {
	w.workersLock.RLock()
	for _, worker := range w.workers {
		done := make(chan bool)
		worker.flush <- done
		<-done
	}
	w.workersLock.RUnlock()
	w.flushMigrateTables(ctx)
	w.flushDeleteStaleTables(ctx)
	return nil
}

func (w *BatchWriter) Close(ctx context.Context) error {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()
	for _, w := range w.workers {
		close(w.ch)
	}
	w.workersWaitGroup.Wait()

	return nil
}

func (w *BatchWriter) worker(ctx context.Context, tableName string, ch <-chan *plugin.MessageInsert, flush <-chan chan bool) {
	sizeBytes := int64(0)
	resources := make([]*plugin.MessageInsert, 0)
	upsertBatch := false
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					w.flush(ctx, tableName, upsertBatch, resources)
				}
				return
			}
			if upsertBatch != r.Upsert {
				w.flush(ctx, tableName, upsertBatch, resources)
				resources = make([]*plugin.MessageInsert, 0)
				sizeBytes = 0
				upsertBatch = r.Upsert
				resources = append(resources, r)
				sizeBytes = util.TotalRecordSize(r.Record)
			} else {
				resources = append(resources, r)
				sizeBytes += util.TotalRecordSize(r.Record)
			}
			if len(resources) >= w.batchSize || sizeBytes+util.TotalRecordSize(r.Record) >= int64(w.batchSizeBytes) {
				w.flush(ctx, tableName, upsertBatch, resources)
				resources = make([]*plugin.MessageInsert, 0)
				sizeBytes = 0
			}
		case <-time.After(w.batchTimeout):
			if len(resources) > 0 {
				w.flush(ctx, tableName, upsertBatch, resources)
				resources = make([]*plugin.MessageInsert, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				w.flush(ctx, tableName, upsertBatch, resources)
				resources = make([]*plugin.MessageInsert, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (w *BatchWriter) flush(ctx context.Context, tableName string, upsertBatch bool, resources []*plugin.MessageInsert) {
	// resources = w.removeDuplicatesByPK(table, resources)
	start := time.Now()
	batchSize := len(resources)
	if err := w.client.WriteTableBatch(ctx, tableName, upsertBatch, resources); err != nil {
		w.logger.Err(err).Str("table", tableName).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
	} else {
		w.logger.Info().Str("table", tableName).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
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

func (w *BatchWriter) flushMigrateTables(ctx context.Context) error {
	if err := w.client.MigrateTables(ctx, w.migrateTableMessages); err != nil {
		return err
	}
	w.migrateTableMessages = w.migrateTableMessages[:0]
	return nil
}

func (w *BatchWriter) flushDeleteStaleTables(ctx context.Context) error {
	if err := w.client.DeleteStale(ctx, w.deleteStaleMessages); err != nil {
		return err
	}
	w.deleteStaleMessages = w.deleteStaleMessages[:0]
	return nil
}

func (w *BatchWriter) flushInsert(ctx context.Context, tableName string) {
	w.workersLock.RLock()
	worker, ok := w.workers[tableName]
	if !ok {
		w.workersLock.RUnlock()
		// no tables to flush
		return
	}
	w.workersLock.RUnlock()
	ch := make(chan bool)
	worker.flush <- ch
	<-ch
}

func (w *BatchWriter) writeAll(ctx context.Context, msgs []plugin.Message) error {
	ch := make(chan plugin.Message, len(msgs))
	for _, msg := range msgs {
		ch <- msg
	}
	close(ch)
	return w.Write(ctx, ch)
}

func (w *BatchWriter) Write(ctx context.Context, msgs <-chan plugin.Message) error {
	for msg := range msgs {
		switch m := msg.(type) {
		case *plugin.MessageDeleteStale:
			if len(w.migrateTableMessages) > 0 {
				if err := w.flushMigrateTables(ctx); err != nil {
					return err
				}
			}
			w.flushInsert(ctx, m.Table.Name)
			w.deleteStaleMessages = append(w.deleteStaleMessages, m)
			if len(w.deleteStaleMessages) > w.batchSize {
				if err := w.flushDeleteStaleTables(ctx); err != nil {
					return err
				}
			}
		case *plugin.MessageInsert:
			if len(w.migrateTableMessages) > 0 {
				if err := w.flushMigrateTables(ctx); err != nil {
					return err
				}
			}
			if len(w.deleteStaleMessages) > 0 {
				if err := w.flushDeleteStaleTables(ctx); err != nil {
					return err
				}
			}
			if err := w.startWorker(ctx, m); err != nil {
				return err
			}
		case *plugin.MessageMigrateTable:
			w.flushInsert(ctx, m.Table.Name)
			if len(w.deleteStaleMessages) > 0 {
				if err := w.flushDeleteStaleTables(ctx); err != nil {
					return err
				}
			}
			w.migrateTableMessages = append(w.migrateTableMessages, m)
			if len(w.migrateTableMessages) > w.batchSize {
				if err := w.flushMigrateTables(ctx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w *BatchWriter) startWorker(ctx context.Context, msg *plugin.MessageInsert) error {
	w.workersLock.RLock()
	md := msg.Record.Schema().Metadata()
	tableName, ok := md.GetValue(schema.MetadataTableName)
	if !ok {
		w.workersLock.RUnlock()
		return fmt.Errorf("table name not found in metadata")
	}
	wr, ok := w.workers[tableName]
	w.workersLock.RUnlock()
	if ok {
		w.workers[tableName].ch <- msg
		return nil
	}
	w.workersLock.Lock()
	ch := make(chan *plugin.MessageInsert)
	flush := make(chan chan bool)
	wr = &worker{
		count: 1,
		ch:    ch,
		flush: flush,
	}
	w.workers[tableName] = wr
	w.workersLock.Unlock()
	w.workersWaitGroup.Add(1)
	go func() {
		defer w.workersWaitGroup.Done()
		w.worker(ctx, tableName, ch, flush)
	}()
	ch <- msg
	return nil
}
