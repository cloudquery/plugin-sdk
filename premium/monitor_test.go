package premium

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type quotaResponse struct {
	hasQuota bool
	err      error
}

func newFakeQuotaMonitor(hasQuota ...quotaResponse) *fakeQuotaMonitor {
	return &fakeQuotaMonitor{responses: hasQuota}
}

type fakeQuotaMonitor struct {
	responses []quotaResponse
	calls     int
}

func (f *fakeQuotaMonitor) HasQuota(_ context.Context) (bool, error) {
	resp := f.responses[f.calls]
	if f.calls < len(f.responses)-1 {
		f.calls++
	}
	return resp.hasQuota, resp.err
}

func (*fakeQuotaMonitor) TeamName() string {
	return "test"
}

func TestWithCancelOnQuotaExceeded_NoInitialQuota(t *testing.T) {
	ctx := context.Background()

	responses := []quotaResponse{
		{false, nil},
	}
	_, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(responses...))

	require.Error(t, err)
}

func TestWithCancelOnQuotaExceeded_NoQuota(t *testing.T) {
	ctx := context.Background()

	responses := []quotaResponse{
		{true, nil},
		{false, nil},
	}
	ctx, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(responses...), WithQuotaCheckPeriod(1*time.Millisecond))
	require.NoError(t, err)

	<-ctx.Done()
	cause := context.Cause(ctx)
	require.ErrorIs(t, ErrNoQuota{team: "test"}, cause)
}

func TestWithCancelOnQuotaCheckConsecutiveFailures(t *testing.T) {
	ctx := context.Background()

	responses := []quotaResponse{
		{true, nil},
		{false, errors.New("test2")},
		{false, errors.New("test3")},
	}
	ctx, err := WithCancelOnQuotaExceeded(ctx,
		newFakeQuotaMonitor(responses...),
		WithQuotaCheckPeriod(1*time.Millisecond),
		WithQuotaMaxConsecutiveFailures(2),
	)
	require.NoError(t, err)
	<-ctx.Done()
	cause := context.Cause(ctx)
	require.Equal(t, "test2\ntest3", cause.Error())
}
