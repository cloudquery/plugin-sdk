package helpers

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// SemaphoreAcauireMax is calling TryAcquire with 1 until it fails and return n-TryAcquireSucess
func TryAcquireMax(ctx context.Context, sem *semaphore.Weighted, n int64) (int64, error) {
	newN := n
	for {
		if newN <= 0 {
			return newN, nil
		}
		// first try will be blocking
		if newN == n {
			if err := sem.Acquire(ctx, 1); err != nil {
				return newN, err
			}
			newN--
			continue
		}
		if !sem.TryAcquire(1) {
			return newN, nil
		}
		newN--
	}
}
