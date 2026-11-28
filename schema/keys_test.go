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

func TestNewKeysRecord_SubsetOfPKs(t *testing.T) {
	// A subset is fine: an event that knows only one of the keys still names a
	// resource well enough to delete by.
	rec, err := NewKeysRecord(keysTestTable(), map[string]any{"account_id": "123456789012"})
	require.NoError(t, err)
	require.EqualValues(t, 1, rec.NumCols())
	require.Equal(t, "account_id", rec.Schema().Field(0).Name)
	require.Equal(t, "account_id=123456789012", KeysString(rec))
}

func TestNewKeysRecord_EmptyValues(t *testing.T) {
	_, err := NewKeysRecord(keysTestTable(), map[string]any{})
	require.ErrorContains(t, err, "no values given")
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

func TestNewKeysRecord_NoIdentityColumns(t *testing.T) {
	table := &Table{
		Name:    "test_table",
		Columns: ColumnList{{Name: "col", Type: arrow.BinaryTypes.String}},
	}
	_, err := NewKeysRecord(table, map[string]any{})
	require.ErrorContains(t, err, "no identity columns")
}

// componentsTestTable is shaped like an AWS table: _cq_id is the only primary
// key and resource identity lives in the primary key components.
func componentsTestTable() *Table {
	table := &Table{
		Name: "test_components",
		Columns: ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
			{Name: "region", Type: arrow.BinaryTypes.String},
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
		},
	}
	AddCqIDs(table)
	return table
}

func TestIdentityColumns(t *testing.T) {
	// Components win when the table declares any.
	table := componentsTestTable()
	require.Equal(t, []string{CqIDColumn.Name}, table.PrimaryKeys())
	require.Equal(t, []string{"account_id", "arn"}, table.IdentityColumns())

	// Otherwise the primary keys are the identity.
	require.Equal(t, []string{"account_id", "arn"}, keysTestTable().PrimaryKeys())
	require.Equal(t, []string{"account_id", "arn"}, keysTestTable().IdentityColumns())

	// Neither: no identity at all.
	require.Empty(t, (&Table{Name: "t", Columns: ColumnList{{Name: "c", Type: arrow.BinaryTypes.String}}}).IdentityColumns())
}

func TestNewKeysRecord_PrimaryKeyComponents(t *testing.T) {
	table := componentsTestTable()

	rec, err := NewKeysRecord(table, map[string]any{
		"arn":        "arn:aws:ec2:i-123",
		"account_id": "123456789012",
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, rec.NumRows())
	require.EqualValues(t, 2, rec.NumCols())
	require.Equal(t, []string{"account_id", "arn"}, []string{rec.Schema().Field(0).Name, rec.Schema().Field(1).Name})
	require.Equal(t, "123456789012", rec.Column(0).ValueStr(0))
	require.Equal(t, "arn:aws:ec2:i-123", rec.Column(1).ValueStr(0))
	require.Equal(t, "account_id=123456789012,arn=arn:aws:ec2:i-123", KeysString(rec))
}

func TestNewKeysRecord_ComponentTableRejectsCqID(t *testing.T) {
	table := componentsTestTable()

	// _cq_id is the table's only primary key, but it is not an identity column.
	_, err := NewKeysRecord(table, map[string]any{
		"account_id":    "123456789012",
		"arn":           "arn:aws:ec2:i-123",
		CqIDColumn.Name: "5c3b1d9e-0000-0000-0000-000000000000",
	})
	require.ErrorContains(t, err, CqIDColumn.Name)
	require.ErrorContains(t, err, "not an identity column")
}

// iamPolicyLikeTable mirrors aws_iam_policies: four identity columns, of which a
// CloudTrail event knows only account_id and arn.
func iamPolicyLikeTable() *Table {
	table := &Table{
		Name: "test_iam_policies",
		Columns: ColumnList{
			{Name: "account_id", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
			{Name: "id", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
			{Name: "input_hash", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
			{Name: "arn", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
			{Name: "description", Type: arrow.BinaryTypes.String},
		},
	}
	AddCqIDs(table)
	return table
}

func TestNewKeysRecord_SubsetOfIdentityColumns(t *testing.T) {
	table := iamPolicyLikeTable()
	require.Equal(t, []string{"account_id", "id", "input_hash", "arn"}, table.IdentityColumns())

	// The event knows 2 of the 4 identity columns.
	rec, err := NewKeysRecord(table, map[string]any{
		"arn":        "arn:aws:iam::123456789012:policy/demo",
		"account_id": "123456789012",
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, rec.NumRows())
	require.EqualValues(t, 2, rec.NumCols())

	// Declaration order, not map order: account_id is declared before arn.
	require.Equal(t, []string{"account_id", "arn"}, []string{rec.Schema().Field(0).Name, rec.Schema().Field(1).Name})
	require.Equal(t, "account_id=123456789012,arn=arn:aws:iam::123456789012:policy/demo", KeysString(rec))
}

func TestNewKeysRecord_SubsetStillRejectsNonIdentityColumn(t *testing.T) {
	_, err := NewKeysRecord(iamPolicyLikeTable(), map[string]any{
		"arn":         "arn:aws:iam::123456789012:policy/demo",
		"description": "not an identity column",
	})
	require.ErrorContains(t, err, "description")
	require.ErrorContains(t, err, "not an identity column")
}

func TestNewKeysRecord_ComponentTableSubset(t *testing.T) {
	rec, err := NewKeysRecord(componentsTestTable(), map[string]any{"arn": "arn:aws:ec2:i-123"})
	require.NoError(t, err)
	require.EqualValues(t, 1, rec.NumCols())
	require.Equal(t, "arn=arn:aws:ec2:i-123", KeysString(rec))
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
