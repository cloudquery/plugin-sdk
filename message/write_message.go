package message

import (
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type writeBaseMessage struct {
}

func (*writeBaseMessage) IsWriteMessage() bool {
	return true
}

type WriteMessage interface {
	GetTable() *schema.Table
	IsWriteMessage() bool
}

type WriteMigrateTable struct {
	writeBaseMessage
	Table        *schema.Table
	ForceMigrate bool
}

func (m WriteMigrateTable) GetTable() *schema.Table {
	return m.Table
}

type WriteInsert struct {
	writeBaseMessage
	Record arrow.Record
}

func (m *WriteInsert) GetTable() *schema.Table {
	table, err := schema.NewTableFromArrowSchema(m.Record.Schema())
	if err != nil {
		panic(err)
	}
	return table
}

type WriteMessages []WriteMessage

type WriteMigrateTables []*WriteMigrateTable

type WriteInserts []*WriteInsert

func (messages WriteMessages) InsertItems() int64 {
	items := int64(0)
	for _, msg := range messages {
		if m, ok := msg.(*WriteInsert); ok {
			items += m.Record.NumRows()
		}
	}
	return items
}

func (messages WriteMessages) InsertMessage() WriteInserts {
	inserts := []*WriteInsert{}
	for _, msg := range messages {
		if m, ok := msg.(*WriteInsert); ok {
			inserts = append(inserts, m)
		}
	}
	return inserts
}

func (m WriteMigrateTables) Exists(tableName string) bool {
	for _, table := range m {
		if table.Table.Name == tableName {
			return true
		}
	}
	return false
}

func (m WriteInserts) Exists(tableName string) bool {
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

func (m WriteInserts) GetRecordsForTable(table *schema.Table) []arrow.Record {
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

// DeleteStale is a pretty specific message which requires the destination to be aware of a CLI use-case
// thus it might be deprecated in the future
// in favour of MessageDelete or MessageRawQuery
// The message indeciates that the destination needs to run something like "DELETE FROM table WHERE _cq_source_name=$1 and sync_time < $2"
type WriteDeleteStale struct {
	writeBaseMessage
	Table      *schema.Table
	SourceName string
	SyncTime   time.Time
}

func (m WriteDeleteStale) GetTable() *schema.Table {
	return m.Table
}
