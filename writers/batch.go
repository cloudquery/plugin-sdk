package writers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/internal/pk"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type Writer interface {
	Write(ctx context.Context, writeOptions plugin.WriteOptions, res <-chan message.Message) error
}

const (
	defaultBatchTimeoutSeconds = 20
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

type BatchWriterClient interface {
	MigrateTables(context.Context, []*message.MigrateTable) error
	WriteTableBatch(ctx context.Context, name string, msgs []*message.Insert) error
	DeleteStale(context.Context, []*message.DeleteStale) error
}

type BatchWriter struct {
	client           BatchWriterClient
	workers          map[string]*worker
	workersLock      sync.RWMutex
	workersWaitGroup sync.WaitGroup

	migrateTableLock     sync.Mutex
	migrateTableMessages []*message.MigrateTable
	deleteStaleLock      sync.Mutex
	deleteStaleMessages  []*message.DeleteStale

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
	ch    chan *message.Insert
	flush chan chan bool
}

func NewBatchWriter(client BatchWriterClient, opts ...Option) (*BatchWriter, error) {
	c := &BatchWriter{
		client:         client,
		workers:        make(map[string]*worker),
		logger:         zerolog.Nop(),
		batchTimeout:   defaultBatchTimeoutSeconds * time.Second,
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.migrateTableMessages = make([]*message.MigrateTable, 0, c.batchSize)
	c.deleteStaleMessages = make([]*message.DeleteStale, 0, c.batchSize)
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
	if err := w.flushMigrateTables(ctx); err != nil {
		return err
	}
	return w.flushDeleteStaleTables(ctx)
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

func (w *BatchWriter) worker(ctx context.Context, tableName string, ch <-chan *message.Insert, flush <-chan chan bool) {
	sizeBytes := int64(0)
	resources := make([]*message.Insert, 0)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					w.flush(ctx, tableName, resources)
				}
				return
			}
			resources = append(resources, r)
			sizeBytes += util.TotalRecordSize(r.Record)

			if len(resources) >= w.batchSize || sizeBytes+util.TotalRecordSize(r.Record) >= int64(w.batchSizeBytes) {
				w.flush(ctx, tableName, resources)
				resources = make([]*message.Insert, 0)
				sizeBytes = 0
			}
		case <-time.After(w.batchTimeout):
			if len(resources) > 0 {
				w.flush(ctx, tableName, resources)
				resources = make([]*message.Insert, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				w.flush(ctx, tableName, resources)
				resources = make([]*message.Insert, 0)
				sizeBytes = 0
			}
			done <- true
		case <-ctx.Done():
			// this means the request was cancelled
			return // after this NO other call will succeed
		}
	}
}

func (w *BatchWriter) flush(ctx context.Context, tableName string, resources []*message.Insert) {
	// resources = w.removeDuplicatesByPK(table, resources)
	start := time.Now()
	batchSize := len(resources)
	if err := w.client.WriteTableBatch(ctx, tableName, resources); err != nil {
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
	w.migrateTableLock.Lock()
	defer w.migrateTableLock.Unlock()
	if len(w.migrateTableMessages) == 0 {
		return nil
	}
	if err := w.client.MigrateTables(ctx, w.migrateTableMessages); err != nil {
		return err
	}
	w.migrateTableMessages = w.migrateTableMessages[:0]
	return nil
}

func (w *BatchWriter) flushDeleteStaleTables(ctx context.Context) error {
	w.deleteStaleLock.Lock()
	defer w.deleteStaleLock.Unlock()
	if len(w.deleteStaleMessages) == 0 {
		return nil
	}
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

func (w *BatchWriter) writeAll(ctx context.Context, msgs []message.Message) error {
	ch := make(chan message.Message, len(msgs))
	for _, msg := range msgs {
		ch <- msg
	}
	close(ch)
	return w.Write(ctx, ch)
}

func (w *BatchWriter) Write(ctx context.Context, msgs <-chan message.Message) error {
	for msg := range msgs {
		switch m := msg.(type) {
		case *message.DeleteStale:
			if err := w.flushMigrateTables(ctx); err != nil {
				return err
			}
			w.flushInsert(ctx, m.Table.Name)
			w.deleteStaleLock.Lock()
			w.deleteStaleMessages = append(w.deleteStaleMessages, m)
			l := len(w.deleteStaleMessages)
			w.deleteStaleLock.Unlock()
			if l > w.batchSize {
				if err := w.flushDeleteStaleTables(ctx); err != nil {
					return err
				}
			}
		case *message.Insert:
			if err := w.flushMigrateTables(ctx); err != nil {
				return err
			}
			if err := w.flushDeleteStaleTables(ctx); err != nil {
				return err
			}
			if err := w.startWorker(ctx, m); err != nil {
				return err
			}
		case *message.MigrateTable:
			w.flushInsert(ctx, m.Table.Name)
			if err := w.flushDeleteStaleTables(ctx); err != nil {
				return err
			}
			w.migrateTableLock.Lock()
			w.migrateTableMessages = append(w.migrateTableMessages, m)
			l := len(w.migrateTableMessages)
			w.migrateTableLock.Unlock()
			if l > w.batchSize {
				if err := w.flushMigrateTables(ctx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w *BatchWriter) startWorker(ctx context.Context, msg *message.Insert) error {
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
		wr.ch <- msg
		return nil
	}
	w.workersLock.Lock()
	ch := make(chan *message.Insert)
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
