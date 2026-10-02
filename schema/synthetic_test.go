package schema

import (
	"encoding/json"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/stretchr/testify/require"
)

func nullableField(dataType arrow.DataType) arrow.Field {
	return arrow.Field{Name: "col", Type: dataType, Nullable: true}
}

func TestSyntheticPairs(t *testing.T) {
	stringList := arrow.ListOf(arrow.BinaryTypes.String)
	tests := []struct {
		name     string
		oldField arrow.Field
		newField arrow.Field
		want     []SyntheticPair
	}{
		{
			name:     "list of strings to json",
			oldField: nullableField(stringList),
			newField: nullableField(types.ExtensionTypes.JSON),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `["env:prod"]`},
				{Case: SyntheticCaseNull, Value: "null"},
				{Case: SyntheticCaseEmptyCollection, Value: `[]`},
				{Case: SyntheticCaseEscaping, Value: `[` + syntheticEscapingString + `]`},
			},
		},
		{
			name:     "json to list of strings",
			oldField: nullableField(types.ExtensionTypes.JSON),
			newField: nullableField(stringList),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `["env:prod"]`},
				{Case: SyntheticCaseNull, Value: "null"},
				{Case: SyntheticCaseEmptyCollection, Value: `[]`},
				{Case: SyntheticCaseEscaping, Value: `[` + syntheticEscapingString + `]`},
			},
		},
		{
			name:     "nested list to json",
			oldField: nullableField(arrow.ListOf(stringList)),
			newField: nullableField(types.ExtensionTypes.JSON),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `[["env:prod"]]`},
				{Case: SyntheticCaseNull, Value: "null"},
				{Case: SyntheticCaseEmptyCollection, Value: `[]`},
				{Case: SyntheticCaseNested, Value: `[["env:prod"]]`},
				{Case: SyntheticCaseEscaping, Value: `[[` + syntheticEscapingString + `]]`},
			},
		},
		{
			name:     "struct to json",
			oldField: nullableField(arrow.StructOf(arrow.Field{Name: "tags", Type: stringList, Nullable: true})),
			newField: nullableField(types.ExtensionTypes.JSON),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `{"tags":["env:prod"]}`},
				{Case: SyntheticCaseNull, Value: "null"},
				{Case: SyntheticCaseNested, Value: `{"tags":["env:prod"]}`},
				{Case: SyntheticCaseEscaping, Value: `{"tags":[` + syntheticEscapingString + `]}`},
			},
		},
		{
			name:     "string to uuid",
			oldField: nullableField(arrow.BinaryTypes.String),
			newField: nullableField(types.ExtensionTypes.UUID),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `"6ba7b810-9dad-11d1-80b4-00c04fd430c8"`},
				{Case: SyntheticCaseNull, Value: "null"},
			},
		},
		{
			name:     "int32 to int64",
			oldField: nullableField(arrow.PrimitiveTypes.Int32),
			newField: nullableField(arrow.PrimitiveTypes.Int64),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `42`},
				{Case: SyntheticCaseNull, Value: "null"},
			},
		},
		{
			name:     "timestamp precision change",
			oldField: nullableField(arrow.FixedWidthTypes.Timestamp_s),
			newField: nullableField(arrow.FixedWidthTypes.Timestamp_us),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `"2024-01-02T03:04:05Z"`},
				{Case: SyntheticCaseNull, Value: "null"},
			},
		},
		{
			name:     "map to map with escaping",
			oldField: nullableField(arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String)),
			newField: nullableField(arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.LargeString)),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `[{"key":"env:prod","value":"env:prod"}]`},
				{Case: SyntheticCaseNull, Value: "null"},
				{Case: SyntheticCaseEmptyCollection, Value: `[]`},
				{Case: SyntheticCaseEscaping, Value: `[{"key":` + syntheticEscapingString + `,"value":` + syntheticEscapingString + `}]`},
			},
		},
		{
			name:     "map of lists is nested",
			oldField: nullableField(arrow.MapOf(arrow.BinaryTypes.String, stringList)),
			newField: nullableField(arrow.MapOf(arrow.BinaryTypes.String, stringList)),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `[{"key":"env:prod","value":["env:prod"]}]`},
				{Case: SyntheticCaseNull, Value: "null"},
				{Case: SyntheticCaseEmptyCollection, Value: `[]`},
				{Case: SyntheticCaseNested, Value: `[{"key":"env:prod","value":["env:prod"]}]`},
				{Case: SyntheticCaseEscaping, Value: `[{"key":` + syntheticEscapingString + `,"value":[` + syntheticEscapingString + `]}]`},
			},
		},
		{
			name:     "not null field has no null pair",
			oldField: arrow.Field{Name: "col", Type: arrow.BinaryTypes.String},
			newField: nullableField(arrow.BinaryTypes.String),
			want: []SyntheticPair{
				{Case: SyntheticCasePopulated, Value: `"env:prod"`},
				{Case: SyntheticCaseEscaping, Value: syntheticEscapingString},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SyntheticPairs(tc.oldField, tc.newField)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestSyntheticPairs_UnableToCompare(t *testing.T) {
	tests := []struct {
		name     string
		oldField arrow.Field
		newField arrow.Field
	}{
		{name: "int to string", oldField: nullableField(arrow.PrimitiveTypes.Int64), newField: nullableField(arrow.BinaryTypes.String)},
		{name: "string to int", oldField: nullableField(arrow.BinaryTypes.String), newField: nullableField(arrow.PrimitiveTypes.Int64)},
		{name: "duration unit change", oldField: nullableField(arrow.FixedWidthTypes.Duration_s), newField: nullableField(arrow.FixedWidthTypes.Duration_ms)},
		{name: "string to list", oldField: nullableField(arrow.BinaryTypes.String), newField: nullableField(arrow.ListOf(arrow.BinaryTypes.String))},
		{name: "bool to int", oldField: nullableField(arrow.FixedWidthTypes.Boolean), newField: nullableField(arrow.PrimitiveTypes.Int64)},
		{name: "fixed size list", oldField: nullableField(arrow.FixedSizeListOf(2, arrow.BinaryTypes.String)), newField: nullableField(arrow.FixedSizeListOf(2, arrow.BinaryTypes.String))},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := SyntheticPairs(tc.oldField, tc.newField)
			require.ErrorIs(t, err, ErrUnableToCompare)
		})
	}
}

