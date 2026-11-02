package queue

import (
	"math/rand"
	"sync"
)

// ConcurrentRandomQueue is a generic, thread-safe queue
// that pops random elements in O(1) time.
type ConcurrentRandomQueue[T any] struct {
	mu     sync.Mutex
	queue  []T
	random *rand.Rand
}

func NewConcurrentRandomQueue[T any](seed int64, capacityHint int) *ConcurrentRandomQueue[T] {
	return &ConcurrentRandomQueue[T]{queue: make([]T, 0, capacityHint), random: rand.New(rand.NewSource(seed))}
}

func (q *ConcurrentRandomQueue[T]) Push(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, item)
}

func (q *ConcurrentRandomQueue[T]) Pop() *T {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return nil
	}
	idx := q.random.Intn(len(q.queue))
	lastIdx := len(q.queue) - 1
	q.queue[idx], q.queue[lastIdx] = q.queue[lastIdx], q.queue[idx]
	item := q.queue[lastIdx]
	q.queue = q.queue[:lastIdx]

	return &item
}
