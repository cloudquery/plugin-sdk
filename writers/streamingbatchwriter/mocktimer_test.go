package streamingbatchwriter

import (
	"time"

	"github.com/cloudquery/plugin-sdk/v4/writers"
)

type mockTicker struct {
	expire chan time.Time
}

func (t *mockTicker) Stop() {
	close(t.expire)
}

func (t *mockTicker) Reset(time.Duration) {}

func (t *mockTicker) Chan() <-chan time.Time {
	return t.expire
}

func newMockTicker() (writers.TickerFunc, chan<- time.Time) {
	expire := make(chan time.Time)
	t := &mockTicker{
		expire: expire,
	}
	return func(time.Duration) writers.Ticker {
		return t
	}, expire
}
