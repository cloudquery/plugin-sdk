package premium

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/marketplacemetering"
	"github.com/aws/aws-sdk-go-v2/service/marketplacemetering/types"
	"github.com/aws/smithy-go"
	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/cloudquery/cloudquery-api-go/auth"
	"github.com/cloudquery/cloudquery-api-go/config"
	"github.com/cloudquery/plugin-sdk/v4/faker"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/premium/mocks"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockTokenClient struct {
	tokenType auth.TokenType
}

func (*MockTokenClient) GetToken() (auth.Token, error) {
	return auth.Token{}, nil
}

func (c *MockTokenClient) GetTokenType() auth.TokenType {
	return c.tokenType
}

func newMockTokenClient(tokenType auth.TokenType) *MockTokenClient {
	return &MockTokenClient{tokenType: tokenType}
}

func TestUsageService_NewUsageClient_Defaults(t *testing.T) {
	err := config.SetConfigHome(t.TempDir())
	require.NoError(t, err)

	err = config.SetValue("team", "config-team")
	require.NoError(t, err)

	uc, err := NewUsageClient(
		plugin.Meta{
			Team: "plugin-team",
			Kind: cqapi.PluginKindSource,
			Name: "vault",
		},
		withTokenClient(newMockTokenClient(auth.BearerToken)),
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
			Kind: cqapi.PluginKindSource,
			Name: "vault",
		},
		WithLogger(logger),
		WithAPIClient(ac),
		withTeamName("override-team-name"),
		WithMaxRetries(10),
		WithMaxWaitTime(120*time.Second),
		WithMaxTimeBetweenFlushes(10*time.Second),
		withTokenClient(newMockTokenClient(auth.BearerToken)),
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

	result, err := usageClient.CheckQuota(ctx)
	require.NoError(t, err)

	assert.False(t, result.HasQuota, "should not have quota")
}

func TestUsageService_HasQuota_WithRowsRemaining(t *testing.T) {
	ctx := context.Background()

	s := createTestServerWithRemainingRows(t, 100)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	result, err := usageClient.CheckQuota(ctx)
	require.NoError(t, err)

	assert.True(t, result.HasQuota, "should have quota")
}

func TestUsageService_Increase_ZeroBatchSize(t *testing.T) {
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

func TestUsageService_IncreaseForTable_ZeroBatchSize(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0))

	tables := 3
	rows := 9999
	for i := 0; i < rows; i++ {
		table := "table:" + strconv.Itoa(i%tables)
		err = usageClient.IncreaseForTable(table, 1)
		require.NoError(t, err)
	}

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, rows, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, rows, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")
}

func TestUsageService_Increase_WithBatchSize(t *testing.T) {
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

func TestUsageService_IncreaseForTable_WithBatchSize(t *testing.T) {
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(uint32(batchSize)))

	tables := 3
	rows := 9999
	for i := 0; i < rows; i++ {
		table := "table:" + strconv.Itoa(i%tables)
		err = usageClient.IncreaseForTable(table, 1)
		require.NoError(t, err)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, rows, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, rows, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")
	assert.True(t, true, s.minExcludingClose() > batchSize, "minimum should be greater than batch size")
}

func TestUsageService_Increase_WithFlushDuration(t *testing.T) {
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

func TestUsageService_IncreaseForTable_WithFlushDuration(t *testing.T) {
	batchSize := 2000

	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(uint32(batchSize)), WithMaxTimeBetweenFlushes(1*time.Millisecond), WithMinTimeBetweenFlushes(0*time.Millisecond))

	tables := 3
	rows := 30
	for i := 0; i < rows; i++ {
		table := "table:" + strconv.Itoa(i%tables)
		err = usageClient.IncreaseForTable(table, 1)
		require.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, rows, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, rows, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")
	assert.True(t, s.minExcludingClose() < batchSize, "we should see updates less than batchsize if ticker is firing")
}

func TestUsageService_Increase_WithMinimumUpdateDuration(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0), WithMinTimeBetweenFlushes(30*time.Second))

	rows := 10000
	for i := 0; i < rows; i++ {
		err = usageClient.Increase(1)
		require.NoError(t, err)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, rows, s.sumOfUpdates(), "total should equal number of updated rows")
	if usageClient.dataOnClose {
		assert.Equal(t, 2, s.numberOfUpdates(), "should only update first time and on close if minimum update duration is set")
	} else {
		assert.Equal(t, 1, s.numberOfUpdates(), "should only update first time if all data was processed before close")
	}
}

