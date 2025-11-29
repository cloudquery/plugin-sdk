package premium

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/marketplacemetering"
	"github.com/aws/aws-sdk-go-v2/service/marketplacemetering/types"
	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/cloudquery/cloudquery-api-go/auth"
	"github.com/cloudquery/cloudquery-api-go/config"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
)

const (
	defaultAPIURL                = "https://api.cloudquery.io"
	defaultBatchLimit            = 1000
	defaultMaxRetries            = 5
	defaultMaxWaitTime           = 60 * time.Second
	defaultMinTimeBetweenFlushes = 10 * time.Second
	defaultMaxTimeBetweenFlushes = 30 * time.Second

	marketplaceDuplicateWaitTime = 1 * time.Second
	marketplaceMinRetries        = 20
)

const (
	UsageIncreaseMethodUnset = iota
	UsageIncreaseMethodTotal
	UsageIncreaseMethodBreakdown
)

const (
	BatchLimitHeader            = "x-cq-batch-limit"
	MinimumUpdateIntervalHeader = "x-cq-minimum-update-interval"
	MaximumUpdateIntervalHeader = "x-cq-maximum-update-interval"
	QueryIntervalHeader         = "x-cq-query-interval"
)

//go:generate mockgen -package=mocks -destination=../premium/mocks/marketplacemetering.go -source=usage.go AWSMarketplaceClientInterface
type AWSMarketplaceClientInterface interface {
	MeterUsage(ctx context.Context, params *marketplacemetering.MeterUsageInput, optFns ...func(*marketplacemetering.Options)) (*marketplacemetering.MeterUsageOutput, error)
}

type TokenClient interface {
	GetToken() (auth.Token, error)
	GetTokenType() auth.TokenType
}

type CheckQuotaResult struct {
	// HasQuota is true if the quota has not been exceeded
	HasQuota bool

	// SuggestedQueryInterval is the suggested interval to wait before querying the API again
	SuggestedQueryInterval time.Duration
}

type QuotaMonitor interface {
	// TeamName returns the team name
	TeamName() string
	// CheckQuota checks if the quota has been exceeded
	CheckQuota(context.Context) (CheckQuotaResult, error)
}

type UsageClient interface {
	QuotaMonitor
	// Increase updates the usage by the given number of rows
	Increase(uint32) error
	// IncreaseForTable updates the usage of a table by the given number of rows
	IncreaseForTable(string, uint32) error
	// Close flushes any remaining rows and closes the quota service
	Close() error
}

type UsageClientOptions func(updater *BatchUpdater)

// WithBatchLimit sets the maximum number of rows to update in a single request
func WithBatchLimit(batchLimit uint32) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.batchLimit = batchLimit
	}
}

// WithMaxTimeBetweenFlushes sets the flush duration - the time at which an update will be triggered even if the batch limit is not reached
func WithMaxTimeBetweenFlushes(maxTimeBetweenFlushes time.Duration) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.maxTimeBetweenFlushes = maxTimeBetweenFlushes
		updater.flushDuration.Reset(maxTimeBetweenFlushes)
	}
}

// WithMinTimeBetweenFlushes sets the minimum time between updates
func WithMinTimeBetweenFlushes(minTimeBetweenFlushes time.Duration) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.minTimeBetweenFlushes = minTimeBetweenFlushes
	}
}

// WithMaxRetries sets the maximum number of retries to update the usage in case of an API error
func WithMaxRetries(maxRetries int) UsageClientOptions {
	return func(updater *BatchUpdater) {
		if maxRetries > 0 {
			updater.maxRetries = maxRetries
		}
	}
}

// WithMaxWaitTime sets the maximum time to wait before retrying a failed update
func WithMaxWaitTime(maxWaitTime time.Duration) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.maxWaitTime = maxWaitTime
	}
}

// WithLogger sets the logger to use - defaults to a no-op logger
func WithLogger(logger zerolog.Logger) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.logger = logger
	}
}

// WithURL sets the API URL to use - defaults to https://api.cloudquery.io
func WithURL(url string) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.url = url
	}
}

// withTeamName sets the team name to use - defaults to the team name from the configuration
func withTeamName(teamName cqapi.TeamName) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.teamName = teamName
	}
}

