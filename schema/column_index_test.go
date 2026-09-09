package schema

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/require"
)

func stringColumns(names ...string) ColumnList {
	cols := make(ColumnList, len(names))
	for i, name := range names {
		cols[i] = Column{Name: name, Type: arrow.BinaryTypes.String}
	}
	return cols
}

func TestColumnIndex(t *testing.T) {
	table := &Table{Name: "test", Columns: stringColumns("a", "b", "c")}

	require.Equal(t, 1, table.ColumnIndex("b"))
	require.Equal(t, -1, table.ColumnIndex("missing"))

	table.BuildColumnIndex()
	require.Equal(t, 0, table.ColumnIndex("a"))
	require.Equal(t, 1, table.ColumnIndex("b"))
	require.Equal(t, 2, table.ColumnIndex("c"))
	require.Equal(t, -1, table.ColumnIndex("missing"))
}

func TestColumnIndexRelations(t *testing.T) {
	table := &Table{
		Name:    "parent",
		Columns: stringColumns("a", "b"),
		Relations: Tables{
			{Name: "child", Columns: stringColumns("c", "d")},
		},
	}
	table.BuildColumnIndex()
	require.Equal(t, 1, table.Relations[0].ColumnIndex("d"))
}

// Columns mutated behind the back of every mutator must not make lookups wrong.
func TestColumnIndexStaleCache(t *testing.T) {
	table := &Table{Name: "test", Columns: stringColumns("a", "b", "c")}
	table.BuildColumnIndex()

	table.Columns = append(stringColumns("z"), table.Columns...)
	require.Equal(t, 0, table.ColumnIndex("z"))
	require.Equal(t, 1, table.ColumnIndex("a"))
	require.Equal(t, 2, table.ColumnIndex("b"))
	require.Equal(t, 3, table.ColumnIndex("c"))

	table.Columns = table.Columns[:1]
	require.Equal(t, 0, table.ColumnIndex("z"))
	require.Equal(t, -1, table.ColumnIndex("c"))
}

func TestColumnIndexInvalidatedByMutators(t *testing.T) {
	t.Run("AddCqIDs", func(t *testing.T) {
		table := &Table{Name: "test", Columns: stringColumns("a", "b")}
		table.BuildColumnIndex()
		AddCqIDs(table)
		require.Nil(t, table.columnIndex)
		require.Equal(t, 2, table.ColumnIndex("a"))
	})

	t.Run("AddCqClientID", func(t *testing.T) {
		table := &Table{Name: "test", Columns: stringColumns("a", "b")}
		table.BuildColumnIndex()
		AddCqClientID(table)
		require.Nil(t, table.columnIndex)
		require.Equal(t, 1, table.ColumnIndex("a"))
	})

	t.Run("OverwriteOrAddColumn", func(t *testing.T) {
		table := &Table{Name: "test", Columns: stringColumns("a", "b")}
		table.BuildColumnIndex()
		table.OverwriteOrAddColumn(&Column{Name: "z", Type: arrow.BinaryTypes.String})
		require.Nil(t, table.columnIndex)
		require.Equal(t, 0, table.ColumnIndex("z"))
		require.Equal(t, 1, table.ColumnIndex("a"))
	})

	t.Run("Copy", func(t *testing.T) {
		table := &Table{Name: "test", Columns: stringColumns("a", "b")}
		table.BuildColumnIndex()
		c := table.Copy(nil)
		require.Nil(t, c.columnIndex)
		require.Equal(t, 1, c.ColumnIndex("b"))
	})
}
