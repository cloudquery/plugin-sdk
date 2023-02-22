package schema

import (
	"testing"

	"github.com/apache/arrow/go/arrow"
	"github.com/google/go-cmp/cmp"
)

func TestCQSchemaToArrow(t *testing.T) {

	expecetdSchema := arrow.NewSchema([]arrow.Field{
		{Name: "_cq_id", Type: &arrow.FixedSizeBinaryType{ByteWidth: 16},
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeUUID,
			})},
		{Name: "_cq_parent_id", Type: &arrow.FixedSizeBinaryType{ByteWidth: 16}, Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeUUID,
			})},
		{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, Nullable: true},
		{Name: "int", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "float", Type: arrow.PrimitiveTypes.Float64, Nullable: true},
		{Name: "uuid", Type: &arrow.FixedSizeBinaryType{ByteWidth: 16}, Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataPrimaryKey: MetadataPrimaryKeyTrue,
				MetadataLogicalType: MetadataLogicalTypeUUID,
			})},
		{Name: "text", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "text_with_null", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "bytea", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "text_array", Type: arrow.ListOf(arrow.BinaryTypes.String), Nullable: true},
		{Name: "text_array_with_null", Type: arrow.ListOf(arrow.BinaryTypes.String), Nullable: true},
		{Name: "int_array", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64), Nullable: true},
		{Name: "timestamp", Type: arrow.FixedWidthTypes.Timestamp_s, Nullable: true},
		{Name: "json", Type: arrow.BinaryTypes.Binary, Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeJSON,
			})},
		{Name: "uuid_array", Type: arrow.ListOf(&arrow.FixedSizeBinaryType{ByteWidth: 16}), Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeUUID,
			})},
		{Name: "inet", Type: arrow.BinaryTypes.Binary, Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeInet,
			})},
		{Name: "inet_array", Type: arrow.ListOf(arrow.BinaryTypes.Binary), Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeInet,
			})},
		{Name: "cidr", Type: arrow.BinaryTypes.Binary, Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeCIDR,
			})},
		{Name: "cidr_array", Type: arrow.ListOf(arrow.BinaryTypes.Binary), Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeCIDR,
			})},
		{Name: "macaddr", Type: arrow.BinaryTypes.Binary, Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeMacAddr,
			})},
		{Name: "macaddr_array", Type: arrow.ListOf(arrow.BinaryTypes.Binary), Nullable: true,
			Metadata: arrow.MetadataFrom(map[string]string{
				MetadataLogicalType: MetadataLogicalTypeMacAddr,
			})},
	}, nil)

	testTable := TestSourceTable("test_table")
	arrowSchema := CQSchemaToArrow(testTable)
	if diff := cmp.Diff(arrowSchema.String(), expecetdSchema.String()); diff != "" {
		t.Errorf(diff)
	}
	if !arrowSchema.Equal(expecetdSchema) {
		t.Errorf("got:\n%v\nwant:\n%v\n", arrowSchema, expecetdSchema)
	}
}