func TestUsageService_IncreaseForTable_WithMinimumUpdateDuration(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(0), WithMinTimeBetweenFlushes(30*time.Second))

	tables := 3
	rows := 9999
	for i := 0; i < rows; i++ {
		table := "table:" + strconv.Itoa(i%tables)
		err = usageClient.IncreaseForTable(table, 1)
		require.NoError(t, err)
	}
	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, rows, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, rows, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")
	if usageClient.dataOnClose {
		assert.Equal(t, 2, s.numberOfUpdates(), "should only update first time and on close if minimum update duration is set")
	} else {
		assert.Equal(t, 1, s.numberOfUpdates(), "should only update first time if all data was processed before close")
	}
}

func TestUsageService_IncreaseForTable_CorrectByTable(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(50))

	tables := 9
	rows := 9999
	for i := 0; i < rows; i++ {
		table := "table:" + strconv.Itoa(i%tables)
		err = usageClient.IncreaseForTable(table, 1)
		require.NoError(t, err)
	}

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, rows, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, rows, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")

	for i := 0; i < tables; i++ {
		assert.Equal(t, 1111, s.tables["table:"+strconv.Itoa(i)].Rows, "table should have correct number of rows")
	}
}

func TestUsageService_AWSMarketplaceDone(t *testing.T) {
	var err error
	ctrl := gomock.NewController(t)

	m := mocks.NewMockAWSMarketplaceClientInterface(ctrl)

	out := marketplacemetering.MeterUsageOutput{}
	inTest := meteringInput{
		marketplacemetering.MeterUsageInput{
			ProductCode:    aws.String("2a8bdkarwqrp0tmo4errl65s7"),
			UsageDimension: aws.String("rows"),
			UsageQuantity:  aws.Int32(int32(0)),
			DryRun:         aws.Bool(true)},
	}
	errTest := smithy.GenericAPIError{Code: "DryRunOperation", Message: "No errors detected in dry run"}

	in := meteringInput{
		MeterUsageInput: marketplacemetering.MeterUsageInput{
			ProductCode:    aws.String("2a8bdkarwqrp0tmo4errl65s7"),
			UsageDimension: aws.String("rows"),
			UsageQuantity:  aws.Int32(20),
			UsageAllocations: []types.UsageAllocation{{
				AllocatedUsageQuantity: aws.Int32(int32(20)),
				Tags: []types.Tag{
					{
						Key:   aws.String("plugin_name"),
						Value: aws.String("vault"),
					},
					{
						Key:   aws.String("plugin_team"),
						Value: aws.String("plugin-team"),
					},
					{
						Key:   aws.String("plugin_kind"),
						Value: aws.String("source"),
					},
				},
			}},
		},
	}
	assert.NoError(t, faker.FakeObject(&out))

	gomock.InOrder(
		m.EXPECT().MeterUsage(gomock.Any(), inTest).Return(&out, &errTest),
		m.EXPECT().MeterUsage(gomock.Any(), in).Return(&out, nil),
	)

	t.Setenv("CQ_AWS_MARKETPLACE_CONTAINER", "true")
	usageClient := newClient(t, nil, WithBatchLimit(50), WithAWSMarketplaceClient(m))

	// This will generate 19,998 rows
	// We expect that there will be 20 rows reported to AWS Marketplace
	rows := 9999
	for i := 0; i < rows; i++ {
		err = usageClient.IncreaseForTable("table", 2)
		require.NoError(t, err)
	}

	err = usageClient.Close()
	require.NoError(t, err)
}

