package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
)

func TestSchemaEncode(t *testing.T) {
	md := arrow.NewMetadata([]string{"true"}, []string{"false"})
	md1 := arrow.NewMetadata([]string{"false"}, []string{"true"})
	schemas := Schemas{
		arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "name", Type: arrow.BinaryTypes.String},
			},
			&md,
		),
		arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "name", Type: arrow.BinaryTypes.String},
			},
			&md1,
		),
	}
	b, err := schemas.Encode()
	if err != nil {
		t.Fatal(err)
	}
	decodedSchemas, err := NewSchemasFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(decodedSchemas) != len(schemas) {
		t.Fatalf("expected %d schemas, got %d", len(schemas), len(decodedSchemas))
	}
	for i := range schemas {
		if !schemas[i].Equal(decodedSchemas[i]) {
			t.Fatalf("expected schema %d to be %v, got %v", i, schemas[i], decodedSchemas[i])
		}
	}
}

func TestRecordToBytesAndNewRecordFromBytes(t *testing.T) {
	md := arrow.NewMetadata([]string{"key"}, []string{"value"})
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "name", Type: arrow.BinaryTypes.String},
		},
		&md,
	)
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer bldr.Release()
	bldr.Field(0).AppendValueFromString("1")
	bldr.Field(1).AppendValueFromString("foo")
	record := bldr.NewRecord()
	b, err := RecordToBytes(record)
	if err != nil {
		t.Fatal(err)
	}
	decodedRecord, err := NewRecordFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	numRows := record.NumRows()
	if numRows != 1 {
		t.Fatalf("expected 1 row, got %d", numRows)
	}
	if diff := RecordDiff(record, decodedRecord); diff != "" {
		t.Fatalf("record differs from expected after NewRecordFromBytes: %v", diff)
	}
}

func TestSchemaToBytesAndNewSchemaFromBytes(t *testing.T) {
	md := arrow.NewMetadata([]string{"key"}, []string{"value"})
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "name", Type: arrow.BinaryTypes.String},
		},
		&md,
	)
	b, err := ToBytes(schema)
	if err != nil {
		t.Fatal(err)
	}
	decodedSchema, err := NewFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !schema.Equal(decodedSchema) {
		t.Fatalf("schema differs from expected after NewSchemaFromBytes. \nBefore: %v,\nAfter: %v", schema, decodedSchema)
	}
}

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
