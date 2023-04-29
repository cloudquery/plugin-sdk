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
	"github.com/cloudquery/plugin-sdk/v2/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/schemav2"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/rs/zerolog"
)

// client is mostly used for testing the destination plugin.
type client struct {
	spec          specs.Destination
	memoryDB      map[string][]arrow.Record
	tables        map[string]*schemav2.Table
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
		tables:   make(map[string]*schemav2.Table),
		spec:     spec,
	}, nil
}

func NewClientErrOnNew(context.Context, zerolog.Logger, specs.Destination) (destination.Client, error) {
	return nil, fmt.Errorf("newTestDestinationMemDBClientErrOnNew")
}

func (c *client) overwrite(table *schemav2.Table, data arrow.Record) {
	pksIndexes := table.PrimaryKeysIndexes()
	for i, row := range c.memoryDB[table.Name] {
		found := true
		for _, pkIndex := range pksIndexes {
			s1 := data.Column(pkIndex).String()
			s2 := row.Column(pkIndex).String()
			if s1 != s2 {
				found = false
			}
		}
		if found {
			c.memoryDB[table.Name] = append(c.memoryDB[table.Name][:i], c.memoryDB[table.Name][i+1:]...)
			c.memoryDB[table.Name] = append(c.memoryDB[table.Name], data)
			return
		}
	}
	c.memoryDB[table.Name] = append(c.memoryDB[table.Name], data)
}

func (c *client) Migrate(_ context.Context, tables schemav2.Tables) error {
	for _, table := range tables {
		memTable := c.memoryDB[table.Name]
		if memTable == nil {
			c.memoryDB[table.Name] = make([]arrow.Record, 0)
			c.tables[table.Name] = table
			continue
		}
		changes := table.GetChanges(c.tables[table.Name])
		// memdb doesn't support any auto-migrate
		if len(changes) == 0 {
			continue
		}
		c.memoryDB[table.Name] = make([]arrow.Record, 0)
		c.tables[table.Name] = table
	}
	return nil
}

func (c *client) Read(_ context.Context, table *schemav2.Table, source string, res chan<- arrow.Record) error {
	if c.memoryDB[table.Name] == nil {
		return nil
	}
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	if sourceColIndex == -1 {
		return fmt.Errorf("table %s doesn't have source column", table.Name)
	}
	var sortedRes []arrow.Record
	c.memoryDBLock.RLock()
	for _, row := range c.memoryDB[table.Name] {
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

func (c *client) Write(ctx context.Context, _ schemav2.Tables, resources <-chan arrow.Record) error {
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
		table := c.tables[tableName]
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[tableName] = append(c.memoryDB[tableName], resource)
		} else {
			c.overwrite(table, resource)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *client) WriteTableBatch(ctx context.Context, table *schemav2.Table, resources []arrow.Record) error {
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
		c.memoryDBLock.Lock()
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[table.Name] = append(c.memoryDB[table.Name], resource)
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

func (c *client) DeleteStale(ctx context.Context, tables schemav2.Tables, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
	}
	return nil
}

func (c *client) deleteStaleTable(_ context.Context, table *schemav2.Table, source string, syncTime time.Time) {
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	syncColIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	var filteredTable []arrow.Record
	for i, row := range c.memoryDB[table.Name] {
		if row.Column(sourceColIndex).(*array.String).Value(0) == source {
			rowSyncTime := row.Column(syncColIndex).(*array.Timestamp).Value(0).ToTime(arrow.Microsecond).UTC()
			if !rowSyncTime.Before(syncTime) {
				filteredTable = append(filteredTable, c.memoryDB[table.Name][i])
			}
		}
	}
	c.memoryDB[table.Name] = filteredTable
}