func TestUsageService_Increase_ErrorOnMixingMethods(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(50))

	assert.Equal(t, usageClient.usageIncreaseMethod, UsageIncreaseMethodUnset, "usage increase method should not be set")

	err = usageClient.Increase(1)
	require.NoError(t, err)

	assert.Equal(t, usageClient.usageIncreaseMethod, UsageIncreaseMethodTotal, "usage increase method should be total")

	err = usageClient.IncreaseForTable("test_table", 1)
	require.ErrorContains(t, err, "mixing usage increase methods is not allowed")

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 1, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, 0, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")
}

func TestUsageService_IncreaseForTable_ErrorOnMixingMethods(t *testing.T) {
	s := createTestServer(t)
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(50))

	assert.Equal(t, usageClient.usageIncreaseMethod, UsageIncreaseMethodUnset, "usage increase method should not be set")

	err = usageClient.IncreaseForTable("test_table", 1)
	require.NoError(t, err)

	assert.Equal(t, usageClient.usageIncreaseMethod, UsageIncreaseMethodBreakdown, "usage increase method should be breakdown")

	err = usageClient.Increase(1)
	require.ErrorContains(t, err, "mixing usage increase methods is not allowed")

	err = usageClient.Close()
	require.NoError(t, err)

	assert.Equal(t, 1, s.sumOfUpdates(), "total should equal number of updated rows")
	assert.Equal(t, 1, s.sumOfTableUpdates(), "breakdown over tables should equal number of updated rows")
}

func TestUsageService_Increase_Header(t *testing.T) {
	s := createTestServer(t)
	s.headers[BatchLimitHeader] = "1000"
	s.headers[MinimumUpdateIntervalHeader] = "60"
	s.headers[MaximumUpdateIntervalHeader] = "120"
	defer s.server.Close()

	apiClient, err := cqapi.NewClientWithResponses(s.server.URL)
	require.NoError(t, err)

	usageClient := newClient(t, apiClient, WithBatchLimit(50))

	// Initial configuration
	assert.Equal(t, uint32(50), usageClient.batchLimit)
	assert.Equal(t, 10*time.Second, usageClient.minTimeBetweenFlushes)
	assert.Equal(t, 30*time.Second, usageClient.maxTimeBetweenFlushes)

	// Generate some usage
	err = usageClient.Increase(100)
	require.NoError(t, err)
	err = usageClient.Close()
	require.NoError(t, err)

	// Check the resulting configuration
	assert.Equal(t, 1, s.numberOfUpdates())
	assert.Equal(t, uint32(1000), usageClient.batchLimit)
	assert.Equal(t, 60*time.Second, usageClient.minTimeBetweenFlushes)
	assert.Equal(t, 120*time.Second, usageClient.maxTimeBetweenFlushes)
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
			Kind: cqapi.PluginKindSource,
			Name: "vault",
		},
		append(ops, withTeamName("team-name"), WithAPIClient(apiClient), withTokenClient(newMockTokenClient(auth.BearerToken)))...)
	require.NoError(t, err)

	return client.(*BatchUpdater)
}

