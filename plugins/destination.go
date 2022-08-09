package plugins

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type DestinationPluginOptions struct {
	Logger zerolog.Logger
}

type DestinationPlugin interface {
	Configure(ctx context.Context, spec specs.DestinationSpec) error
	CreateTables(ctx context.Context, table []*schema.Table) error
	Save(ctx context.Context, resources []*schema.Resource) error
	GetExampleConfig(ctx context.Context) string
}
