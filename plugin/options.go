package plugin

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v3/schema"
)

type GetTables func(ctx context.Context, c Client) (schema.Tables, error)

type Option func(*Plugin)

// WithDynamicTableOption allows the plugin to return list of tables after call to New
func WithDynamicTableOption(getDynamicTables GetTables) Option {
	return func(p *Plugin) {
		p.getDynamicTables = getDynamicTables
	}
}

// WithNoInternalColumns won't add internal columns (_cq_id, _cq_parent_cq_id) to the plugin tables
func WithNoInternalColumns() Option {
	return func(p *Plugin) {
		p.internalColumns = false
	}
}

func WithUnmanaged() Option {
	return func(p *Plugin) {
		p.unmanaged = true
	}
}

// WithTitleTransformer allows the plugin to control how table names get turned into titles for the
// generated documentation.
func WithTitleTransformer(t func(*schema.Table) string) Option {
	return func(p *Plugin) {
		p.titleTransformer = t
	}
}


func WithStaticTables(tables schema.Tables) Option {
	return func(p *Plugin) {
		p.staticTables = tables
	}
}