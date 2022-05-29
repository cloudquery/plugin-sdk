package stats

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elliotchance/orderedmap"
	"github.com/hashicorp/go-hclog"
	"github.com/segmentio/stats/v4"
)

type stat struct {
	start    time.Time
	duration time.Duration
	stopped  bool
}

type durationLogger struct {
	logger            hclog.Logger
	trackedOperations *orderedmap.OrderedMap
	mu                sync.Mutex
}

type Options struct {
	tick    time.Duration
	handler stats.Handler
}

func NewClockWithObserve(name string, tags ...stats.Tag) *stats.Clock {
	// The default clock doesn't send a measurement on start (only on stop)
	// We want both on start AND stop, so we wrap the ClockAt method
	cl := stats.DefaultEngine.ClockAt(name, time.Now(), tags...)
	stats.DefaultEngine.Observe(name, time.Duration(0), tags...)
	return cl
}

// This is executed in the context of the calling method
// We would like to keep track of still running operations, and completed operations durations
// HandleMeasures can be called by `NewClockWithObserve` which indicates a "start" of an operation
// Or by `clock.Stop` which indicates a "stop" of an operation
// We pass the measurements to a channel and periodically aggregate the data and print a hearbeat log
func (h *durationLogger) HandleMeasures(t time.Time, measures ...stats.Measure) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, m := range measures {
		id, stamp := getMeasurementDetails(m.Fields[0].Name, m.Tags)
		if stamp {
			item, ok := h.trackedOperations.Get(id)
			if ok {
				h.trackedOperations.Set(id, stat{start: item.(stat).start, duration: m.Fields[0].Value.Duration(), stopped: true})
			}
		} else {
			h.trackedOperations.Set(id, stat{start: t})
		}
	}
}

// This is executed in the context of the tick go routine
// We reported all stopped operations onces with their duration
// Still running operations are reported on each tick, until they finish
func (h *durationLogger) Flush() {
	h.mu.Lock()
	defer h.mu.Unlock()

	var durationReported []string
	for el := h.trackedOperations.Front(); el != nil; el = el.Next() {
		id := el.Key
		stat := el.Value.(stat)
		if stat.stopped {
			// `clock.Stop` was called, so we log the total duration and remove the operation from future logs
			durationReported = append(durationReported, id.(string))
			h.logger.Debug("heartbeat", "id", id, "duration", formatSeconds(stat.duration))
		} else {
			// `clock.Stop` was not called, so the operation is still running
			// We log the duration since the start of the operation
			h.logger.Debug("heartbeat", "id", id, "running_for", formatSeconds(time.Since(stat.start)))
		}
	}

	for _, id := range durationReported {
		h.trackedOperations.Delete(id)
	}
}

func Start(ctx context.Context, logger hclog.Logger, options ...func(*Options)) {
	stats.DefaultEngine.Prefix = ""

	opts := &Options{tick: time.Minute, handler: newHandler(logger)}
	for _, o := range options {
		o(opts)
	}

	logger.Debug("starting stats collector heartbeat", "tick", opts.tick)

	stats.Register(opts.handler)

	go func() {
		ticker := time.NewTicker(opts.tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats.Flush()
			}
		}
	}()
}

func WithTick(tick time.Duration) func(*Options) {
	return func(opts *Options) {
		opts.tick = tick
	}
}

func WithHandler(handler stats.Handler) func(*Options) {
	return func(opts *Options) {
		opts.handler = handler
	}
}

func Flush() {
	stats.Flush()
}

func newHandler(logger hclog.Logger) stats.Handler {
	return &durationLogger{logger: logger, trackedOperations: orderedmap.NewOrderedMap()}
}

func getMeasurementDetails(name string, tags []stats.Tag) (string, bool) {
	var stamp = false
	s := make([]string, 0, len(tags))
	s = append(s, name)
	for _, t := range tags {
		// stamp is added on `clock.Stop()`
		// we want that both `clock.Start()` and `clock.Stop()` have the same map id
		// `clock.Stop()` adds a tag named `stamp` so we don't add it to the id
		if t.Name != "stamp" {
			s = append(s, t.Name, t.Value)
		} else {
			stamp = true
		}
	}
	return strings.Join(s, ":"), stamp
}

func formatSeconds(duration time.Duration) string {
	return fmt.Sprintf("%ds", int64(duration.Seconds()))
}
