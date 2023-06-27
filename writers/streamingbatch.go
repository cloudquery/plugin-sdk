package writers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type StreamingBatchWriterClient interface {
	MigrateTables(context.Context, []*message.MigrateTable) error
	DeleteStale(context.Context, []*message.DeleteStale) error
	WriteTable(context.Context, <-chan *message.Insert) error
}

type StreamingBatchWriter struct {
	client           StreamingBatchWriterClient
	workers          map[string]*streamingbatchworker
	workersLock      sync.RWMutex
	workersWaitGroup sync.WaitGroup

	migrateTableLock     sync.Mutex
	migrateTableMessages []*message.MigrateTable
	deleteStaleLock      sync.Mutex
	deleteStaleMessages  []*message.DeleteStale

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSizeRows  int64
	batchSizeBytes int64
}

type StreamingBatchWriterOption func(*StreamingBatchWriter)

func WithStreamingBatchWriterLogger(logger zerolog.Logger) StreamingBatchWriterOption {
	return func(p *StreamingBatchWriter) {
		p.logger = logger
	}
}

func WithStreamingBatchWriterBatchTimeout(timeout time.Duration) StreamingBatchWriterOption {
	return func(p *StreamingBatchWriter) {
		p.batchTimeout = timeout
	}
}

func WithStreamingBatchWriterBatchSizeRows(size int64) StreamingBatchWriterOption {
	return func(p *StreamingBatchWriter) {
		p.batchSizeRows = size
	}
}

func WithStreamingBatchWriterBatchSizeBytes(size int64) StreamingBatchWriterOption {
	return func(p *StreamingBatchWriter) {
		p.batchSizeBytes = size
	}
}

type streamingbatchworker struct {
	count int
	ch    chan *message.Insert
	flush chan chan bool
}

func NewStreamingBatchWriter(client StreamingBatchWriterClient, opts ...StreamingBatchWriterOption) (*StreamingBatchWriter, error) {
	c := &StreamingBatchWriter{
		client:       client,
		workers:      make(map[string]*streamingbatchworker),
		logger:       zerolog.Nop(),
		batchTimeout: defaultBatchTimeoutSeconds * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.migrateTableMessages = make([]*message.MigrateTable, 0, c.batchSizeRows)
	c.deleteStaleMessages = make([]*message.DeleteStale, 0, c.batchSizeRows)
	return c, nil
}

func (w *StreamingBatchWriter) Flush(ctx context.Context) error {
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

func (w *StreamingBatchWriter) stopWorkers() {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()
	for _, w := range w.workers {
		close(w.ch)
	}
	w.workersWaitGroup.Wait()
	w.workers = make(map[string]*streamingbatchworker)
}

func (w *StreamingBatchWriter) worker(ctx context.Context, tableName string, ch <-chan *message.Insert, errCh chan<- error, flush <-chan chan bool) {
	var sizeBytes, sizeRows int64
	opened := false

	var (
		clientCh    chan *message.Insert
		clientErrCh chan error
	)

	doOpen := func() {
		clientCh = make(chan *message.Insert)
		clientErrCh = make(chan error)
		go func() {
			clientErrCh <- w.client.WriteTable(ctx, clientCh)
			close(clientErrCh)
		}()
	}
	doClose := func() {
		if opened {
			close(clientCh)
			if err := <-clientErrCh; err != nil {
				errCh <- fmt.Errorf("WriteTable failed on %s: %w", tableName, err)
			}
		}
		opened = false
		sizeBytes, sizeRows = 0, 0
	}

	for {
		select {
		case r, ok := <-ch:
			if !ok {
				doClose()
				return
			}

			if (w.batchSizeRows > 0 && sizeRows >= w.batchSizeRows) || (w.batchSizeBytes > 0 && sizeBytes+util.TotalRecordSize(r.Record) >= w.batchSizeBytes) {
				doClose()
			}

			if !opened {
				doOpen()
				opened = true
			}

			clientCh <- r
			sizeRows++
			sizeBytes += util.TotalRecordSize(r.Record)
		case <-time.After(w.batchTimeout):
			if sizeRows > 0 {
				doClose()
			}
		case done := <-flush:
			if sizeRows > 0 {
				doClose()
			}
			done <- true
		}
	}
}

func (w *StreamingBatchWriter) flushMigrateTables(ctx context.Context) error {
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

func (w *StreamingBatchWriter) flushDeleteStaleTables(ctx context.Context) error {
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

func (w *StreamingBatchWriter) flushInsert(_ context.Context, tableName string) {
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

func (w *StreamingBatchWriter) Write(ctx context.Context, msgs <-chan message.Message) error {
	errCh := make(chan error)

	go func() {
		for err := range errCh {
			w.logger.Err(err).Msg("error from StreamingBatchWriter")
		}
	}()

	hasWorkers := false
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
			if w.batchSizeRows > 0 && int64(l) > w.batchSizeRows {
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
			hasWorkers = true
			if err := w.startWorker(ctx, errCh, m); err != nil {
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
			if w.batchSizeRows > 0 && int64(l) > w.batchSizeRows {
				if err := w.flushMigrateTables(ctx); err != nil {
					return err
				}
			}
		}
	}

	if err := w.Flush(ctx); err != nil {
		return err
	}

	if hasWorkers {
		w.stopWorkers()
	}
	close(errCh)
	return nil
}

func (w *StreamingBatchWriter) startWorker(ctx context.Context, errCh chan<- error, msg *message.Insert) error {
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
	wr = &streamingbatchworker{
		count: 1,
		ch:    ch,
		flush: flush,
	}
	w.workers[tableName] = wr
	w.workersLock.Unlock()
	w.workersWaitGroup.Add(1)
	go func() {
		defer w.workersWaitGroup.Done()
		w.worker(ctx, tableName, ch, errCh, flush)
	}()
	ch <- msg
	return nil
}
