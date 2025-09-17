package schema

import (
	"fmt"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/stretchr/testify/require"
)

func TestFindEmptyColumns(t *testing.T) {
	table := TestTable("test", TestSourceOptions{})
	tg := NewTestDataGenerator(0)
	record := tg.Generate(table, GenTestDataOptions{
		MaxRows:  1,
		NullRows: true,
	})
	v := FindEmptyColumns(table, []arrow.RecordBatch{record})
	require.NotEmpty(t, v)
	require.Len(t, v, len(table.Columns)-1) // exclude "id"
}

func TestFindEmptyColumnsNotEmpty(t *testing.T) {
	table := TestTable("test", TestSourceOptions{})
	tg := NewTestDataGenerator(0)
	record := tg.Generate(table, GenTestDataOptions{
		MaxRows:  1,
		NullRows: false,
	})
	v := FindEmptyColumns(table, []arrow.RecordBatch{record})
	require.Empty(t, v)
}

func TestFindEmptyColumnsJSON(t *testing.T) {
	table := &Table{
		Name: "test",
		Columns: ColumnList{
			{Name: "json", Type: types.ExtensionTypes.JSON},
		},
	}
	sc := table.ToArrowSchema()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	err := bldr.Field(0).UnmarshalJSON([]byte(`[{}]`))
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal json for column: %v", err))
	}
	records := []arrow.RecordBatch{bldr.NewRecordBatch()}
	bldr.Release()

	v := FindEmptyColumns(table, records)
	require.NotEmpty(t, v)
	require.Len(t, v, 1)
}
