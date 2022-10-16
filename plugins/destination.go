package plugins

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type NewDestinationClientFunc func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error)

type DestinationStats struct {
	// Errors number of errors / failed writes
	Errors uint64
	// Writes number of successful writes
	Writes uint64
}

type DestinationClient interface {
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, res <-chan *schema.DestinationResource) error
	Stats() DestinationStats
	DeleteStale(ctx context.Context, tables []string, sourceName string, syncTime time.Time) error
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
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
}

func (s *DestinationStats) Add(other *DestinationStats) {
	if other == nil {
		return
	}
	atomic.AddUint64(&s.Errors, other.Errors)
	atomic.AddUint64(&s.Writes, other.Writes)
}

const writeWorkers = 2

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

func (p *DestinationPlugin) Stats() DestinationStats {
	return p.client.Stats()
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
	return p.client.Migrate(ctx, tables)
}

func (p *DestinationPlugin) Write(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, res <-chan *schema.DestinationResource) error {
	ch := make(chan *schema.DestinationResource)
	SetDestinationManagedCqColumns(tables)

	eg, ctx := errgroup.WithContext(ctx)
	// given most destination plugins writing in batch we are using a worker pool to write in parallel
	// it might not generalize well and we might need to move it to each destination plugin implementation.
	for i := 0; i < writeWorkers; i++ {
		eg.Go(func() error {
			return p.client.Write(ctx, ch)
		})
	}
	for {
		select {
		case <-ctx.Done():
			res = nil
		case r, ok := <-res:
			if ok {
				r.Data = append([]interface{}{sourceName, syncTime}, r.Data...)
				ch <- r
			} else {
				res = nil
			}
		}
		if res == nil {
			break
		}
	}

	close(ch)
	if err := eg.Wait(); err != nil {
		return err
	}
	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		if err := p.DeleteStale(ctx, tables.TableNames(), sourceName, syncTime); err != nil {
			return err
		}
	}
	return nil
}

func (p *DestinationPlugin) DeleteStale(ctx context.Context, tables []string, sourceName string, syncTime time.Time) error {
	return p.client.DeleteStale(ctx, tables, sourceName, syncTime)
}

func (p *DestinationPlugin) Close(ctx context.Context) error {
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
