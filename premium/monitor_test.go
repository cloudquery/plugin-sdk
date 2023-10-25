package premium

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newFakeQuotaMonitor(hasQuota ...bool) *fakeQuotaMonitor {
	return &fakeQuotaMonitor{hasQuota: hasQuota}
}

type fakeQuotaMonitor struct {
	hasQuota []bool
	calls    int
}

func (f *fakeQuotaMonitor) HasQuota(_ context.Context) (bool, error) {
	hasQuota := f.hasQuota[f.calls]
	if f.calls < len(f.hasQuota)-1 {
		f.calls++
	}
	return hasQuota, nil
}

func TestWithCancelOnQuotaExceeded_NoInitialQuota(t *testing.T) {
	ctx := context.Background()

	_, _, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(false))

	require.Error(t, err)
}

func TestWithCancelOnQuotaExceeded_NoQuota(t *testing.T) {
	ctx := context.Background()

	ctx, _, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(true, false), WithQuotaCheckPeriod(1*time.Millisecond))
	require.NoError(t, err)

	<-ctx.Done()
}

func TestWithCancelOnQuotaExceeded_HasQuotaCanceled(t *testing.T) {
	ctx := context.Background()

	ctx, cancel, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(true, true, true), WithQuotaCheckPeriod(1*time.Millisecond))
	require.NoError(t, err)
	cancel()

	<-ctx.Done()
}