// WithAPIClient sets the API client to use - defaults to a client using a bearer token generated from the refresh token stored in the configuration
func WithAPIClient(apiClient *cqapi.ClientWithResponses) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.apiClient = apiClient
	}
}

// WithAWSMarketplaceClient sets the AWS Marketplace client to use - defaults to marketplacemetering.NewFromConfig()
func WithAWSMarketplaceClient(awsMarketplaceClient AWSMarketplaceClientInterface) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.awsMarketplaceClient = awsMarketplaceClient
	}
}

// withTokenClient sets the token client to use - defaults to auth.NewTokenClient(). Used in tests to mock the token client
func withTokenClient(tokenClient TokenClient) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.tokenClient = tokenClient
	}
}

var (
	_ UsageClient = (*BatchUpdater)(nil)
	_ UsageClient = (*NoOpUsageClient)(nil)
)

type BatchUpdater struct {
	logger      zerolog.Logger
	url         string
	apiClient   *cqapi.ClientWithResponses
	tokenClient TokenClient

	awsMarketplaceClient AWSMarketplaceClientInterface

	// Plugin details
	teamName       cqapi.TeamName
	pluginMeta     plugin.Meta
	installationID string

	// Configuration
	batchLimit            uint32
	maxRetries            int
	maxWaitTime           time.Duration
	minTimeBetweenFlushes time.Duration
	maxTimeBetweenFlushes time.Duration

	// State
	sync.Mutex
	flushDuration       *time.Ticker
	rows                uint32
	tables              map[string]uint32
	lastUpdateTime      time.Time
	triggerUpdate       chan struct{}
	done                chan struct{}
	closeError          chan error
	isClosed            bool
	dataOnClose         bool
	usageIncreaseMethod int

	// Testing
	timeFunc func() time.Time
}

func NewUsageClient(meta plugin.Meta, ops ...UsageClientOptions) (UsageClient, error) {
	u := &BatchUpdater{
		logger: zerolog.Nop(),
		url:    defaultAPIURL,

		pluginMeta: meta,

		batchLimit:            defaultBatchLimit,
		minTimeBetweenFlushes: defaultMinTimeBetweenFlushes,
		maxTimeBetweenFlushes: defaultMaxTimeBetweenFlushes,
		maxRetries:            defaultMaxRetries,
		maxWaitTime:           defaultMaxWaitTime,
		flushDuration:         time.NewTicker(defaultMaxTimeBetweenFlushes),
		triggerUpdate:         make(chan struct{}),
		done:                  make(chan struct{}),
		closeError:            make(chan error),
		timeFunc:              time.Now,

		tables: map[string]uint32{},
	}
	for _, op := range ops {
		op(u)
	}

	if meta.SkipUsageClient {
		u.logger.Debug().Msg("Disabling usage client")
		return &NoOpUsageClient{
			TeamNameValue: u.teamName,
		}, nil
	}
	// If user wants to use the AWS Marketplace for billing, don't even try to communicate with CQ API
	if isAWSMarketplace() {
		err := u.setupAWSMarketplace()
		return u, err
	}

	if u.tokenClient == nil {
		u.tokenClient = auth.NewTokenClient()
	}

	// Fail early if the token is not set
	if _, err := u.tokenClient.GetToken(); err != nil {
		return nil, err
	}

	// Create a default api client if none was provided
	if u.apiClient == nil {
		retryClient := retryablehttp.NewClient()
		retryClient.Logger = nil
		retryClient.RetryMax = u.maxRetries
		retryClient.RetryWaitMax = u.maxWaitTime

		var err error
		u.apiClient, err = cqapi.NewClientWithResponses(u.url,
			cqapi.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
				token, err := u.tokenClient.GetToken()
				if err != nil {
					return fmt.Errorf("failed to get token: %w", err)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
				return nil
			}),
			cqapi.WithHTTPClient(retryClient.StandardClient()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create api client: %w", err)
		}
	}

	// Set team name from configuration if not provided
	if u.teamName == "" {
		teamName, err := u.getTeamNameByTokenType(u.tokenClient.GetTokenType())
		if err != nil {
			return nil, fmt.Errorf("failed to get team name: %w", err)
		}
		u.teamName = teamName
	}
	u.installationID = determineInstallationID()

	u.backgroundUpdater()

	return u, nil
}

