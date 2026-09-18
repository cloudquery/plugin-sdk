package message

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func deleteRecordTestTable() *schema.Table {
	grandChild := &schema.Table{
		Name:    "test_grandchild",
		Columns: schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
	}
	child1 := &schema.Table{
		Name:      "test_child_one",
		Columns:   schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
		Relations: schema.Tables{grandChild},
	}
	child2 := &schema.Table{
		Name:    "test_child_two",
		Columns: schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
	}
	table := &schema.Table{
		Name: "test_parent",
		Columns: schema.ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKey: true},
			{Name: "region", Type: arrow.BinaryTypes.String},
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true},
		},
		Relations: schema.Tables{child1, child2},
	}
	return table.Copy(nil)
}

func TestNewDeleteRecordByKeys(t *testing.T) {
	table := deleteRecordTestTable()
	keys, err := schema.NewKeysRecord(table, map[string]any{"account_id": "123456789012", "id": int64(42)})
	require.NoError(t, err)

	msg, err := NewDeleteRecordByKeys(table, keys)
	require.NoError(t, err)
	require.Equal(t, "test_parent", msg.TableName)

	require.Len(t, msg.WhereClause, 1)
	group := msg.WhereClause[0]
	require.Equal(t, "AND", group.GroupingType)
	require.Len(t, group.Predicates, 2)

	require.Equal(t, "EQ", group.Predicates[0].Operator)
	require.Equal(t, "account_id", group.Predicates[0].Column)
	require.EqualValues(t, 1, group.Predicates[0].Record.NumCols())
	require.EqualValues(t, 1, group.Predicates[0].Record.NumRows())
	require.Equal(t, "account_id", group.Predicates[0].Record.Schema().Field(0).Name)
	require.Equal(t, "123456789012", group.Predicates[0].Record.Column(0).ValueStr(0))

	require.Equal(t, "EQ", group.Predicates[1].Operator)
	require.Equal(t, "id", group.Predicates[1].Column)
	require.Equal(t, "id", group.Predicates[1].Record.Schema().Field(0).Name)
	require.Equal(t, "42", group.Predicates[1].Record.Column(0).ValueStr(0))

	require.ElementsMatch(t, TableRelations{
		{TableName: "test_child_one", ParentTable: "test_parent"},
		{TableName: "test_grandchild", ParentTable: "test_child_one"},
		{TableName: "test_child_two", ParentTable: "test_parent"},
	}, msg.TableRelations)
}

func TestNewDeleteRecordByKeys_NoRelations(t *testing.T) {
	table := &schema.Table{
		Name:    "test_leaf",
		Columns: schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
	}
	keys, err := schema.NewKeysRecord(table, map[string]any{"arn": "arn:aws:ec2:i-123"})
	require.NoError(t, err)

	msg, err := NewDeleteRecordByKeys(table, keys)
	require.NoError(t, err)
	require.Empty(t, msg.TableRelations)
	require.Len(t, msg.WhereClause[0].Predicates, 1)
}

func TestNewDeleteRecordByKeys_ColumnMismatch(t *testing.T) {
	table := deleteRecordTestTable()
	other := &schema.Table{
		Name:    "test_other",
		Columns: schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
	}
	keys, err := schema.NewKeysRecord(other, map[string]any{"arn": "arn:aws:ec2:i-123"})
	require.NoError(t, err)

	_, err = NewDeleteRecordByKeys(table, keys)
	require.ErrorContains(t, err, "test_parent")
}

func TestNewDeleteRecordByKeys_TooManyRows(t *testing.T) {
	table := &schema.Table{
		Name:    "test_leaf",
		Columns: schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
	}

	sc := arrow.NewSchema([]arrow.Field{{Name: "arn", Type: arrow.BinaryTypes.String, Nullable: true}}, nil)
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	defer bldr.Release()
	bldr.Field(0).(*array.StringBuilder).AppendValues([]string{"a", "b"}, nil)
	twoRows := bldr.NewRecordBatch()

	_, err := NewDeleteRecordByKeys(table, twoRows)
	require.ErrorContains(t, err, "exactly one row")
}

func TestNewDeleteRecordByKeys_NilKeys(t *testing.T) {
	table := &schema.Table{
		Name:    "test_leaf",
		Columns: schema.ColumnList{{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true}},
	}
	_, err := NewDeleteRecordByKeys(table, nil)
	require.Error(t, err)
}
