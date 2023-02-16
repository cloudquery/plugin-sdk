package memdb

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/plugins/destination"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

// client is mostly used for testing the destination plugin.
type client struct {
	schema.DefaultTransformer
	spec          specs.Destination
	memoryDB      map[string][][]any
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

func GetNewClient(options ...Option) destination.NewClientFunc {
	c := &client{
		memoryDB:     make(map[string][][]any),
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
		memoryDB: make(map[string][][]any),
		tables:   make(map[string]*schema.Table),
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
func (c *client) overwrite(table *schema.Table, data []any) {
	pks := table.PrimaryKeys()
	pksIndex := make([]int, len(pks))
	for i := range pks {
		pksIndex[i] = table.Columns.Index(pks[i])
	}
	for i, row := range c.memoryDB[table.Name] {
		found := true
		for _, pkIndex := range pksIndex {
			if !row[pkIndex].(schema.CQType).Equal(data[pkIndex].(schema.CQType)) {
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

func (c *client) Migrate(_ context.Context, tables schema.Tables) error {
	for _, table := range tables {
		memTable := c.memoryDB[table.Name]
		if memTable == nil {
			c.memoryDB[table.Name] = make([][]any, 0)
			c.tables[table.Name] = table
			continue
		}
		changes := table.GetChanges(c.tables[table.Name])
		// memdb doesn't support any auto-migrate
		if changes == nil {
			continue
		}
		c.memoryDB[table.Name] = make([][]any, 0)
		c.tables[table.Name] = table
	}
	return nil
}

func (c *client) Read(_ context.Context, table *schema.Table, source string, res chan<- []any) error {
	if c.memoryDB[table.Name] == nil {
		return nil
	}
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	var sortedRes [][]any
	c.memoryDBLock.RLock()
	for _, row := range c.memoryDB[table.Name] {
		if row[sourceColIndex].(*schema.Text).Str == source {
			sortedRes = append(sortedRes, row)
		}
	}
	c.memoryDBLock.RUnlock()

	for _, row := range sortedRes {
		res <- row
	}
	return nil
}

func (c *client) Write(ctx context.Context, tables schema.Tables, resources <-chan *destination.ClientResource) error {
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
		if c.spec.WriteMode == specs.WriteModeAppend {
			c.memoryDB[resource.TableName] = append(c.memoryDB[resource.TableName], resource.Data)
		} else {
			c.overwrite(tables.Get(resource.TableName), resource.Data)
		}
		c.memoryDBLock.Unlock()
	}
	return nil
}

func (c *client) WriteTableBatch(ctx context.Context, table *schema.Table, resources [][]any) error {
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

func (c *client) DeleteStale(ctx context.Context, tables schema.Tables, source string, syncTime time.Time) error {
	for _, table := range tables {
		c.deleteStaleTable(ctx, table, source, syncTime)
		if err := c.DeleteStale(ctx, table.Relations, source, syncTime); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) deleteStaleTable(_ context.Context, table *schema.Table, source string, syncTime time.Time) {
	sourceColIndex := table.Columns.Index(schema.CqSourceNameColumn.Name)
	syncColIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	var filteredTable [][]any
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
