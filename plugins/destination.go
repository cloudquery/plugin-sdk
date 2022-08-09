package plugins

import (
	"context"

	"github.com/cloudquery/cq-plugin-sdk/schema"
	"github.com/cloudquery/cq-plugin-sdk/spec"
	"github.com/rs/zerolog"
)

type DestinationPluginOptions struct {
	Logger zerolog.Logger
}

type DestinationPlugin interface {
	Configure(ctx context.Context, spec spec.DestinationSpec) error
	CreateTables(ctx context.Context, table []*schema.Table) error
	Save(ctx context.Context, resources []*schema.Resource) error
	GetExampleConfig(ctx context.Context) string
}
