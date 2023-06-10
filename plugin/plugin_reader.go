package plugin

import (
	"context"
	"fmt"

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

type ReadOnlyClient interface {
	Tables(ctx context.Context) (schema.Tables, error)
	Sync(ctx context.Context, options SyncOptions, res chan<- Message) error
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

func (p *Plugin) syncAll(ctx context.Context, options SyncOptions) (Messages, error) {
	var err error
	ch := make(chan Message)
	go func() {
		defer close(ch)
		err = p.Sync(ctx, options, ch)
	}()
	// nolint:prealloc
	var resources []Message
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, err
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, options SyncOptions, res chan<- Message) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	// p.syncTime = options.SyncTime
	// startTime := time.Now()

	if err := p.client.Sync(ctx, options, res); err != nil {
		return fmt.Errorf("failed to sync unmanaged client: %w", err)
	}

	// p.logger.Info().Uint64("resources", p.metrics.TotalResources()).Uint64("errors", p.metrics.TotalErrors()).Uint64("panics", p.metrics.TotalPanics()).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}
