package queue

import (
	"testing"
	"time"
)

func TestActiveWorkSignalWaitReturnsAfterEarlierDone(t *testing.T) {
	s := newActiveWorkSignal()
	s.Add()
	s.Add()
	s.Done()

	waited := make(chan struct{})
	go func() {
		s.Wait()
		close(waited)
	}()

	select {
	case <-waited:
	case <-time.After(5 * time.Second):
		t.Fatal("Wait missed a Done that happened before it was called")
	}
}
