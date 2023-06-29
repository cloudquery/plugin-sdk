// Package streamingbatchwriter provides a writers.Writer implementation that writes to a client that implements the streamingbatchwriter.Client interface.
//
// Write messages are sent to the client with three separate methods: MigrateTable, WriteTable, and DeleteStale. Each method is called separate goroutines.
// Message types are processed in blocks: Receipt of a new message type will cause the previous message type processing to end (if it exists) which is signalled
// to the handler by closing the channel. The handler should return after processing all messages.
//
// For Insert messages (handled by WriteTable) each table creates separate goroutine. Number of goroutines is limited by the number of tables.
// Thus, each WriteTable invocation is for a single table (all messages sent to WriteTable are guaranteed to be for the same table).
//
// After a 'batch' is complete, the channel is closed. The handler is expected to block until the channel is closed and to keep processing in a streaming fashion.
//
// Batches are considered complete when:
// 1. The batch timeout is reached
// 2. The batch size is reached
// 3. The batch size in bytes is reached
// 4. A different message type is received
//
// Each handler can get invoked multiple times as new batches are processed.
// Handlers get invoked only if there's a message of that type at hand: First message of the batch is immediately available in the channel.
package streamingbatchwriter

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

// Client is the interface that must be implemented by the client of StreamingBatchWriter.
type Client interface {
	// MigrateTable should block and handle WriteMigrateTable messages until the channel is closed.
	MigrateTable(context.Context, <-chan *message.WriteMigrateTable) error

	// DeleteStale should block and handle WriteDeleteStale messages until the channel is closed.
	DeleteStale(context.Context, <-chan *message.WriteDeleteStale) error

	// WriteTable should block and handle writes to a single table until the channel is closed. Table metadata can be found in the first WriteInsert message.
	// The channel is closed when all inserts in the batch have been sent. New batches, if any, will be sent on a new call to WriteTable.
	WriteTable(context.Context, <-chan *message.WriteInsert) error
}

type StreamingBatchWriter struct {
	client Client

	insertWorkers    map[string]*streamingWorkerManager[*message.WriteInsert]
	migrateWorker    *streamingWorkerManager[*message.WriteMigrateTable]
	deleteWorker     *streamingWorkerManager[*message.WriteDeleteStale]
	workersLock      sync.RWMutex
	workersWaitGroup sync.WaitGroup

	lastMsgType writers.MsgType

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSizeRows  int64
	batchSizeBytes int64

	timerFn timerFn
}

type timerFn func(timeout time.Duration) <-chan time.Time

// Assert at compile-time that StreamingBatchWriter implements the Writer interface
var _ writers.Writer = (*StreamingBatchWriter)(nil)

type Option func(*StreamingBatchWriter)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *StreamingBatchWriter) {
		p.logger = logger
	}
}

func WithBatchTimeout(timeout time.Duration) Option {
	return func(p *StreamingBatchWriter) {
		p.batchTimeout = timeout
	}
}

func WithBatchSizeRows(size int64) Option {
	return func(p *StreamingBatchWriter) {
		p.batchSizeRows = size
	}
}

func WithBatchSizeBytes(size int64) Option {
	return func(p *StreamingBatchWriter) {
		p.batchSizeBytes = size
	}
}

func withTimerFn(timer timerFn) Option {
	return func(p *StreamingBatchWriter) {
		p.timerFn = timer
	}
}

