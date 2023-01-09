package destination

import (
	"sync/atomic"
)

type Metrics struct {
	// Errors number of errors / failed writes
	Errors uint64
	// Writes number of successful writes
	Writes uint64
}

func (mx *Metrics) Get() Metrics {
	return Metrics{
		Errors: atomic.LoadUint64(&mx.Errors),
		Writes: atomic.LoadUint64(&mx.Writes),
	}
}

func (mx *Metrics) Failed(amount int) {
	atomic.AddUint64(&mx.Errors, uint64(amount))
}

func (mx *Metrics) Wrote(amount int) {
	atomic.AddUint64(&mx.Writes, uint64(amount))
}
