package source

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/backend"
	"github.com/cloudquery/plugin-sdk/caser"
	"github.com/cloudquery/plugin-sdk/internal/backends/local"
	"github.com/cloudquery/plugin-sdk/internal/backends/nop"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

type Options struct {
	Backend backend.Backend
}

type NewExecutionClientFunc func(context.Context, zerolog.Logger, specs.Source, Options) (schema.ClientMeta, error)

type SourceUnmanagedClient interface {
	schema.ClientMeta
	Sync(ctx context.Context, metrics *Metrics, res chan<- *schema.Resource) error
}

// Plugin is the base structure required to pass to sdk.serve
// We take a declarative approach to API here similar to Cobra
type Plugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	name string
	// Version of the plugin
	version string
	// Called upon configure call to validate and init configuration
	newExecutionClient NewExecutionClientFunc
	// dynamic table function if specified
	getDynamicTables GetTables
	// Tables is all tables supported by this source plugin
	tables schema.Tables
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
	client schema.ClientMeta
	// sessionTables are the
	sessionTables schema.Tables
	// backend is the backend used to store the cursor state
	backend backend.Backend
	// spec is the spec the client was initialized with
	spec specs.Source
	// NoInternalColumns if set to true will not add internal columns to tables such as _cq_id and _cq_parent_id
	// useful for sources such as PostgreSQL and other databases
	internalColumns bool
	// unmanaged if set to true then the plugin will call Sync directly and not use the scheduler
	unmanaged bool
}

const (
	maxAllowedDepth = 4
)

// Add internal columns
func addInternalColumns(tables []*schema.Table) error {
	for _, table := range tables {
		if c := table.Column("_cq_id"); c != nil {
			return fmt.Errorf("table %s already has column _cq_id", table.Name)
		}
		cqID := schema.CqIDColumn
		if len(table.PrimaryKeys()) == 0 {
			cqID.CreationOptions.PrimaryKey = true
		}
		table.Columns = append([]schema.Column{cqID, schema.CqParentIDColumn}, table.Columns...)
		if err := addInternalColumns(table.Relations); err != nil {
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

// NewPlugin returns a new plugin with a given name, version, tables, newExecutionClient
// and additional options.
func NewPlugin(name string, version string, tables []*schema.Table, newExecutionClient NewExecutionClientFunc, options ...Option) *Plugin {
	p := Plugin{
		name:               name,
		version:            version,
		tables:             tables,
		newExecutionClient: newExecutionClient,
		metrics:            &Metrics{TableClient: make(map[string]map[string]*TableClientMetrics)},
		caser:              caser.New(),
		internalColumns: 	true,
	}
	for _, opt := range options {
		opt(&p)
	}
	setParents(p.tables, nil)
	if err := transformTables(p.tables); err != nil {
		panic(err)
	}
	if p.internalColumns {
		if err := addInternalColumns(p.tables); err != nil {
			panic(err)
		}
	}
	if err := p.validate(); err != nil {
		panic(err)
	}
	p.maxDepth = maxDepth(p.tables)
	if p.maxDepth > maxAllowedDepth {
		panic(fmt.Errorf("max depth of tables is %d, max allowed is %d", p.maxDepth, maxAllowedDepth))
	}
	return &p
}

func (p *Plugin) SetLogger(logger zerolog.Logger) {
	p.logger = logger.With().Str("module", p.name+"-src").Logger()
}

// Tables returns all tables supported by this source plugin
func (p *Plugin) Tables() schema.Tables {
	return p.tables
}

func (p *Plugin) HasDynamicTables() bool {
	return p.getDynamicTables != nil
}

func (p *Plugin) GetDynamicTables() schema.Tables {
	return p.sessionTables
}

// TablesForSpec returns all tables supported by this source plugin that match the given spec.
// It validates the tables part of the spec and will return an error if it is found to be invalid.
// This is deprecated method
func (p *Plugin) TablesForSpec(spec specs.Source) (schema.Tables, error) {
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}
	tables, err := p.tables.FilterDfs(spec.Tables, spec.SkipTables, spec.SkipDependentTables)
	if err != nil {
		return nil, fmt.Errorf("failed to filter tables: %w", err)
	}
	return tables, nil
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return p.name
}

// Version returns the version of this plugin
func (p *Plugin) Version() string {
	return p.version
}

func (p *Plugin) Metrics() *Metrics {
	return p.metrics
}

func (p *Plugin) Init(ctx context.Context, spec specs.Source) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()

	var err error
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		return fmt.Errorf("invalid spec: %w", err)
	}
	p.spec = spec

	switch spec.Backend {
	case specs.BackendNone:
		p.backend = nop.New()
	case specs.BackendLocal:
		p.backend, err = local.New(spec)
		if err != nil {
			return fmt.Errorf("failed to initialize local backend: %w", err)
		}
	default:
		return fmt.Errorf("unknown backend: %s", spec.Backend)
	}

	tables := p.tables
	if p.getDynamicTables != nil {
		p.client, err = p.newExecutionClient(ctx, p.logger, spec, Options{Backend: p.backend})
		if err != nil {
			return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
		}
		tables, err = p.getDynamicTables(ctx, p.client)
		if err != nil {
			return fmt.Errorf("failed to get dynamic tables: %w", err)
		}

		tables, err = tables.FilterDfs(spec.Tables, spec.SkipTables, spec.SkipDependentTables)
		if err != nil {
			return fmt.Errorf("failed to filter tables: %w", err)
		}
		if len(tables) == 0 {
			return fmt.Errorf("no tables to sync - please check your spec 'tables' and 'skip_tables' settings")
		}

		setParents(tables, nil)
		if err := transformTables(tables); err != nil {
			return err
		}
		if p.internalColumns {
			if err := addInternalColumns(tables); err != nil {
				return err
			}
		}
		if err := p.validate(); err != nil {
			return err
		}
		p.maxDepth = maxDepth(tables)
		if p.maxDepth > maxAllowedDepth {
			return fmt.Errorf("max depth of tables is %d, max allowed is %d", p.maxDepth, maxAllowedDepth)
		}
	} else {
		tables, err = tables.FilterDfs(spec.Tables, spec.SkipTables, spec.SkipDependentTables)
		if err != nil {
			return fmt.Errorf("failed to filter tables: %w", err)
		}
	}

	p.sessionTables = tables
	return nil
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *Plugin) Sync(ctx context.Context, res chan<- *schema.Resource) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()

	if p.client == nil {
		var err error
		p.client, err = p.newExecutionClient(ctx, p.logger, p.spec, Options{Backend: p.backend})
		if err != nil {
			return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
		}
	}

	startTime := time.Now()
	if p.unmanaged {
		unmanagedClient := p.client.(SourceUnmanagedClient)
		if err := unmanagedClient.Sync(ctx, p.metrics, res); err != nil {
			return fmt.Errorf("failed to sync unmanaged client: %w", err)
		}
	} else {
		switch p.spec.Scheduler {
		case specs.SchedulerDFS:
			p.syncDfs(ctx, p.spec, p.client, p.sessionTables, res)
		case specs.SchedulerRoundRobin:
			p.syncRoundRobin(ctx, p.spec, p.client, p.sessionTables, res)
		default:
			return fmt.Errorf("unknown scheduler %s. Options are: %v", p.spec.Scheduler, specs.AllSchedulers.String())
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
