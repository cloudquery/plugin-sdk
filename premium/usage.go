package premium

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/cloudquery/cloudquery-api-go/auth"
	"github.com/cloudquery/cloudquery-api-go/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultAPIURL                = "https://api.cloudquery.io"
	defaultBatchLimit            = 1000
	defaultMaxRetries            = 5
	defaultMaxWaitTime           = 60 * time.Second
	defaultMinTimeBetweenFlushes = 10 * time.Second
	defaultMaxTimeBetweenFlushes = 30 * time.Second
)

type TokenClient interface {
	GetToken() (auth.Token, error)
	GetTokenType() auth.TokenType
}

type QuotaMonitor interface {
	// TeamName returns the team name
	TeamName() string
	// HasQuota returns true if the quota has not been exceeded
	HasQuota(context.Context) (bool, error)
}

type UsageClient interface {
	QuotaMonitor
	// Increase updates the usage by the given number of rows
	Increase(uint32) error
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
		updater.maxRetries = maxRetries
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

// withTokenClient sets the token client to use - defaults to auth.NewTokenClient(). Used in tests to mock the token client
func withTokenClient(tokenClient TokenClient) UsageClientOptions {
	return func(updater *BatchUpdater) {
		updater.tokenClient = tokenClient
	}
}

func WithPluginTeam(pluginTeam string) cqapi.PluginTeam {
	return pluginTeam
}

func WithPluginKind(pluginKind string) cqapi.PluginKind {
	return cqapi.PluginKind(pluginKind)
}

func WithPluginName(pluginName string) cqapi.PluginName {
	return pluginName
}

var _ UsageClient = (*BatchUpdater)(nil)

type BatchUpdater struct {
	logger      zerolog.Logger
	url         string
	apiClient   *cqapi.ClientWithResponses
	tokenClient TokenClient

	// Plugin details
	teamName   cqapi.TeamName
	pluginTeam cqapi.PluginTeam
	pluginKind cqapi.PluginKind
	pluginName cqapi.PluginName

	// Configuration
	batchLimit            uint32
	maxRetries            int
	maxWaitTime           time.Duration
	minTimeBetweenFlushes time.Duration
	maxTimeBetweenFlushes time.Duration

	// State
	lastUpdateTime time.Time
	rowsToUpdate   atomic.Uint32
	triggerUpdate  chan struct{}
	done           chan struct{}
	closeError     chan error
	isClosed       bool
}

func NewUsageClient(pluginTeam cqapi.PluginTeam, pluginKind cqapi.PluginKind, pluginName cqapi.PluginName, ops ...UsageClientOptions) (*BatchUpdater, error) {
	u := &BatchUpdater{
		logger: zerolog.Nop(),
		url:    defaultAPIURL,

		pluginTeam: pluginTeam,
		pluginKind: pluginKind,
		pluginName: pluginName,

		batchLimit:            defaultBatchLimit,
		minTimeBetweenFlushes: defaultMinTimeBetweenFlushes,
		maxTimeBetweenFlushes: defaultMaxTimeBetweenFlushes,
		maxRetries:            defaultMaxRetries,
		maxWaitTime:           defaultMaxWaitTime,
		triggerUpdate:         make(chan struct{}),
		done:                  make(chan struct{}),
		closeError:            make(chan error),
	}
	for _, op := range ops {
		op(u)
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
		ac, err := cqapi.NewClientWithResponses(u.url, cqapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			token, err := u.tokenClient.GetToken()
			if err != nil {
				return fmt.Errorf("failed to get token: %w", err)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return nil
		}))
		if err != nil {
			return nil, fmt.Errorf("failed to create api client: %w", err)
		}
		u.apiClient = ac
	}

	// Set team name from configuration if not provided
	if u.teamName == "" {
		teamName, err := u.getTeamNameByTokenType(u.tokenClient.GetTokenType())
		if err != nil {
			return nil, fmt.Errorf("failed to get team name: %w", err)
		}
		u.teamName = teamName
	}

	u.backgroundUpdater()

	return u, nil
}

func (u *BatchUpdater) Increase(rows uint32) error {
	if rows <= 0 {
		return fmt.Errorf("rows must be greater than zero got %d", rows)
	}

	if u.isClosed {
		return fmt.Errorf("usage updater is closed")
	}

	u.rowsToUpdate.Add(rows)

	// Trigger an update unless an update is already in process
	select {
	case u.triggerUpdate <- struct{}{}:
	default:
		return nil
	}

	return nil
}

func (u *BatchUpdater) TeamName() string {
	return u.teamName
}

func (u *BatchUpdater) HasQuota(ctx context.Context) (bool, error) {
	u.logger.Debug().Str("url", u.url).Str("team", u.teamName).Str("pluginTeam", u.pluginTeam).Str("pluginKind", string(u.pluginKind)).Str("pluginName", u.pluginName).Msg("checking quota")
	usage, err := u.apiClient.GetTeamPluginUsageWithResponse(ctx, u.teamName, u.pluginTeam, u.pluginKind, u.pluginName)
	if err != nil {
		return false, fmt.Errorf("failed to get usage: %w", err)
	}
	if usage.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("failed to get usage: %s", usage.Status())
	}

	hasQuota := usage.JSON200.RemainingRows == nil || *usage.JSON200.RemainingRows > 0
	return hasQuota, nil
}

