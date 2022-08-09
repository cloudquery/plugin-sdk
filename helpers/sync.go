package helpers

import (
	"golang.org/x/sync/semaphore"
)

// SemaphoreAcauireMax is calling TryAcquire with 1 until it fails and return n-TryAcquireSucess
func TryAcquireMax(sem *semaphore.Weighted, n int64) int64 {
	newN := n
	for {
		if newN <= 0 {
			return newN
		}
		if !sem.TryAcquire(1) {
			return newN
		}
		newN--
	}
}
