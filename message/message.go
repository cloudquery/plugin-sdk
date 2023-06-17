package message

import (
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type Message interface {
	GetTable() *schema.Table
}

type MigrateTable struct {
	Table *schema.Table
}

func (m MigrateTable) GetTable() *schema.Table {
	return m.Table
}

type Insert struct {
	Record arrow.Record
	Upsert bool
}

func (m *Insert) GetTable() *schema.Table {
	table, err := schema.NewTableFromArrowSchema(m.Record.Schema())
	if err != nil {
		panic(err)
	}
	return table
}

// DeleteStale is a pretty specific message which requires the destination to be aware of a CLI use-case
// thus it might be deprecated in the future
// in favour of MessageDelete or MessageRawQuery
// The message indeciates that the destination needs to run something like "DELETE FROM table WHERE _cq_source_name=$1 and sync_time < $2"
type DeleteStale struct {
	Table      *schema.Table
	SourceName string
	SyncTime   time.Time
}

func (m DeleteStale) GetTable() *schema.Table {
	return m.Table
}

type Messages []Message

type MigrateTables []*MigrateTable

type Inserts []*Insert

func (messages Messages) InsertItems() int64 {
	items := int64(0)
	for _, msg := range messages {
		switch m := msg.(type) {
		case *Insert:
			items += m.Record.NumRows()
		}
	}
	return items
}

func (messages Messages) InsertMessage() Inserts {
	inserts := []*Insert{}
	for _, msg := range messages {
		switch m := msg.(type) {
		case *Insert:
			inserts = append(inserts, m)
		}
	}
	return inserts
}

func (m MigrateTables) Exists(tableName string) bool {
	for _, table := range m {
		if table.Table.Name == tableName {
			return true
		}
	}
	return false
}

func (m Inserts) Exists(tableName string) bool {
	for _, insert := range m {
		md := insert.Record.Schema().Metadata()
		tableNameMeta, ok := md.GetValue(schema.MetadataTableName)
		if !ok {
			continue
		}
		if tableNameMeta == tableName {
			return true
		}
	}
	return false
}

func (m Inserts) GetRecordsForTable(table *schema.Table) []arrow.Record {
	res := []arrow.Record{}
	for _, insert := range m {
		md := insert.Record.Schema().Metadata()
		tableNameMeta, ok := md.GetValue(schema.MetadataTableName)
		if !ok {
			continue
		}
		if tableNameMeta == table.Name {
			res = append(res, insert.Record)
		}
	}
	return res
}
