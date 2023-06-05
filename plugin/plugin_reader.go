package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/rs/zerolog"
)

type SyncOptions struct {
	Tables            []string
	SkipTables        []string
	Concurrency       int64
	Scheduler         Scheduler
	DeterministicCQID bool
	// SyncTime if specified then this will be add to every table as _sync_time column
	SyncTime time.Time
	// If spceified then this will be added to every table as _source_name column
	SourceName   string
	StateBackend state.Client
}

type ReadOnlyClient interface {
	Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error
	Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error
	Close(ctx context.Context) error
}

type NewReadOnlyClientFunc func(context.Context, zerolog.Logger, any) (ReadOnlyClient, error)

// NewReadOnlyPlugin returns a new CloudQuery Plugin with the given name, version and implementation.
// this plugin will only support read operations. For ReadWrite plugin use NewPlugin.
func NewReadOnlyPlugin(name string, version string, newClient NewReadOnlyClientFunc, options ...Option) *Plugin {
	newClientWrapper := func(ctx context.Context, logger zerolog.Logger, any any) (Client, error) {
		readOnlyClient, err := newClient(ctx, logger, any)
		if err != nil {
			return nil, err
		}
		wrapperClient := struct {
			ReadOnlyClient
			UnimplementedWriter
		}{
			ReadOnlyClient: readOnlyClient,
		}
		return wrapperClient, nil
	}
	return NewPlugin(name, version, newClientWrapper, options...)
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

func (p *Plugin) syncAll(ctx context.Context, options SyncOptions) ([]arrow.Record, error) {
	var err error
	ch := make(chan arrow.Record)
	go func() {
		defer close(ch)
		err = p.Sync(ctx, options, ch)
	}()
	// nolint:prealloc
	var resources []arrow.Record
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, err
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	p.syncTime = options.SyncTime
	startTime := time.Now()

	if err := p.client.Sync(ctx, options, res); err != nil {
		return fmt.Errorf("failed to sync unmanaged client: %w", err)
	}

	p.logger.Info().Uint64("resources", p.metrics.TotalResources()).Uint64("errors", p.metrics.TotalErrors()).Uint64("panics", p.metrics.TotalPanics()).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}
