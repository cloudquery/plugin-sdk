package premium

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type remainingQuotaResponse struct {
	remainingRows *int64
	err           error
}

func newFakeQuotaMonitor(hasQuota ...remainingQuotaResponse) *fakeQuotaMonitor {
	return &fakeQuotaMonitor{responses: hasQuota}
}

type fakeQuotaMonitor struct {
	responses []remainingQuotaResponse
	calls     int
}

func (*fakeQuotaMonitor) HasQuota(_ context.Context) (bool, error) {
	return true, nil
}

func (f *fakeQuotaMonitor) RemainingRows(_ context.Context) (*int64, error) {
	resp := f.responses[f.calls]
	if f.calls < len(f.responses)-1 {
		f.calls++
	}
	return resp.remainingRows, resp.err
}

func (*fakeQuotaMonitor) TeamName() string {
	return "test"
}

func TestWithCancelOnQuotaExceeded_NoInitialQuota(t *testing.T) {
	ctx := context.Background()

	responses := []remainingQuotaResponse{
		{int64ptr(0), nil},
	}
	_, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(responses...))

	require.Error(t, err)
}

func TestWithCancelOnQuotaExceeded_NoQuota(t *testing.T) {
	ctx := context.Background()

	responses := []remainingQuotaResponse{
		{int64ptr(1000), nil},
		{int64ptr(0), nil},
	}
	ctx, err := WithCancelOnQuotaExceeded(ctx, newFakeQuotaMonitor(responses...), WithQuotaCheckPeriod(1*time.Millisecond))
	require.NoError(t, err)

	<-ctx.Done()
	cause := context.Cause(ctx)
	require.ErrorIs(t, ErrNoQuota{team: "test"}, cause)
}

func TestWithCancelOnQuotaCheckConsecutiveFailures(t *testing.T) {
	ctx := context.Background()

	responses := []remainingQuotaResponse{
		{int64ptr(1000), nil},
		{int64ptr(0), errors.New("test2")},
		{int64ptr(0), errors.New("test3")},
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

func int64ptr(i int64) *int64 {
	return &i
}
