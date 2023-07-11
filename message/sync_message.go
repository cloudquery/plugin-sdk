package message

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"golang.org/x/exp/slices"
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

func (m SyncInserts) GetRecords() []arrow.Record {
	res := make([]arrow.Record, len(m))
	for i := range m {
		res[i] = m[i].Record
	}
	return res
}

// Get all records for a single table
func (m SyncInserts) GetRecordsForTable(table *schema.Table) []arrow.Record {
	res := make([]arrow.Record, 0, len(m))
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

// Get all records for an  array of tables including any table relations
func (m SyncInserts) GetRecordsForTables(tables schema.Tables) []arrow.Record {
	allTables := tables.FlattenTables()
	res := make([]arrow.Record, 0, len(m))
	tableNames := make([]string, 0, len(allTables))
	for i, tableName := range allTables {
		tableNames[i] = tableName.Name
	}
	for _, insert := range m {
		md := insert.Record.Schema().Metadata()
		tableNameMeta, ok := md.GetValue(schema.MetadataTableName)
		if !ok {
			continue
		}
		for _, tableName := range tableNames {
			if tableNameMeta == tableName {
				res = append(res, insert.Record)
				break
			}
		}
	}
	return slices.Clip(res)
}
