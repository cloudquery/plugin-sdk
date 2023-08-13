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

type NewClientOptions struct {
	NoConnection bool
}

type NewClientFunc func(context.Context, zerolog.Logger, []byte, NewClientOptions) (Client, error)

type Client interface {
	SourceClient
	DestinationClient
}

type UnimplementedDestination struct{}

func (UnimplementedDestination) Write(context.Context, <-chan message.WriteMessage) error {
	return ErrNotImplemented
}

func (UnimplementedDestination) Read(context.Context, *schema.Table, chan<- arrow.Record) error {
	return ErrNotImplemented
}

type UnimplementedSource struct{}

func (UnimplementedSource) Sync(context.Context, SyncOptions, chan<- message.SyncMessage) error {
	return ErrNotImplemented
}

func (UnimplementedSource) Tables(context.Context, TableOptions) (schema.Tables, error) {
	return nil, ErrNotImplemented
}

// Plugin is the base structure required to pass to sdk.serve
// We take a declarative approach to API here similar to Cobra
type Plugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	name string
	// Version of the plugin
	version string
	// Title of the plugin as appears in CloudQuery registry
	title string
	// Short description of the plugin as appears in CloudQuery registry
	shortDescription string
	// Long description of the plugin as appears in CloudQuery registry
	description string
	// categories of the plugin as appears in CloudQuery registry
	categories []string
	// targets to build plugin for
	targets []BuildTarget
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

// NewPlugin returns a new CloudQuery Plugin with the given name, version and implementation.
// Depending on the options, it can be a write-only plugin, read-only plugin, or both.
func NewPlugin(name string, version string, newClient NewClientFunc, options ...Option) *Plugin {
	p := Plugin{
		name:            name,
		version:         version,
		internalColumns: true,
		newClient:       newClient,
		title: 				 	 name,
		categories: 		[]string{},
		targets: buildTargets,
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

func (p *Plugin) Title() string {
	return p.title
}

func (p *Plugin) Description() string {
	return p.description
}

func (p *Plugin) ShortDescription() string {
	return p.shortDescription
}

func (p *Plugin) Categories() []string {
	return p.categories
}

func (p *Plugin) Targets() []BuildTarget {
	return p.targets
}

func (p *Plugin) SetLogger(logger zerolog.Logger) {
	p.logger = logger.With().Str("module", p.name+"-src").Logger()
}

func (p *Plugin) Tables(ctx context.Context, options TableOptions) (schema.Tables, error) {
	if p.client == nil {
		return nil, fmt.Errorf("plugin not initialized")
	}
	tables, err := p.client.Tables(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	return tables, nil
}

// Init initializes the plugin with the given spec.
func (p *Plugin) Init(ctx context.Context, spec []byte, options NewClientOptions) error {
	if !p.mu.TryLock() {
		return fmt.Errorf("plugin already in use")
	}
	defer p.mu.Unlock()
	var err error
	p.client, err = p.newClient(ctx, p.logger, spec, options)
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
