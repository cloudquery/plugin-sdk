package batchsender

import (
	"sync"
	"testing"
	"time"
)

// This test verifies there is no data race between Send() and the timer-triggered flush.
func TestSend_ConcurrentWithTimerFlush(_ *testing.T) {
	// The race occurs when:
	//   1. a Send() call schedules a timer via time.AfterFunc
	//   2. the timer fires and calls flush() on a separate goroutine
	//   3. another Send() reads bs.items concurrently.
	//
	// To trigger this, we send items from multiple goroutines with delays around batchTimeout so the timer fires between Sends.
	var mu sync.Mutex
	var received []any

	const numGoroutines = 5
	const sendsPerGoroutine = 20

	bs := NewBatchSender(func(items any) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, items)
	})

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range sendsPerGoroutine {
				bs.Send("item")
				time.Sleep(batchTimeout + 10*time.Millisecond)
			}
		}()
	}

	wg.Wait()
	bs.Close()
}
