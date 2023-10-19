package premium

import (
	"context"
	"encoding/json"
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

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(0))

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

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(0))

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

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(0))

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

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(uint32(batchSize)))

	for i := 0; i < 10000; i++ {
		err = usageClient.Increase(ctx, 1)
		require.NoError(t, err)
	}
	err = usageClient.Close(ctx)
	require.NoError(t, err)

	assert.Equal(t, 10000, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.True(t, true, s.minExcludingClose() > batchSize, "minimum should be greater than batch size")
}

func TestUsageService_WithTicker(t *testing.T) {
	ctx := context.Background()
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(uint32(batchSize)), WithTickerDuration(1))

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

func TestUsageService_NoUpdates(t *testing.T) {
	ctx := context.Background()

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(0))

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

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(0))

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

	usageClient := NewUsageClient(ctx, apiClient, "myteam", "mnorbury-team", "source", "vault", WithBatchLimit(0))

	// Close the service first
	err = usageClient.Close(ctx)
	require.NoError(t, err)

	err = usageClient.Increase(ctx, 10)
	require.Error(t, err, "should not be able to update closed service")

	assert.Equal(t, 0, s.numberOfUpdates(), "total number of updates should be zero")
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
