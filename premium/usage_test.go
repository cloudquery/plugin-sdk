package premium

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/cloudquery/cloudquery-api-go/auth"
	"github.com/cloudquery/cloudquery-api-go/config"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockTokenClient struct {
}

func (*MockTokenClient) GetToken() (auth.Token, error) {
	return auth.Token{}, nil
}

func (*MockTokenClient) GetTokenType() auth.TokenType {
	return auth.BearerToken
}

func TestUsageService_NewUsageClient_Defaults(t *testing.T) {
	err := config.SetConfigHome(t.TempDir())
	require.NoError(t, err)

	err = config.SetValue("team", "config-team")
	require.NoError(t, err)

	uc, err := NewUsageClient(
		plugin.Meta{
			Team: "plugin-team",
			Kind: cqapi.Source,
			Name: "vault",
		},
		withTokenClient(&MockTokenClient{}),
	)
	require.NoError(t, err)

	bu := uc.(*BatchUpdater)

	assert.NotNil(t, bu.apiClient)
	assert.Equal(t, "config-team", bu.teamName)
	assert.Equal(t, zerolog.Nop(), bu.logger)
	assert.Equal(t, 5, bu.maxRetries)
	assert.Equal(t, 60*time.Second, bu.maxWaitTime)
	assert.Equal(t, 30*time.Second, bu.maxTimeBetweenFlushes)
}

func TestUsageService_NewUsageClient_Override(t *testing.T) {
	ac, err := cqapi.NewClientWithResponses("http://localhost")
	require.NoError(t, err)

	logger := zerolog.New(zerolog.NewTestWriter(t))

	uc, err := NewUsageClient(
		plugin.Meta{
			Team: "plugin-team",
			Kind: cqapi.Source,
			Name: "vault",
		},
		WithLogger(logger),
		WithAPIClient(ac),
		withTeamName("override-team-name"),
		WithMaxRetries(10),
		WithMaxWaitTime(120*time.Second),
		WithMaxTimeBetweenFlushes(10*time.Second),
		withTokenClient(&MockTokenClient{}),
	)
	require.NoError(t, err)

	bu := uc.(*BatchUpdater)

	assert.Equal(t, ac, bu.apiClient)
	assert.Equal(t, "override-team-name", bu.teamName)
	assert.Equal(t, logger, bu.logger)
	assert.Equal(t, 10, bu.maxRetries)
	assert.Equal(t, 120*time.Second, bu.maxWaitTime)
	assert.Equal(t, 10*time.Second, bu.maxTimeBetweenFlushes)
}

func TestUsageService_HasQuota_NoRowsRemaining(t *testing.T) {
	ctx := context.Background()

	s := createTestServerWithRemainingRows(t, 0)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

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

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	hasQuota, err := usageClient.HasQuota(ctx)
	require.NoError(t, err)

	assert.True(t, hasQuota, "should have quota")
}

func TestUsageService_ZeroBatchSize(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(1)
		require.NoError(t, err)
	}

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
}

func TestUsageService_WithBatchSize(t *testing.T) {
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(uint32(batchSize)))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(1)
		require.NoError(t, err)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.True(t, true, s.minExcludingClose() > batchSize, "minimum should be greater than batch size")
}

func TestUsageService_WithFlushDuration(t *testing.T) {
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(uint32(batchSize)), WithMaxTimeBetweenFlushes(1*time.Millisecond), WithMinTimeBetweenFlushes(0*time.Millisecond))

	for i := 0; i < 10; i++ {
		err = usageClient.Increase(10)
		require.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 100, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.True(t, s.minExcludingClose() < batchSize, "we should see updates less than batchsize if ticker is firing")
}

func TestUsageService_WithMinimumUpdateDuration(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0), WithMinTimeBetweenFlushes(30*time.Second))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(1)
		require.NoError(t, err)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, 2, s.numberOfUpdates(), "should only update first time and on close if minimum update duration is set")
}

func TestUsageService_NoUpdates(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 0, s.numberOfUpdates(), "total number of updates should be zero")
}

func TestUsageService_UpdatesWithZeroRows(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	err = usageClient.Increase(0)
	require.Error(t, err, "should not be able to update with zero rows")

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 0, s.numberOfUpdates(), "total number of updates should be zero")
}

func TestUsageService_TriggerQuotaUpdate_WithNoRowsRemaining(t *testing.T) {
	ctx := context.Background()

	s := createTestServerWithRemainingRows(t, 100, 0)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	newCtx, err := WithCancelOnQuotaExceeded(ctx, usageClient, WithQuotaCheckPeriod(30*time.Second))
	require.NoError(t, err)

	go func() {
		for i := 0; i < 100; i++ {
			_ = usageClient.Increase(50)
		}
	}()

	<-newCtx.Done()
	err = context.Cause(newCtx)
	require.ErrorIs(t, ErrNoQuota{team: "team-name"}, err)
}

