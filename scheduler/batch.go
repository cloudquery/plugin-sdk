package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"github.com/rs/zerolog"
)

type (
	BatchSettings struct {
		MaxRows int
		Timeout time.Duration
	}

	BatchOption func(settings *BatchSettings)
)

func WithBatchOptions(options ...BatchOption) Option {
	return func(s *Scheduler) {
		if s.batchSettings == nil {
			s.batchSettings = new(BatchSettings)
		}
		for _, o := range options {
			o(s.batchSettings)
		}
	}
}

func WithBatchMaxRows(rows int) BatchOption {
	return func(s *BatchSettings) {
		s.MaxRows = rows
	}
}

func WithBatchTimeout(timeout time.Duration) BatchOption {
	return func(s *BatchSettings) {
		s.Timeout = timeout
	}
}

func (s *BatchSettings) getBatcher(ctx context.Context, res chan<- message.SyncMessage, logger zerolog.Logger) batcherInterface {
	if s.Timeout > 0 && s.MaxRows > 1 {
		return &batcher{
			done:    ctx.Done(),
			res:     res,
			maxRows: s.MaxRows,
			timeout: s.Timeout,
			logger:  logger.With().Int("max_rows", s.MaxRows).Dur("timeout_ms", s.Timeout).Logger(),
		}
	}

	return &nopBatcher{res: res}
}

type batcherInterface interface {
	process(res *schema.Resource)
	close()
}

type nopBatcher struct {
	res chan<- message.SyncMessage
}

func (n *nopBatcher) process(resource *schema.Resource) {
	n.res <- &message.SyncInsert{Record: resource.GetValues().ToArrowRecord(resource.Table.ToArrowSchema())}
}

func (*nopBatcher) close() {}

var _ batcherInterface = (*nopBatcher)(nil)

type batcher struct {
	done <-chan struct{}

	res chan<- message.SyncMessage

	maxRows int
	timeout time.Duration

	// using sync primitives by value here implies that batcher is to be used by pointer only
	// workers is a sync.Map rather than a map + mutex pair
	// because worker allocation & lookup falls into one of the sync.Map use-cases,
	// namely, ever-growing cache (write once, read many times).
	workers sync.Map // k = table name, v = *worker
	wg      sync.WaitGroup

	logger zerolog.Logger
}

type worker struct {
	ch               chan *schema.Resource
	flush            chan chan struct{}
	curRows, maxRows int
	builder          *array.RecordBuilder // we can reuse that
	res              chan<- message.SyncMessage

	// debug logging
	tableName string
	logger    *zerolog.Logger
}

// send must be called on len(rows) > 0
func (w *worker) send() {
	w.logger.Debug().Str("table", w.tableName).Int("rows", w.curRows).Msg("send")
	w.res <- &message.SyncInsert{Record: w.builder.NewRecord()}
	// we need to reserve here as NewRecord (& underlying NewArray calls) reset the memory
	w.builder.Reserve(w.maxRows)
	w.curRows = 0 // reset
}

func (w *worker) work(done <-chan struct{}, timeout time.Duration) {
	ticker := writers.NewTicker(timeout)
	defer ticker.Stop()
	tickerCh := ticker.Chan()

	for {
		select {
		case r, ok := <-w.ch:
			if !ok {
				if w.curRows > 0 {
					w.send()
				}
				return
			}

			// append to builder
			scalar.AppendToRecordBuilder(w.builder, r.GetValues())
			w.curRows++
			// check if we need to flush
			if w.maxRows > 0 && w.curRows == w.maxRows {
				w.send()
				ticker.Reset(timeout)
			}

		case <-tickerCh:
			if w.curRows > 0 {
				w.send()
			}

		case ch := <-w.flush:
			if w.curRows > 0 {
				w.send()
				ticker.Reset(timeout)
			}
			close(ch)

		case <-done:
			// this means the request was cancelled
			return // after this NO other call will succeed
		}
	}
}

func (b *batcher) process(res *schema.Resource) {
	table := res.Table
	// already running worker
	v, loaded := b.workers.Load(table.Name)
	if loaded {
		v.(*worker).ch <- res
		return
	}

	// we alloc only ch here, as it may be needed right away
	// for instance, if another goroutine will get the value allocated by us
	wr := &worker{ch: make(chan *schema.Resource, 5)} // 5 is quite enough
	v, loaded = b.workers.LoadOrStore(table.Name, wr)
	if loaded {
		// means that the worker was already in tne sync.Map, so we just discard the wr value
		close(wr.ch)          // for GC
		v.(*worker).ch <- res // send res to the already allocated worker
		return
	}

	// fill in the required data
	// start wr
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		// fill in the worker fields
		wr.flush = make(chan chan struct{})
		wr.maxRows = b.maxRows
		wr.builder = array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
		wr.res = b.res
		wr.builder.Reserve(b.maxRows)
		wr.logger = &b.logger
		wr.tableName = table.Name

		// start processing
		wr.work(b.done, b.timeout)
	}()

	wr.ch <- res
}

func (b *batcher) close() {
	b.workers.Range(func(_, v any) bool {
		close(v.(*worker).ch)
		return true
	})
	b.wg.Wait()
}
