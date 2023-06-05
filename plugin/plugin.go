package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type NewClientFunc func(context.Context, zerolog.Logger, any) (Client, error)

type Client interface {
	Tables(ctx context.Context) (schema.Tables, error)
	Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error
	Migrate(ctx context.Context, tables schema.Tables, migrateMode MigrateOptions) error
	Write(ctx context.Context, options WriteOptions, res <-chan arrow.Record) error
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

func (UnimplementedWriter) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	return fmt.Errorf("not implemented")
}

type UnimplementedSync struct{}

func (UnimplementedSync) Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error {
	return fmt.Errorf("not implemented")
}

func (UnimplementedSync) Tables(ctx context.Context) (schema.Tables, error) {
	return nil, fmt.Errorf("not implemented")
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
	// Logger to call, this logger is passed to the serve.Serve Client, if not defined Serve will create one instead.
	logger zerolog.Logger
	// maxDepth is the max depth of tables
	maxDepth uint64
	// mu is a mutex that limits the number of concurrent init/syncs (can only be one at a time)
	mu sync.Mutex
	// client is the initialized session client
	client Client
	// spec is the spec the client was initialized with
	spec any
	// NoInternalColumns if set to true will not add internal columns to tables such as _cq_id and _cq_parent_id
	// useful for sources such as PostgreSQL and other databases
	internalColumns bool
	// titleTransformer allows the plugin to control how table names get turned into titles for generated documentation
	titleTransformer func(*schema.Table) string
}

const (
	maxAllowedDepth = 4
)

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
		name:             name,
		version:          version,
		internalColumns:  true,
		titleTransformer: DefaultTitleTransformer,
		newClient:        newClient,
	}
	for _, opt := range options {
		opt(&p)
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

func (p *Plugin) Tables(ctx context.Context) (schema.Tables, error) {
	tables, err := p.client.Tables(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	setParents(tables, nil)
	if err := transformTables(tables); err != nil {
		return nil, err
	}
	if p.internalColumns {
		if err := p.addInternalColumns(tables); err != nil {
			return nil, err
		}
	}
	p.maxDepth = maxDepth(tables)
	if p.maxDepth > maxAllowedDepth {
		return nil, fmt.Errorf("max depth of tables is %d, max allowed is %d", p.maxDepth, maxAllowedDepth)
	}
	if err := p.validate(tables); err != nil {
		return nil, err
	}
	return tables, nil
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
