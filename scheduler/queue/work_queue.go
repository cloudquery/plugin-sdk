package queue

import (
	"context"
	"sync"
)

// workQueue holds the work units waiting for a worker and counts the ones in progress.
//
// The sync is done when no work unit is queued and none is in progress. Workers queue
// relations while their table is in progress, so an empty queue alone does not mean done.
//
//   - Push queues a work unit.
//   - Next takes a work unit off the queue and marks it in progress.
//   - Done marks a work unit taken by Next as finished.
type workQueue struct {
	mu         sync.Mutex
	queue      *ConcurrentRandomQueue[WorkUnit]
	inProgress int
	// Buffered with capacity 1 so a change that happens before Next starts to wait is kept rather than lost.
	changed chan struct{}
}

func newWorkQueue(seed int64, capacityHint int) *workQueue {
	return &workQueue{
		queue:   NewConcurrentRandomQueue[WorkUnit](seed, capacityHint),
		changed: make(chan struct{}, 1),
	}
}

func (q *workQueue) Push(item WorkUnit) {
	q.mu.Lock()
	q.queue.Push(item)
	q.mu.Unlock()
	q.notify()
}

// Next blocks until a work unit is queued and returns it marked in progress.
// It returns nil when no work unit is queued or in progress, or when ctx is done.
func (q *workQueue) Next(ctx context.Context) *WorkUnit {
	for ctx.Err() == nil {
		q.mu.Lock()
		item := q.queue.Pop()
		if item != nil {
			q.inProgress++
		}
		idle := item == nil && q.inProgress == 0
		q.mu.Unlock()

		if item != nil {
			return item
		}
		if idle {
			return nil
		}
		select {
		case <-q.changed:
		case <-ctx.Done():
		}
	}
	return nil
}

func (q *workQueue) Done() {
	q.mu.Lock()
	q.inProgress--
	q.mu.Unlock()
	q.notify()
}

func (q *workQueue) notify() {
	select {
	case q.changed <- struct{}{}:
	default:
	}
}
