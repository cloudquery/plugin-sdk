package premium

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUsageService_HasQuota_NoRowsRemaining(t *testing.T) {
	ctx := context.Background()

	s := createTestServerWithRemainingRows(t, 0)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0))

	hasQuota, err := usageClient.HasQuota(ctx)
	require.NoError(t, err)

	assert.False(t, hasQuota, "should not have quota")
}

func TestUsageService_HasQuota_WithRowsRemaining(t *testing.T) {
	ctx := context.Background()

	s := createTestServerWithRemainingRows(t, 100)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0))

	hasQuota, err := usageClient.HasQuota(ctx)
	require.NoError(t, err)

	assert.True(t, hasQuota, "should have quota")
}

func TestUsageService_ZeroBatchSize(t *testing.T) {
	ctx := context.Background()

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(ctx, 1)
		require.NoError(t, err)
	}

	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
}

func TestUsageService_WithBatchSize(t *testing.T) {
	ctx := context.Background()
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(uint32(batchSize)))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(ctx, 1)
		require.NoError(t, err)
	}
	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.True(t, true, s.minExcludingClose() > batchSize, "minimum should be greater than batch size")
}

func TestUsageService_WithFlushDuration(t *testing.T) {
	ctx := context.Background()
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(uint32(batchSize)), WithMaxTimeBetweenFlushes(1*time.Millisecond), WithMinTimeBetweenFlushes(0*time.Millisecond))

	for i := 0; i < 10; i++ {
		err = usageClient.Increase(ctx, 10)
		require.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}
	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 100, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.True(t, s.minExcludingClose() < batchSize, "we should see updates less than batchsize if ticker is firing")
}

func TestUsageService_WithMinimumUpdateDuration(t *testing.T) {
	ctx := context.Background()

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0), WithMinTimeBetweenFlushes(30*time.Second))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(ctx, 1)
		require.NoError(t, err)
	}
	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, 2, s.numberOfUpdates(), "should only update first time and on close if minimum update duration is set")
}

func TestUsageService_NoUpdates(t *testing.T) {
	ctx := context.Background()

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0))

	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 0, s.numberOfUpdates(), "total number of updates should be zero")
}

func TestUsageService_UpdatesWithZeroRows(t *testing.T) {
	ctx := context.Background()

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0))

	err = usageClient.Increase(ctx, 0)
	require.Error(t, err, "should not be able to update with zero rows")

	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 0, s.numberOfUpdates(), "total number of updates should be zero")
}

func TestUsageService_ShouldNotUpdateClosedService(t *testing.T) {
	ctx := context.Background()

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(ctx, apiClient, WithBatchLimit(0))

	// Close the service first
	err = usageClient.Close(ctx)
	require.NoError(t, err)

	err = usageClient.Increase(ctx, 10)
	require.Error(t, err, "should not be able to update closed service")

	assert.Equal(t, 0, s.numberOfUpdates(), "total number of updates should be zero")
}

func TestUsageService_CalculateRetryDuration_Exp(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		headers         http.Header
		retry           int
		expectedSeconds int
		ops             func(client *BatchUpdater)
	}{
		{
			name:            "first retry",
			statusCode:      http.StatusOK,
			headers:         http.Header{},
			retry:           0,
			expectedSeconds: 1,
		},
		{
			name:            "second retry",
			statusCode:      http.StatusOK,
			headers:         http.Header{},
			retry:           1,
			expectedSeconds: 2,
		},
		{
			name:            "third retry",
			statusCode:      http.StatusOK,
			headers:         http.Header{},
			retry:           2,
			expectedSeconds: 4,
		},
		{
			name:            "fourth retry",
			statusCode:      http.StatusOK,
			headers:         http.Header{},
			retry:           3,
			expectedSeconds: 8,
		},
		{
			name:            "should max out at max wait time",
			statusCode:      http.StatusOK,
			headers:         http.Header{},
			retry:           10,
			expectedSeconds: 30,
			ops: func(client *BatchUpdater) {
				client.maxWaitTime = 30 * time.Second
			},
		},
	}

	for _, tt := range tests {
		usageClient := newClient(context.Background(), nil)
		if tt.ops != nil {
			tt.ops(usageClient)
		}
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, err := usageClient.calculateRetryDuration(tt.statusCode, tt.headers, time.Now(), tt.retry)
			require.NoError(t, err)

			assert.InDeltaf(t, tt.expectedSeconds, retryDuration.Seconds(), 0.1, "retry duration should be %d seconds", tt.expectedSeconds)
		})
	}
}

