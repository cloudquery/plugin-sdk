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
)

type batcher struct {
	ctx     context.Context
	ctxDone <-chan struct{}

	res chan<- message.SyncMessage

	maxRows int
	timeout time.Duration

	// using sync primitives by value here implies that batcher is to be used by pointer only
	// workers is a sync.Map rather than a map + mutex pair
	// because worker allocation & lookup falls into one of the sync.Map use-cases,
	// namely, ever-growing cache (write once, read many times).
	workers sync.Map // k = table name, v = *worker
	wg      sync.WaitGroup
}

type worker struct {
	ch               chan *schema.Resource
	flush            chan chan struct{}
	curRows, maxRows int
	builder          *array.RecordBuilder // we can reuse that
	res              chan<- message.SyncMessage
}

// send must be called on len(rows) > 0
func (w *worker) send() {
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
			if !ok || r.TableDone() {
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

		// start processing
		wr.work(b.ctxDone, b.timeout)
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

func newBatcher(ctx context.Context, res chan<- message.SyncMessage, maxRows int, timeout time.Duration) *batcher {
	return &batcher{
		ctx:     ctx,
		ctxDone: ctx.Done(),
		res:     res,
		maxRows: maxRows,
		timeout: timeout,
	}
}
