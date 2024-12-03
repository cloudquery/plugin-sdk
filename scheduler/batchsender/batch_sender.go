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
	sendFn    func(any)
	items     []any
	timer     *time.Timer
	itemsLock sync.Mutex
}

func NewBatchSender(sendFn func(any)) *BatchSender {
	return &BatchSender{sendFn: sendFn}
}

func (bs *BatchSender) Send(item any) {
	if bs.timer != nil {
		bs.timer.Stop()
	}

	items := helpers.InterfaceSlice(item)

	// If item is already a slice, send it directly
	// together with the current batch
	if len(items) > 1 {
		bs.flush(items...)
		return
	}

	// Otherwise, add item to the current batch
	bs.appendToBatch(items...)

	// If the current batch has reached the batch size, send it
	if len(bs.items) >= batchSize {
		bs.flush()
		return
	}

	// Otherwise, start a timer to send the current batch after the batch timeout
	bs.timer = time.AfterFunc(batchTimeout, func() { bs.flush() })
}

func (bs *BatchSender) appendToBatch(items ...any) {
	bs.itemsLock.Lock()
	defer bs.itemsLock.Unlock()

	bs.items = append(bs.items, items...)
}

func (bs *BatchSender) flush(items ...any) {
	bs.itemsLock.Lock()
	defer bs.itemsLock.Unlock()

	bs.items = append(bs.items, items...)

	bs.sendFn(bs.items)
	bs.items = nil
}
