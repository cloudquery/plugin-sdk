package metrics

import (
	"sync"
	"time"
)

type durationMeasurement struct {
	startTime time.Time
	started   bool
	duration  time.Duration
	sem       sync.Mutex
}

func (dm *durationMeasurement) Start(start time.Time) {
	// If we have already started, don't start again. This can happen for relational tables that are resolved multiple times (per parent resource)
	dm.sem.Lock()
	defer dm.sem.Unlock()
	if dm.started {
		return
	}

	dm.started = true
	dm.startTime = start
}

// End calculates, updates and returns the delta duration for updating OTEL counters.
func (dm *durationMeasurement) End(end time.Time) time.Duration {
	var delta time.Duration
	newDuration := end.Sub(dm.startTime)

	dm.sem.Lock()
	defer dm.sem.Unlock()

	delta = newDuration - dm.duration
	dm.duration = newDuration
	return delta
}
