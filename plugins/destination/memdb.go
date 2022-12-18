package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

// memDBClient is mostly used for testing the destination plugin.
type memDBClient struct {
	schema.DefaultTransformer
	spec          specs.Destination
	memoryDB      map[string][][]interface{}
	errOnWrite    bool
	blockingWrite bool
}

type memDBClientOption func(*memDBClient)

func withMemDBErrOnWrite() memDBClientOption {
	return func(c *memDBClient) {
		c.errOnWrite = true
	}
}

func withManagedMemDBBlockingWrite() memDBClientOption {
	return func(c *memDBClient) {
		c.blockingWrite = true
	}
}

func getNewManagedMemDBClient(options ...memDBClientOption) NewClientFunc {
	c := &memDBClient{
		memoryDB: make(map[string][][]interface{}),
	}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, specs.Destination) (Client, error) {
		return c, nil
	}
}

func newMemDBClient(_ context.Context, _ zerolog.Logger, spec specs.Destination) (Client, error) {
	return &memDBClient{
		memoryDB: make(map[string][][]interface{}),
		spec: spec,
	}, nil
}

func newManagedMemDBClientErrOnNew(context.Context, zerolog.Logger, specs.Destination) (Client, error) {
	return nil, fmt.Errorf("newTestDestinationMemDBClientErrOnNew")
}

func (*memDBClient) ReverseTransformValues(_ *schema.Table, values []interface{}) (schema.CQTypes, error) {
	res := make(schema.CQTypes, len(values))
	for i, v := range values {
		res[i] = v.(schema.CQType)
	}
	return res, nil
}
func (c *memDBClient) overwrite(table *schema.Table, data []interface{}) {
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


func (c *memDBClient) Migrate(_ context.Context, tables schema.Tables) error {
	for _, table := range tables {
		if c.memoryDB[table.Name] == nil {
			c.memoryDB[table.Name] = make([][]interface{}, 0)
		}
	}
	return nil
}

func (c *memDBClient) Read(_ context.Context, table *schema.Table, source string, res chan<- []interface{}) error {
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

func (c *memDBClient) Write(ctx context.Context, tables schema.Tables, resources <-chan *ClientResource) error {
	if c.errOnWrite {
		return fmt.Errorf("errOnWrite")
	}
	if c.blockingWrite {
		<-ctx.Done()
		if c.errOnWrite {
			return fmt.Errorf("errOnWrite")
		}
		return nil
	}
	for resource := range resources {
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[resource.TableName] = append(c.memoryDB[resource.TableName], resource.Data)
		} else {
			c.overwrite(tables.Get(resource.TableName), resource.Data)
		}
	}
	return nil
}

func (c *memDBClient) WriteTableBatch(ctx context.Context, table *schema.Table, resources [][]interface{}) error {
	if c.errOnWrite {
		return fmt.Errorf("errOnWrite")
	}
	if c.blockingWrite {
		<-ctx.Done()
		if c.errOnWrite {
			return fmt.Errorf("errOnWrite")
		}
		return nil
	}
	for _, resource := range resources {
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[table.Name] = append(c.memoryDB[table.Name], resource)
		} else {
			c.overwrite(table, resource)
		}
	}
	return nil
}

func (*memDBClient) Metrics() Metrics {
	return Metrics{}
}

func (c *memDBClient) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *memDBClient) DeleteStale(ctx context.Context, tables schema.Tables, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
		if err := c.DeleteStale(ctx, table.Relations, source, syncTime); err != nil {
			return err
		}
	}
	return nil
}

func (c *memDBClient) deleteStaleTable(_ context.Context, table *schema.Table, source string, syncTime time.Time) {
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	syncColIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	var filteredTable [][]interface{}
	for i, row := range c.memoryDB[table.Name] {
		if row[sourceColIndex].(*schema.Text).Str == source {
			rowSyncTime := row[syncColIndex].(*schema.Timestamptz)
			if !rowSyncTime.Time.UTC().Before(syncTime) {
				filteredTable = append(filteredTable, c.memoryDB[table.Name][i])
			}
		}
	}
	c.memoryDB[table.Name] = filteredTable
}
