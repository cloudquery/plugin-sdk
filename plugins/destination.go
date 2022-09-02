package plugins

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type DestinationNewExecutionClientFunc func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error)

type DestinationClient interface {
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, table string, data map[string]interface{}) error
	SetLogger(logger zerolog.Logger)
}

type DestinationOption func(*DestinationPlugin)

func WithDestinationExampleConfig(exampleConfig string) DestinationOption {
	return func(p *DestinationPlugin) {
		p.exampleConfig = exampleConfig
	}
}

func WithDestinationLogger(logger zerolog.Logger) DestinationOption {
	return func(p *DestinationPlugin) {
		p.logger = logger
	}
}

// DestinationPlugin is the base structure required by calls to serve.Serve.
type DestinationPlugin struct {
	// name of plugin i.e aws, gcp, azure, etc
	name string
	// version of the plugin
	version string
	// example config for the plugin
	exampleConfig string
	// called on init to create a new DestinationClient
	newExecutionClient DestinationNewExecutionClientFunc
	// logger that should be used by the plugin
	logger zerolog.Logger
	// client returned by call to newExecutionClient
	client DestinationClient
}

func (p *DestinationPlugin) Name() string {
	return p.name
}

func (p *DestinationPlugin) Version() string {
	return p.version
}

func (p *DestinationPlugin) ExampleConfig() string {
	return p.exampleConfig
}

func (p *DestinationPlugin) Migrate(ctx context.Context, tables schema.Tables) error {
	return p.client.Migrate(ctx, tables)
}

func (p *DestinationPlugin) Write(ctx context.Context, table string, data map[string]interface{}) error {
	return p.client.Write(ctx, table, data)
}

func (p *DestinationPlugin) SetLogger(logger zerolog.Logger) {
	p.logger = logger
	p.client.SetLogger(logger)
}

func NewDestinationPlugin(name string, version string, newExecutionClient DestinationNewExecutionClientFunc, opts ...DestinationOption) *DestinationPlugin {
	p := DestinationPlugin{
		name:               name,
		version:            version,
		exampleConfig:      "",
		logger:             zerolog.Logger{},
		newExecutionClient: newExecutionClient,
	}
	for _, opt := range opts {
		opt(&p)
	}
	if err := p.validate(); err != nil {
		panic(err)
	}
	return &p
}

func (p *DestinationPlugin) validate() error {
	if p.newExecutionClient == nil {
		return fmt.Errorf("newExecutionClient function not defined for source plugin: " + p.name)
	}

	if p.name == "" {
		return errors.New("plugin name should not be empty")
	}

	return nil
}