func TestSyntheticRecords(t *testing.T) {
	oldField := arrow.Field{Name: "tags", Type: arrow.ListOf(arrow.BinaryTypes.String), Nullable: true}
	newField := arrow.Field{Name: "tags", Type: types.ExtensionTypes.JSON, Nullable: true}
	pairs, err := SyntheticPairs(oldField, newField)
	require.NoError(t, err)

	oldRecord, newRecord, err := SyntheticRecords(oldField, newField, pairs)
	require.NoError(t, err)
	defer oldRecord.Release()
	defer newRecord.Release()

	require.True(t, oldRecord.Schema().Equal(arrow.NewSchema([]arrow.Field{oldField}, nil)))
	require.True(t, newRecord.Schema().Equal(arrow.NewSchema([]arrow.Field{newField}, nil)))
	require.EqualValues(t, len(pairs), oldRecord.NumRows())
	require.EqualValues(t, len(pairs), newRecord.NumRows())

	for i, pair := range pairs {
		oldValue, err := json.Marshal(oldRecord.Column(0).GetOneForMarshal(i))
		require.NoError(t, err)
		newValue, err := json.Marshal(newRecord.Column(0).GetOneForMarshal(i))
		require.NoError(t, err)
		require.JSONEq(t, pair.Value, string(oldValue), "old value for case %s", pair.Case)
		require.JSONEq(t, pair.Value, string(newValue), "new value for case %s", pair.Case)
	}
	populated, err := json.Marshal(newRecord.Column(0).GetOneForMarshal(0))
	require.NoError(t, err)
	require.Equal(t, `["env:prod"]`, string(populated))
}

func TestSyntheticPairs_SameTypeIsComparable(t *testing.T) {
	for _, column := range TestTable("test", TestSourceOptions{}).Columns {
		t.Run(column.Name, func(t *testing.T) {
			field := column.ToArrowField()
			pairs, err := SyntheticPairs(field, field)
			require.NoError(t, err)
			oldRecord, newRecord, err := SyntheticRecords(field, field, pairs)
			require.NoError(t, err)
			oldRecord.Release()
			newRecord.Release()
		})
	}
}
