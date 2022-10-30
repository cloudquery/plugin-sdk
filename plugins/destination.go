package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type NewDestinationClientFunc func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error)

type DestinationClient interface {
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, table string, data map[string]interface{}) error
	Metrics() DestinationMetrics
	DeleteStale(ctx context.Context, tables string, sourceName string, syncTime time.Time) error
	Close(ctx context.Context) error
}

type DestinationPlugin struct {
	// Name of destination plugin i.e postgresql,snowflake
	name string
	// Version of the destination plugin
	version string
	// Called upon configure call to validate and init configuration
	newDestinationClient NewDestinationClientFunc
	// initialized destination client
	client DestinationClient
	// spec the client was initialized with
	spec specs.Destination
	// tables destination was last migrated with
	tables schema.Tables
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
}

type WriteSummary struct {
	SuccessWrites uint64
	FailedWrites  uint64
	FailedDeletes uint64
}

func NewDestinationPlugin(name string, version string, newDestinationClient NewDestinationClientFunc) *DestinationPlugin {
	p := &DestinationPlugin{
		name:                 name,
		version:              version,
		newDestinationClient: newDestinationClient,
	}
	return p
}

func (p *DestinationPlugin) Name() string {
	return p.name
}

func (p *DestinationPlugin) Version() string {
	return p.version
}

func (p *DestinationPlugin) Metrics() DestinationMetrics {
	return p.client.Metrics()
}

// we need lazy loading because we want to be able to initialize after
func (p *DestinationPlugin) Init(ctx context.Context, logger zerolog.Logger, spec specs.Destination) error {
	var err error
	p.logger = logger
	p.spec = spec
	p.client, err = p.newDestinationClient(ctx, logger, spec)
	if err != nil {
		return err
	}
	return nil
}

// we implement all DestinationClient functions so we can hook into pre-post behavior
func (p *DestinationPlugin) Migrate(ctx context.Context, tables schema.Tables) error {
	SetDestinationManagedCqColumns(tables)

	if p.client == nil {
		return fmt.Errorf("destination client not initialized")
	}
	p.tables = tables
	return p.client.Migrate(ctx, tables)
}

func (p *DestinationPlugin) Write(ctx context.Context, sourceName string, syncTime time.Time, res <-chan *schema.Resource) *WriteSummary {
	if p.client == nil {
		return nil
	}
	summary := WriteSummary{}
	for r := range res {
		r.Data[schema.CqSourceNameColumn.Name] = sourceName
		r.Data[schema.CqSyncTimeColumn.Name] = syncTime
		err := p.client.Write(ctx, r.TableName, r.Data)
		if err != nil {
			summary.FailedWrites++
			p.logger.Error().Str("table", r.TableName).Err(err).Msgf("failed to write to destination")
		} else {
			summary.SuccessWrites++
		}
	}
	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		failedDeletes := p.DeleteStale(ctx, p.tables.TableNames(), sourceName, syncTime)
		summary.FailedDeletes = failedDeletes
	}
	return &summary
}

func (p *DestinationPlugin) DeleteStale(ctx context.Context, tables []string, sourceName string, syncTime time.Time) uint64 {
	if p.client == nil {
		return 0
	}
	failedDeletes := uint64(0)
	for _, t := range tables {
		if err := p.client.DeleteStale(ctx, t, sourceName, syncTime); err != nil {
			p.logger.Error().Err(err).Msgf("failed to delete stale records")
			failedDeletes++
		}
	}
	return failedDeletes
}

func (p *DestinationPlugin) Close(ctx context.Context) error {
	if p.client == nil {
		return fmt.Errorf("destination client not initialized")
	}
	return p.client.Close(ctx)
}

// Overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func SetDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)

		SetDestinationManagedCqColumns(table.Relations)
	}
}
