package streamingbatchwriter

import "time"

type mockTimer struct {
	expire chan time.Time
}

func (t *mockTimer) timer(time.Duration) (<-chan time.Time, func()) {
	return t.expire, func() {}
}

func newMockTimer() (timerFn, chan time.Time) {
	expire := make(chan time.Time)
	t := &mockTimer{
		expire: expire,
	}
	return t.timer, expire
}
