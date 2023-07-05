package writers

import (
	"time"
)

type TickerFunc func(interval time.Duration) (ch <-chan time.Time, done func())

func NewTicker(interval time.Duration) (<-chan time.Time, func()) {
	if interval <= 0 {
		return nil, nop
	}
	t := time.NewTicker(interval)
	return t.C, t.Stop
}

func nop() {}
