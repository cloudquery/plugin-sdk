package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

// TestDestinationMemDBClient is mostly used for testing.
type TestDestinationMemDBClient struct {
	schema.DefaultTransformer
	spec          specs.Destination
	memoryDB      map[string][][]interface{}
	errOnWrite    bool
	blockingWrite bool
}

type TestDestinationOption func(*TestDestinationMemDBClient)

func withErrOnWrite() TestDestinationOption {
	return func(c *TestDestinationMemDBClient) {
		c.errOnWrite = true
	}
}

func withBlockingWrite() TestDestinationOption {
	return func(c *TestDestinationMemDBClient) {
		c.blockingWrite = true
	}
}

func getNewTestDestinationMemDBClient(options ...TestDestinationOption) NewDestinationClientFunc {
	c := &TestDestinationMemDBClient{
		memoryDB: make(map[string][][]interface{}),
	}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
		return c, nil
	}
}

func NewTestDestinationMemDBClient(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
	return &TestDestinationMemDBClient{
		memoryDB: make(map[string][][]interface{}),
	}, nil
}

func newTestDestinationMemDBClientErrOnNew(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
	return nil, fmt.Errorf("newTestDestinationMemDBClientErrOnNew")
}

func (*TestDestinationMemDBClient) ReverseTransformValues(_ *schema.Table, values []interface{}) (schema.CQTypes, error) {
	res := make(schema.CQTypes, len(values))
	for i, v := range values {
		res[i] = v.(schema.CQType)
	}
	return res, nil
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

func (c *TestDestinationMemDBClient) PreWrite(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) (interface{}, error) {
	return nil, nil
}

func (c *TestDestinationMemDBClient) PostWrite(ctx context.Context, writeClient interface{}, tables schema.Tables, sourceName string, syncTime time.Time) error {
	return nil
}

func (c *TestDestinationMemDBClient) WriteTableBatch(ctx context.Context, _ interface{}, table *schema.Table, resources [][]interface{}) error {
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
	if c.spec.WriteMode == specs.WriteModeAppend {
		c.memoryDB[table.Name] = append(c.memoryDB[table.Name], resources...)
	} else {
		for _, resource := range resources {
			c.overwrite(table, resource)
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
