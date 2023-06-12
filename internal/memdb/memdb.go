package memdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
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
	}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, any) (plugin.Client, error) {
		return c, nil
	}
}

func NewMemDBClient(_ context.Context, _ zerolog.Logger, spec any) (plugin.Client, error) {
	return &client{
		memoryDB: make(map[string][]arrow.Record),
		tables:   make(map[string]*schema.Table),
	}, nil
}

func NewMemDBClientErrOnNew(context.Context, zerolog.Logger, any) (plugin.Client, error) {
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

func (c *client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- plugin.Message) error {
	c.memoryDBLock.RLock()

	for tableName := range c.memoryDB {
		if !plugin.IsTable(tableName, options.Tables, options.SkipTables) {
			continue
		}
		for _, row := range c.memoryDB[tableName] {
			res <- &plugin.MessageInsert{
				Record: row,
				Upsert: false,
			}
		}
	}
	c.memoryDBLock.RUnlock()
	return nil
}

func (c *client) Tables(ctx context.Context) (schema.Tables, error) {
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

func (c *client) Write(ctx context.Context, options plugin.WriteOptions, msgs <-chan plugin.Message) error {
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
		case *plugin.MessageCreateTable:
			c.migrate(ctx, msg.Table)
		case *plugin.MessageDeleteStale:
			c.deleteStale(ctx, msg)
		case *plugin.MessageInsert:
			sc := msg.Record.Schema()
			tableName, ok := sc.Metadata().GetValue(schema.MetadataTableName)
			if !ok {
				return fmt.Errorf("table name not found in schema metadata")
			}
			table := c.tables[tableName]
			if msg.Upsert {
				c.overwrite(table, msg.Record)
			} else {
				c.memoryDB[tableName] = append(c.memoryDB[tableName], msg.Record)
			}
		}

		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *client) Close(context.Context) error {
	c.memoryDB = nil
	return nil
}

func (c *client) deleteStale(_ context.Context, msg *plugin.MessageDeleteStale) {
	var filteredTable []arrow.Record
	tableName := msg.Table.Name
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
