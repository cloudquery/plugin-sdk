package loggedlimiter

import (
	"context"
	"iter"
	"sync"
	"time"
)

type LoggedLimiter struct {
	Capacity int // Number of workers

	ch      chan struct{}
	logLock sync.Mutex // Protects access to current
	logs    []LogEntry
}

type LogEntry struct {
	Time         time.Time
	UsedCapacity int
}

func New(capacity int) *LoggedLimiter {
	limiter := &LoggedLimiter{
		Capacity: capacity,
		ch:       make(chan struct{}, capacity), // Channel capacity is fixed
	}

	return limiter
}

func (l *LoggedLimiter) Acquire(ctx context.Context, _ int) error {
	select {
	case l.ch <- struct{}{}:
		l._log()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *LoggedLimiter) Release(_ int) {
	<-l.ch // Free up a slot
	l._log()
}

func (l *LoggedLimiter) _log() {
	l.logLock.Lock()
	defer l.logLock.Unlock()
	l.logs = append(l.logs, LogEntry{Time: time.Now(), UsedCapacity: len(l.ch)})
}

func (l *LoggedLimiter) Logs() iter.Seq[LogEntry] {
	return func(yield func(LogEntry) bool) {
		for i := range l.logs {
			l.logLock.Lock()
			log := l.logs[i]
			l.logLock.Unlock()
			if !yield(log) {
				return
			}
		}
	}
}
