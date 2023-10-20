package premium

import (
	"context"
	"fmt"
	cqapi "github.com/cloudquery/cloudquery-api-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBatchLimit  = 1000
	defaultDurationMS  = 10000
	defaultMaxRetries  = 5
	defaultMaxWaitTime = 60 * time.Second
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

// WithTickerDuration sets the duration between updates if the number of rows to update have not reached the batch limit
func WithTickerDuration(durationms int) UpdaterOptions {
	return func(updater *BatchUpdater) {
		updater.tickerDuration = durationms
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

	teamName   string
	pluginTeam string
	pluginKind string
	pluginName string

	batchLimit     uint32
	tickerDuration int
	maxRetries     int
	maxWaitTime    time.Duration
	rowsToUpdate   atomic.Uint32
	triggerUpdate  chan struct{}
	done           chan struct{}
	wg             *sync.WaitGroup
	isClosed       bool
}

func NewUsageClient(ctx context.Context, apiClient *cqapi.ClientWithResponses, teamName, pluginTeam, pluginKind, pluginName string, ops ...UpdaterOptions) *BatchUpdater {
	u := &BatchUpdater{
		apiClient: apiClient,

		teamName:   teamName,
		pluginTeam: pluginTeam,
		pluginKind: pluginKind,
		pluginName: pluginName,

		batchLimit:     defaultBatchLimit,
		tickerDuration: defaultDurationMS,
		maxRetries:     defaultMaxRetries,
		maxWaitTime:    defaultMaxWaitTime,
		triggerUpdate:  make(chan struct{}),
		done:           make(chan struct{}),
		wg:             &sync.WaitGroup{},
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
	u.wg.Wait()

	return nil
}

func (u *BatchUpdater) backgroundUpdater(ctx context.Context) {
	started := make(chan struct{})
	u.wg.Add(1)

	duration := time.Duration(u.tickerDuration) * time.Millisecond
	ticker := time.NewTicker(duration)

	go func() {
		defer u.wg.Done()
		started <- struct{}{}
		for {
			select {
			case <-u.triggerUpdate:
				rowsToUpdate := u.rowsToUpdate.Load()
				if rowsToUpdate < u.batchLimit {
					// Not enough rows to update
					continue
				}
				log.Info().Msgf("updating usage: %d", rowsToUpdate)
				if err := u.updateUsageWithRetryAndBackoff(ctx, rowsToUpdate); err != nil {
					log.Error().Err(err).Msg("failed to update usage")
					// TODO: what to do with an update error
					continue
				}
				u.rowsToUpdate.Add(-rowsToUpdate)
			case <-ticker.C:
				rowsToUpdate := u.rowsToUpdate.Load()
				if rowsToUpdate == 0 {
					continue
				}
				if err := u.updateUsageWithRetryAndBackoff(ctx, rowsToUpdate); err != nil {
					log.Error().Err(err).Msg("failed to update usage")
					// TODO: what to do with an update error
					continue
				}
				u.rowsToUpdate.Add(-rowsToUpdate)
			case <-u.done:
				remainingRows := u.rowsToUpdate.Load()
				if remainingRows != 0 {
					log.Info().Msgf("updating usage: %d", remainingRows)
					if err := u.updateUsageWithRetryAndBackoff(ctx, remainingRows); err != nil {
						log.Error().Err(err).Msg("failed to update usage")
					}
					u.rowsToUpdate.Add(-remainingRows)
				}
				log.Info().Msg("background updater exiting")
				return
			}
		}
	}()
	<-started
}

func (u *BatchUpdater) updateUsageWithRetryAndBackoff(ctx context.Context, numberToUpdate uint32) error {
	var retryDelay time.Duration

	for retry := 0; retry < u.maxRetries; retry++ {
		startTime := time.Now()

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
		if resp.StatusCode() == http.StatusOK {
			return nil
		}

		retryDelay = time.Duration(1<<retry) * time.Second
		if retryDelay > u.maxWaitTime {
			retryDelay = u.maxWaitTime
		}

		sleepDuration := retryDelay - time.Since(startTime)
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
	return fmt.Errorf("failed to update usage: max retries exceeded")
}