func TestUsageService_CalculateRetryDuration_ServerBackPressure(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		headers         http.Header
		retry           int
		expectedSeconds int
		ops             func(client *BatchUpdater)
		wantErr         error
	}{
		{
			name:            "should use exponential backoff on 200",
			statusCode:      http.StatusOK,
			headers:         http.Header{},
			retry:           0,
			expectedSeconds: 1,
		},
		{
			name:            "should use exponential backoff on 429 if no retry-after header",
			statusCode:      http.StatusTooManyRequests,
			headers:         http.Header{},
			retry:           1,
			expectedSeconds: 2,
		},
		{
			name:            "should use retry-after header if present on 429",
			statusCode:      http.StatusTooManyRequests,
			headers:         http.Header{"Retry-After": []string{"5"}},
			retry:           0,
			expectedSeconds: 5,
		},
		{
			name:       "should raise an error if the server wants us to wait longer than max wait time",
			statusCode: http.StatusTooManyRequests,
			headers:    http.Header{"Retry-After": []string{"40"}},
			retry:      0,
			ops: func(client *BatchUpdater) {
				client.maxWaitTime = 30 * time.Second
			},
			wantErr: errors.New("retry-after header exceeds max wait time: 40s > 30s"),
		},
	}

	for _, tt := range tests {
		usageClient := newClient(context.Background(), nil)
		if tt.ops != nil {
			tt.ops(usageClient)
		}
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, err := usageClient.calculateRetryDuration(tt.statusCode, tt.headers, time.Now(), tt.retry)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			}

			assert.InDeltaf(t, tt.expectedSeconds, retryDuration.Seconds(), 0.1, "retry duration should be %d seconds", tt.expectedSeconds)
		})
	}
}

func newClient(ctx context.Context, apiClient *cqapi.ClientWithResponses, ops ...UpdaterOptions) *BatchUpdater {
	return NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", ops...)
}

func createTestServerWithRemainingRows(t *testing.T, remainingRows int) *testStage {
	stage := testStage{
		remainingRows: remainingRows,
		update:        make([]int, 0),
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			if _, err := fmt.Fprintf(w, `{"remaining_rows": %d}`, stage.remainingRows); err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "POST" {
			dec := json.NewDecoder(r.Body)
			var req cqapi.IncreaseTeamPluginUsageJSONRequestBody
			err := dec.Decode(&req)
			require.NoError(t, err)

			stage.update = append(stage.update, req.Rows)

			w.WriteHeader(http.StatusOK)
			return
		}
	})

	stage.server = httptest.NewServer(handler)

	return &stage
}

func createTestServer(t *testing.T) *testStage {
	return createTestServerWithRemainingRows(t, 0)
}

type testStage struct {
	server *httptest.Server

	remainingRows int
	update        []int
}

func (s *testStage) numberOfUpdates() int {
	return len(s.update)
}

func (s *testStage) sumOfUpdates() int {
	sum := 0
	for _, val := range s.update {
		sum += val
	}
	return sum
}

func (s *testStage) minExcludingClose() int {
	m := math.MaxInt
	for i := 0; i < len(s.update); i++ {
		if s.update[i] < m {
			m = s.update[i]
		}
	}
	return m
}
