package streamingbatchwriter

import (
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/writers"
)

type mockTicker struct {
	expire  chan time.Time
	stopped sync.Once
}

func (t *mockTicker) Stop() {
	t.stopped.Do(func() {
		close(t.expire)
	})
}

func (t *mockTicker) Tick() {
	t.expire <- time.Now()
}

func (*mockTicker) Reset(time.Duration) {}

func (t *mockTicker) Chan() <-chan time.Time {
	return t.expire
}

func newMockTicker() (writers.TickerFunc, func()) {
	expire := make(chan time.Time)
	t := &mockTicker{
		expire: expire,
	}
	return func(time.Duration) writers.Ticker {
		return t
	}, t.Tick
}
