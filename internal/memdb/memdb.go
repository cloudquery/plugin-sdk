package memdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
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

type Option func(*client)

type Spec struct {
}

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

func GetNewClient(options ...Option) plugin.NewClientFunc {
	c := &client{
		memoryDB:     make(map[string][]arrow.Record),
		memoryDBLock: sync.RWMutex{},
		tables:       make(map[string]*schema.Table),
	}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return c, nil
	}
}

func NewMemDBClient(ctx context.Context, l zerolog.Logger, spec []byte, options plugin.NewClientOptions) (plugin.Client, error) {
	return GetNewClient()(ctx, l, spec, options)
}

func NewMemDBClientErrOnNew(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
	return nil, fmt.Errorf("newTestDestinationMemDBClientErrOnNew")
}

func (c *client) overwrite(table *schema.Table, record arrow.Record) {
	for i := int64(0); i < record.NumRows(); i++ {
		c.overwriteRow(table, record.NewSlice(i, i+1))
	}
}

func (c *client) overwriteRow(table *schema.Table, data arrow.Record) {
	tableName := table.Name
	pksIndex := table.PrimaryKeysIndexes()
	if len(pksIndex) == 0 {
		c.memoryDB[tableName] = append(c.memoryDB[tableName], data)
		return
	}

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

func (*client) ID() string {
	return "testDestinationMemDB"
}

func (*client) GetSpec() any {
	return &Spec{}
}

func (c *client) Read(_ context.Context, table *schema.Table, res chan<- arrow.Record) error {
	c.memoryDBLock.RLock()
	defer c.memoryDBLock.RUnlock()

	tableName := table.Name
	// we iterate over records in reverse here because we don't set an expectation
	// of ordering on plugins, and we want to make sure that the tests are not
	// dependent on the order of insertion either.
	rows := c.memoryDB[tableName]
	for i := len(rows) - 1; i >= 0; i-- {
		res <- rows[i]
	}
	return nil
}

func (c *client) Sync(_ context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	c.memoryDBLock.RLock()

	for tableName := range c.memoryDB {
		if !plugin.MatchesTable(tableName, options.Tables, options.SkipTables) {
			continue
		}
		for _, row := range c.memoryDB[tableName] {
			res <- &message.SyncInsert{
				Record: row,
			}
		}
	}
	c.memoryDBLock.RUnlock()
	return nil
}

func (c *client) Tables(context.Context, plugin.TableOptions) (schema.Tables, error) {
	tables := make(schema.Tables, 0, len(c.tables))
	for _, table := range c.tables {
		tables = append(tables, table)
	}
	return tables, nil
}

func (c *client) migrate(_ context.Context, table *schema.Table) {
	tableName := table.Name
	memTable := c.memoryDB[tableName]
	if memTable == nil {
		c.memoryDB[tableName] = make([]arrow.Record, 0)
		c.tables[tableName] = table
		return
	}

	changes := table.GetChanges(c.tables[tableName])
	// memdb doesn't support any auto-migrate
	if changes == nil {
		return
	}
	c.memoryDB[tableName] = make([]arrow.Record, 0)
	c.tables[tableName] = table
}

func (c *client) Write(ctx context.Context, msgs <-chan message.WriteMessage) error {
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

	for msg := range msgs {
		c.memoryDBLock.Lock()

		switch msg := msg.(type) {
		case *message.WriteMigrateTable:
			c.migrate(ctx, msg.Table)
		case *message.WriteDeleteStale:
			c.deleteStale(ctx, msg)
		case *message.WriteInsert:
			sc := msg.Record.Schema()
			tableName, ok := sc.Metadata().GetValue(schema.MetadataTableName)
			if !ok {
				return fmt.Errorf("table name not found in schema metadata")
			}
			table := c.tables[tableName]
			c.overwrite(table, msg.Record)
		}

		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *client) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *client) deleteStale(_ context.Context, msg *message.WriteDeleteStale) {
	var filteredTable []arrow.Record
	tableName := msg.TableName
	for i, row := range c.memoryDB[tableName] {
		sc := row.Schema()
		indices := sc.FieldIndices(schema.CqSourceNameColumn.Name)
		if len(indices) == 0 {
			continue
		}
		sourceColIndex := indices[0]
		indices = sc.FieldIndices(schema.CqSyncTimeColumn.Name)
		if len(indices) == 0 {
			continue
		}
		syncColIndex := indices[0]

		if row.Column(sourceColIndex).(*array.String).Value(0) == msg.SourceName {
			rowSyncTime := row.Column(syncColIndex).(*array.Timestamp).Value(0).ToTime(arrow.Microsecond).UTC()
			if !rowSyncTime.Before(msg.SyncTime) {
				filteredTable = append(filteredTable, c.memoryDB[tableName][i])
			}
		}
	}
	c.memoryDB[tableName] = filteredTable
}