func createTestServerWithRemainingRows(t *testing.T, remainingRows int) *testStage {
	stage := testStage{
		remainingRows: remainingRows,
		headers:       make(map[string]string),
		update:        make([]int, 0),
		tables: make(map[string]struct {
			Name string
			Rows int
		}),
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

			stage.mu.Lock()
			defer stage.mu.Unlock()
			stage.update = append(stage.update, req.Rows)

			if req.Tables != nil {
				for _, table := range *req.Tables {
					tbl, ok := stage.tables[table.Name]

					if !ok {
						stage.tables[table.Name] = struct {
							Name string
							Rows int
						}{Name: table.Name, Rows: table.Rows}
						continue
					}

					tbl.Rows += table.Rows
					stage.tables[table.Name] = tbl
				}
			}

			// Set response headers
			for k, v := range stage.headers {
				w.Header().Set(k, v)
			}

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

	headers map[string]string

	remainingRows int
	update        []int
	tables        map[string]struct {
		Name string
		Rows int
	}
	mu sync.RWMutex
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

func (s *testStage) sumOfTableUpdates() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sum := 0
	for _, val := range s.tables {
		sum += val.Rows
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

func Test_UsageClientInit_FromManagedSyncsAPIKeys(t *testing.T) {
	type testCase struct {
		name string
		envs map[string]string
		err  string
	}
	testCases := []testCase{
		{
			name: "sync run API key with team name",
			envs: map[string]string{
				auth.EnvVarCloudQueryAPIKey: "cqsr_api_key",
				"_CQ_TEAM_NAME":             "cqrn_team_name",
			},
		},
		{
			name: "sync run API key no team name",
			envs: map[string]string{
				auth.EnvVarCloudQueryAPIKey: "cqsr_api_key",
			},
			err: "failed to get team name: _CQ_TEAM_NAME environment variable not set",
		},
		{
			name: "sync test connection API key with team name",
			envs: map[string]string{
				auth.EnvVarCloudQueryAPIKey: "cqstc_api_key",
				"_CQ_TEAM_NAME":             "cqstc_team_name",
			},
		},
		{
			name: "sync test connection API key no team name",
			envs: map[string]string{
				auth.EnvVarCloudQueryAPIKey: "cqstc_api_key",
			},
			err: "failed to get team name: _CQ_TEAM_NAME environment variable not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			_, err := NewUsageClient(
				plugin.Meta{
					Team: "plugin-team",
					Kind: cqapi.PluginKindSource,
					Name: "test",
				},
			)
			if tc.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_UsageClientInit_UnknownTokenType(t *testing.T) {
	type testCase struct {
		name string
		envs map[string]string
		err  string
	}
	testCases := []testCase{
		{
			name: "unknown API key with team name",
			envs: map[string]string{
				"_CQ_TEAM_NAME": "team_name",
			},
		},
		{
			name: "unknown API key no team name",
			envs: map[string]string{},
			err:  "unsupported token type:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			_, err := NewUsageClient(
				plugin.Meta{
					Team: "plugin-team",
					Kind: cqapi.PluginKindSource,
					Name: "test",
				},
				withTokenClient(newMockTokenClient(math.MaxInt)),
			)
			if tc.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type meteringInput struct {
	marketplacemetering.MeterUsageInput
}

func (mi meteringInput) Matches(x any) bool {
	testInput, ok := x.(*marketplacemetering.MeterUsageInput)
	if !ok {
		return false
	}

	if aws.ToString(testInput.ProductCode) != aws.ToString(mi.ProductCode) {
		return false
	}
	if aws.ToString(testInput.UsageDimension) != aws.ToString(mi.UsageDimension) {
		return false
	}
	if aws.ToInt32(testInput.UsageQuantity) != aws.ToInt32(mi.UsageQuantity) {
		return false
	}
	if aws.ToBool(testInput.DryRun) != aws.ToBool(mi.DryRun) {
		return false
	}
	return true
}

func (mi meteringInput) String() string {
	return fmt.Sprintf("{ProductCode:%s UsageDimension:%s UsageQuantity:%d}", aws.ToString(mi.ProductCode), aws.ToString(mi.UsageDimension), aws.ToInt32(mi.UsageQuantity))
}

func TestRoundDown(t *testing.T) {
	cases := []struct {
		x    uint32
		unit uint32
		want uint32
	}{
		{1000, 1000, 1000},
		{2000, 1000, 2000},
		{2001, 1000, 2000},
		{2999, 1000, 2000},
		{0, 1000, 0},
		{123, 1000, 0},
		{2500, 1000, 2000},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, roundDown(c.x, c.unit))
	}
}

func TestRoundUp(t *testing.T) {
	cases := []struct {
		x    uint32
		unit uint32
		want uint32
	}{
		{1000, 1000, 1000},
		{2000, 1000, 2000},
		{2001, 1000, 3000},
		{2999, 1000, 3000},
		{0, 1000, 0},
		{123, 1000, 1000},
		{2500, 1000, 3000},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, roundUp(c.x, c.unit))
	}
}
