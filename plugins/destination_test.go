package plugins

import (
	"context"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)


type testDestinationClient struct {
	schema.DefaultTransformer
	DefaultReverseTransformer
	spec  specs.Destination
	memoryDB map[string][][]interface{}
}

func newTestDestinationClient(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
	return &testDestinationClient{
		memoryDB: make(map[string][][]interface{}),
	}, nil
}

func (c *testDestinationClient) overwrite(table *schema.Table, data []interface{}) error {
	pks := table.PrimaryKeys()
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
			return nil
		}
	}
	c.memoryDB[table.Name] = append(c.memoryDB[table.Name], data)
	return nil
}

func (c *testDestinationClient) Initialize(_ context.Context, spec specs.Destination) error {
	c.spec = spec
	return nil
}
func (c *testDestinationClient) Migrate(ctx context.Context, tables schema.Tables) error {
	for _, table := range tables {
		if c.memoryDB[table.Name] == nil {
			c.memoryDB[table.Name] = make([][]interface{}, 0)
		}
	}
	return nil
}

func (c *testDestinationClient) Read(ctx context.Context, table *schema.Table, source string, res chan<- []interface{}) error {
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

func (c *testDestinationClient) Write(_ context.Context, tables schema.Tables, resources <-chan *ClientResource) error {
	for resource := range resources {
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[resource.TableName] = append(c.memoryDB[resource.TableName], resource.Data)
		} else {
			c.overwrite(tables.Get(resource.TableName), resource.Data)
		}
	}
	return nil
}

func (*testDestinationClient) Metrics() DestinationMetrics {
	return DestinationMetrics{}
}

func (c *testDestinationClient) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *testDestinationClient) DeleteStale(ctx context.Context,tables schema.Tables, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
		c.DeleteStale(ctx, table.Relations, source, syncTime)
	}
	return nil
}

func (c *testDestinationClient) deleteStaleTable(_ context.Context,table *schema.Table, source string, syncTime time.Time) {
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

func TestDestinationPlugin(t *testing.T) {
	p := NewDestinationPlugin("test", "development", newTestDestinationClient)
	DestinationPluginTestSuite(t, p, specs.Destination{
		WriteMode: specs.WriteModeAppend,
	})
}