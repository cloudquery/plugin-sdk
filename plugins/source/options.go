package source

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
)

type GetTables func(ctx context.Context, c schema.ClientMeta) (schema.Tables, error)

type Option func(*Plugin)

// WithDynamicTableOption allows the plugin to return list of tables after call to New
func WithDynamicTableOption(getDynamicTables GetTables) Option {
	return func(p *Plugin) {
		p.getDynamicTables = getDynamicTables
	}
}
