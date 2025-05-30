package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RecordDiff(l arrow.Record, r arrow.Record) string {
	var sb strings.Builder
	if l.NumCols() != r.NumCols() {
		return fmt.Sprintf("different number of columns: %d vs %d", l.NumCols(), r.NumCols())
	}
	if l.NumRows() != r.NumRows() {
		return fmt.Sprintf("different number of rows: %d vs %d", l.NumRows(), r.NumRows())
	}
	for i := 0; i < int(l.NumCols()); i++ {
		edits, err := array.Diff(l.Column(i), r.Column(i))
		if err != nil {
			panic(fmt.Sprintf("left: %v, right: %v, error: %v", l.Column(i).DataType(), r.Column(i).DataType(), err))
		}
		diff := edits.UnifiedDiff(l.Column(i), r.Column(i))
		if diff != "" {
			sb.WriteString(l.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func buildTestRecord(withClientIDValue string) arrow.Record {
	testFields := []arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "name", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64, Nullable: true},
		{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, Nullable: true},
		{Name: "uuid", Type: types.UUID, Nullable: true},
	}
	if withClientIDValue != "" {
		testFields = append(testFields, CqClientIDColumn.ToArrowField())
	}
	schema := arrow.NewSchema(testFields, nil)

	testValuesCount := 10
	builders := []array.Builder{
		array.NewInt64Builder(memory.DefaultAllocator),
		array.NewStringBuilder(memory.DefaultAllocator),
		array.NewFloat64Builder(memory.DefaultAllocator),
		array.NewBooleanBuilder(memory.DefaultAllocator),
		types.NewUUIDBuilder(memory.DefaultAllocator),
	}
	for _, builder := range builders {
		builder.Reserve(testValuesCount)
		switch b := builder.(type) {
		case *array.Int64Builder:
			for i := range testValuesCount {
				b.Append(int64(i))
			}
		case *array.StringBuilder:
			for i := range testValuesCount {
				b.AppendString(fmt.Sprintf("test%d", i))
			}
		case *array.Float64Builder:
			for i := range testValuesCount {
				b.Append(float64(i))
			}
		case *array.BooleanBuilder:
			for i := range testValuesCount {
				b.Append(i%2 == 0)
			}
		case *types.UUIDBuilder:
			for i := range testValuesCount {
				b.Append(newUUID(uuid.NameSpaceURL, []byte(fmt.Sprintf("test%d", i))))
			}
		}
	}
	if withClientIDValue != "" {
		builder := array.NewStringBuilder(memory.DefaultAllocator)
		builder.Reserve(testValuesCount)
		for range testValuesCount {
			builder.AppendString(withClientIDValue)
		}
		builders = append(builders, builder)
	}
	values := lo.Map(builders, func(builder array.Builder, _ int) arrow.Array {
		return builder.NewArray()
	})
	return array.NewRecord(schema, values, int64(testValuesCount))
}

func TestAddInternalColumnsToRecord(t *testing.T) {
	tests := []struct {
		name               string
		record             arrow.Record
		cqClientIDValue    string
		expectedNewColumns int64
	}{
		{
			name:               "add _cq_id,_cq_parent_id,_cq_client_id",
			record:             buildTestRecord(""),
			cqClientIDValue:    "new_client_id",
			expectedNewColumns: 3,
		},
		{
			name:               "add cq_client_id,cq_id replace existing _cq_client_id",
			record:             buildTestRecord("existing_client_id"),
			cqClientIDValue:    "new_client_id",
			expectedNewColumns: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddInternalColumnsToRecord(tt.record, tt.cqClientIDValue)
			require.NoError(t, err)
			require.Equal(t, tt.record.NumRows(), got.NumRows())
			require.Equal(t, tt.record.NumCols()+tt.expectedNewColumns, got.NumCols())

			gotSchema := got.Schema()
			cqIDFields := gotSchema.FieldIndices(CqIDColumn.Name)
			require.Len(t, cqIDFields, 1)

			cqParentIDFields := gotSchema.FieldIndices(CqParentIDColumn.Name)
			require.Len(t, cqParentIDFields, 1)

			cqClientIDFields := gotSchema.FieldIndices(CqClientIDColumn.Name)
			require.Len(t, cqClientIDFields, 1)

			cqIDArray := got.Column(cqIDFields[0])
			require.Equal(t, types.UUID, cqIDArray.DataType())
			require.Equal(t, tt.record.NumRows(), int64(cqIDArray.Len()))

			cqParentIDArray := got.Column(cqParentIDFields[0])
			require.Equal(t, types.UUID, cqParentIDArray.DataType())
			require.Equal(t, tt.record.NumRows(), int64(cqParentIDArray.Len()))

			cqClientIDArray := got.Column(cqClientIDFields[0])
			require.Equal(t, arrow.BinaryTypes.String, cqClientIDArray.DataType())
			require.Equal(t, tt.record.NumRows(), int64(cqClientIDArray.Len()))

			for i := range cqIDArray.Len() {
				cqID := cqIDArray.GetOneForMarshal(i).(uuid.UUID)
				require.NotEmpty(t, cqID)
			}
			for i := range cqParentIDArray.Len() {
				cqParentID := cqParentIDArray.GetOneForMarshal(i)
				require.Nil(t, cqParentID)
			}
			for i := range cqClientIDArray.Len() {
				cqClientID := cqClientIDArray.GetOneForMarshal(i).(string)
				require.Equal(t, tt.cqClientIDValue, cqClientID)
			}
		})
	}
}

func TestCQIDHashingConsistency(t *testing.T) {
	record := buildTestRecord("")
	got, err := AddInternalColumnsToRecord(record, "")
	require.NoError(t, err)
	cqIDFields := got.Schema().FieldIndices(CqIDColumn.Name)
	require.Len(t, cqIDFields, 1)
	// we are now using an internal version of the official SHA1 module
	// d8f3b1de-8c63-5a0e-a1aa-19e9b5311c24 is the expected hash value from the official SHA1 module that we expect in our implementation
	assert.Equal(t, "d8f3b1de-8c63-5a0e-a1aa-19e9b5311c24", got.Column(cqIDFields[0]).ValueStr(0))
}
