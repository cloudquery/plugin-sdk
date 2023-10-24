package premium

import (
	"context"
	"fmt"
	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	defaultBatchLimit            = 1000
	defaultMaxRetries            = 5
	defaultMaxWaitTime           = 60 * time.Second
	defaultMinTimeBetweenFlushes = 10 * time.Second
	defaultMaxTimeBetweenFlushes = 30 * time.Second
)

type UsageClient interface {
	// Increase updates the usage by the given number of rows
	Increase(context.Context, uint32)
	// HasQuota returns true if the quota has not been exceeded
	HasQuota(context.Context) (bool, error)
	// Close flushes any remaining rows and closes the quota service
	Close() error
}

type UpdaterOptions func(updater *BatchUpdater)

// WithBatchLimit sets the maximum number of rows to update in a single request
func WithBatchLimit(batchLimit uint32) UpdaterOptions {
	return func(updater *BatchUpdater) {
		updater.batchLimit = batchLimit
	}
}

// WithMaxTimeBetweenFlushes sets the flush duration - the time at which an update will be triggered even if the batch limit is not reached
func WithMaxTimeBetweenFlushes(maxTimeBetweenFlushes time.Duration) UpdaterOptions {
	return func(updater *BatchUpdater) {
		updater.maxTimeBetweenFlushes = maxTimeBetweenFlushes
	}
}

// WithMinTimeBetweenFlushes sets the minimum time between updates
func WithMinTimeBetweenFlushes(minTimeBetweenFlushes time.Duration) UpdaterOptions {
	return func(updater *BatchUpdater) {
		updater.minTimeBetweenFlushes = minTimeBetweenFlushes
	}
}

// WithMaxRetries sets the maximum number of retries to update the usage in case of an API error
func WithMaxRetries(maxRetries int) UpdaterOptions {
	return func(updater *BatchUpdater) {
		updater.maxRetries = maxRetries
	}
}

// WithMaxWaitTime sets the maximum time to wait before retrying a failed update
func WithMaxWaitTime(maxWaitTime time.Duration) UpdaterOptions {
	return func(updater *BatchUpdater) {
		updater.maxWaitTime = maxWaitTime
	}
}

type BatchUpdater struct {
	apiClient *cqapi.ClientWithResponses

	// Plugin details
	teamName   string
	pluginTeam string
	pluginKind string
	pluginName string

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

func NewUsageClient(ctx context.Context, apiClient *cqapi.ClientWithResponses, teamName, pluginTeam, pluginKind, pluginName string, ops ...UpdaterOptions) *BatchUpdater {
	u := &BatchUpdater{
		apiClient: apiClient,

		teamName:   teamName,
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

	u.backgroundUpdater(ctx)

	return u
}

func (u *BatchUpdater) Increase(_ context.Context, rows uint32) error {
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

func (u *BatchUpdater) HasQuota(ctx context.Context) (bool, error) {
	usage, err := u.apiClient.GetTeamPluginUsageWithResponse(ctx, u.teamName, u.pluginTeam, cqapi.PluginKind(u.pluginKind), u.pluginName)
	if err != nil {
		return false, fmt.Errorf("failed to get usage: %w", err)
	}
	if usage.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("failed to get usage: %s", usage.Status())
	}
	return *usage.JSON200.RemainingRows > 0, nil
}

func (u *BatchUpdater) Close(_ context.Context) error {
	u.isClosed = true

	close(u.done)

	return <-u.closeError
}

func (u *BatchUpdater) backgroundUpdater(ctx context.Context) {
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
		queryStartTime := time.Now()

		resp, err := u.apiClient.IncreaseTeamPluginUsageWithResponse(ctx, u.teamName, cqapi.IncreaseTeamPluginUsageJSONRequestBody{
			RequestId:  uuid.New(),
			PluginTeam: u.pluginTeam,
			PluginKind: cqapi.PluginKind(u.pluginKind),
			PluginName: u.pluginName,
			Rows:       int(numberToUpdate),
		})
		if err != nil {
			return fmt.Errorf("failed to update usage: %w", err)
		}
		if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
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
	if retryableStatusCode(statusCode) {
		retryAfter := headers.Get("Retry-After")
		if retryAfter != "" {
			retryDelay, err := time.ParseDuration(retryAfter + "s")
			if err != nil {
				return 0, fmt.Errorf("failed to parse retry-after header: %w", err)
			}
			return retryDelay, nil
		}
	}

	baseRetry := min(time.Duration(1<<retry)*time.Second, u.maxWaitTime)
	jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
	retryDelay := baseRetry + jitter
	return retryDelay - time.Since(queryStartTime), nil
}

func retryableStatusCode(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable
}
