package premium

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ErrNoQuota struct {
	team string
}

func (e ErrNoQuota) Error() string {
	return fmt.Sprintf("You have reached this plugin's usage limit for the month, please visit https://cloud.cloudquery.io/teams/%s/billing to upgrade your plan or increase the limit.", e.team)
}

const DefaultQuotaCheckInterval = 30 * time.Second
const DefaultMaxQuotaFailures = 10 // 5 minutes

type quotaChecker struct {
	qm                     QuotaMonitor
	duration               time.Duration
	maxConsecutiveFailures int
}

type QuotaCheckOption func(*quotaChecker)

// WithQuotaCheckPeriod controls the time interval between quota checks
func WithQuotaCheckPeriod(duration time.Duration) QuotaCheckOption {
	return func(m *quotaChecker) {
		m.duration = duration
	}
}

// WithQuotaMaxConsecutiveFailures controls the number of consecutive failed quota checks before the context is cancelled
func WithQuotaMaxConsecutiveFailures(n int) QuotaCheckOption {
	return func(m *quotaChecker) {
		m.maxConsecutiveFailures = n
	}
}

// WithCancelOnQuotaExceeded monitors the quota usage at intervals defined by duration and cancels the context if the quota is exceeded
func WithCancelOnQuotaExceeded(ctx context.Context, qm QuotaMonitor, ops ...QuotaCheckOption) (context.Context, error) {
	m := quotaChecker{
		qm:                     qm,
		duration:               DefaultQuotaCheckInterval,
		maxConsecutiveFailures: DefaultMaxQuotaFailures,
	}
	for _, op := range ops {
		op(&m)
	}

	if err := m.checkInitialQuota(ctx); err != nil {
		return ctx, err
	}

	newCtx := m.startQuotaMonitor(ctx)

	return newCtx, nil
}

func (qc quotaChecker) checkInitialQuota(ctx context.Context) error {
	hasQuota, err := qc.qm.HasQuota(ctx)
	if err != nil {
		return err
	}
	if hasQuota.SuggestedQueryInterval > 0 {
		qc.duration = hasQuota.SuggestedQueryInterval
	}

	if !hasQuota.HasQuota {
		return ErrNoQuota{team: qc.qm.TeamName()}
	}

	return nil
}

func (qc quotaChecker) startQuotaMonitor(ctx context.Context) context.Context {
	newCtx, cancelWithCause := context.WithCancelCause(ctx)
	go func() {
		ticker := time.NewTicker(qc.duration)
		consecutiveFailures := 0
		var hasQuotaErrors error
		for {
			select {
			case <-newCtx.Done():
				return
			case <-ticker.C:
				hasQuota, err := qc.qm.HasQuota(newCtx)
				if err != nil {
					consecutiveFailures++
					hasQuotaErrors = errors.Join(hasQuotaErrors, err)
					if consecutiveFailures >= qc.maxConsecutiveFailures {
						cancelWithCause(hasQuotaErrors)
						return
					}
					continue
				}
				if hasQuota.SuggestedQueryInterval > 0 && qc.duration != hasQuota.SuggestedQueryInterval {
					qc.duration = hasQuota.SuggestedQueryInterval
					ticker.Stop()
					ticker = time.NewTicker(qc.duration)
				}
				consecutiveFailures = 0
				hasQuotaErrors = nil
				if !hasQuota.HasQuota {
					cancelWithCause(ErrNoQuota{team: qc.qm.TeamName()})
					return
				}
			}
		}
	}()
	return newCtx
}
