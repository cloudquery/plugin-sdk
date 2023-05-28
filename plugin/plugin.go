package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/backend"
	"github.com/cloudquery/plugin-sdk/v3/caser"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"

	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v0"
)

type Options struct {
	Backend backend.Backend
}

type NewExecutionClientFunc func(context.Context, zerolog.Logger, specs.Source, Options) (schema.ClientMeta, error)

type NewClientFunc func(context.Context, zerolog.Logger, pbPlugin.Spec) (Client, error)

type UnmanagedClient interface {
	schema.ClientMeta
	Sync(ctx context.Context, metrics *Metrics, syncSpec pbPlugin.SyncSpec, res chan<- *schema.Resource) error
}

type Client interface {
	Sync(ctx context.Context, metrics *Metrics, res chan<- *schema.Resource) error
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, tables schema.Tables, res <-chan arrow.Record) error
	DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error
	Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error
}

type UnimplementedWriter struct{}

func (UnimplementedWriter) WriteTableBatch(context.Context, *schema.Table, []arrow.Record) error {
	return fmt.Errorf("not implemented")
}

type UnimplementedSync struct{}

func (UnimplementedSync) Sync(ctx context.Context, metrics *Metrics, res chan<- *schema.Resource) error {
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
	// backend is the backend used to store the cursor state
	backend backend.Backend
	// spec is the spec the client was initialized with
	spec pbPlugin.Spec
	// NoInternalColumns if set to true will not add internal columns to tables such as _cq_id and _cq_parent_id
	// useful for sources such as PostgreSQL and other databases
	internalColumns bool
	// unmanaged if set to true then the plugin will call Sync directly and not use the scheduler
	unmanaged bool
	// titleTransformer allows the plugin to control how table names get turned into titles for generated documentation
	titleTransformer func(*schema.Table) string
	syncTime         time.Time
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
			return resource.Set(c.Name, p.spec.Name)
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

func NewPlugin(name string, version string, newClient NewClientFunc, options ...Option) *Plugin {
	p := Plugin{
		name:             name,
		version:          version,
		internalColumns:    true,
		caser:              caser.New(),
		titleTransformer:   DefaultTitleTransformer,
		newClient: 				newClient,
	}
	for _, opt := range options {
		opt(&p)
	}
	if p.staticTables != nil {
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

// Tables returns all tables supported by this source plugin
func (p *Plugin) StaticTables() schema.Tables {
	return p.staticTables
}

func (p *Plugin) HasDynamicTables() bool {
	return p.getDynamicTables != nil
}

func (p *Plugin) DynamicTables() schema.Tables {
	return p.sessionTables
}

func (p *Plugin) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error {
	return p.client.Read(ctx, table, sourceName, res)
}

func (p *Plugin) Metrics() *Metrics {
	return p.metrics
}

func (p *Plugin) Init(ctx context.Context, spec pbPlugin.Spec) error {
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

func (p *Plugin) Migrate(ctx context.Context, tables schema.Tables) error {
	return p.client.Migrate(ctx, tables)
}

func (p *Plugin) writeUnmanaged(ctx context.Context, _ specs.Source, tables schema.Tables, _ time.Time, res <-chan arrow.Record) error {
	return p.client.Write(ctx, tables, res)
}

func (p *Plugin) Write(ctx context.Context, sourceSpec pbPlugin.Spec, tables schema.Tables, syncTime time.Time, res <-chan arrow.Record) error {
	syncTime = syncTime.UTC()
	if err := p.client.Write(ctx, tables, res); err != nil {
		return err
	}
	if p.spec.WriteSpec.WriteMode == pbPlugin.WRITE_MODE_WRITE_MODE_OVERWRITE_DELETE_STALE {
		tablesToDelete := tables
		if sourceSpec.BackendSpec != nil {
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

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, syncTime time.Time, syncSpec pbPlugin.SyncSpec, res chan<- *schema.Resource) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	p.syncTime = syncTime

	startTime := time.Now()
	if p.unmanaged {
		unmanagedClient := p.client.(UnmanagedClient)
		if err := unmanagedClient.Sync(ctx, p.metrics, syncSpec, res); err != nil {
			return fmt.Errorf("failed to sync unmanaged client: %w", err)
		}
	} else {
		switch syncSpec.Scheduler {
		case pbPlugin.SyncSpec_SCHEDULER_DFS:
			p.syncDfs(ctx, syncSpec, p.client, p.sessionTables, res)
		case pbPlugin.SyncSpec_SCHEDULER_ROUND_ROBIN:
			p.syncRoundRobin(ctx, syncSpec, p.client, p.sessionTables, res)
		default:
			return fmt.Errorf("unknown scheduler %s. Options are: %v", syncSpec.Scheduler, specs.AllSchedulers.String())
		}
	}

	p.logger.Info().Uint64("resources", p.metrics.TotalResources()).Uint64("errors", p.metrics.TotalErrors()).Uint64("panics", p.metrics.TotalPanics()).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}

func (p *Plugin) Close(ctx context.Context) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	if p.backend != nil {
		err := p.backend.Close(ctx)
		if err != nil {
			return fmt.Errorf("failed to close backend: %w", err)
		}
		p.backend = nil
	}
	return nil
}
