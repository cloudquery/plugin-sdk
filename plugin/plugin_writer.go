package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type WriteOptions struct {
	// WriteMode is the mode to write to the database
	WriteMode WriteMode
	// Predefined tables are available if tables are known at the start of the write
	Tables schema.Tables
}

type MigrateOptions struct {
	// MigrateMode is the mode to migrate the database
	MigrateMode MigrateMode
}

func (p *Plugin) Migrate(ctx context.Context, tables schema.Tables, options MigrateOptions) error {
	if p.client == nil {
		return fmt.Errorf("plugin is not initialized")
	}
	return p.client.Migrate(ctx, tables, options)
}

// this function is currently used mostly for testing so it's not a public api
func (p *Plugin) writeOne(ctx context.Context, options WriteOptions, resource arrow.Record) error {
	resources := []arrow.Record{resource}
	return p.writeAll(ctx, options, resources)
}

// this function is currently used mostly for testing so it's not a public api
func (p *Plugin) writeAll(ctx context.Context, options WriteOptions, resources []arrow.Record) error {
	ch := make(chan arrow.Record, len(resources))
	for _, resource := range resources {
		ch <- resource
	}
	close(ch)
	tables := make(schema.Tables, 0)
	tableNames := make(map[string]struct{})
	for _, resource := range resources {
		sc := resource.Schema()
		tableMD := sc.Metadata()
		name, found := tableMD.GetValue(schema.MetadataTableName)
		if !found {
			return fmt.Errorf("missing table name")
		}
		if _, ok := tableNames[name]; ok {
			continue
		}
		table, err := schema.NewTableFromArrowSchema(resource.Schema())
		if err != nil {
			return err
		}
		tables = append(tables, table)
		tableNames[table.Name] = struct{}{}
	}
	options.Tables = tables
	return p.Write(ctx, options, ch)
}

func (p *Plugin) Write(ctx context.Context, options WriteOptions, res <-chan arrow.Record) error {
	if err := p.client.Write(ctx, options, res); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	syncTime = syncTime.UTC()
	return p.client.DeleteStale(ctx, tables, sourceName, syncTime)
}
