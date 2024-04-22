package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
)

type batcher struct {
	res chan<- message.SyncMessage

	size    int
	timeout time.Duration

	// using sync primitives by value here implies that batcher is to be used by pointer only
	workers sync.Map // k = table name, v = *worker
	wg      sync.WaitGroup
}

type worker struct {
	ch      chan *schema.Resource
	flush   chan chan struct{}
	rows    schema.Resources
	builder *array.RecordBuilder // we can reuse that
	res     chan<- message.SyncMessage
}

// send must be called on len(rows) > 0
func (w *worker) send() {
	w.builder.Reserve(len(w.rows)) // prealloc

	for _, row := range w.rows {
		scalar.AppendToRecordBuilder(w.builder, row.GetValues())
	}

	w.rows = w.rows[:0]
	w.res <- &message.SyncInsert{Record: w.builder.NewRecord()}
}

func (w *worker) work(ctx context.Context, size int, timeout time.Duration) {
	ticker := writers.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case r, ok := <-w.ch:
			if !ok {
				if len(w.rows) > 0 {
					w.send()
				}
				return
			}

			w.rows = append(w.rows, r)
			if size > 0 && len(w.rows) == size {
				w.send()
				ticker.Reset(timeout)
			}

		case <-ticker.Chan():
			if len(w.rows) > 0 {
				w.send()
			}

		case done := <-w.flush:
			if len(w.rows) > 0 {
				w.send()
				ticker.Reset(timeout)
			}
			close(done)

		case <-ctx.Done():
			// this means the request was cancelled
			return // after this NO other call will succeed
		}
	}
}

func (b *batcher) worker(ctx context.Context, res *schema.Resource) {
	table := res.Table
	v, loaded := b.workers.Load(table.Name)
	if !loaded {
		newWorker := &worker{
			ch:      make(chan *schema.Resource),
			flush:   make(chan chan struct{}),
			rows:    make(schema.Resources, 0, b.size),
			builder: array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema()),
			res:     b.res,
		}
		// we use LoadOrStore as there may be other goroutine doing the same
		v, loaded = b.workers.LoadOrStore(table.Name, newWorker)
		if loaded {
			// discard newWorker
			close(newWorker.ch)
			close(newWorker.flush)
			newWorker.builder.Release()
		}
	}
	wr := v.(*worker)
	if !loaded {
		// now loaded basically determines if we need to call the goroutine.
		// Note that loaded will be false only for the goroutine that actually successfully stores the worker
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			wr.work(ctx, b.size, b.timeout)
		}()
	}

	wr.ch <- res
}

func (b *batcher) close() {
	b.workers.Range(func(_, v any) bool {
		close(v.(*worker).ch)
		return true
	})
	b.wg.Wait()
}

func newBatcher(res chan<- message.SyncMessage, size int, timeout time.Duration) *batcher {
	return &batcher{
		res:     res,
		size:    size,
		timeout: timeout,
	}
}