func (u *BatchUpdater) setupAWSMarketplace() error {
	ctx := context.TODO()
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	// This allows us to be able to inject a mock client for testing
	if u.awsMarketplaceClient == nil {
		u.awsMarketplaceClient = marketplacemetering.NewFromConfig(cfg)
	}
	u.teamName = "AWS_MARKETPLACE"
	// This needs to be larger than normal, because we can only send a single usage record per second (from each compute node)
	u.batchLimit = 1000000000

	u.minTimeBetweenFlushes = 1 * time.Minute
	u.maxRetries = max(u.maxRetries, marketplaceMinRetries)
	u.backgroundUpdater()

	_, err = u.awsMarketplaceClient.MeterUsage(ctx, &marketplacemetering.MeterUsageInput{
		ProductCode:    aws.String(awsMarketplaceProductCode()),
		Timestamp:      aws.Time(u.timeFunc()),
		UsageDimension: aws.String("rows"),
		UsageQuantity:  aws.Int32(int32(0)),
		DryRun:         aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed dry run invocation with error: %w", err)
	}
	return nil
}

func isAWSMarketplace() bool {
	return os.Getenv("CQ_AWS_MARKETPLACE_CONTAINER") == "true"
}

func awsMarketplaceProductCode() string {
	if os.Getenv("CQ_AWS_MARKETPLACE_CONTAINER") == "true" {
		return "2a8bdkarwqrp0tmo4errl65s7"
	}
	return ""
}

func (u *BatchUpdater) Increase(rows uint32) error {
	if u.usageIncreaseMethod == UsageIncreaseMethodBreakdown {
		return errors.New("mixing usage increase methods is not allowed, use IncreaseForTable instead")
	}

	if rows <= 0 {
		return fmt.Errorf("rows must be greater than zero got %d", rows)
	}

	if u.isClosed {
		return errors.New("usage updater is closed")
	}

	u.Lock()
	defer u.Unlock()

	if u.usageIncreaseMethod == UsageIncreaseMethodUnset {
		u.usageIncreaseMethod = UsageIncreaseMethodTotal
	}
	u.rows += rows

	// Trigger an update unless an update is already in process
	select {
	case u.triggerUpdate <- struct{}{}:
	default:
	}

	return nil
}

func (u *BatchUpdater) IncreaseForTable(table string, rows uint32) error {
	if u.usageIncreaseMethod == UsageIncreaseMethodTotal {
		return errors.New("mixing usage increase methods is not allowed, use Increase instead")
	}

	if rows <= 0 {
		return fmt.Errorf("rows must be greater than zero got %d", rows)
	}

	if u.isClosed {
		return errors.New("usage updater is closed")
	}

	u.Lock()
	defer u.Unlock()

	if u.usageIncreaseMethod == UsageIncreaseMethodUnset {
		u.usageIncreaseMethod = UsageIncreaseMethodBreakdown
	}

	u.tables[table] += rows
	u.rows += rows

	// Trigger an update unless an update is already in process
	select {
	case u.triggerUpdate <- struct{}{}:
	default:
	}

	return nil
}

func (u *BatchUpdater) TeamName() string {
	return u.teamName
}

func (u *BatchUpdater) CheckQuota(ctx context.Context) (CheckQuotaResult, error) {
	if u.awsMarketplaceClient != nil {
		return CheckQuotaResult{HasQuota: true}, nil
	}
	u.logger.Debug().Str("url", u.url).Str("team", u.teamName).Str("pluginTeam", u.pluginMeta.Team).Str("pluginKind", string(u.pluginMeta.Kind)).Str("pluginName", u.pluginMeta.Name).Msg("checking quota")
	usage, err := u.apiClient.GetTeamPluginUsageWithResponse(ctx, u.teamName, u.pluginMeta.Team, u.pluginMeta.Kind, u.pluginMeta.Name)
	if err != nil {
		return CheckQuotaResult{HasQuota: false}, fmt.Errorf("failed to get usage: %w", err)
	}
	if usage.StatusCode() != http.StatusOK {
		if u.tokenClient.GetTokenType() == auth.APIKey && usage.StatusCode() == http.StatusForbidden {
			u.logger.Warn().Msg("API Key may have expired. Please see the CloudQuery Console see the expiration status.")
		}
		return CheckQuotaResult{HasQuota: false}, fmt.Errorf("failed to get usage: %s", usage.Status())
	}

	res := CheckQuotaResult{
		HasQuota: usage.JSON200.RemainingRows == nil || *usage.JSON200.RemainingRows > 0,
	}
	if usage.HTTPResponse == nil {
		return res, nil
	}
	if headerValue := usage.HTTPResponse.Header.Get(QueryIntervalHeader); headerValue != "" {
		interval, err := strconv.ParseUint(headerValue, 10, 32)
		if interval > 0 {
			res.SuggestedQueryInterval = time.Duration(interval) * time.Second
		} else {
			u.logger.Warn().Err(err).Str(QueryIntervalHeader, headerValue).Msg("failed to parse query interval")
		}
	}
	return res, nil
}

func (u *BatchUpdater) Close() error {
	u.isClosed = true

	close(u.done)

	return <-u.closeError
}

func (u *BatchUpdater) getTableUsage() (usage []cqapi.UsageIncreaseTablesInner, total uint32) {
	u.Lock()
	defer u.Unlock()

	for key, value := range u.tables {
		if value == 0 {
			continue
		}
		usage = append(usage, cqapi.UsageIncreaseTablesInner{
			Name: key,
			Rows: int(value),
		})
	}

	return usage, u.rows
}

func (u *BatchUpdater) subtractTableUsageForAWSMarketplace(total uint32) {
	for table := range u.tables {
		tableTotal := u.tables[table]
		if tableTotal < 1 {
			continue
		}
		if tableTotal >= total {
			u.tables[table] -= total
			// we can return early because we have subtracted enough rows
			return
		}
		u.tables[table] = 0
		total -= tableTotal
	}
}
func (u *BatchUpdater) subtractTableUsage(usage []cqapi.UsageIncreaseTablesInner, total uint32) {
	u.Lock()
	defer u.Unlock()
	if u.awsMarketplaceClient != nil {
		u.subtractTableUsageForAWSMarketplace(total)
		return
	}

	for _, table := range usage {
		u.tables[table.Name] -= uint32(table.Rows)
	}

	u.rows -= total
}

func (u *BatchUpdater) backgroundUpdater() {
	ctx := context.Background()
	started := make(chan struct{})

	go func() {
		started <- struct{}{}
		for {
			select {
			case <-u.triggerUpdate:
				// If we are using AWS Marketplace, we should only report the usage at the end of the sync
				if u.awsMarketplaceClient != nil {
					continue
				}
				if time.Since(u.lastUpdateTime) < u.minTimeBetweenFlushes {
					// Not enough time since last update
					continue
				}

				tables, totals := u.getTableUsage()

				if totals < u.batchLimit {
					// Not enough rows to update
					continue
				}

				if err := u.updateUsageWithRetryAndBackoff(ctx, totals, tables); err != nil {
					u.logger.Warn().Err(err).Msg("failed to update usage")
					continue
				}
				u.subtractTableUsage(tables, totals)

			case <-u.flushDuration.C:
				// If we are using AWS Marketplace, we should only report the usage at the end of the sync
				if u.awsMarketplaceClient != nil {
					continue
				}
				if time.Since(u.lastUpdateTime) < u.minTimeBetweenFlushes {
					// Not enough time since last update
					continue
				}

				tables, totals := u.getTableUsage()

				if totals == 0 {
					continue
				}

				if err := u.updateUsageWithRetryAndBackoff(ctx, totals, tables); err != nil {
					u.logger.Warn().Err(err).Msg("failed to update usage")
					continue
				}
				u.subtractTableUsage(tables, totals)

			case <-u.done:
				tables, totals := u.getTableUsage()
				if totals != 0 {
					u.dataOnClose = true
					// To allow us to round up the total in the last batch we need to save the original total
					// to use in the last subtractTableUsage
					originalTotals := totals

					// If we are using AWS Marketplace, we need to round up to the nearest 1000
					if u.awsMarketplaceClient != nil {
						totals = roundUp(totals, 1000)
					}
					if err := u.updateUsageWithRetryAndBackoff(ctx, totals, tables); err != nil {
						u.closeError <- err
						return
					}
					u.subtractTableUsage(tables, originalTotals)
				}
				u.closeError <- nil
				return
			}
		}
	}()
	<-started
}

func (u *BatchUpdater) reportUsageToAWSMarketplace(ctx context.Context, rows uint32) error {
	// AWS marketplace requires usage to be reported as groups of 1000
	rows /= 1000
	usage := []types.UsageAllocation{{
		AllocatedUsageQuantity: aws.Int32(int32(rows)),
		Tags: []types.Tag{
			{
				Key:   aws.String("plugin_name"),
				Value: aws.String(u.pluginMeta.Name),
			},
			{
				Key:   aws.String("plugin_team"),
				Value: aws.String(u.pluginMeta.Team),
			},
			{
				Key:   aws.String("plugin_kind"),
				Value: aws.String(string(u.pluginMeta.Kind)),
			},
		},
	}}
	// Timestamp + UsageDimension + UsageQuantity are required fields and must be unique
	// since Timestamp only maintains a granularity of seconds, we need to ensure our batch size is large enough
	_, err := u.awsMarketplaceClient.MeterUsage(ctx, &marketplacemetering.MeterUsageInput{
		// Product code is a unique identifier for a product in AWS Marketplace
		// Each product is given a unique product code when it is listed in AWS Marketplace
		// in the future we can have multiple product codes for container or AMI based listings
		ProductCode:      aws.String(awsMarketplaceProductCode()),
		Timestamp:        aws.Time(u.timeFunc()),
		UsageDimension:   aws.String("rows"),
		UsageAllocations: usage,
		UsageQuantity:    aws.Int32(int32(rows)),
	})
	if err != nil {
		return fmt.Errorf("failed to update usage: %w", err)
	}
	return nil
}

func (u *BatchUpdater) updateMarketplaceUsage(ctx context.Context, rows uint32) error {
	var lastErr error
	for retry := 0; retry < u.maxRetries; retry++ {
		u.logger.Debug().Int("try", retry).Int("max_retries", u.maxRetries).Uint32("rows", rows).Msg("updating usage")

		lastErr = u.reportUsageToAWSMarketplace(ctx, rows)
		if lastErr == nil {
			u.logger.Debug().Int("try", retry).Uint32("rows", rows).Msg("usage updated")
			return nil
		}

		var de *types.DuplicateRequestException
		if !errors.As(lastErr, &de) {
			return fmt.Errorf("failed to update usage: %w", lastErr)
		}
		u.logger.Debug().Err(lastErr).Int("try", retry).Uint32("rows", rows).Msg("usage update failed due to duplicate request")

		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(marketplaceDuplicateWaitTime + jitter)
	}
	return fmt.Errorf("failed to update usage: max retries exceeded: %w", lastErr)
}

func (u *BatchUpdater) updateUsageWithRetryAndBackoff(ctx context.Context, rows uint32, tables []cqapi.UsageIncreaseTablesInner) error {
	// If the AWS Marketplace client is set, use it to track usage
	if u.awsMarketplaceClient != nil {
		return u.updateMarketplaceUsage(ctx, rows)
	}

	u.logger.Debug().Str("url", u.url).Uint32("rows", rows).Msg("updating usage")
	payload := cqapi.IncreaseTeamPluginUsageJSONRequestBody{
		RequestId:  uuid.New(),
		PluginTeam: u.pluginMeta.Team,
		PluginKind: u.pluginMeta.Kind,
		PluginName: u.pluginMeta.Name,
		Rows:       int(rows),
	}
	if len(u.installationID) > 0 {
		payload.InstallationID = &u.installationID
	}
	if len(tables) > 0 {
		payload.Tables = &tables
	}

	resp, err := u.apiClient.IncreaseTeamPluginUsageWithResponse(ctx, u.teamName, payload)
	if err == nil && resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		u.logger.Debug().Str("url", u.url).Int("status_code", resp.StatusCode()).Uint32("rows", rows).Msg("usage updated")
		u.lastUpdateTime = u.timeFunc().UTC()
		if resp.HTTPResponse != nil {
			u.updateConfigurationFromHeaders(resp.HTTPResponse.Header)
		}
		return nil
	}

	return fmt.Errorf("failed to update usage: %w", err)
}

