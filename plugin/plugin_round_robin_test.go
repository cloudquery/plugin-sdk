package plugin

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v0"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/rs/zerolog"
)

type testPluginClient struct {
	memoryDB      map[string][]arrow.Record
	tables        map[string]*schema.Table
	memoryDBLock  sync.RWMutex
}

type testPluginSpec struct {
	ConnectionString string `json:"connection_string"`
}

func (c *testPluginClient) Sync(ctx context.Context, metrics *Metrics, res chan<- *schema.Resource) error {
	return nil
}

func (c *testPluginClient) Migrate(ctx context.Context, tables schema.Tables) error {
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
	return nil
}

func (c *testPluginClient) Write(ctx context.Context, tables schema.Tables, resources <-chan arrow.Record) error {
	for resource := range resources {
		c.memoryDBLock.Lock()
		sc := resource.Schema()
		tableName, ok := sc.Metadata().GetValue(schema.MetadataTableName)
		if !ok {
			return fmt.Errorf("table name not found in schema metadata")
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

func (c *testPluginClient) overwrite(table *schema.Table, data arrow.Record) {
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

func (c *testPluginClient) deleteStaleTable(_ context.Context, table *schema.Table, source string, syncTime time.Time) {
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

func (c *testPluginClient) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	return nil
}

func (c *testPluginClient) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error {
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
		if arr.(*array.String).Value(0) == sourceName {
			sortedRes = append(sortedRes, row)
		}
	}
	c.memoryDBLock.RUnlock()

	for _, row := range sortedRes {
		res <- row
	}
	return nil
}

func NewTestPluginClient(context.Context, zerolog.Logger, pbPlugin.Spec) (Client, error) {
	return &testPluginClient{
		memoryDB: make(map[string][]arrow.Record),
		tables:   make(map[string]*schema.Table),
	}, nil
}

func TestPluginRoundRobin(t *testing.T) {
	p := NewPlugin("test", "v0.0.0", NewTestPluginClient)
}