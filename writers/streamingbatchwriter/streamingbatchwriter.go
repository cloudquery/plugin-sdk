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
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/internal/batch"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"github.com/rs/zerolog"
)

// Client is the interface that must be implemented by the client of StreamingBatchWriter.
type Client interface {
	// MigrateTable should block and handle WriteMigrateTable messages until the channel is closed or an error is returned.
	MigrateTable(context.Context, <-chan *message.WriteMigrateTable) error

	// DeleteStale should block and handle WriteDeleteStale messages until the channel is closed or an error is returned.
	DeleteStale(context.Context, <-chan *message.WriteDeleteStale) error

	// DeleteRecords should block and handle WriteDeleteRecord messages until the channel is closed or an error is returned.
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

func (w *StreamingBatchWriter) Flush(context.Context) error {
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
	return nil // not checked below
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
	return nil // not checked below
}

func (w *StreamingBatchWriter) Write(ctx context.Context, msgs <-chan message.WriteMessage) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		defer w.Close(ctx)
		for msg := range msgs {
			msgType := writers.MsgID(msg)
			if w.lastMsgType != writers.MsgTypeUnset && w.lastMsgType != msgType {
				_ = w.Flush(ctx)
			}
			w.lastMsgType = msgType
			if err := w.startWorker(ctx, errCh, msg); err != nil {
				errCh <- err
			}
		}
	}()

	var errs []error
	for err := range errCh {
		if err != nil {
			w.logger.Error().Err(err).Msg("error in streaming batch writer")
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// startWorker starts a worker for the given message type and table, or uses the existing worker if one is already running for that table.
// It returns an immediate error if the message type is not supported or the table name cannot be determined from the message.
// Errors from the running worker are sent to the errCh channel.
func (w *StreamingBatchWriter) startWorker(ctx context.Context, errCh chan<- error, msg message.WriteMessage) error {
	var tableName string

	if mi, ok := msg.(*message.WriteInsert); ok {
		md := mi.Record.Schema().Metadata()
		tableName, ok = md.GetValue(schema.MetadataTableName)
		if !ok {
			return errors.New("table name not found in metadata")
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

		w.migrateWorker = &streamingWorkerManager[*message.WriteMigrateTable]{
			ch:        make(chan *message.WriteMigrateTable),
			writeFunc: w.client.MigrateTable,

			flush: make(chan chan bool),
			errCh: errCh,

			tableName:    tableName,
			limit:        batch.CappedAt(0, w.batchSizeRows),
			batchTimeout: w.batchTimeout,
			tickerFn:     w.tickerFn,
			failed:       &atomic.Bool{},
		}

		w.workersWaitGroup.Add(1)
		go w.migrateWorker.run(ctx, &w.workersWaitGroup)
		w.migrateWorker.ch <- m

		return nil
	case *message.WriteDeleteStale:
		w.workersLock.Lock()
		defer w.workersLock.Unlock()

		if w.deleteStaleWorker != nil {
			w.deleteStaleWorker.ch <- m
			return nil
		}

		w.deleteStaleWorker = &streamingWorkerManager[*message.WriteDeleteStale]{
			ch:        make(chan *message.WriteDeleteStale),
			writeFunc: w.client.DeleteStale,
			tableName: tableName,

			flush: make(chan chan bool),
			errCh: errCh,

			limit:        batch.CappedAt(0, w.batchSizeRows),
			batchTimeout: w.batchTimeout,
			tickerFn:     w.tickerFn,
			failed:       &atomic.Bool{},
		}

		w.workersWaitGroup.Add(1)
		go w.deleteStaleWorker.run(ctx, &w.workersWaitGroup)
		w.deleteStaleWorker.ch <- m

		return nil
	case *message.WriteInsert:
		w.workersLock.RLock()
		worker, ok := w.insertWorkers[tableName]
		w.workersLock.RUnlock()
		if ok {
			worker.ch <- m
			return nil
		}

		w.workersLock.Lock()
		activeWorker, ok := w.insertWorkers[tableName]
		if ok {
			w.workersLock.Unlock()
			// some other goroutine could have already added the worker just send the message to it
			activeWorker.ch <- m
			return nil
		}

		worker = &streamingWorkerManager[*message.WriteInsert]{
			ch:        make(chan *message.WriteInsert),
			writeFunc: w.client.WriteTable,
			tableName: tableName,

			flush: make(chan chan bool),
			errCh: errCh,

			limit:        batch.CappedAt(w.batchSizeBytes, w.batchSizeRows),
			batchTimeout: w.batchTimeout,
			tickerFn:     w.tickerFn,
			failed:       &atomic.Bool{},
		}

		w.insertWorkers[tableName] = worker
		w.workersLock.Unlock()

		w.workersWaitGroup.Add(1)
		go worker.run(ctx, &w.workersWaitGroup)
		worker.ch <- m

		return nil
	case *message.WriteDeleteRecord:
		w.workersLock.Lock()
		defer w.workersLock.Unlock()

		if w.deleteRecordWorker != nil {
			w.deleteRecordWorker.ch <- m
			return nil
		}

		// TODO: flush all workers for nested tables as well (See https://github.com/cloudquery/plugin-sdk/issues/1296)
		w.deleteRecordWorker = &streamingWorkerManager[*message.WriteDeleteRecord]{
			ch:        make(chan *message.WriteDeleteRecord),
			writeFunc: w.client.DeleteRecords,
			tableName: tableName,

			flush: make(chan chan bool),
			errCh: errCh,

			limit:        batch.CappedAt(w.batchSizeBytes, w.batchSizeRows),
			batchTimeout: w.batchTimeout,
			tickerFn:     w.tickerFn,
			failed:       &atomic.Bool{},
		}

		w.workersWaitGroup.Add(1)
		go w.deleteRecordWorker.run(ctx, &w.workersWaitGroup)
		w.deleteRecordWorker.ch <- m

		return nil
	default:
		return fmt.Errorf("unhandled message type: %T", msg)
	}
}

// streamingWorkerManager manages a worker that processes messages of type T for table tableName.
type streamingWorkerManager[T message.WriteMessage] struct {
	ch        chan T
	writeFunc func(context.Context, <-chan T) error
	tableName string

	flush chan chan bool
	errCh chan<- error

	limit        *batch.Cap
	batchTimeout time.Duration
	tickerFn     writers.TickerFunc
	failed       *atomic.Bool
	workerWg     sync.WaitGroup

	inputCh chan T
	mu      sync.Mutex // protects inputCh
}

func (s *streamingWorkerManager[T]) closeFlush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.inputCh != nil {
		close(s.inputCh)
		s.inputCh = nil
		s.limit.Reset()
	}
}

func (s *streamingWorkerManager[T]) send(ctx context.Context, data T) {
	// Don't create new goroutines if we're in a failed state
	if s.failed.Load() {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.inputCh != nil {
		select {
		case <-ctx.Done():
			return
		case s.inputCh <- data:
		}
		return
	}

	s.inputCh = make(chan T)
	s.workerWg.Add(1)

	// start consuming our new channel
	go func(ch chan T) {
		defer s.workerWg.Done()
		defer func() {
			if msg := recover(); msg != nil {
				switch v := msg.(type) {
				case error:
					s.errCh <- fmt.Errorf("panic processing %s: %w [recovered]", s.tableName, v)
				default:
					s.errCh <- fmt.Errorf("panic processing %s: %v [recovered]", s.tableName, msg)
				}
			}
		}()
		defer func() { // modified closeFlush
			s.mu.Lock()
			defer s.mu.Unlock()
			if s.inputCh == ch { // only close if we're still the active channel
				close(s.inputCh)
				s.inputCh = nil
			}
		}()

		err := s.writeFunc(ctx, ch)
		if err != nil {
			s.failed.Store(true)
			go func() {
				// nolint:revive
				for range ch { // drain the channel to avoid deadlock
				}
			}()
			select {
			case <-ctx.Done():
				return
			default:
				s.errCh <- fmt.Errorf("handler failed on %s: %w", s.tableName, err)
			}
		}
	}(s.inputCh)

	select {
	case <-ctx.Done():
		return
	case s.inputCh <- data:
	}
}

func (s *streamingWorkerManager[T]) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.workerWg.Wait()
	defer s.closeFlush()

	ticker := s.tickerFn(s.batchTimeout)
	defer ticker.Stop()

	tickerCh, ctxDone := ticker.Chan(), ctx.Done()

	for {
		select {
		case r, ok := <-s.ch:
			if !ok {
				return
			}
			if ins, ok := any(r).(*message.WriteInsert); ok {
				add, toFlush, rest := batch.SliceRecord(ins.Record, s.limit)
				if add != nil {
					s.limit.AddSlice(add)
					s.send(ctx, any(&message.WriteInsert{Record: add.RecordBatch}).(T))
				}
				if len(toFlush) > 0 || rest != nil || s.limit.ReachedLimit() {
					// flush current batch
					s.closeFlush()
					ticker.Reset(s.batchTimeout)
				}
				for _, sliceToFlush := range toFlush {
					s.limit.AddRows(sliceToFlush.NumRows())
					s.send(ctx, any(&message.WriteInsert{Record: sliceToFlush}).(T))
					s.closeFlush()
					ticker.Reset(s.batchTimeout)
				}

				// set the remainder
				if rest != nil {
					s.limit.AddSlice(rest)
					s.send(ctx, any(&message.WriteInsert{Record: rest.RecordBatch}).(T))
				}
			} else {
				s.send(ctx, r)
				s.limit.AddRows(1)
				if s.limit.ReachedLimit() {
					s.closeFlush()
					ticker.Reset(s.batchTimeout)
				}
			}

		case <-tickerCh:
			if s.limit.Rows() > 0 {
				s.closeFlush()
			}
		case done := <-s.flush:
			if s.limit.Rows() > 0 {
				s.closeFlush()
				ticker.Reset(s.batchTimeout)
			}
			done <- true
		case <-ctxDone:
			// this means the request was cancelled
			return // after this NO other call will succeed
		}
	}
}
