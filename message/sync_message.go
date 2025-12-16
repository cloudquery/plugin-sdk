package message

import (
	"slices"

	"github.com/apache/arrow-go/v18/arrow"
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
	Record arrow.RecordBatch
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

func (messages SyncMessages) GetInserts() SyncInserts {
	inserts := make(SyncInserts, 0, len(messages))
	for _, msg := range messages {
		if m, ok := msg.(*SyncInsert); ok {
			inserts = append(inserts, m)
		}
	}
	return slices.Clip(inserts)
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

func (m SyncInserts) GetRecords() []arrow.RecordBatch {
	res := make([]arrow.RecordBatch, len(m))
	for i := range m {
		res[i] = m[i].Record
	}
	return res
}

// Get all records for a single table
func (m SyncInserts) GetRecordsForTable(table *schema.Table) []arrow.RecordBatch {
	res := make([]arrow.RecordBatch, 0, len(m))
	for _, insert := range m {
		md := insert.Record.Schema().Metadata()
		tableNameMeta, ok := md.GetValue(schema.MetadataTableName)
		if !ok || tableNameMeta != table.Name {
			continue
		}
		res = append(res, insert.Record)
	}
	return slices.Clip(res)
}

type SyncDeleteRecord struct {
	syncBaseMessage
	// TODO: Instead of using this struct we should derive the DeletionKeys and parent/child relation from the schema.Table itself
	DeleteRecord
}

func (m SyncDeleteRecord) GetTable() *schema.Table {
	return &schema.Table{Name: m.TableName}
}

type SyncError struct {
	syncBaseMessage
	TableName string
	Error     string
}

func (e SyncError) GetTable() *schema.Table {
	return &schema.Table{Name: e.TableName}
}
