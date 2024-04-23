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

	"github.com/apache/arrow/go/v15/arrow/util"
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

	// DeleteRecords should block and handle WriteDeleteRecord messages until the channel is closed.
	DeleteRecords(context.Context, <-chan *message.WriteDeleteRecord) error

	// WriteTable should block and handle writes to a single table until the channel is closed. Table metadata can be found in the first WriteInsert message.
	// The channel is closed when all inserts in the batch have been sent. New batches, if any, will be sent on a new call to WriteTable.
	WriteTable(context.Context, <-chan *message.WriteInsert) error
}

type StreamingBatchWriter struct {
	client Client

	insertWorkers      map[string]*streamingWorkerManager[*message.WriteInsert]
	migrateWorker      *streamingWorkerManager[*message.WriteMigrateTable]
	deleteStaleWorker  *streamingWorkerManager[*message.WriteDeleteStale]
	deleteRecordWorker *streamingWorkerManager[*message.WriteDeleteRecord]

	workersLock      sync.RWMutex
	workersWaitGroup sync.WaitGroup

	lastMsgType writers.MsgType

	logger         zerolog.Logger
	batchTimeout   time.Duration
	batchSizeRows  int64
	batchSizeBytes int64

	tickerFn writers.TickerFunc
}

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

func withTickerFn(tickerFn writers.TickerFunc) Option {
	return func(p *StreamingBatchWriter) {
		p.tickerFn = tickerFn
	}
}

