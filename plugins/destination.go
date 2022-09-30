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
	Initialize(ctx context.Context, spec specs.Destination) error
	Migrate(ctx context.Context, tables schema.Tables) error
	WriteRow(ctx context.Context, table string, data map[string]interface{}) error
	Close(ctx context.Context) error
	SetLogger(logger zerolog.Logger)
}
