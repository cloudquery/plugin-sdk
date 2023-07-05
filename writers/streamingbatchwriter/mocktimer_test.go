package streamingbatchwriter

import (
	"time"

	"github.com/cloudquery/plugin-sdk/v4/writers"
)

type mockTimer struct {
	expire chan time.Time
}

func (t *mockTimer) timer(time.Duration) (<-chan time.Time, func()) {
	return t.expire, t.close
}

func (t *mockTimer) close() {
	close(t.expire)
}

func newMockTimer() (writers.TickerFunc, chan time.Time) {
	expire := make(chan time.Time)
	t := &mockTimer{
		expire: expire,
	}
	return t.timer, expire
}
