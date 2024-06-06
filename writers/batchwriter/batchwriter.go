package batchwriter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v16/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"github.com/rs/zerolog"
)

type Client interface {
	MigrateTables(context.Context, message.WriteMigrateTables) error
	WriteTableBatch(ctx context.Context, name string, messages message.WriteInserts) error
	DeleteStale(context.Context, message.WriteDeleteStales) error
	DeleteRecord(context.Context, message.WriteDeleteRecords) error
}

type BatchWriter struct {
	client           Client
	workers          map[string]*worker
	workersLock      sync.RWMutex
	workersWaitGroup sync.WaitGroup

	migrateTableLock     sync.Mutex
	migrateTableMessages message.WriteMigrateTables
	deleteStaleLock      sync.Mutex
	deleteStaleMessages  message.WriteDeleteStales
	deleteRecordLock     sync.Mutex
	deleteRecordMessages message.WriteDeleteRecords

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSize      int64
	batchSizeBytes int64
}

// Assert at compile-time that BatchWriter implements the Writer interface
var _ writers.Writer = (*BatchWriter)(nil)

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
		p.batchSize = int64(size)
	}
}

func WithBatchSizeBytes(size int) Option {
	return func(p *BatchWriter) {
		p.batchSizeBytes = int64(size)
	}
}

type worker struct {
	ch    chan *message.WriteInsert
	flush chan chan bool
}

const (
	defaultBatchTimeoutSeconds = 20
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

func New(client Client, opts ...Option) (*BatchWriter, error) {
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
	c.migrateTableMessages = make([]*message.WriteMigrateTable, 0, c.batchSize)
	c.deleteStaleMessages = make([]*message.WriteDeleteStale, 0, c.batchSize)
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

func (w *BatchWriter) Close(context.Context) error {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()
	for _, w := range w.workers {
		close(w.ch)
	}
	w.workersWaitGroup.Wait()

	return nil
}

func (w *BatchWriter) worker(ctx context.Context, tableName string, ch <-chan *message.WriteInsert, flush <-chan chan bool) {
	var bytes, rows int64
	resources := make([]*message.WriteInsert, 0, w.batchSize) // at least we have 1 row per record

	ticker := writers.NewTicker(w.batchTimeout)
	defer ticker.Stop()

	tickerCh, ctxDone := ticker.Chan(), ctx.Done()

	send := func() {
		w.flushTable(ctx, tableName, resources)
		clear(resources)
		resources = resources[:0]
		bytes, rows = 0, 0
	}
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if rows > 0 {
					w.flushTable(ctx, tableName, resources)
				}
				return
			}

			recordRows, recordBytes := r.Record.NumRows(), util.TotalRecordSize(r.Record)
			if (w.batchSize > 0 && rows+recordRows > w.batchSize) ||
				(w.batchSizeBytes > 0 && bytes+recordBytes > w.batchSizeBytes) {
				if rows == 0 {
					// New record overflows batch by itself.
					// Flush right away.
					// TODO: slice
					resources = append(resources, r)
					send()
					ticker.Reset(w.batchTimeout)
					continue
				}
				// rows > 0
				send()
				ticker.Reset(w.batchTimeout)
			}
			if recordRows > 0 {
				// only save records with rows
				resources = append(resources, r)
				rows += recordRows
				bytes += recordBytes
			}

		case <-tickerCh:
			if rows > 0 {
				send()
			}
		case done := <-flush:
			if rows > 0 {
				send()
				ticker.Reset(w.batchTimeout)
			}
			done <- true
		case <-ctxDone:
			// this means the request was cancelled
			return // after this NO other call will succeed
		}
	}
}

func (w *BatchWriter) flushTable(ctx context.Context, tableName string, resources []*message.WriteInsert) {
	// resources = w.removeDuplicatesByPK(table, resources)
	start := time.Now()
	batchSize := len(resources)
	if err := w.client.WriteTableBatch(ctx, tableName, resources); err != nil {
		w.logger.Err(err).Str("table", tableName).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
	} else {
		w.logger.Info().Str("table", tableName).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
	}
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

func (w *BatchWriter) flushDeleteRecordTables(ctx context.Context) error {
	w.deleteRecordLock.Lock()
	defer w.deleteRecordLock.Unlock()
	if len(w.deleteRecordMessages) == 0 {
		return nil
	}
	if err := w.client.DeleteRecord(ctx, w.deleteRecordMessages); err != nil {
		return err
	}
	w.deleteRecordMessages = w.deleteRecordMessages[:0]
	return nil
}

func (w *BatchWriter) flushInsert(tableName string) {
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

func (w *BatchWriter) writeAll(ctx context.Context, msgs []message.WriteMessage) error {
	ch := make(chan message.WriteMessage, len(msgs))
	for _, msg := range msgs {
		ch <- msg
	}
	close(ch)
	return w.Write(ctx, ch)
}

func (w *BatchWriter) Write(ctx context.Context, msgs <-chan message.WriteMessage) error {
	for msg := range msgs {
		switch m := msg.(type) {
		case *message.WriteDeleteStale:
			if err := w.flushMigrateTables(ctx); err != nil {
				return err
			}
			w.flushInsert(m.TableName)
			w.deleteStaleLock.Lock()
			w.deleteStaleMessages = append(w.deleteStaleMessages, m)
			l := int64(len(w.deleteStaleMessages))
			w.deleteStaleLock.Unlock()
			if w.batchSize > 0 && l > w.batchSize {
				if err := w.flushDeleteStaleTables(ctx); err != nil {
					return err
				}
			}
		case *message.WriteDeleteRecord:
			if err := w.flushMigrateTables(ctx); err != nil {
				return err
			}
			if err := w.flushDeleteStaleTables(ctx); err != nil {
				return err
			}
			// Ensure all related workers are flushed
			for _, rel := range m.TableRelations {
				w.flushInsert(rel.TableName)
			}
			w.deleteRecordLock.Lock()
			w.deleteRecordMessages = append(w.deleteRecordMessages, m)
			l := int64(len(w.deleteRecordMessages))
			w.deleteRecordLock.Unlock()
			if w.batchSize > 0 && l > w.batchSize {
				if err := w.flushDeleteRecordTables(ctx); err != nil {
					return err
				}
			}
		case *message.WriteInsert:
			if err := w.flushMigrateTables(ctx); err != nil {
				return err
			}
			if err := w.flushDeleteStaleTables(ctx); err != nil {
				return err
			}
			if err := w.startWorker(ctx, m); err != nil {
				return err
			}
		case *message.WriteMigrateTable:
			w.flushInsert(m.Table.Name)
			if err := w.flushDeleteStaleTables(ctx); err != nil {
				return err
			}
			w.migrateTableLock.Lock()
			w.migrateTableMessages = append(w.migrateTableMessages, m)
			l := int64(len(w.migrateTableMessages))
			w.migrateTableLock.Unlock()
			if w.batchSize > 0 && l > w.batchSize {
				if err := w.flushMigrateTables(ctx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w *BatchWriter) startWorker(_ context.Context, msg *message.WriteInsert) error {
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
	ch := make(chan *message.WriteInsert)
	flush := make(chan chan bool)
	wr = &worker{
		ch:    ch,
		flush: flush,
	}
	w.workers[tableName] = wr
	w.workersLock.Unlock()
	w.workersWaitGroup.Add(1)
	go func() {
		defer w.workersWaitGroup.Done()
		// TODO: we need to create a cancellable context that then can be cancelled via
		// w.cancelWorkers()
		w.worker(context.Background(), tableName, ch, flush)
	}()
	ch <- msg
	return nil
}
