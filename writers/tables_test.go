package writers

import (
	"errors"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func TestFailedTables(t *testing.T) {
	table := &schema.Table{Name: "inserted", Columns: schema.ColumnList{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	defer bldr.Release()
	insert := &message.WriteInsert{Record: bldr.NewRecordBatch()}
	failed := errors.New("write failed")

	err := FailedTables(failed, []message.WriteMessage{
		insert,
		&message.WriteMigrateTable{Table: &schema.Table{Name: "migrated"}},
		&message.WriteDeleteStale{TableName: "stale"},
		&message.WriteDeleteRecord{DeleteRecord: message.DeleteRecord{TableName: "deleted"}},
	})
	require.ErrorIs(t, err, failed)
	require.Equal(t, []string{"deleted", "inserted", "migrated", "stale"}, message.FailedTables(err))

	require.NoError(t, FailedTables(nil, []message.WriteMessage{insert}))

	noName := array.NewRecordBuilder(memory.DefaultAllocator, arrow.NewSchema([]arrow.Field{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}, nil))
	defer noName.Release()
	err = FailedTables(failed, []message.WriteMessage{insert, &message.WriteInsert{Record: noName.NewRecordBatch()}})
	require.Equal(t, failed, err, "a message without a table may belong to any table")
}
