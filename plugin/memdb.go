package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

// client is mostly used for testing the destination plugin.
type client struct {
	memoryDB      map[string][]arrow.Record
	tables        map[string]*schema.Table
	memoryDBLock  sync.RWMutex
	errOnWrite    bool
	blockingWrite bool
}

type MemDBOption func(*client)

func WithErrOnWrite() MemDBOption {
	return func(c *client) {
		c.errOnWrite = true
	}
}

func WithBlockingWrite() MemDBOption {
	return func(c *client) {
		c.blockingWrite = true
	}
}

func GetNewClient(options ...MemDBOption) NewClientFunc {
	c := &client{
		memoryDB:     make(map[string][]arrow.Record),
		memoryDBLock: sync.RWMutex{},
	}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, any) (Client, error) {
		return c, nil
	}
}

func NewMemDBClient(_ context.Context, _ zerolog.Logger, spec any) (Client, error) {
	return &client{
		memoryDB: make(map[string][]arrow.Record),
		tables:   make(map[string]*schema.Table),
	}, nil
}

func NewMemDBClientErrOnNew(context.Context, zerolog.Logger, []byte) (Client, error) {
	return nil, fmt.Errorf("newTestDestinationMemDBClientErrOnNew")
}

func (c *client) overwrite(table *schema.Table, data arrow.Record) {
	pksIndex := table.PrimaryKeysIndexes()
	tableName := table.Name
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

func (c *client) ID() string {
	return "testDestinationMemDB"
}

func (c *client) NewManagedSyncClient(context.Context, SyncOptions) (ManagedSyncClient, error) {
	return nil, fmt.Errorf("not supported")
}

func (c *client) Sync(ctx context.Context, options SyncOptions, res chan<- arrow.Record) error {
	c.memoryDBLock.RLock()
	for tableName := range c.memoryDB {
		for _, row := range c.memoryDB[tableName] {
			res <- row
		}
	}
	c.memoryDBLock.RUnlock()
	return nil
}

func (c *client) Migrate(_ context.Context, tables schema.Tables, migrateMode MigrateMode) error {
	for _, table := range tables {
		tableName := table.Name
		memTable := c.memoryDB[tableName]
		if memTable == nil {
			c.memoryDB[tableName] = make([]arrow.Record, 0)
			c.tables[tableName] = table
			continue
		}

		changes := table.GetChanges(c.tables[tableName])
		// memdb doesn't support any auto-migrate
		if changes == nil {
			continue
		}
		c.memoryDB[tableName] = make([]arrow.Record, 0)
		c.tables[tableName] = table
	}
	return nil
}

func (c *client) Read(_ context.Context, table *schema.Table, source string, res chan<- arrow.Record) error {
	tableName := table.Name
	if c.memoryDB[tableName] == nil {
		return nil
	}
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	if sourceColIndex == -1 {
		return fmt.Errorf("table %s doesn't have source column", tableName)
	}
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

func (c *client) Write(ctx context.Context, _ schema.Tables, writeMode WriteMode, resources <-chan arrow.Record) error {
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
		sc := resource.Schema()
		tableName, ok := sc.Metadata().GetValue(schema.MetadataTableName)
		if !ok {
			return fmt.Errorf("table name not found in schema metadata")
		}
		table := c.tables[tableName]
		if writeMode == WriteModeAppend {
			c.memoryDB[tableName] = append(c.memoryDB[tableName], resource)
		} else {
			c.overwrite(table, resource)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *client) WriteTableBatch(ctx context.Context, table *schema.Table, writeMode WriteMode, resources []arrow.Record) error {
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
	tableName := table.Name
	for _, resource := range resources {
		c.memoryDBLock.Lock()
		if writeMode == WriteModeAppend {
			c.memoryDB[tableName] = append(c.memoryDB[tableName], resource)
		} else {
			c.overwrite(table, resource)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (*client) Metrics() Metrics {
	return Metrics{}
}

func (c *client) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *client) DeleteStale(ctx context.Context, tables schema.Tables, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
	}
	return nil
}

func (c *client) deleteStaleTable(_ context.Context, table *schema.Table, source string, syncTime time.Time) {
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	syncColIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	tableName := table.Name
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
