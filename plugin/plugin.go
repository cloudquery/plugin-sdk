package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

const (
	defaultBatchTimeoutSeconds = 20
	defaultBatchSize           = 10000
	defaultBatchSizeBytes      = 5 * 1024 * 1024 // 5 MiB
)

type NewClientFunc func(context.Context, zerolog.Logger, any) (Client, error)

type ManagedSyncClient interface {
	ID() string
}

type Client interface {
	Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error
	Migrate(ctx context.Context, tables schema.Tables, migrateMode MigrateMode) error
	WriteTableBatch(ctx context.Context, table *schema.Table, writeMode WriteMode, data []arrow.Record) error
	Write(ctx context.Context, tables schema.Tables, writeMode WriteMode, res <-chan arrow.Record) error
	DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error
	Close(ctx context.Context) error
}

type UnimplementedWriter struct{}

func (UnimplementedWriter) Migrate(ctx context.Context, tables schema.Tables, migrateMode MigrateMode) error {
	return fmt.Errorf("not implemented")
}

func (UnimplementedWriter) Write(ctx context.Context, tables schema.Tables, writeMode WriteMode, res <-chan arrow.Record) error {
	return fmt.Errorf("not implemented")
}

func (UnimplementedWriter) WriteTableBatch(ctx context.Context, table *schema.Table, writeMode WriteMode, data []arrow.Record) error {
	return fmt.Errorf("not implemented")
}

func (UnimplementedWriter) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	return fmt.Errorf("not implemented")
}

type UnimplementedSync struct{}

func (UnimplementedSync) Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error {
	return fmt.Errorf("not implemented")
}

type UnimplementedRead struct{}

func (UnimplementedRead) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error {
	return fmt.Errorf("not implemented")
}

// Plugin is the base structure required to pass to sdk.serve
// We take a declarative approach to API here similar to Cobra
type Plugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	name string
	// Version of the plugin
	version string
	// Called upon init call to validate and init configuration
	newClient NewClientFunc
	// dynamic table function if specified
	getDynamicTables GetTables
	// Tables are static tables that defined in compile time by the plugin
	staticTables schema.Tables
	// status sync metrics
	metrics *Metrics
	// Logger to call, this logger is passed to the serve.Serve Client, if not defined Serve will create one instead.
	logger zerolog.Logger
	// resourceSem is a semaphore that limits the number of concurrent resources being fetched
	resourceSem *semaphore.Weighted
	// tableSem is a semaphore that limits the number of concurrent tables being fetched
	tableSems []*semaphore.Weighted
	// maxDepth is the max depth of tables
	maxDepth uint64
	// caser
	caser *caser.Caser
	// mu is a mutex that limits the number of concurrent init/syncs (can only be one at a time)
	mu sync.Mutex

	// client is the initialized session client
	client Client
	// sessionTables are the
	sessionTables schema.Tables
	// spec is the spec the client was initialized with
	spec any
	// NoInternalColumns if set to true will not add internal columns to tables such as _cq_id and _cq_parent_id
	// useful for sources such as PostgreSQL and other databases
	internalColumns bool
	// unmanagedSync if set to true then the plugin will call Sync directly and not use the scheduler
	unmanagedSync bool
	// titleTransformer allows the plugin to control how table names get turned into titles for generated documentation
	titleTransformer  func(*schema.Table) string
	syncTime          time.Time
	sourceName        string
	deterministicCQId bool

	managedWriter bool
	workers       map[string]*worker
	workersLock   *sync.Mutex

	batchTimeout          time.Duration
	defaultBatchSize      int
	defaultBatchSizeBytes int
}

const (
	maxAllowedDepth = 4
)

// Add internal columns
func (p *Plugin) addInternalColumns(tables []*schema.Table) error {
	for _, table := range tables {
		if c := table.Column("_cq_id"); c != nil {
			return fmt.Errorf("table %s already has column _cq_id", table.Name)
		}
		cqID := schema.CqIDColumn
		if len(table.PrimaryKeys()) == 0 {
			cqID.PrimaryKey = true
		}
		cqSourceName := schema.CqSourceNameColumn
		cqSyncTime := schema.CqSyncTimeColumn
		cqSourceName.Resolver = func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			return resource.Set(c.Name, p.sourceName)
		}
		cqSyncTime.Resolver = func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			return resource.Set(c.Name, p.syncTime)
		}

		table.Columns = append([]schema.Column{cqSourceName, cqSyncTime, cqID, schema.CqParentIDColumn}, table.Columns...)
		if err := p.addInternalColumns(table.Relations); err != nil {
			return err
		}
	}
	return nil
}

// Set parent links on relational tables
func setParents(tables schema.Tables, parent *schema.Table) {
	for _, table := range tables {
		table.Parent = parent
		setParents(table.Relations, table)
	}
}

// Apply transformations to tables
func transformTables(tables schema.Tables) error {
	for _, table := range tables {
		if table.Transform != nil {
			if err := table.Transform(table); err != nil {
				return fmt.Errorf("failed to transform table %s: %w", table.Name, err)
			}
		}
		if err := transformTables(table.Relations); err != nil {
			return err
		}
	}
	return nil
}

func maxDepth(tables schema.Tables) uint64 {
	var depth uint64
	if len(tables) == 0 {
		return 0
	}
	for _, table := range tables {
		newDepth := 1 + maxDepth(table.Relations)
		if newDepth > depth {
			depth = newDepth
		}
	}
	return depth
}

// NewPlugin returns a new CloudQuery Plugin with the given name, version and implementation.
// Depending on the options, it can be write only plugin, read only plugin or both.
func NewPlugin(name string, version string, newClient NewClientFunc, options ...Option) *Plugin {
	p := Plugin{
		name:                  name,
		version:               version,
		internalColumns:       true,
		caser:                 caser.New(),
		titleTransformer:      DefaultTitleTransformer,
		newClient:             newClient,
		metrics:               &Metrics{TableClient: make(map[string]map[string]*TableClientMetrics)},
		workers:               make(map[string]*worker),
		workersLock:           &sync.Mutex{},
		batchTimeout:          time.Duration(defaultBatchTimeoutSeconds) * time.Second,
		defaultBatchSize:      defaultBatchSize,
		defaultBatchSizeBytes: defaultBatchSizeBytes,
	}
	for _, opt := range options {
		opt(&p)
	}
	if p.staticTables != nil {
		setParents(p.staticTables, nil)
		if err := transformTables(p.staticTables); err != nil {
			panic(err)
		}
		if p.internalColumns {
			if err := p.addInternalColumns(p.staticTables); err != nil {
				panic(err)
			}
		}
		p.maxDepth = maxDepth(p.staticTables)
		if p.maxDepth > maxAllowedDepth {
			panic(fmt.Errorf("max depth of tables is %d, max allowed is %d", p.maxDepth, maxAllowedDepth))
		}
		if err := p.validate(p.staticTables); err != nil {
			panic(err)
		}
	}

	return &p
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return p.name
}

// Version returns the version of this plugin
func (p *Plugin) Version() string {
	return p.version
}

func (p *Plugin) SetLogger(logger zerolog.Logger) {
	p.logger = logger.With().Str("module", p.name+"-src").Logger()
}

func (p *Plugin) Metrics() *Metrics {
	return p.metrics
}

// Init initializes the plugin with the given spec.
func (p *Plugin) Init(ctx context.Context, spec any) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	var err error
	p.client, err = p.newClient(ctx, p.logger, spec)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	p.spec = spec

	return nil
}

func (p *Plugin) Close(ctx context.Context) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	return p.client.Close(ctx)
}
