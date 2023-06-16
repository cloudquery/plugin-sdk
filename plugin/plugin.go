package plugin

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

var ErrNotImplemented = fmt.Errorf("not implemented")

type NewClientFunc func(context.Context, zerolog.Logger, []byte) (Client, error)

type Client interface {
	SourceClient
	DestinationClient
}

type UnimplementedDestination struct{}

func (UnimplementedDestination) Write(ctx context.Context, options WriteOptions, res <-chan message.Message) error {
	return ErrNotImplemented
}

func (UnimplementedDestination) Read(ctx context.Context, table *schema.Table, res chan<- arrow.Record) error {
	return fmt.Errorf("not implemented")
}

type UnimplementedSource struct{}

func (UnimplementedSource) Sync(ctx context.Context, options SyncOptions, res chan<- message.Message) error {
	return ErrNotImplemented
}

func (UnimplementedSource) Tables(ctx context.Context) (schema.Tables, error) {
	return nil, ErrNotImplemented
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
	// mu is a mutex that limits the number of concurrent init/syncs (can only be one at a time)
	mu sync.Mutex
	// client is the initialized session client
	client Client
	// spec is the spec the client was initialized with
	spec any
	// NoInternalColumns if set to true will not add internal columns to tables such as _cq_id and _cq_parent_id
	// useful for sources such as PostgreSQL and other databases
	internalColumns bool
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
// Depending on the options, it can be a write-only plugin, read-only plugin, or both.
func NewPlugin(name string, version string, newClient NewClientFunc, options ...Option) *Plugin {
	p := Plugin{
		name:            name,
		version:         version,
		internalColumns: true,
		newClient:       newClient,
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
	if p.client == nil {
		return nil, fmt.Errorf("plugin not initialized")
	}
	tables, err := p.client.Tables(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	return tables, nil
}

// Init initializes the plugin with the given spec.
func (p *Plugin) Init(ctx context.Context, spec []byte) error {
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
	if p.client == nil {
		return nil
	}
	return p.client.Close(ctx)
}
