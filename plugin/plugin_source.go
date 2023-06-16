package plugin

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/glob"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/rs/zerolog"
)

type SyncOptions struct {
	Tables            []string
	SkipTables        []string
	Concurrency       int64
	DeterministicCQID bool
	StateBackend      state.Client
}

type SourceClient interface {
	GetSpec() any
	Close(ctx context.Context) error
	Tables(ctx context.Context) (schema.Tables, error)
	Sync(ctx context.Context, options SyncOptions, res chan<- message.Message) error
}

func MatchesTable(name string, includeTablesPattern []string, skipTablesPattern []string) bool {
	for _, pattern := range skipTablesPattern {
		if glob.Glob(pattern, name) {
			return false
		}
	}
	for _, pattern := range includeTablesPattern {
		if glob.Glob(pattern, name) {
			return true
		}
	}
	return false
}

type NewSourceClientFunc func(context.Context, zerolog.Logger, any) (SourceClient, error)

// NewSourcePlugin returns a new CloudQuery Plugin with the given name, version and implementation.
// Source plugins only support read operations. For Read & Write plugin use NewPlugin.
func NewSourcePlugin(name string, version string, newClient NewSourceClientFunc, options ...Option) *Plugin {
	newClientWrapper := func(ctx context.Context, logger zerolog.Logger, spec []byte) (Client, error) {
		sourceClient, err := newClient(ctx, logger, spec)
		if err != nil {
			return nil, err
		}
		wrapperClient := struct {
			SourceClient
			UnimplementedDestination
		}{
			SourceClient: sourceClient,
		}
		return wrapperClient, nil
	}
	return NewPlugin(name, version, newClientWrapper, options...)
}

func (p *Plugin) readAll(ctx context.Context, table *schema.Table) ([]arrow.Record, error) {
	var err error
	ch := make(chan arrow.Record)
	go func() {
		defer close(ch)
		err = p.client.Read(ctx, table, ch)
	}()
	// nolint:prealloc
	var records []arrow.Record
	for record := range ch {
		records = append(records, record)
	}
	return records, err
}

func (p *Plugin) SyncAll(ctx context.Context, options SyncOptions) (message.Messages, error) {
	var err error
	ch := make(chan message.Message)
	go func() {
		defer close(ch)
		err = p.Sync(ctx, options, ch)
	}()
	// nolint:prealloc
	var resources []message.Message
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, err
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, options SyncOptions, res chan<- message.Message) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	// startTime := time.Now()

	if err := p.client.Sync(ctx, options, res); err != nil {
		return fmt.Errorf("failed to sync unmanaged client: %w", err)
	}

	// p.logger.Info().Uint64("resources", p.metrics.TotalResources()).Uint64("errors", p.metrics.TotalErrors()).Uint64("panics", p.metrics.TotalPanics()).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}
