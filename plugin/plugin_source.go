package plugin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/glob"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type BackendOptions struct {
	TableName  string
	Connection string
	// SupportsCheckpoints is set when the CLI persists state from message.SyncCheckpoint,
	// so state clients send their writes on the sync stream instead of writing to the backend.
	SupportsCheckpoints bool
}

type Shard struct {
	Num   int32
	Total int32
}

type SyncOptions struct {
	Tables              []string
	SkipTables          []string
	SkipDependentTables bool
	DeterministicCQID   bool
	BackendOptions      *BackendOptions
	Shard               *Shard
	// Refresh, when non-empty, limits the sync to the listed resources: the plugin
	// re-fetches each one and emits an Insert for the resources it finds and a
	// DeleteRecord for the ones it does not. `Tables` is ignored and no
	// MigrateTable message is emitted.
	Refresh []RefreshResource
}

// Op is the kind of change that happened to a resource.
type Op int

const (
	// OpUpsert means the resource was created or updated and should be re-fetched.
	OpUpsert Op = iota
	// OpDelete means the resource no longer exists and should be deleted downstream.
	OpDelete
)

func (o Op) String() string {
	switch o {
	case OpUpsert:
		return "upsert"
	case OpDelete:
		return "delete"
	default:
		return fmt.Sprintf("Op(%d)", int(o))
	}
}

// Change is a single resource-level change observed by a WatchClient.
type Change struct {
	// Table is the name of the table the changed resource belongs to.
	Table string
	// Keys is a one-row record whose columns are a non-empty subset of the
	// table's identity columns: its primary key components when it declares any,
	// otherwise its primary keys. An event that knows only part of a resource's
	// identity, such as its ARN, supplies just that part. Build it with
	// schema.NewKeysRecord.
	Keys arrow.RecordBatch
	Op   Op
	// EventID identifies this change in the source's event stream. It is echoed
	// back by Ack to checkpoint progress.
	EventID string
	// Time is when the change happened at the source.
	Time time.Time
}

// WatchOptions configures a Watch stream.
type WatchOptions struct {
	Tables         []string
	SkipTables     []string
	BackendOptions *BackendOptions
	Shard          *Shard
	// Drain asks Watch to return once the sources are caught up instead of
	// following the stream indefinitely.
	Drain bool
}

// Position tells Watch where to start reading from.
type Position struct {
	// Since starts the stream at the given time. When nil, Watch resumes from the
	// stored checkpoint.
	Since *time.Time
}

// Ready is the result of WatchReady: the state of the stored checkpoint.
type Ready struct {
	HasCheckpoint  bool
	CheckpointTime time.Time
}

// RefreshResource identifies one resource to re-fetch during a refresh sync.
type RefreshResource struct {
	Table string
	// Keys is a one-row record whose columns are a non-empty subset of the
	// table's identity columns: its primary key components when it declares any,
	// otherwise its primary keys. An event that knows only part of a resource's
	// identity, such as its ARN, supplies just that part. Build it with
	// schema.NewKeysRecord.
	Keys arrow.RecordBatch
	Op   Op
}

// WatchClient is an optional interface a source client can implement to support
// event-based syncing. It is type-asserted on the client, like OnBeforeSender.
//
// Call order and concurrency, which implementations must honour:
//
//   - WatchReady is called first and on its own. Watch and Ack are not called
//     until it has returned.
//   - Watch and Ack are then called concurrently, from different goroutines, for
//     the lifetime of the stream. An Ack routinely arrives while Watch is blocked
//     sending a change or waiting for the next one, so Watch must not be the only
//     thing that can make progress.
//   - Acks are serialised with respect to each other, but not with respect to
//     Watch. Any state the two share - a checkpoint, a state client, a cursor -
//     must be guarded by the implementation. Neither the SDK nor the gRPC server
//     takes a lock around these calls, and Plugin.Sync can run a refresh on the
//     same client at the same time.
//   - Close may be called while a watch is still running, so it must be safe
//     against a concurrent Ack.
type WatchClient interface {
	// WatchReady prepares the watch and reports the state of the stored checkpoint.
	// It is called once, before Watch and Ack.
	WatchReady(ctx context.Context, opts WatchOptions) (Ready, error)
	// Watch streams changes to the given channel. It returns nil once the sources
	// are drained (when WatchOptions.Drain is set) or when ctx is done. It runs
	// concurrently with Ack.
	Watch(ctx context.Context, pos Position, changes chan<- Change) error
	// Ack checkpoints every change up to and including eventID. It is called from
	// a different goroutine than Watch, and may be called while Watch is blocked.
	Ack(ctx context.Context, eventID string) error
}

type SourceClient interface {
	Close(ctx context.Context) error
	Tables(ctx context.Context, options TableOptions) (schema.Tables, error)
	Sync(ctx context.Context, options SyncOptions, res chan<- message.SyncMessage) error
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
	newClientWrapper := func(ctx context.Context, logger zerolog.Logger, spec []byte, _ NewClientOptions) (Client, error) {
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

func (p *Plugin) SyncAll(ctx context.Context, options SyncOptions) (message.SyncMessages, error) {
	var err error
	ch := make(chan message.SyncMessage)
	go func() {
		defer close(ch)
		err = p.Sync(ctx, options, ch)
	}()
	// nolint:prealloc
	var resources message.SyncMessages
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, err
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, options SyncOptions, res chan<- message.SyncMessage) error {
	if !p.mu.TryLock() {
		return errors.New("plugin already in use")
	}
	defer p.mu.Unlock()
	if p.client == nil {
		return errors.New("plugin not initialized. call Init() first")
	}

	if err := p.client.Sync(ctx, options, res); err != nil {
		return fmt.Errorf("failed to sync unmanaged client: %w", err)
	}

	return nil
}

// watchClient returns the client as a WatchClient, or ErrNotImplemented if the
// client does not support watching.
func (p *Plugin) watchClient() (WatchClient, error) {
	if p.client == nil {
		return nil, errors.New("plugin not initialized. call Init() first")
	}
	c, ok := p.client.(WatchClient)
	if !ok {
		return nil, ErrNotImplemented
	}
	return c, nil
}

// WatchReady prepares a watch stream and reports the state of the stored checkpoint.
// It returns ErrNotImplemented if the client does not implement WatchClient.
//
// Unlike Sync, the watch methods do not take the plugin lock: a watch stream is
// long-lived and must not block refresh syncs running alongside it.
func (p *Plugin) WatchReady(ctx context.Context, opts WatchOptions) (Ready, error) {
	c, err := p.watchClient()
	if err != nil {
		return Ready{}, err
	}
	ready, err := c.WatchReady(ctx, opts)
	if err != nil {
		return Ready{}, fmt.Errorf("failed to prepare watch: %w", err)
	}
	return ready, nil
}

// Watch streams changes from the client to the given channel until it is drained
// or ctx is done. It returns ErrNotImplemented if the client does not implement
// WatchClient.
func (p *Plugin) Watch(ctx context.Context, pos Position, changes chan<- Change) error {
	c, err := p.watchClient()
	if err != nil {
		return err
	}
	if err := c.Watch(ctx, pos, changes); err != nil {
		return fmt.Errorf("failed to watch: %w", err)
	}
	return nil
}

// Ack checkpoints every change up to and including eventID. It returns
// ErrNotImplemented if the client does not implement WatchClient.
func (p *Plugin) Ack(ctx context.Context, eventID string) error {
	c, err := p.watchClient()
	if err != nil {
		return err
	}
	if err := c.Ack(ctx, eventID); err != nil {
		return fmt.Errorf("failed to ack event %q: %w", eventID, err)
	}
	return nil
}
