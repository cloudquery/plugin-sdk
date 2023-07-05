package writers

import (
	"time"
)

type TickerFunc func(time.Duration) Ticker

type Ticker interface {
	Stop()
	Reset(d time.Duration)
	Chan() <-chan time.Time
}

func NewTicker(interval time.Duration) Ticker {
	if interval <= 0 {
		return nopTicker{}
	}
	return &ticker{time.NewTicker(interval)}
}

type ticker struct {
	*time.Ticker
}

func (t *ticker) Chan() <-chan time.Time {
	return t.C
}

type nopTicker struct{}

func (nopTicker) Stop() {}

func (nopTicker) Reset(_ time.Duration) {}

func (nopTicker) Chan() <-chan time.Time {
	return nil
}
