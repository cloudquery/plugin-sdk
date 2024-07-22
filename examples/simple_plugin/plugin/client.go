package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/cloudquery/plugin-sdk/examples/simple_plugin/client"
	"github.com/cloudquery/plugin-sdk/examples/simple_plugin/services"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/rs/zerolog"
)

type Client struct {
	logger    zerolog.Logger
	config    client.Spec
	tables    schema.Tables
	scheduler *scheduler.Scheduler
}

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	schedulerClient := &client.TestClient{
		Logger: c.logger,
		Spec:   c.config,
	}

	return c.scheduler.Sync(ctx, schedulerClient, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (c *Client) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (*Client) Close(_ context.Context) error {
	return nil
}

func (*Client) Write(context.Context, <-chan message.WriteMessage) error {
	// Not implemented, just used for testing destination packaging
	return nil
}

func (*Client) Read(context.Context, *schema.Table, chan<- arrow.Record) error {
	// Not implemented, just used for testing destination packaging
	return nil
}

func (*Client) Transform(_ context.Context, _ <-chan arrow.Record, _ chan<- arrow.Record) error {
	// Not implemented, just used for testing destination packaging
	return nil
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
	}
	return tables
}
