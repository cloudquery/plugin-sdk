package batchwriter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"github.com/rs/zerolog"
)

type Client interface {
	MigrateTables(context.Context, []*message.WriteMigrateTable) error
	WriteTableBatch(ctx context.Context, name string, msgs []*message.WriteInsert) error
	DeleteStale(context.Context, []*message.WriteDeleteStale) error
}

type BatchWriter struct {
	client           Client
	workers          map[string]*worker
	workersLock      sync.RWMutex
	workersWaitGroup sync.WaitGroup

	migrateTableLock     sync.Mutex
	migrateTableMessages []*message.WriteMigrateTable
	deleteStaleLock      sync.Mutex
	deleteStaleMessages  []*message.WriteDeleteStale

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSize      int
	batchSizeBytes int
}

// Assert at compile-time that BatchWriter implements the Writer interface
var _ writers.Writer = (*BatchWriter)(nil)

type Option func(*BatchWriter)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *BatchWriter) {
		p.logger = logger
	}
}

type worker struct {
	count int
	ch    chan *message.WriteInsert
	flush chan chan bool
}

func New(client Client, batchSize, batchSizeBytes int, batchTimeout time.Duration, opts ...Option) (*BatchWriter, error) {
	c := &BatchWriter{
		client:         client,
		workers:        make(map[string]*worker),
		logger:         zerolog.Nop(),
		batchSize:      batchSize,
		batchSizeBytes: batchSizeBytes,
		batchTimeout:   batchTimeout,
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
	sizeBytes := int64(0)
	resources := make([]*message.WriteInsert, 0, w.batchSize)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					w.flushTable(ctx, tableName, resources)
				}
				return
			}

			if (w.batchSize > 0 && len(resources) >= w.batchSize) || (w.batchSizeBytes > 0 && sizeBytes+util.TotalRecordSize(r.Record) >= int64(w.batchSizeBytes)) {
				w.flushTable(ctx, tableName, resources)
				resources, sizeBytes = resources[:0], 0
			}

			resources = append(resources, r)
			sizeBytes += util.TotalRecordSize(r.Record)
		case <-time.After(w.batchTimeout):
			if len(resources) > 0 {
				w.flushTable(ctx, tableName, resources)
				resources, sizeBytes = resources[:0], 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				w.flushTable(ctx, tableName, resources)
				resources, sizeBytes = resources[:0], 0
			}
			done <- true
		case <-ctx.Done():
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

// func (*BatchWriter) removeDuplicatesByPK(table *schema.Table, resources []*message.Insert) []*message.Insert {
// 	pkIndices := table.PrimaryKeysIndexes()
// 	// special case where there's no PK at all
// 	if len(pkIndices) == 0 {
// 		return resources
// 	}

// 	pks := make(map[string]struct{}, len(resources))
// 	res := make([]*message.Insert, 0, len(resources))
// 	for _, r := range resources {
// 		if r.Record.NumRows() > 1 {
// 			panic(fmt.Sprintf("record with more than 1 row: %d", r.Record.NumRows()))
// 		}
// 		key := pk.String(r.Record)
// 		_, ok := pks[key]
// 		if !ok {
// 			pks[key] = struct{}{}
// 			res = append(res, r)
// 			continue
// 		}
// 		// duplicate, release
// 		r.Release()
// 	}

// 	return res
// }

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
			w.flushInsert(m.Table.Name)
			w.deleteStaleLock.Lock()
			w.deleteStaleMessages = append(w.deleteStaleMessages, m)
			l := len(w.deleteStaleMessages)
			w.deleteStaleLock.Unlock()
			if w.batchSize > 0 && l > w.batchSize {
				if err := w.flushDeleteStaleTables(ctx); err != nil {
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
			l := len(w.migrateTableMessages)
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

func (w *BatchWriter) startWorker(ctx context.Context, msg *message.WriteInsert) error {
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
