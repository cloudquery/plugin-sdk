package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/examples/simple_plugin/client"
	"github.com/cloudquery/plugin-sdk/examples/simple_plugin/services"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/rs/zerolog"
)

const (
	// checkpointKey holds the id of the last acked event.
	checkpointKey = "events:checkpoint"
	// checkpointTimeKey holds when that event was acked, as RFC3339.
	checkpointTimeKey = "events:checkpoint_time"
)

var (
	_ plugin.WatchClient  = (*Client)(nil)
	_ scheduler.Refresher = (*Client)(nil)
)

type Client struct {
	logger    zerolog.Logger
	config    client.Spec
	tables    schema.Tables
	scheduler *scheduler.Scheduler

	// watchMu guards stateClient and checkpoint. Watch and Ack run concurrently
	// for the lifetime of a watch stream, and Close can land on top of either.
	watchMu     sync.Mutex
	stateClient state.Client
	checkpoint  string
}

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	if len(options.Refresh) > 0 {
		return c.scheduler.Refresh(ctx, c.tables, c, options.Refresh, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
	}

	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	return c.scheduler.Sync(ctx, c.schedulerClient(), tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (c *Client) schedulerClient() *client.TestClient {
	return &client.TestClient{
		Logger: c.logger,
		Spec:   c.config,
	}
}

// Refresh answers a refresh sync from the `refresh_data` rows in the spec. A row
// matches when every primary key column of the table holds the requested value.
func (c *Client) Refresh(_ context.Context, table *schema.Table, keys arrow.RecordBatch) (schema.ClientMeta, []any, error) {
	want := schema.KeysString(keys)
	for _, row := range c.config.RefreshData[table.Name] {
		rowKeys, err := rowKeysRecord(table, row)
		if err != nil {
			return nil, nil, fmt.Errorf("refresh_data row of table %q: %w", table.Name, err)
		}
		if schema.KeysString(rowKeys) != want {
			continue
		}
		return c.schedulerClient().WithClientID(rowClientID(row)), []any{row}, nil
	}
	return nil, nil, scheduler.ErrNotFound
}

// rowKeysRecord builds the keys record of a `refresh_data` row so it can be
// compared with the keys of the resource being refreshed.
func rowKeysRecord(table *schema.Table, row map[string]any) (arrow.RecordBatch, error) {
	values := make(map[string]any, len(table.PrimaryKeys()))
	for _, name := range table.PrimaryKeys() {
		value, ok := row[name]
		if !ok {
			return nil, fmt.Errorf("missing primary key %q", name)
		}
		values[name] = value
	}
	return schema.NewKeysRecord(table, values)
}

// rowClientID reports which multiplexed client a `refresh_data` row belongs to.
func rowClientID(row map[string]any) int {
	switch id := row["client_id"].(type) {
	case float64:
		return int(id)
	case int:
		return id
	default:
		return 0
	}
}

// WatchReady connects to the state backend and reports the stored checkpoint.
func (c *Client) WatchReady(ctx context.Context, opts plugin.WatchOptions) (plugin.Ready, error) {
	stateClient, err := state.NewConnectedClient(ctx, opts.BackendOptions)
	if err != nil {
		return plugin.Ready{}, fmt.Errorf("failed to connect to the state backend: %w", err)
	}

	eventID, err := stateClient.GetKey(ctx, checkpointKey)
	if err != nil {
		return plugin.Ready{}, fmt.Errorf("failed to read %s: %w", checkpointKey, err)
	}
	c.setWatchState(stateClient, eventID)
	if eventID == "" {
		return plugin.Ready{}, nil
	}

	ready := plugin.Ready{HasCheckpoint: true}
	at, err := stateClient.GetKey(ctx, checkpointTimeKey)
	if err != nil {
		return plugin.Ready{}, fmt.Errorf("failed to read %s: %w", checkpointTimeKey, err)
	}
	if at != "" {
		ready.CheckpointTime, err = time.Parse(time.RFC3339, at)
		if err != nil {
			return plugin.Ready{}, fmt.Errorf("failed to parse %s %q: %w", checkpointTimeKey, at, err)
		}
	}
	return ready, nil
}

// Watch replays the changes configured in the spec and returns. Everything is
// replayed when the caller asked for a point in time or when there is no
// checkpoint; otherwise only the changes after the checkpointed event are.
func (c *Client) Watch(ctx context.Context, pos plugin.Position, changes chan<- plugin.Change) error {
	pending := c.config.Changes
	if checkpoint := c.watchCheckpoint(); pos.Since == nil && checkpoint != "" {
		pending = changesAfter(pending, checkpoint)
	}

	for _, change := range pending {
		table := c.tables.Get(change.Table)
		if table == nil {
			return fmt.Errorf("change %q refers to unknown table %q", change.EventID, change.Table)
		}
		keys, err := schema.NewKeysRecord(table, change.Keys)
		if err != nil {
			return fmt.Errorf("failed to build the keys of change %q: %w", change.EventID, err)
		}
		select {
		case changes <- plugin.Change{
			Table:   change.Table,
			Keys:    keys,
			Op:      change.PluginOp(),
			EventID: change.EventID,
			Time:    time.Now().UTC(),
		}:
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

// changesAfter returns the changes that follow the given event id. An id that is
// not in the list leaves every change pending.
func changesAfter(changes []client.Change, eventID string) []client.Change {
	for i, change := range changes {
		if change.EventID == eventID {
			return changes[i+1:]
		}
	}
	return changes
}

// setWatchState records what WatchReady found.
func (c *Client) setWatchState(stateClient state.Client, checkpoint string) {
	c.watchMu.Lock()
	defer c.watchMu.Unlock()
	c.stateClient = stateClient
	c.checkpoint = checkpoint
}

// watchCheckpoint returns the last acked event id.
func (c *Client) watchCheckpoint() string {
	c.watchMu.Lock()
	defer c.watchMu.Unlock()
	return c.checkpoint
}

// Ack checkpoints every change up to and including eventID. It runs on a
// different goroutine than Watch, so it takes the lock for the whole write:
// acks are serialised with each other anyway, and Close must not pull the state
// client out from under one.
func (c *Client) Ack(ctx context.Context, eventID string) error {
	c.watchMu.Lock()
	defer c.watchMu.Unlock()

	if c.stateClient == nil {
		return errors.New("watch not started: call WatchReady first")
	}
	if err := c.stateClient.SetKey(ctx, checkpointKey, eventID); err != nil {
		return fmt.Errorf("failed to store %s: %w", checkpointKey, err)
	}
	if err := c.stateClient.SetKey(ctx, checkpointTimeKey, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("failed to store %s: %w", checkpointTimeKey, err)
	}
	if err := c.stateClient.Flush(ctx); err != nil {
		return fmt.Errorf("failed to flush the state backend: %w", err)
	}
	c.checkpoint = eventID
	return nil
}

func (c *Client) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

// Close releases the state backend. It can be called while a watch is still
// running, so it waits for an in-flight Ack instead of closing underneath it.
func (c *Client) Close(_ context.Context) error {
	c.watchMu.Lock()
	defer c.watchMu.Unlock()

	if c.stateClient == nil {
		return nil
	}
	err := c.stateClient.Close()
	c.stateClient = nil
	return err
}

func (*Client) Write(context.Context, <-chan message.WriteMessage) error {
	// Not implemented, just used for testing destination packaging
	return nil
}

func (*Client) Read(context.Context, *schema.Table, chan<- arrow.RecordBatch) error {
	// Not implemented, just used for testing destination packaging
	return nil
}

func (*Client) Transform(_ context.Context, _ <-chan arrow.RecordBatch, _ chan<- arrow.RecordBatch) error {
	// Not implemented, just used for testing destination packaging
	return nil
}

func (*Client) TransformSchema(_ context.Context, _ *arrow.Schema) (*arrow.Schema, error) {
	// Not implemented, just used for testing destination packaging
	return nil, nil
}

func Configure(_ context.Context, logger zerolog.Logger, spec []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	if opts.NoConnection {
		return &Client{
			logger: logger,
			tables: getTables(),
		}, nil
	}

	config := &client.Spec{}
	if err := json.Unmarshal(spec, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}
	config.SetDefaults()
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate spec: %w", err)
	}

	return &Client{
		logger: logger,
		config: *config,
		scheduler: scheduler.NewScheduler(
			scheduler.WithLogger(logger),
		),
		tables: getTables(),
	}, nil
}

func TestConnection(_ context.Context, _ zerolog.Logger, spec []byte) error {
	config := &client.Spec{}
	if err := json.Unmarshal(spec, config); err != nil {
		return plugin.NewTestConnError("INVALID_SPEC", fmt.Errorf("failed to unmarshal spec: %w", err))
	}

	config.SetDefaults()
	if err := config.Validate(); err != nil {
		return plugin.NewTestConnError("INVALID_SPEC", fmt.Errorf("failed to validate spec: %w", err))
	}
	return nil
}

func getTables() schema.Tables {
	tables := schema.Tables{
		services.TestSomeTable(),
	}
	if err := transformers.TransformTables(tables); err != nil {
		panic(err)
	}
	for _, t := range tables {
		schema.AddCqIDs(t)
		schema.AddCqClientID(t)
	}
	return tables
}
