package helpers

import (
	"context"
	"time"
)

// Sleep pauses for the given duration or aborts immediately if the given context is canceled
func Sleep(ctx context.Context, dur time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(dur):
		return nil
	}
}