func New(client Client, opts ...Option) (*StreamingBatchWriter, error) {
	c := &StreamingBatchWriter{
		client:         client,
		insertWorkers:  make(map[string]*streamingWorkerManager[*message.WriteInsert]),
		logger:         zerolog.Nop(),
		batchTimeout:   writers.DefaultBatchTimeoutSeconds * time.Second,
		batchSizeRows:  writers.DefaultBatchSize,
		batchSizeBytes: writers.DefaultBatchSizeBytes,
		timerFn:        timer,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (w *StreamingBatchWriter) Flush(_ context.Context) error {
	w.workersLock.RLock()
	if w.migrateWorker != nil {
		done := make(chan bool)
		w.migrateWorker.flush <- done
		<-done
	}
	if w.deleteWorker != nil {
		done := make(chan bool)
		w.deleteWorker.flush <- done
		<-done
	}
	for _, worker := range w.insertWorkers {
		done := make(chan bool)
		worker.flush <- done
		<-done
	}
	w.workersLock.RUnlock()
	return nil
}

func (w *StreamingBatchWriter) stopWorkers() {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()
	for _, w := range w.insertWorkers {
		close(w.ch)
	}
	if w.migrateWorker != nil {
		close(w.migrateWorker.ch)
	}
	if w.deleteWorker != nil {
		close(w.deleteWorker.ch)
	}
	w.workersWaitGroup.Wait()

	w.insertWorkers = make(map[string]*streamingWorkerManager[*message.WriteInsert])
	w.migrateWorker = nil
	w.deleteWorker = nil
}

func (w *StreamingBatchWriter) Write(ctx context.Context, msgs <-chan message.WriteMessage) error {
	errCh := make(chan error)

	go func() {
		for err := range errCh {
			w.logger.Err(err).Msg("error from StreamingBatchWriter")
		}
	}()

	hasWorkers := false

	for msg := range msgs {
		msgType := writers.MsgID(msg)
		if w.lastMsgType != msgType {
			if err := w.Flush(ctx); err != nil {
				return err
			}
		}
		w.lastMsgType = msgType
		hasWorkers = true
		if err := w.startWorker(ctx, errCh, msg); err != nil {
			return err
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

func (w *StreamingBatchWriter) startWorker(ctx context.Context, errCh chan<- error, msg message.WriteMessage) error {
	var tableName string

	if mi, ok := msg.(*message.WriteInsert); ok {
		md := mi.Record.Schema().Metadata()
		tableName, ok = md.GetValue(schema.MetadataTableName)
		if !ok {
			return fmt.Errorf("table name not found in metadata")
		}
	} else {
		tableName = msg.GetTable().Name
	}

	switch m := msg.(type) {
	case *message.WriteMigrateTable:
		w.workersLock.Lock()
		defer w.workersLock.Unlock()
		if w.migrateWorker != nil {
			w.migrateWorker.ch <- m
			return nil
		}
		ch := make(chan *message.WriteMigrateTable)
		flush := make(chan chan bool)
		w.migrateWorker = &streamingWorkerManager[*message.WriteMigrateTable]{
			ch:        ch,
			writeFunc: w.client.MigrateTable,

			flush: flush,
			errCh: errCh,

			batchSizeRows: w.batchSizeRows,
			batchTimeout:  w.batchTimeout,
			timerFn:       w.timerFn,
		}

		w.workersWaitGroup.Add(1)
		go w.migrateWorker.run(ctx, &w.workersWaitGroup, tableName)
		w.migrateWorker.ch <- m
		return nil
	case *message.WriteDeleteStale:
		w.workersLock.Lock()
		defer w.workersLock.Unlock()
		if w.deleteWorker != nil {
			w.deleteWorker.ch <- m
			return nil
		}
		ch := make(chan *message.WriteDeleteStale)
		flush := make(chan chan bool)
		w.deleteWorker = &streamingWorkerManager[*message.WriteDeleteStale]{
			ch:        ch,
			writeFunc: w.client.DeleteStale,

			flush: flush,
			errCh: errCh,

			batchSizeRows: w.batchSizeRows,
			batchTimeout:  w.batchTimeout,
		}

		w.workersWaitGroup.Add(1)
		go w.deleteWorker.run(ctx, &w.workersWaitGroup, tableName)
		w.deleteWorker.ch <- m
		return nil
	case *message.WriteInsert:
		w.workersLock.RLock()
		wr, ok := w.insertWorkers[tableName]
		w.workersLock.RUnlock()
		if ok {
			wr.ch <- m
			return nil
		}

		ch := make(chan *message.WriteInsert)
		flush := make(chan chan bool)
		wr = &streamingWorkerManager[*message.WriteInsert]{
			ch:        ch,
			writeFunc: w.client.WriteTable,

			flush: flush,
			errCh: errCh,

			batchSizeRows:  w.batchSizeRows,
			batchSizeBytes: w.batchSizeBytes,
			batchTimeout:   w.batchTimeout,
			timerFn:        w.timerFn,
		}
		w.workersLock.Lock()
		w.insertWorkers[tableName] = wr
		w.workersLock.Unlock()

		w.workersWaitGroup.Add(1)
		go wr.run(ctx, &w.workersWaitGroup, tableName)
		ch <- m
		return nil

	default:
		return fmt.Errorf("unhandled message type: %T", msg)
	}
}

type streamingWorkerManager[T message.WriteMessage] struct {
	ch        chan T
	writeFunc func(context.Context, <-chan T) error

	flush chan chan bool
	errCh chan<- error

	batchSizeRows  int64
	batchSizeBytes int64
	batchTimeout   time.Duration
	timerFn        timerFn
}

func (s *streamingWorkerManager[T]) run(ctx context.Context, wg *sync.WaitGroup, tableName string) {
	defer wg.Done()
	var (
		clientCh            chan T
		clientErrCh         chan error
		open                bool
		sizeBytes, sizeRows int64
	)

	ensureOpened := func() {
		if open {
			return
		}

		clientCh = make(chan T)
		clientErrCh = make(chan error, 1)
		go func() {
			defer close(clientErrCh)
			defer func() {
				if err := recover(); err != nil {
					clientErrCh <- fmt.Errorf("panic: %v", err)
				}
			}()
			clientErrCh <- s.writeFunc(ctx, clientCh)
		}()
		open = true
	}
	closeFlush := func() {
		if open {
			close(clientCh)
			if err := <-clientErrCh; err != nil {
				s.errCh <- fmt.Errorf("handler failed on %s: %w", tableName, err)
			}
		}
		open = false
		sizeBytes, sizeRows = 0, 0
	}
	defer closeFlush()

	for {
		select {
		case r, ok := <-s.ch:
			if !ok {
				return
			}

			var recSize int64
			if ins, ok := any(r).(*message.WriteInsert); ok {
				recSize = util.TotalRecordSize(ins.Record)
			}

			if (s.batchSizeRows > 0 && sizeRows >= s.batchSizeRows) || (s.batchSizeBytes > 0 && sizeBytes+recSize >= s.batchSizeBytes) {
				closeFlush()
			}

			ensureOpened()
			clientCh <- r
			sizeRows++
			sizeBytes += recSize
		case <-s.timerFn(s.batchTimeout):
			if sizeRows > 0 {
				closeFlush()
			}
		case done := <-s.flush:
			if sizeRows > 0 {
				closeFlush()
			}
			done <- true
		}
	}
}

func timer(timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return nil
	}
	return time.After(timeout)
}