// updateConfigurationFromHeaders updates the configuration based on the headers returned by the API
func (u *BatchUpdater) updateConfigurationFromHeaders(header http.Header) {
	if headerValue := header.Get(BatchLimitHeader); headerValue != "" {
		newBatchLimit, err := strconv.ParseUint(headerValue, 10, 32)
		if newBatchLimit > 0 {
			u.batchLimit = uint32(newBatchLimit)
		} else {
			u.logger.Warn().Err(err).Str(BatchLimitHeader, headerValue).Msg("failed to parse batch limit")
		}
	}

	if headerValue := header.Get(MinimumUpdateIntervalHeader); headerValue != "" {
		newInterval, err := strconv.ParseInt(headerValue, 10, 32)
		if newInterval > 0 {
			u.minTimeBetweenFlushes = time.Duration(newInterval) * time.Second
		} else {
			u.logger.Warn().Err(err).Str(MinimumUpdateIntervalHeader, headerValue).Msg("failed to parse minimum update interval")
		}
	}

	if headerValue := header.Get(MaximumUpdateIntervalHeader); headerValue != "" {
		newInterval, err := strconv.ParseInt(headerValue, 10, 32)
		if newInterval > 0 {
			newMaxTimeBetweenFlushes := time.Duration(newInterval) * time.Second
			if u.maxTimeBetweenFlushes != newMaxTimeBetweenFlushes {
				u.maxTimeBetweenFlushes = newMaxTimeBetweenFlushes
				u.flushDuration.Reset(u.maxTimeBetweenFlushes)
			}
		} else {
			u.logger.Warn().Err(err).Str(MaximumUpdateIntervalHeader, headerValue).Msg("failed to parse maximum update interval")
		}
	}
}

