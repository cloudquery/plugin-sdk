package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (p *Plugin) Migrate(ctx context.Context, tables schema.Tables, migrateMode MigrateMode) error {
	if p.client == nil {
		return fmt.Errorf("plugin is not initialized")
	}
	return p.client.Migrate(ctx, tables, migrateMode)
}

// this function is currently used mostly for testing so it's not a public api
func (p *Plugin) writeOne(ctx context.Context, sourceName string, syncTime time.Time, writeMode WriteMode, resource arrow.Record) error {
	resources := []arrow.Record{resource}
	return p.writeAll(ctx, sourceName, syncTime, writeMode, resources)
}

// this function is currently used mostly for testing so it's not a public api
func (p *Plugin) writeAll(ctx context.Context, sourceName string, syncTime time.Time, writeMode WriteMode, resources []arrow.Record) error {
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
	return p.Write(ctx, sourceName, tables, syncTime, writeMode, ch)
}

func (p *Plugin) Write(ctx context.Context, sourceName string, tables schema.Tables, syncTime time.Time, writeMode WriteMode, res <-chan arrow.Record) error {
	syncTime = syncTime.UTC()
	if p.managedWriter {
		if err := p.writeManagedTableBatch(ctx, tables, writeMode, res); err != nil {
			return err
		}
	} else {
		if err := p.client.Write(ctx, tables, writeMode, res); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	syncTime = syncTime.UTC()
	return p.client.DeleteStale(ctx, tables, sourceName, syncTime)
}
