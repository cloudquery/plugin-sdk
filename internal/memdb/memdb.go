package memdb

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/cloudquery/plugin-sdk/plugins/destination"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

// client is mostly used for testing the destination plugin.
type client struct {
	schema.DefaultTransformer
	spec          specs.Destination
	memoryDB      map[string][]arrow.Record
	tables        map[string]*arrow.Schema
	memoryDBLock  sync.RWMutex
	errOnWrite    bool
	blockingWrite bool
}

type Option func(*client)

func WithErrOnWrite() Option {
	return func(c *client) {
		c.errOnWrite = true
	}
}

func WithBlockingWrite() Option {
	return func(c *client) {
		c.blockingWrite = true
	}
}

func GetNewClient(options ...Option) destination.NewClientFunc {
	c := &client{
		memoryDB:     make(map[string][]arrow.Record),
		memoryDBLock: sync.RWMutex{},
	}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, specs.Destination) (destination.Client, error) {
		return c, nil
	}
}

func getTestLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func NewClient(_ context.Context, _ zerolog.Logger, spec specs.Destination) (destination.Client, error) {
	return &client{
		memoryDB: make(map[string][]arrow.Record),
		tables:   make(map[string]*arrow.Schema),
		spec:     spec,
	}, nil
}

func NewClientErrOnNew(context.Context, zerolog.Logger, specs.Destination) (destination.Client, error) {
	return nil, fmt.Errorf("newTestDestinationMemDBClientErrOnNew")
}

func (*client) ReverseTransformValues(_ *schema.Table, values []any) (schema.CQTypes, error) {
	res := make(schema.CQTypes, len(values))
	for i, v := range values {
		res[i] = v.(schema.CQType)
	}
	return res, nil
}

func (c *client) overwrite(table *arrow.Schema, data arrow.Record) {
	pksIndex := schema.PrimaryKeyIndices(table)
	tableName := schema.TableName(table)
	for i, row := range c.memoryDB[tableName] {
		found := true
		for _, pkIndex := range pksIndex {
			s1 := data.Column(pkIndex).String()
			s2 := row.Column(pkIndex).String()
			if s1 != s2 {
				found = false
			}
		}
		if found {
			c.memoryDB[tableName] = append(c.memoryDB[tableName][:i], c.memoryDB[tableName][i+1:]...)
			c.memoryDB[tableName] = append(c.memoryDB[tableName], data)
			return
		}
	}
	c.memoryDB[tableName] = append(c.memoryDB[tableName], data)
}

func (c *client) Migrate(_ context.Context, tables schema.Schemas) error {
	for _, table := range tables {
		tableName := schema.TableName(table)
		memTable := c.memoryDB[tableName]
		if memTable == nil {
			c.memoryDB[tableName] = make([]arrow.Record, 0)
			c.tables[tableName] = table
			continue
		}
		changes := schema.GetSchemaChanges(table, c.tables[tableName])
		// memdb doesn't support any auto-migrate
		if changes == nil {
			continue
		}
		c.memoryDB[tableName] = make([]arrow.Record, 0)
		c.tables[tableName] = table
	}
	return nil
}

func (c *client) Read(_ context.Context, table *arrow.Schema, source string, res chan<- arrow.Record) error {
	tableName := schema.TableName(table)
	if c.memoryDB[tableName] == nil {
		return nil
	}
	indices := table.FieldIndices(schema.CqSourceNameColumn.Name)
	if len(indices) == 0 {
		return fmt.Errorf("table %s doesn't have source column", tableName)
	}
	sourceColIndex := indices[0]
	var sortedRes []arrow.Record
	c.memoryDBLock.RLock()
	for _, row := range c.memoryDB[tableName] {
		arr := row.Column(sourceColIndex)
		if arr.(*array.String).Value(0) == source {
			sortedRes = append(sortedRes, row)
		}
	}
	c.memoryDBLock.RUnlock()

	for _, row := range sortedRes {
		res <- row
	}
	return nil
}

func (c *client) Write(ctx context.Context, tables schema.Schemas, resources <-chan arrow.Record) error {
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
		c.memoryDBLock.Lock()
		tableName, err := schema.TableNameFromSchema(resource.Schema())
		if err != nil {
			return err
		}
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[tableName] = append(c.memoryDB[tableName], resource)
		} else {
			c.overwrite(resource.Schema(), resource)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *client) WriteTableBatch(ctx context.Context, table *arrow.Schema, resources []arrow.Record) error {
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
	tableName := schema.TableName(table)
	for _, resource := range resources {
		c.memoryDBLock.Lock()
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[tableName] = append(c.memoryDB[tableName], resource)
		} else {
			c.overwrite(table, resource)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (*client) Metrics() destination.Metrics {
	return destination.Metrics{}
}

func (c *client) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *client) DeleteStale(ctx context.Context, tables schema.Schemas, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
	}
	return nil
}

func (c *client) deleteStaleTable(_ context.Context, table *arrow.Schema, source string, syncTime time.Time) {
	sourceColIndex := table.FieldIndices(schema.CqSourceNameColumn.Name)[0]
	syncColIndex := table.FieldIndices(schema.CqSyncTimeColumn.Name)[0]
	tableName := schema.TableName(table)
	var filteredTable []arrow.Record
	for i, row := range c.memoryDB[tableName] {
		if row.Column(sourceColIndex).(*array.String).Value(0) == source {
			rowSyncTime := row.Column(syncColIndex).(*array.Timestamp).Value(0).ToTime(arrow.Microsecond).UTC()
			if !rowSyncTime.Before(syncTime) {
				filteredTable = append(filteredTable, c.memoryDB[tableName][i])
			}
		}
	}
	c.memoryDB[tableName] = filteredTable
}