func (u *BatchUpdater) Close() error {
	u.isClosed = true

	close(u.done)

	return <-u.closeError
}

func (u *BatchUpdater) backgroundUpdater() {
	ctx := context.Background()
	started := make(chan struct{})

	flushDuration := time.NewTicker(u.maxTimeBetweenFlushes)

	go func() {
		started <- struct{}{}
		for {
			select {
			case <-u.triggerUpdate:
				if time.Since(u.lastUpdateTime) < u.minTimeBetweenFlushes {
					// Not enough time since last update
					continue
				}

				rowsToUpdate := u.rowsToUpdate.Load()
				if rowsToUpdate < u.batchLimit {
					// Not enough rows to update
					continue
				}
				if err := u.updateUsageWithRetryAndBackoff(ctx, rowsToUpdate); err != nil {
					log.Warn().Err(err).Msg("failed to update usage")
					continue
				}
				u.rowsToUpdate.Add(-rowsToUpdate)
			case <-flushDuration.C:
				if time.Since(u.lastUpdateTime) < u.minTimeBetweenFlushes {
					// Not enough time since last update
					continue
				}
				rowsToUpdate := u.rowsToUpdate.Load()
				if rowsToUpdate == 0 {
					continue
				}
				if err := u.updateUsageWithRetryAndBackoff(ctx, rowsToUpdate); err != nil {
					log.Warn().Err(err).Msg("failed to update usage")
					continue
				}
				u.rowsToUpdate.Add(-rowsToUpdate)
			case <-u.done:
				remainingRows := u.rowsToUpdate.Load()
				if remainingRows != 0 {
					if err := u.updateUsageWithRetryAndBackoff(ctx, remainingRows); err != nil {
						u.closeError <- err
						return
					}
					u.rowsToUpdate.Add(-remainingRows)
				}
				u.closeError <- nil
				return
			}
		}
	}()
	<-started
}

func (u *BatchUpdater) updateUsageWithRetryAndBackoff(ctx context.Context, numberToUpdate uint32) error {
	for retry := 0; retry < u.maxRetries; retry++ {
		u.logger.Debug().Str("url", u.url).Int("try", retry).Int("max_retries", u.maxRetries).Uint32("rows", numberToUpdate).Msg("updating usage")
		queryStartTime := time.Now()

		resp, err := u.apiClient.IncreaseTeamPluginUsageWithResponse(ctx, u.teamName, cqapi.IncreaseTeamPluginUsageJSONRequestBody{
			RequestId:  uuid.New(),
			PluginTeam: u.pluginTeam,
			PluginKind: u.pluginKind,
			PluginName: u.pluginName,
			Rows:       int(numberToUpdate),
		})
		if err != nil {
			return fmt.Errorf("failed to update usage: %w", err)
		}
		if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
			u.logger.Debug().Str("url", u.url).Int("try", retry).Int("status_code", resp.StatusCode()).Uint32("rows", numberToUpdate).Msg("usage updated")
			u.lastUpdateTime = time.Now().UTC()
			return nil
		}

		retryDuration, err := u.calculateRetryDuration(resp.StatusCode(), resp.HTTPResponse.Header, queryStartTime, retry)
		if err != nil {
			return fmt.Errorf("failed to calculate retry duration: %w", err)
		}
		if retryDuration > 0 {
			time.Sleep(retryDuration)
		}
	}
	return fmt.Errorf("failed to update usage: max retries exceeded")
}

// calculateRetryDuration calculates the duration to sleep relative to the query start time before retrying an update
func (u *BatchUpdater) calculateRetryDuration(statusCode int, headers http.Header, queryStartTime time.Time, retry int) (time.Duration, error) {
	if !retryableStatusCode(statusCode) {
		return 0, fmt.Errorf("non-retryable status code: %d", statusCode)
	}

	// Check if we have a retry-after header
	retryAfter := headers.Get("Retry-After")
	if retryAfter != "" {
		retryDelay, err := time.ParseDuration(retryAfter + "s")
		if err != nil {
			return 0, fmt.Errorf("failed to parse retry-after header: %w", err)
		}
		return retryDelay, nil
	}

	// Calculate exponential backoff
	baseRetry := min(time.Duration(1<<retry)*time.Second, u.maxWaitTime)
	jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
	retryDelay := baseRetry + jitter
	return retryDelay - time.Since(queryStartTime), nil
}

func retryableStatusCode(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable
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
			return "", fmt.Errorf("team name not set. Hint: use `cloudquery switch <team>`")
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
		return "", fmt.Errorf("unsupported token type: %v", tokenType)
	}
}
