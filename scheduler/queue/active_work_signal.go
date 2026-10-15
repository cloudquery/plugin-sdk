package queue

import (
	"sync/atomic"
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
//   - When the dispatcher takes a task off the queue to hand it to a worker, call `Add()` (like a WaitGroup).
//     Marking the task active before the handoff ensures the idle check cannot fire while
//     the task is in flight between the queue and a worker.
//   - When a worker finishes a task, call `Done()` (like a WaitGroup)
//
// - If the queue is empty, check `IsIdle()` to check if no workers are active.
// - If workers are still active, call `Wait()` to block until one of them finishes.
type activeWorkSignal struct {
	activeWorkUnitCount atomic.Int32
	isStarted           atomic.Bool
	// Buffered with capacity 1 so a Done that happens before Wait is kept rather than lost.
	workDone chan struct{}
}

func newActiveWorkSignal() *activeWorkSignal {
	return &activeWorkSignal{workDone: make(chan struct{}, 1)}
}

// Add means the dispatcher has taken a task off the queue for a worker.
func (s *activeWorkSignal) Add() {
	s.activeWorkUnitCount.Add(1)
	s.isStarted.Store(true)
}

// Done means a worker has finished working on a task.
//
// Wake up the work queuing goroutine: the worker might have queued new tasks, or all work might have finished.
func (s *activeWorkSignal) Done() {
	s.activeWorkUnitCount.Add(-1)
	select {
	case s.workDone <- struct{}{}:
	default:
	}
}

// IsIdle returns true if no workers are active. If queue is empty and workers idle, done!
func (s *activeWorkSignal) IsIdle() bool {
	return s.isStarted.Load() && s.activeWorkUnitCount.Load() <= 0
}

// Wait blocks until a worker finishes a task. It returns immediately if one has
// finished since the previous Wait returned.
func (s *activeWorkSignal) Wait() {
	<-s.workDone
}
