package plugins

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

// type DestinationOption func(*DestinationPlugin)

// type WriteFunc func(ctx context.Context, spec specs.DestinationSpec, tables []*schema.Table, resources <-chan *schema.Resource) error

// type DestinationPlugin struct {
// 	// Name of the plugin.
// 	name string
// 	// Version is the version of the plugin.
// 	version string
// 	// JsonSchema for specific source plugin spec
// 	jsonSchema string
// 	// ExampleConfig is the example configuration for this plugin
// 	exampleConfig string
// 	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
// 	logger zerolog.Logger
// 	// Write is a function that get a stream of resources and write them to the configured destination
// 	// with the configured mode.
// 	Write WriteFunc
// }

type DestinationPlugin interface {
	Name() string
	Version() string
	ExampleConfig() string
	Initialize(ctx context.Context, spec specs.Destination) error
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, resources *schema.Resource) error
	SetLogger(logger zerolog.Logger)
}

// func WithExampleConfig(exampleConfig string) DestinationOption {
// 	return func(p *DestinationPlugin) {
// 		p.exampleConfig = exampleConfig
// 	}
// }

// func WithJsonSchema(jsonSchema string) DestinationOption {
// 	return func(p *DestinationPlugin) {
// 		p.jsonSchema = jsonSchema
// 	}
// }

// func WithLogger(logger zerolog.Logger) DestinationOption {
// 	return func(p *DestinationPlugin) {
// 		p.logger = logger
// 	}
// }

// func NewDestinationClient(name string, version string, writeFunc WriteFunc, opts ...DestinationOption) *DestinationPlugin {
// 	p := DestinationPlugin{
// 		name:    name,
// 		version: version,
// 	}
// 	for _, opt := range opts {
// 		opt(&p)
// 	}
// 	return &p
// }

// func (p *DestinationPlugin) GetExampleConfig() string {
// 	return p.exampleConfig
// }
