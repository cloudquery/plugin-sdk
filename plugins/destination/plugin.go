package destination

import (
	"context"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type writerType int

const (
	unmanaged writerType = iota
	managed
)

const (
	defaultBatchTimeoutSeconds = 20
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

type NewClientFunc func(context.Context, zerolog.Logger, specs.Destination) (Client, error)

type ManagedWriter interface {
	WriteTableBatch(ctx context.Context, table *schema.Table, data [][]any) error
}

type UnimplementedManagedWriter struct{}

type UnmanagedWriter interface {
	Write(ctx context.Context, tables schema.Tables, res <-chan *ClientResource) error
	Metrics() Metrics
}

type UnimplementedUnmanagedWriter struct{}

func (*UnimplementedManagedWriter) WriteTableBatch(context.Context, *schema.Table, [][]any) error {
	panic("WriteTableBatch not implemented")
}

func (*UnimplementedUnmanagedWriter) Write(context.Context, schema.Tables, <-chan *ClientResource) error {
	panic("Write not implemented")
}

func (*UnimplementedUnmanagedWriter) Metrics() Metrics {
	panic("Metrics not implemented")
}

type Client interface {
	schema.CQTypeTransformer
	ReverseTransformValues(table *schema.Table, values []any) (schema.CQTypes, error)
	Migrate(ctx context.Context, tables schema.Tables) error
	Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- []any) error
	ManagedWriter
	UnmanagedWriter
	DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error
	Close(ctx context.Context) error
}

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
	newClient  NewClientFunc
	writerType writerType
	// initialized destination client
	client Client
	// spec the client was initialized with
	spec specs.Destination
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger

	// This is in use if the user passed a managed client
	metrics     map[string]*Metrics
	metricsLock *sync.RWMutex

	workers     map[string]*worker
	workersLock *sync.Mutex

	batchTimeout          time.Duration
	defaultBatchSize      int
	defaultBatchSizeBytes int
}

func WithManagedWriter() Option {
	return func(p *Plugin) {
		p.writerType = managed
	}
}

func WithBatchTimeout(seconds int) Option {
	return func(p *Plugin) {
		p.batchTimeout = time.Duration(seconds) * time.Second
	}
}

func WithDefaultBatchSize(defaultBatchSize int) Option {
	return func(p *Plugin) {
		p.defaultBatchSize = defaultBatchSize
	}
}

func WithDefaultBatchSizeBytes(defaultBatchSizeBytes int) Option {
	return func(p *Plugin) {
		p.defaultBatchSizeBytes = defaultBatchSizeBytes
	}
}

// NewPlugin creates a new destination plugin
func NewPlugin(name string, version string, newClientFunc NewClientFunc, opts ...Option) *Plugin {
	p := &Plugin{
		name:                  name,
		version:               version,
		newClient:             newClientFunc,
		metrics:               make(map[string]*Metrics),
		metricsLock:           &sync.RWMutex{},
		workers:               make(map[string]*worker),
		workersLock:           &sync.Mutex{},
		batchTimeout:          time.Duration(defaultBatchTimeoutSeconds) * time.Second,
		defaultBatchSize:      defaultBatchSize,
		defaultBatchSizeBytes: defaultBatchSizeBytes,
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
	switch p.writerType {
	case unmanaged:
		return p.client.Metrics()
	case managed:
		metrics := Metrics{}
		p.metricsLock.RLock()
		for _, m := range p.metrics {
			metrics.Errors += m.Errors
			metrics.Writes += m.Writes
		}
		p.metricsLock.RUnlock()
		return metrics
	default:
		panic("unknown client type")
	}
}

// we need lazy loading because we want to be able to initialize after
func (p *Plugin) Init(ctx context.Context, logger zerolog.Logger, spec specs.Destination) error {
	var err error
	p.logger = logger
	p.spec = spec
	p.spec.SetDefaults(p.defaultBatchSize, p.defaultBatchSizeBytes)
	p.client, err = p.newClient(ctx, logger, p.spec)
	if err != nil {
		return err
	}
	return nil
}

// we implement all DestinationClient functions so we can hook into pre-post behavior
func (p *Plugin) Migrate(ctx context.Context, tables schema.Tables) error {
	SetDestinationManagedCqColumns(tables)
	setCqIDColumnOptionsForTables(tables)
	p.setPKsForTables(tables)
	return p.client.Migrate(ctx, tables)
}

func (p *Plugin) readAll(ctx context.Context, table *schema.Table, sourceName string) ([]schema.CQTypes, error) {
	var readErr error
	ch := make(chan schema.CQTypes)
	go func() {
		defer close(ch)
		readErr = p.Read(ctx, table, sourceName, ch)
	}()
	// nolint:prealloc
	var resources []schema.CQTypes
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, readErr
}

func (p *Plugin) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- schema.CQTypes) error {
	SetDestinationManagedCqColumns(schema.Tables{table})
	ch := make(chan []any)
	var err error
	go func() {
		defer close(ch)
		err = p.client.Read(ctx, table, sourceName, ch)
	}()
	for resource := range ch {
		r, err := p.client.ReverseTransformValues(table, resource)
		if err != nil {
			return err
		}
		res <- r
	}
	return err
}

