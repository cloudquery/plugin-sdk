package schema

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/require"
)

func keysTestTable() *Table {
	return &Table{
		Name: "test_table",
		Columns: ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKey: true},
			{Name: "not_a_key", Type: arrow.BinaryTypes.String},
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true},
		},
	}
}

func TestNewKeysRecord_StringPK(t *testing.T) {
	table := &Table{
		Name: "test_table",
		Columns: ColumnList{
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true},
			{Name: "other", Type: arrow.PrimitiveTypes.Int64},
		},
	}

	rec, err := NewKeysRecord(table, map[string]any{"arn": "arn:aws:ec2:i-123"})
	require.NoError(t, err)
	require.EqualValues(t, 1, rec.NumRows())
	require.EqualValues(t, 1, rec.NumCols())
	require.Equal(t, "arn", rec.Schema().Field(0).Name)
	require.Equal(t, "arn:aws:ec2:i-123", rec.Column(0).ValueStr(0))
}

func TestNewKeysRecord_CompositePKOrder(t *testing.T) {
	table := &Table{
		Name: "test_table",
		Columns: ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKey: true},
			{Name: "region", Type: arrow.BinaryTypes.String},
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true},
		},
	}

	// Map iteration order must not affect the record's column order: it follows
	// table.PrimaryKeys().
	rec, err := NewKeysRecord(table, map[string]any{"id": int64(42), "account_id": "123456789012"})
	require.NoError(t, err)
	require.EqualValues(t, 2, rec.NumCols())
	require.Equal(t, []string{"account_id", "id"}, []string{rec.Schema().Field(0).Name, rec.Schema().Field(1).Name})
	require.Equal(t, "123456789012", rec.Column(0).ValueStr(0))
	require.Equal(t, "42", rec.Column(1).ValueStr(0))
}

func TestNewKeysRecord_MissingPK(t *testing.T) {
	_, err := NewKeysRecord(keysTestTable(), map[string]any{"account_id": "123456789012"})
	require.ErrorContains(t, err, "arn")
	require.ErrorContains(t, err, "test_table")
}

func TestNewKeysRecord_ExtraColumn(t *testing.T) {
	_, err := NewKeysRecord(keysTestTable(), map[string]any{
		"account_id": "123456789012",
		"arn":        "arn:aws:ec2:i-123",
		"not_a_key":  "nope",
	})
	require.ErrorContains(t, err, "not_a_key")
}

func TestNewKeysRecord_UnknownColumn(t *testing.T) {
	_, err := NewKeysRecord(keysTestTable(), map[string]any{
		"account_id":  "123456789012",
		"arn":         "arn:aws:ec2:i-123",
		"nonexistent": "nope",
	})
	require.ErrorContains(t, err, "nonexistent")
}

func TestNewKeysRecord_NoPrimaryKeys(t *testing.T) {
	table := &Table{
		Name:    "test_table",
		Columns: ColumnList{{Name: "col", Type: arrow.BinaryTypes.String}},
	}
	_, err := NewKeysRecord(table, map[string]any{})
	require.ErrorContains(t, err, "no primary keys")
}

func TestKeysString(t *testing.T) {
	table := &Table{
		Name: "test_table",
		Columns: ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKey: true},
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true},
		},
	}

	rec, err := NewKeysRecord(table, map[string]any{"account_id": "123456789012", "id": int64(42)})
	require.NoError(t, err)
	require.Equal(t, "account_id=123456789012,id=42", KeysString(rec))
}

func TestKeysString_Null(t *testing.T) {
	table := &Table{
		Name: "test_table",
		Columns: ColumnList{
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKey: true},
		},
	}

	rec, err := NewKeysRecord(table, map[string]any{"arn": nil})
	require.NoError(t, err)
	require.Equal(t, "arn=(null)", KeysString(rec))
}

func TestKeysString_Empty(t *testing.T) {
	require.Equal(t, "", KeysString(nil))
}
