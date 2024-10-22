package queue

import (
	"sync"
	"sync/atomic"
	"time"
)

// activeWorkSignal is a thread-safe coordinator for awaiting a worker pool
// that relies on a queue that might be temporarily empty.
//
// If queue is empty and workers idle, done!
//
// If the queue is empty but a worker is working on a task, we must wait and check
// if it's empty after that worker finishes. That's why we need this.
//
// Use it like this:
//
// - When a worker picks up a task, call `Add()` (like a WaitGroup)
// - When a worker finishes a task, call `Done()` (like a WaitGroup)
//
// - If the queue is empty, check `IsIdle()` to check if no workers are active.
// - If workers are still active, call `Wait()` to block until state changes.
type activeWorkSignal struct {
	countChangeSignal   *sync.Cond
	activeWorkUnitCount *atomic.Int32
	isStarted           *atomic.Bool
}

func newActiveWorkSignal() *activeWorkSignal {
	return &activeWorkSignal{
		countChangeSignal:   sync.NewCond(&sync.Mutex{}),
		activeWorkUnitCount: &atomic.Int32{},
		isStarted:           &atomic.Bool{},
	}
}

// Add means a worker has started working on a task.
//
// Wake up the work queuing goroutine.
func (s *activeWorkSignal) Add() {
	s.activeWorkUnitCount.Add(1)
	s.isStarted.Store(true)
	s.countChangeSignal.Signal()
}

// Done means a worker has finished working on a task.
//
// If the count became zero, wake up the work queuing goroutine (might have finished).
func (s *activeWorkSignal) Done() {
	s.activeWorkUnitCount.Add(-1)
	s.countChangeSignal.Signal()
}

// IsIdle returns true if no workers are active. If queue is empty and workers idle, done!
func (s *activeWorkSignal) IsIdle() bool {
	return s.isStarted.Load() && s.activeWorkUnitCount.Load() <= 0
}

// Wait blocks until the count of active workers changes.
func (s *activeWorkSignal) Wait() {
	// A race condition is possible when the last active table asynchronously
	// queues a relation. The table finishes (calling `.Done()`) a moment
	// before the queue receives the `.Push()`. At this point, the queue is
	// empty and there are no active workers.
	//
	// A moment later, the queue receives the `.Push()` and queues a new task.
	//
	// This is a very infrequent case according to tests, but it happens.
	time.Sleep(10 * time.Millisecond)

	if s.activeWorkUnitCount.Load() <= 0 {
		return
	}
	s.countChangeSignal.L.Lock()
	defer s.countChangeSignal.L.Unlock()
	s.countChangeSignal.Wait()
}