const (
	defaultBatchTimeoutSeconds = 20
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

func New(client Client, opts ...Option) (*StreamingBatchWriter, error) {
	c := &StreamingBatchWriter{
		client:         client,
		insertWorkers:  make(map[string]*streamingWorkerManager[*message.WriteInsert]),
		logger:         zerolog.Nop(),
		batchTimeout:   defaultBatchTimeoutSeconds * time.Second,
		batchSizeRows:  defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
		tickerFn:       writers.NewTicker,
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
	if w.deleteStaleWorker != nil {
		done := make(chan bool)
		w.deleteStaleWorker.flush <- done
		<-done
	}
	if w.deleteRecordWorker != nil {
		done := make(chan bool)
		w.deleteRecordWorker.flush <- done
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

func (w *StreamingBatchWriter) Close(context.Context) error {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()
	for _, w := range w.insertWorkers {
		close(w.ch)
	}
	if w.migrateWorker != nil {
		close(w.migrateWorker.ch)
	}
	if w.deleteStaleWorker != nil {
		close(w.deleteStaleWorker.ch)
	}
	if w.deleteRecordWorker != nil {
		close(w.deleteRecordWorker.ch)
	}
	w.workersWaitGroup.Wait()

	w.insertWorkers = make(map[string]*streamingWorkerManager[*message.WriteInsert])
	w.migrateWorker = nil
	w.deleteStaleWorker = nil
	w.deleteRecordWorker = nil
	w.lastMsgType = writers.MsgTypeUnset

	return nil
}

func (w *StreamingBatchWriter) Write(ctx context.Context, msgs <-chan message.WriteMessage) error {
	errCh := make(chan error)
	defer close(errCh)

	go func() {
		for err := range errCh {
			w.logger.Err(err).Msg("error from StreamingBatchWriter")
		}
	}()

	for msg := range msgs {
		msgType := writers.MsgID(msg)
		if w.lastMsgType != writers.MsgTypeUnset && w.lastMsgType != msgType {
			if err := w.Flush(ctx); err != nil {
				return err
			}
		}
		w.lastMsgType = msgType
		if err := w.startWorker(ctx, errCh, msg); err != nil {
			return err
		}
	}

	return w.Close(ctx)
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
		w.migrateWorker = &streamingWorkerManager[*message.WriteMigrateTable]{
			ch:        ch,
			writeFunc: w.client.MigrateTable,

			flush: make(chan chan bool),
			errCh: errCh,

			batchSizeRows: w.batchSizeRows,
			batchTimeout:  w.batchTimeout,
			tickerFn:      w.tickerFn,
		}

		w.workersWaitGroup.Add(1)
		go w.migrateWorker.run(ctx, &w.workersWaitGroup, tableName)
		w.migrateWorker.ch <- m
		return nil
	case *message.WriteDeleteStale:
		w.workersLock.Lock()
		defer w.workersLock.Unlock()
		if w.deleteStaleWorker != nil {
			w.deleteStaleWorker.ch <- m
			return nil
		}
		ch := make(chan *message.WriteDeleteStale)
		w.deleteStaleWorker = &streamingWorkerManager[*message.WriteDeleteStale]{
			ch:        ch,
			writeFunc: w.client.DeleteStale,

			flush: make(chan chan bool),
			errCh: errCh,

			batchSizeRows: w.batchSizeRows,
			batchTimeout:  w.batchTimeout,
			tickerFn:      w.tickerFn,
		}

		w.workersWaitGroup.Add(1)
		go w.deleteStaleWorker.run(ctx, &w.workersWaitGroup, tableName)
		w.deleteStaleWorker.ch <- m
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
		wr = &streamingWorkerManager[*message.WriteInsert]{
			ch:        ch,
			writeFunc: w.client.WriteTable,

			flush: make(chan chan bool),
			errCh: errCh,

			batchSizeRows:  w.batchSizeRows,
			batchSizeBytes: w.batchSizeBytes,
			batchTimeout:   w.batchTimeout,
			tickerFn:       w.tickerFn,
		}
		w.workersLock.Lock()
		wrOld, ok := w.insertWorkers[tableName]
		if ok {
			w.workersLock.Unlock()
			// some other goroutine could have already added the worker
			// just send the message to it & discard our allocated worker
			wrOld.ch <- m
			return nil
		}
		w.insertWorkers[tableName] = wr
		w.workersLock.Unlock()

		w.workersWaitGroup.Add(1)
		go wr.run(ctx, &w.workersWaitGroup, tableName)
		ch <- m
		return nil
	case *message.WriteDeleteRecord:
		w.workersLock.Lock()
		defer w.workersLock.Unlock()
		if w.deleteRecordWorker != nil {
			w.deleteRecordWorker.ch <- m
			return nil
		}
		ch := make(chan *message.WriteDeleteRecord)
		// TODO: flush all workers for nested tables as well (See https://github.com/cloudquery/plugin-sdk/issues/1296)
		w.deleteRecordWorker = &streamingWorkerManager[*message.WriteDeleteRecord]{
			ch:        ch,
			writeFunc: w.client.DeleteRecords,

			flush: make(chan chan bool),
			errCh: errCh,

			batchSizeRows: w.batchSizeRows,
			batchTimeout:  w.batchTimeout,
			tickerFn:      w.tickerFn,
		}

		w.workersWaitGroup.Add(1)
		go w.deleteRecordWorker.run(ctx, &w.workersWaitGroup, tableName)
		w.deleteRecordWorker.ch <- m
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
	tickerFn       writers.TickerFunc
}

func (s *streamingWorkerManager[T]) run(ctx context.Context, wg *sync.WaitGroup, tableName string) {
	defer wg.Done()
	var (
		clientCh    chan T
		clientErrCh chan error
		open        bool
		bytes, rows = writers.NewCapped(s.batchSizeBytes), writers.NewCapped(s.batchSizeRows)
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
		bytes.Reset()
		rows.Reset()
	}
	defer closeFlush()

	ticker := s.tickerFn(s.batchTimeout)
	defer ticker.Stop()
	for {
		select {
		case r, ok := <-s.ch:
			if !ok {
				return
			}

			var recSize int64
			rowSize := int64(1) // at least 1 row for messages without records
			if ins, ok := any(r).(*message.WriteInsert); ok {
				recSize = util.TotalRecordSize(ins.Record)
				rowSize = ins.Record.NumRows()
			}

			ensureOpened()
			clientCh <- r
			rows.Add(rowSize)
			bytes.Add(recSize)

			if rows.ReachedLimit() || bytes.ReachedLimit() {
				closeFlush()
				ticker.Reset(s.batchTimeout)
			}
		case <-ticker.Chan():
			if rows.Current() > 0 {
				closeFlush()
			}
		case done := <-s.flush:
			if rows.Current() > 0 {
				closeFlush()
				ticker.Reset(s.batchTimeout)
			}
			done <- true
		}
	}
}
