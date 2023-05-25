package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/rs/zerolog"
)

type NewClientFunc func(context.Context, zerolog.Logger, specs.Destination) (Client, error)

type UnmanagedWriter interface {
	Write(context.Context, specs.Source, schema.Tables, time.Time, <-chan arrow.Record) error
	Metrics() Metrics
}

type ManagedWriter interface {
	WriteTableBatch(context.Context, specs.Source, *schema.Table, time.Time, []arrow.Record) error
}

type Client interface {
	Migrate(ctx context.Context, tables schema.Tables) error
	Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error
	ManagedWriter
	UnmanagedWriter
	DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error
	Close(ctx context.Context) error
}

type BatchingWriter interface {
	UnmanagedWriter
}

type (
	BatchingWriterFunc     func(ManagedWriter) BatchingWriter
	BatchingWriterFuncFunc func(*specs.Destination) BatchingWriterFunc
)

type ClientResource struct {
	TableName string
	Data      []any
}

type Option func(*Plugin)

type Plugin struct {
	// Name of destination plugin i.e postgresql,snowflake
	name string
	// Version of the destination plugin
	version string
	// Called upon configure call to validate and init configuration
	newClient NewClientFunc
	// initialized destination client
	client Client
	// spec the client was initialized with
	spec specs.Destination
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
	// batchingWriter to use instead of passing data to client's Writer
	batchingWriter BatchingWriter
	// batchingWriterFuncFunc is to create the batchingWriter during Client initialization
	batchingWriterFuncFunc BatchingWriterFuncFunc
}

func WithManagedWriter(f BatchingWriterFuncFunc) Option {
	return func(p *Plugin) {
		p.batchingWriterFuncFunc = f
	}
}

// NewPlugin creates a new destination plugin
func NewPlugin(name string, version string, newClientFunc NewClientFunc, opts ...Option) *Plugin {
	p := &Plugin{
		name:      name,
		version:   version,
		newClient: newClientFunc,
	}
	if newClientFunc == nil {
		// we do this check because we only call this during runtime later on so it can fail
		// before the server starts
		panic("newClientFunc can't be nil")
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Plugin) Name() string {
	return p.name
}

func (p *Plugin) Version() string {
	return p.version
}

func (p *Plugin) Metrics() Metrics {
	if p.batchingWriter == nil {
		return p.client.Metrics()
	}
	return p.batchingWriter.Metrics()
}

// we need lazy loading because we want to be able to initialize after
func (p *Plugin) Init(ctx context.Context, logger zerolog.Logger, spec specs.Destination) error {
	p.logger = logger
	p.spec = spec

	var bwf BatchingWriterFunc
	if p.batchingWriterFuncFunc != nil {
		bwf = p.batchingWriterFuncFunc(&p.spec)
	}

	var err error
	p.client, err = p.newClient(ctx, logger, p.spec)
	if err != nil {
		return err
	}
	if bwf != nil {
		p.batchingWriter = bwf(p.client)
	}

	return nil
}

// we implement all DestinationClient functions so we can hook into pre-post behavior
func (p *Plugin) Migrate(ctx context.Context, tables schema.Tables) error {
	if err := checkDestinationColumns(tables); err != nil {
		return err
	}
	return p.client.Migrate(ctx, tables)
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

// this function is currently used mostly for testing, so it's not a public api
func (p *Plugin) writeOne(ctx context.Context, sourceSpec specs.Source, syncTime time.Time, resource arrow.Record) error {
	resources := []arrow.Record{resource}
	return p.writeAll(ctx, sourceSpec, syncTime, resources)
}

// this function is currently used mostly for testing, so it's not a public api
func (p *Plugin) writeAll(ctx context.Context, sourceSpec specs.Source, syncTime time.Time, resources []arrow.Record) error {
	ch := make(chan arrow.Record, len(resources))
	for _, resource := range resources {
		ch <- resource
	}
	close(ch)
	tables := make(schema.Tables, 0)
	tableNames := make(map[string]struct{})
	for _, resource := range resources {
		sc := resource.Schema()
		tableMD := sc.Metadata()
		name, found := tableMD.GetValue(schema.MetadataTableName)
		if !found {
			return fmt.Errorf("missing table name")
		}
		if _, ok := tableNames[name]; ok {
			continue
		}
		table, err := schema.NewTableFromArrowSchema(resource.Schema())
		if err != nil {
			return err
		}
		tables = append(tables, table)
		tableNames[table.Name] = struct{}{}
	}
	return p.Write(ctx, sourceSpec, tables, syncTime, ch)
}

func (p *Plugin) Write(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan arrow.Record) error {
	syncTime = syncTime.UTC()
	err := checkDestinationColumns(tables)
	if err != nil {
		return err
	}

	if p.batchingWriter == nil {
		err = p.client.Write(ctx, sourceSpec, tables, syncTime, res)
	} else {
		err = p.batchingWriter.Write(ctx, sourceSpec, tables, syncTime, res)
	}
	if err != nil {
		return err
	}
	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		tablesToDelete := tables
		if sourceSpec.Backend != specs.BackendNone {
			tablesToDelete = make(schema.Tables, 0, len(tables))
			for _, t := range tables {
				if !t.IsIncremental {
					tablesToDelete = append(tablesToDelete, t)
				}
			}
		}
		if err := p.DeleteStale(ctx, tablesToDelete, sourceSpec.Name, syncTime); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	syncTime = syncTime.UTC()
	return p.client.DeleteStale(ctx, tables, sourceName, syncTime)
}

func (p *Plugin) Close(ctx context.Context) error {
	return p.client.Close(ctx)
}

// BatchWriter returns the current batching writer or nil, used in testing
func (p *Plugin) BatchingWriter() BatchingWriter {
	return p.batchingWriter
}

func checkDestinationColumns(tables schema.Tables) error {
	for _, table := range tables {
		if table.Columns.Index(schema.CqSourceNameColumn.Name) == -1 {
			return fmt.Errorf("table %s is missing column %s. please consider upgrading source plugin", table.Name, schema.CqSourceNameColumn.Name)
		}
		if table.Columns.Index(schema.CqSyncTimeColumn.Name) == -1 {
			return fmt.Errorf("table %s is missing column %s. please consider upgrading source plugin", table.Name, schema.CqSourceNameColumn.Name)
		}
		column := table.Columns.Get(schema.CqIDColumn.Name)
		if column != nil {
			if !column.NotNull {
				return fmt.Errorf("column %s.%s cannot be nullable. please consider upgrading source plugin", table.Name, schema.CqIDColumn.Name)
			}
			if !column.Unique {
				return fmt.Errorf("column %s.%s must be unique. please consider upgrading source plugin", table.Name, schema.CqIDColumn.Name)
			}
		}
	}
	return nil
}
