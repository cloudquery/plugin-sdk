package batchsender

import (
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/helpers"
)

const (
	batchSize    = 100
	batchTimeout = 100 * time.Millisecond
)

// BatchSender is a helper struct that batches items and sends them in batches of batchSize or after batchTimeout.
//
// - If item is already a slice, it will be sent directly
// - Otherwise, it will be added to the current batch
// - If the current batch has reached the batch size, it will be sent immediately
// - Otherwise, a timer will be started to send the current batch after the batch timeout
type BatchSender struct {
	sendFn func(any)
	items  []any
	timer  *time.Timer
	mu     sync.Mutex
}

func NewBatchSender(sendFn func(any)) *BatchSender {
	return &BatchSender{sendFn: sendFn}
}

func (bs *BatchSender) Send(item any) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.timer != nil {
		bs.timer.Stop()
	}

	items := helpers.InterfaceSlice(item)

	// If item is already a slice, send it directly
	// together with the current batch
	if len(items) > 1 {
		bs.flushLocked(items...)
		return
	}

	// Otherwise, add item to the current batch
	bs.items = append(bs.items, items...)

	// If the current batch has reached the batch size, send it
	if len(bs.items) >= batchSize {
		bs.flushLocked()
		return
	}

	// Otherwise, start a timer to send the current batch after the batch timeout
	bs.timer = time.AfterFunc(batchTimeout, func() {
		bs.mu.Lock()
		defer bs.mu.Unlock()
		bs.flushLocked()
	})
}

// flushLocked sends all buffered items. Must be called with bs.mu held.
func (bs *BatchSender) flushLocked(items ...any) {
	bs.items = append(bs.items, items...)

	if len(bs.items) == 0 {
		return
	}

	bs.sendFn(bs.items)
	bs.items = nil
}

func (bs *BatchSender) Close() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.timer != nil {
		bs.timer.Stop()
	}
	bs.flushLocked()
}