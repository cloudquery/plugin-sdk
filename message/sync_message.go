package message

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type syncBaseMessage struct {
}

func (*syncBaseMessage) IsSyncMessage() bool {
	return true
}

type SyncMessage interface {
	GetTable() *schema.Table
	IsSyncMessage() bool
}

type SyncMigrateTable struct {
	syncBaseMessage
	Table *schema.Table
}

func (m SyncMigrateTable) GetTable() *schema.Table {
	return m.Table
}

type SyncInsert struct {
	syncBaseMessage
	Record arrow.Record
}

func (m *SyncInsert) GetTable() *schema.Table {
	table, err := schema.NewTableFromArrowSchema(m.Record.Schema())
	if err != nil {
		panic(err)
	}
	return table
}

type SyncMessages []SyncMessage

type SyncMigrateTables []*SyncMigrateTable

type SyncInserts []*SyncInsert

func (messages SyncMessages) InsertItems() int64 {
	items := int64(0)
	for _, msg := range messages {
		if m, ok := msg.(*SyncInsert); ok {
			items += m.Record.NumRows()
		}
	}
	return items
}

func (messages SyncMessages) InsertMessage() SyncInserts {
	inserts := []*SyncInsert{}
	for _, msg := range messages {
		if m, ok := msg.(*SyncInsert); ok {
			inserts = append(inserts, m)
		}
	}
	return inserts
}

func (m SyncMigrateTables) Exists(tableName string) bool {
	for _, table := range m {
		if table.Table.Name == tableName {
			return true
		}
	}
	return false
}

func (m SyncInserts) Exists(tableName string) bool {
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

func (m SyncInserts) GetRecordsForTable(table *schema.Table) []arrow.Record {
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