// this function is currently used mostly for testing so it's not a public api
func (p *Plugin) writeOne(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, resource schema.DestinationResource) error {
	resources := []schema.DestinationResource{resource}
	return p.writeAll(ctx, sourceSpec, tables, syncTime, resources)
}

// this function is currently used mostly for testing so it's not a public api
func (p *Plugin) writeAll(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, resources []schema.DestinationResource) error {
	ch := make(chan schema.DestinationResource, len(resources))
	for _, resource := range resources {
		ch <- resource
	}
	close(ch)
	return p.Write(ctx, sourceSpec, tables, syncTime, ch)
}

func (p *Plugin) Write(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan schema.DestinationResource) error {
	syncTime = syncTime.UTC()
	SetDestinationManagedCqColumns(tables)
	p.setPKsForTables(tables)
	switch p.writerType {
	case unmanaged:
		if err := p.writeUnmanaged(ctx, sourceSpec, tables, syncTime, res); err != nil {
			return err
		}
	case managed:
		if err := p.writeManagedTableBatch(ctx, sourceSpec, tables, syncTime, res); err != nil {
			return err
		}
	default:
		panic("unknown client type")
	}
	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		tablesToDelete := tables
		if sourceSpec.Backend != specs.BackendNone {
			include := func(t *schema.Table) bool {
				return true
			}
			exclude := func(t *schema.Table) bool {
				return t.IsIncremental
			}
			tablesToDelete = tables.FilterDfsFunc(include, exclude, sourceSpec.SkipDependentTables)
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

// Overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func SetDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)
		SetDestinationManagedCqColumns(table.Relations)
	}
}

// this is for backward compatibility for sources that didn't update the SDK yet
// TODO: remove this in the future once all sources have updated the SDK
func setCqIDColumnOptionsForTables(tables []*schema.Table) {
	for _, table := range tables {
		for i, c := range table.Columns {
			if c.Name == schema.CqIDColumn.Name {
				table.Columns[i].CreationOptions.NotNull = true
				table.Columns[i].CreationOptions.Unique = true
			}
		}
		setCqIDColumnOptionsForTables(table.Relations)
	}
}

func (p *Plugin) setPKsForTables(tables schema.Tables) {
	if p.spec.PKMode == specs.PKModeCQID {
		setCQIDAsPrimaryKeysForTables(tables)
	}
}
func setCQIDAsPrimaryKeysForTables(tables schema.Tables) {
	for _, table := range tables {
		for i, col := range table.Columns {
			table.Columns[i].CreationOptions.PrimaryKey = col.Name == schema.CqIDColumn.Name
		}
		setCQIDAsPrimaryKeysForTables(table.Relations)
	}
}