func TestUsageService_ShouldNotUpdateClosedService(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	// Close the service first
	err = usageClient.Close()
	require.NoError(t, err)

	err = usageClient.Increase(10)
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
			statusCode:      http.StatusServiceUnavailable,
			headers:         http.Header{},
			retry:           0,
			expectedSeconds: 1,
		},
		{
			name:            "second retry",
			statusCode:      http.StatusServiceUnavailable,
			headers:         http.Header{},
			retry:           1,
			expectedSeconds: 2,
		},
		{
			name:            "third retry",
			statusCode:      http.StatusServiceUnavailable,
			headers:         http.Header{},
			retry:           2,
			expectedSeconds: 4,
		},
		{
			name:            "fourth retry",
			statusCode:      http.StatusServiceUnavailable,
			headers:         http.Header{},
			retry:           3,
			expectedSeconds: 8,
		},
		{
			name:            "should max out at max wait time",
			statusCode:      http.StatusServiceUnavailable,
			headers:         http.Header{},
			retry:           10,
			expectedSeconds: 30,
			ops: func(client *BatchUpdater) {
				client.maxWaitTime = 30 * time.Second
			},
		},
	}

	for _, tt := range tests {
		usageClient := newClient(t, nil)
		if tt.ops != nil {
			tt.ops(usageClient)
		}
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, err := usageClient.calculateRetryDuration(tt.statusCode, tt.headers, time.Now(), tt.retry)
			require.NoError(t, err)

			assert.InDeltaf(t, tt.expectedSeconds, retryDuration.Seconds(), 1, "retry duration should be %d seconds", tt.expectedSeconds)
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
			name:            "should use exponential backoff on 503 and no header",
			statusCode:      http.StatusServiceUnavailable,
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
	}

	for _, tt := range tests {
		usageClient := newClient(t, nil)
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

			assert.InDeltaf(t, tt.expectedSeconds, retryDuration.Seconds(), 1, "retry duration should be %d seconds", tt.expectedSeconds)
		})
	}
}

func newClient(t *testing.T, apiClient *cqapi.ClientWithResponses, ops ...UsageClientOptions) *BatchUpdater {
	client, err := NewUsageClient(
		plugin.Meta{
			Team: "plugin-team",
			Kind: cqapi.Source,
			Name: "vault",
		},
		append(ops, withTeamName("team-name"), WithAPIClient(apiClient), withTokenClient(&MockTokenClient{}))...)
	require.NoError(t, err)

	return client.(*BatchUpdater)
}

func createTestServerWithRemainingRows(t *testing.T, remainingRows ...int) *testStage {
	stage := &testStage{
		remainingRows: remainingRows,
		update:        make([]int, 0),
		calls:         make(map[string]int),
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			stage.mu.Lock()
			defer stage.mu.Unlock()
			stage.calls[r.Method]++
		}()
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			methodCalls := stage.calls[r.Method]
			if methodCalls >= len(stage.remainingRows) {
				methodCalls = len(stage.remainingRows) - 1
			}
			remainingRows := stage.remainingRows[methodCalls]
			if _, err := fmt.Fprintf(w, `{"remaining_rows": %d}`, remainingRows); err != nil {
				t.Fatal(err)
			}
			return
		}
		if r.Method == "POST" {
			dec := json.NewDecoder(r.Body)
			var req cqapi.IncreaseTeamPluginUsageJSONRequestBody
			err := dec.Decode(&req)
			require.NoError(t, err)

			stage.mu.Lock()
			defer stage.mu.Unlock()
			stage.update = append(stage.update, req.Rows)

			w.WriteHeader(http.StatusOK)
			return
		}
	})

	stage.server = httptest.NewServer(handler)

	return stage
}

func createTestServer(t *testing.T) *testStage {
	return createTestServerWithRemainingRows(t, 0)
}

type testStage struct {
	server *httptest.Server

	remainingRows []int
	update        []int
	calls         map[string]int
	mu            sync.RWMutex
}

func (s *testStage) numberOfUpdates() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.update)
}

func (s *testStage) sumOfUpdates() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sum := 0
	for _, val := range s.update {
		sum += val
	}
	return sum
}

func (s *testStage) minExcludingClose() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := math.MaxInt
	for i := 0; i < len(s.update); i++ {
		if s.update[i] < m {
			m = s.update[i]
		}
	}
	return m
}
