package plugin

import (
	"context"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type GetTables func(ctx context.Context, c Client) (schema.Tables, error)

type Option func(*Plugin)

// WithDynamicTable allows the plugin to return list of tables after call to New
func WithDynamicTable(getDynamicTables GetTables) Option {
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

func WithManagedWriter() Option {
	return func(p *Plugin) {
		p.managedWriter = true
	}
}

func WithBatchTimeout(seconds int) Option {
	return func(p *Plugin) {
		p.batchTimeout = time.Duration(seconds) * time.Second
	}
}

func WithDefaultBatchSize(defaultBatchSize int) Option {
	return func(p *Plugin) {
		p.defaultBatchSize = defaultBatchSize
	}
}

func WithDefaultBatchSizeBytes(defaultBatchSizeBytes int) Option {
	return func(p *Plugin) {
		p.defaultBatchSizeBytes = defaultBatchSizeBytes
	}
}
