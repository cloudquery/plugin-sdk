package premium

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type ErrNoQuota struct {
	team string
}

func (e ErrNoQuota) Error() string {
	return fmt.Sprintf("You have reached this plugin's usage limit for the month, please visit https://cloudquery.io/teams/%s/billing to upgrade your plan or increase the limit.", e.team)
}

const DefaultQuotaCheckInterval = 30 * time.Second
const DefaultMaxQuotaFailures = 10 // 5 minutes

type quotaChecker struct {
	qm                     QuotaMonitor
	duration               time.Duration
	maxConsecutiveFailures int
	remainingRows          *atomic.Int64
	triggerQuotaCheck      chan struct{}
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
		triggerQuotaCheck:      make(chan struct{}),
	}
	for _, op := range ops {
		op(&m)
	}

	if err := m.checkInitialQuota(ctx); err != nil {
		return ctx, err
	}

	if batchUpdater, ok := qm.(*BatchUpdater); ok {
		batchUpdater.quotaChecker = &m
	}

	newCtx := m.startQuotaMonitor(ctx)

	return newCtx, nil
}

func (qc *quotaChecker) checkInitialQuota(ctx context.Context) error {
	if err := qc.refreshRemainingRows(ctx); err != nil {
		return err
	}

	if qc.remainingRows != nil && qc.remainingRows.Load() <= 0 {
		return ErrNoQuota{team: qc.qm.TeamName()}
	}

	return nil
}

func (qc *quotaChecker) startQuotaMonitor(ctx context.Context) context.Context {
	newCtx, cancelWithCause := context.WithCancelCause(ctx)
	go func() {
		ticker := time.NewTicker(qc.duration)
		consecutiveFailures := 0
		var hasQuotaErrors error
		for {
			select {
			case <-newCtx.Done():
				return
			case <-qc.triggerQuotaCheck:
				// Attempt to refresh the remaining rows immediately when triggered
				if err := qc.refreshRemainingRows(newCtx); err != nil {
					// Assume we have no quota if we can't refresh the remaining rows as this case is only triggered when
					// think we have exhausted the quota we knew about
					cancelWithCause(ErrNoQuota{team: qc.qm.TeamName()})
				}
				// Check if we have quota after refreshing the remaining rows - this covers the case where more quota has
				// been added after the initial check
				if qc.remainingRows != nil && qc.remainingRows.Load() <= 0 {
					cancelWithCause(ErrNoQuota{team: qc.qm.TeamName()})
					return
				}
			case <-ticker.C:
				err := qc.refreshRemainingRows(newCtx)
				if err != nil {
					consecutiveFailures++
					hasQuotaErrors = errors.Join(hasQuotaErrors, err)
					if consecutiveFailures >= qc.maxConsecutiveFailures {
						cancelWithCause(hasQuotaErrors)
						return
					}
					continue
				}
				consecutiveFailures = 0
				hasQuotaErrors = nil
				if qc.remainingRows != nil && qc.remainingRows.Load() <= 0 {
					cancelWithCause(ErrNoQuota{team: qc.qm.TeamName()})
					return
				}
			}
		}
	}()
	return newCtx
}

func (qc *quotaChecker) refreshRemainingRows(ctx context.Context) error {
	remainingRows, err := qc.qm.RemainingRows(ctx)
	if err != nil {
		return err
	}

	if remainingRows != nil {
		if qc.remainingRows == nil {
			qc.remainingRows = new(atomic.Int64)
		}
		qc.remainingRows.Store(*remainingRows)
	}

	return nil
}