func (u *BatchUpdater) getTeamNameByTokenType(tokenType auth.TokenType) (string, error) {
	switch tokenType {
	case auth.BearerToken:
		teamName, err := config.GetValue("team")
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("config file for reading team name not found (%w). Hint: use `cloudquery login` and/or `cloudquery switch <team>`", err)
		} else if err != nil {
			return "", fmt.Errorf("failed to get team name from config: %w", err)
		}
		if teamName == "" {
			return "", errors.New("team name not set. Hint: use `cloudquery switch <team>`")
		}
		return teamName, nil
	case auth.APIKey:
		resp, err := u.apiClient.ListTeamsWithResponse(context.Background(), &cqapi.ListTeamsParams{})
		if err != nil {
			return "", fmt.Errorf("failed to list teams for API key: %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			return "", fmt.Errorf("failed to list teams for API key, status code: %s", resp.Status())
		}
		if len(resp.JSON200.Items) != 1 {
			return "", fmt.Errorf("expected to find exactly one team for API key, found %d", len(resp.JSON200.Items))
		}
		return resp.JSON200.Items[0].Name, nil
	default:
		team := os.Getenv("_CQ_TEAM_NAME")
		if team == "" {
			switch tokenType {
			case auth.SyncRunAPIKey, auth.SyncTestConnectionAPIKey:
				return "", errors.New("_CQ_TEAM_NAME environment variable not set")
			}
			return "", fmt.Errorf("unsupported token type: %v", tokenType)
		}
		return team, nil
	}
}

func determineInstallationID() string {
	return os.Getenv("_CQ_INSTALLATION_ID")
}

type NoOpUsageClient struct {
	TeamNameValue string
}

func (n *NoOpUsageClient) TeamName() string {
	return n.TeamNameValue
}

func (NoOpUsageClient) CheckQuota(_ context.Context) (CheckQuotaResult, error) {
	return CheckQuotaResult{HasQuota: true}, nil
}

func (NoOpUsageClient) Increase(_ uint32) error {
	return nil
}

func (NoOpUsageClient) IncreaseForTable(_ string, _ uint32) error {
	return nil
}

func (NoOpUsageClient) Close() error {
	return nil
}

func roundDown(x, unit uint32) uint32 {
	return x - (x % unit)
}

func roundUp(x, unit uint32) uint32 {
	if x%unit == 0 {
		return x
	}
	return x + (unit - x%unit)
}
