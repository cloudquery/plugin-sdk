package plugins

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type DestinationPlugin interface {
	Name() string
	Version() string
	ExampleConfig() string
	Initialize(ctx context.Context, spec specs.Destination) error
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, table string, data map[string]interface{}) error
	SetLogger(logger zerolog.Logger)
}
