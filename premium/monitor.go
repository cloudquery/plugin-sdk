package premium

import (
	"context"
	"errors"
	"time"
)

var ErrNoQuota = errors.New("no remaining quota for the month, please increase your usage limit if you want to continue syncing this plugin")

const DefaultQuotaCheckInterval = 30 * time.Second

type quotaChecker struct {
	qm       QuotaMonitor
	duration time.Duration
}

type QuotaCheckOption func(*quotaChecker)

// WithQuotaCheckPeriod the time interval between quota checks
func WithQuotaCheckPeriod(duration time.Duration) QuotaCheckOption {
	return func(m *quotaChecker) {
		m.duration = duration
	}
}

// WithCancelOnQuotaExceeded monitors the quota usage at intervals defined by duration and cancels the context if the quota is exceeded
func WithCancelOnQuotaExceeded(ctx context.Context, qm QuotaMonitor, ops ...QuotaCheckOption) (context.Context, func(), error) {
	m := quotaChecker{
		qm:       qm,
		duration: DefaultQuotaCheckInterval,
	}
	for _, op := range ops {
		op(&m)
	}

	if err := m.checkInitialQuota(ctx); err != nil {
		return ctx, nil, err
	}

	ctx, cancel := m.startQuotaMonitor(ctx)

	return ctx, cancel, nil
}

func (qc quotaChecker) checkInitialQuota(ctx context.Context) error {
	hasQuota, err := qc.qm.HasQuota(ctx)
	if err != nil {
		return err
	}

	if !hasQuota {
		return ErrNoQuota
	}

	return nil
}

func (qc quotaChecker) startQuotaMonitor(ctx context.Context) (context.Context, func()) {
	newCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		ticker := time.NewTicker(qc.duration)
		for {
			select {
			case <-newCtx.Done():
				return
			case <-ticker.C:
				hasQuota, err := qc.qm.HasQuota(newCtx)
				if err != nil {
					continue
				}
				if !hasQuota {
					return
				}
			}
		}
	}()
	return newCtx, cancel
}
