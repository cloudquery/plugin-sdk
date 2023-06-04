package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/google/uuid"
)



type SyncOptions struct {
	Tables            []string
	SkipTables        []string
	Concurrency       int64
	Scheduler         Scheduler
	DeterministicCQID bool
	StateBackend      state.Client
}

// Tables returns all tables supported by this source plugin
func (p *Plugin) StaticTables() schema.Tables {
	return p.staticTables
}

func (p *Plugin) HasDynamicTables() bool {
	return p.getDynamicTables != nil
}

func (p *Plugin) DynamicTables() schema.Tables {
	return p.sessionTables
}

func (p *Plugin) readAll(ctx context.Context, table *schema.Table, sourceName string) ([]arrow.Record, error) {
	var readErr error
	ch := make(chan arrow.Record)
	go func() {
		defer close(ch)
		readErr = p.Read(ctx, table, sourceName, ch)
	}()
	// nolint:prealloc
	var resources []arrow.Record
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, readErr
}

func (p *Plugin) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error {
	return p.client.Read(ctx, table, sourceName, res)
}

func (p *Plugin) Acknowledge(ctx context.Context, recordUUID uuid.UUID) error {
	return nil
}

func (p *Plugin) syncAll(ctx context.Context, sourceName string, syncTime time.Time, options SyncOptions) ([]arrow.Record, error) {
	var err error
	ch := make(chan arrow.Record)
	go func() {
		defer close(ch)
		err = p.Sync(ctx, sourceName, syncTime, options, ch)
	}()
	// nolint:prealloc
	var resources []arrow.Record
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, err
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, sourceName string, syncTime time.Time, options SyncOptions, res chan<- arrow.Record) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	p.syncTime = syncTime
	startTime := time.Now()

	if p.unmanagedSync {
		if err := p.client.Sync(ctx, options, res); err != nil {
			return fmt.Errorf("failed to sync unmanaged client: %w", err)
		}
	} else {
		if err := p.managedSync(ctx, sourceName, syncTime, options, res); err != nil {
			return fmt.Errorf("failed to sync managed client: %w", err)
		}
	}

	p.logger.Info().Uint64("resources", p.metrics.TotalResources()).Uint64("errors", p.metrics.TotalErrors()).Uint64("panics", p.metrics.TotalPanics()).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}
