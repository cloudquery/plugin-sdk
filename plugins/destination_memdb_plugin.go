package plugins

import (
	"context"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

// TestDestinationMemDBClient is mostly used for testing.
type TestDestinationMemDBClient struct {
	schema.DefaultTransformer
	DefaultReverseTransformer
	spec     specs.Destination
	memoryDB map[string][][]interface{}
}

func NewTestDestinationMemDBClient(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
	return &TestDestinationMemDBClient{
		memoryDB: make(map[string][][]interface{}),
	}, nil
}

func (c *TestDestinationMemDBClient) overwrite(table *schema.Table, data []interface{}) {
	pks := table.PrimaryKeys()
	//nolint:prealloc
	var pksIndex []int
	for _, pk := range pks {
		pksIndex = append(pksIndex, table.Columns.Index(pk))
	}
	for i, row := range c.memoryDB[table.Name] {
		found := true
		for _, pkIndex := range pksIndex {
			if row[pkIndex] != data[pkIndex] {
				found = false
			}
		}
		if found {
			c.memoryDB[table.Name][i] = data
			return
		}
	}
	c.memoryDB[table.Name] = append(c.memoryDB[table.Name], data)
}

func (c *TestDestinationMemDBClient) Initialize(_ context.Context, spec specs.Destination) error {
	c.spec = spec
	return nil
}
func (c *TestDestinationMemDBClient) Migrate(_ context.Context, tables schema.Tables) error {
	for _, table := range tables {
		if c.memoryDB[table.Name] == nil {
			c.memoryDB[table.Name] = make([][]interface{}, 0)
		}
	}
	return nil
}

func (c *TestDestinationMemDBClient) Read(_ context.Context, table *schema.Table, source string, res chan<- []interface{}) error {
	if c.memoryDB[table.Name] == nil {
		return nil
	}
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	for _, row := range c.memoryDB[table.Name] {
		if row[sourceColIndex].(*schema.Text).Str == source {
			res <- row
		}
	}

	return nil
}

func (c *TestDestinationMemDBClient) Write(_ context.Context, tables schema.Tables, resources <-chan *ClientResource) error {
	for resource := range resources {
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[resource.TableName] = append(c.memoryDB[resource.TableName], resource.Data)
		} else {
			c.overwrite(tables.Get(resource.TableName), resource.Data)
		}
	}
	return nil
}

func (*TestDestinationMemDBClient) Metrics() DestinationMetrics {
	return DestinationMetrics{}
}

func (c *TestDestinationMemDBClient) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *TestDestinationMemDBClient) DeleteStale(ctx context.Context, tables schema.Tables, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
		if err := c.DeleteStale(ctx, table.Relations, source, syncTime); err != nil {
			return err
		}
	}
	return nil
}

func (c *TestDestinationMemDBClient) deleteStaleTable(_ context.Context, table *schema.Table, source string, syncTime time.Time) {
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	syncColIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	for i, row := range c.memoryDB[table.Name] {
		if row[sourceColIndex].(*schema.Text).Str == source {
			rowSyncTime := row[syncColIndex].(*schema.Timestamptz)
			if rowSyncTime.Time.Before(syncTime) {
				c.memoryDB[table.Name] = append(c.memoryDB[table.Name][:i], c.memoryDB[table.Name][i+1:]...)
			}
		}
	}
}
