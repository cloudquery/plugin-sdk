package streamingbatchwriter

import "time"

type mockTimer struct {
	expire chan time.Time
}

func (t *mockTimer) timer(d time.Duration) <-chan time.Time {
	return t.expire
}

func newMockTimer() (timerFn, chan time.Time) {
	expire := make(chan time.Time)
	t := &mockTimer{
		expire: expire,
	}
	return t.timer, expire
}
