package message

import (
	"slices"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type writeBaseMessage struct{}

func (*writeBaseMessage) IsWriteMessage() bool { return true }

type WriteMessage interface {
	GetTable() *schema.Table
	IsWriteMessage() bool
}

type WriteMessages []WriteMessage

func (messages WriteMessages) InsertItems() int64 {
	items := int64(0)
	for _, msg := range messages {
		if m, ok := msg.(*WriteInsert); ok {
			items += m.Record.NumRows()
		}
	}
	return items
}

func (messages WriteMessages) GetInserts() WriteInserts {
	inserts := make(WriteInserts, 0, len(messages))
	for _, msg := range messages {
		if m, ok := msg.(*WriteInsert); ok {
			inserts = append(inserts, m)
		}
	}
	return slices.Clip(inserts)
}

type WriteMigrateTable struct {
	writeBaseMessage
	Table        *schema.Table
	MigrateForce bool
}

func (m WriteMigrateTable) GetTable() *schema.Table { return m.Table }

type WriteMigrateTables []*WriteMigrateTable

func (m WriteMigrateTables) Exists(tableName string) bool {
	return slices.ContainsFunc(m, func(msg *WriteMigrateTable) bool {
		return msg.Table.Name == tableName
	})
}

func (m WriteMigrateTables) GetMessageByTable(tableName string) *WriteMigrateTable {
	for _, msg := range m {
		if msg.Table.Name == tableName {
			return msg
		}
	}
	return nil
}

type WriteInsert struct {
	writeBaseMessage
	Record arrow.RecordBatch
}

func (m *WriteInsert) GetTable() *schema.Table {
	table, err := schema.NewTableFromArrowSchema(m.Record.Schema())
	if err != nil {
		panic(err)
	}
	return table
}

type WriteInserts []*WriteInsert

func (m WriteInserts) Exists(tableName string) bool {
	return slices.ContainsFunc(m, func(msg *WriteInsert) bool {
		tableNameMeta, ok := msg.Record.Schema().Metadata().GetValue(schema.MetadataTableName)
		return ok && tableNameMeta == tableName
	})
}

func (m WriteInserts) GetRecords() []arrow.RecordBatch {
	res := make([]arrow.RecordBatch, len(m))
	for i := range m {
		res[i] = m[i].Record
	}
	return res
}

func (m WriteInserts) GetRecordsForTable(table *schema.Table) []arrow.RecordBatch {
	res := make([]arrow.RecordBatch, 0, len(m))
	for _, insert := range m {
		tableNameMeta, ok := insert.Record.Schema().Metadata().GetValue(schema.MetadataTableName)
		if !ok || tableNameMeta != table.Name {
			continue
		}
		res = append(res, insert.Record)
	}
	return slices.Clip(res)
}

// WriteDeleteStale is a pretty specific message which requires the destination to be aware of a CLI use-case
// thus it might be deprecated in the future
// in favour of MessageDelete or MessageRawQuery
// The message indicates that the destination needs to run something like "DELETE FROM table WHERE _cq_source_name=$1 and sync_time < $2"
type WriteDeleteStale struct {
	writeBaseMessage
	TableName  string
	SourceName string
	SyncTime   time.Time
}

func (m WriteDeleteStale) GetTable() *schema.Table {
	return &schema.Table{Name: m.TableName}
}

type WriteDeleteStales []*WriteDeleteStale

func (m WriteDeleteStales) Exists(tableName string) bool {
	return slices.ContainsFunc(m, func(msg *WriteDeleteStale) bool {
		return msg.TableName == tableName
	})
}

type TableRelation struct {
	TableName   string
	ParentTable string
}

type TableRelations []TableRelation

type Predicate struct {
	Operator string
	Column   string
	Record   arrow.RecordBatch
}

type Predicates []Predicate

type PredicateGroup struct {
	// This will be AND or OR
	GroupingType string
	Predicates   Predicates
}

type PredicateGroups []PredicateGroup

type DeleteRecord struct {
	TableName      string
	TableRelations TableRelations
	WhereClause    PredicateGroups
	SyncTime       time.Time
}

type WriteDeleteRecord struct {
	writeBaseMessage
	DeleteRecord
}

func (m WriteDeleteRecord) GetTable() *schema.Table {
	return &schema.Table{Name: m.TableName}
}

type WriteDeleteRecords []*WriteDeleteRecord
