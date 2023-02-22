package schema

import (
	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
)

const (
	MetadataPrimaryKey     = "cq:extension:primary_key"
	MetadataPrimaryKeyTrue = "true"

	MetadataLogicalType        = "cq:extension:logical_type"
	MetadataLogicalTypeUUID    = "UUID"
	MetadataLogicalTypeJSON    = "JSON"
	MetadataLogicalTypeInet    = "Inet"
	MetadataLogicalTypeCIDR    = "CIDR"
	MetadataLogicalTypeMacAddr = "MacAddr"
)

func CQColumnToArrowField(col *Column) arrow.Field {
	var typ arrow.DataType
	metadata := make(map[string]string)

	switch col.Type {
	case TypeBool:
		typ = arrow.FixedWidthTypes.Boolean
	case TypeInt:
		typ = arrow.PrimitiveTypes.Int64
	case TypeFloat:
		typ = arrow.PrimitiveTypes.Float64
	case TypeUUID:
		typ = &arrow.FixedSizeBinaryType{ByteWidth: 16}
		metadata[MetadataLogicalType] = MetadataLogicalTypeUUID
	case TypeString:
		typ = arrow.BinaryTypes.String
	case TypeByteArray:
		typ = arrow.BinaryTypes.Binary
	case TypeStringArray:
		typ = arrow.ListOf(arrow.BinaryTypes.String)
	case TypeIntArray:
		typ = arrow.ListOf(arrow.PrimitiveTypes.Int64)
	case TypeTimestamp:
		typ = arrow.FixedWidthTypes.Timestamp_s
	case TypeJSON:
		typ = arrow.BinaryTypes.Binary
		metadata[MetadataLogicalType] = MetadataLogicalTypeJSON
	case TypeUUIDArray:
		typ = arrow.ListOf(&arrow.FixedSizeBinaryType{ByteWidth: 16})
		metadata[MetadataLogicalType] = MetadataLogicalTypeUUID
	case TypeInet:
		typ = arrow.BinaryTypes.Binary
		metadata[MetadataLogicalType] = MetadataLogicalTypeInet
	case TypeInetArray:
		typ = arrow.ListOf(arrow.BinaryTypes.Binary)
		metadata[MetadataLogicalType] = MetadataLogicalTypeInet
	case TypeCIDR:
		typ = arrow.BinaryTypes.Binary
		metadata[MetadataLogicalType] = MetadataLogicalTypeCIDR
	case TypeCIDRArray:
		typ = arrow.ListOf(arrow.BinaryTypes.Binary)
		metadata[MetadataLogicalType] = MetadataLogicalTypeCIDR
	case TypeMacAddr:
		typ = arrow.BinaryTypes.Binary
		metadata[MetadataLogicalType] = MetadataLogicalTypeMacAddr
	case TypeMacAddrArray:
		typ = arrow.ListOf(arrow.BinaryTypes.Binary)
		metadata[MetadataLogicalType] = MetadataLogicalTypeMacAddr
	default:
		panic("unknown type " + typ.Name())
	}
	if col.CreationOptions.PrimaryKey {
		metadata[MetadataPrimaryKey] = MetadataPrimaryKeyTrue
	}
	return arrow.Field{
		Name:     col.Name,
		Type:     typ,
		Nullable: !col.CreationOptions.NotNull,
		Metadata: arrow.MetadataFrom(metadata),
	}
}

func CQSchemaToArrow(table *Table) *arrow.Schema {
	fields := make([]arrow.Field, 0, len(table.Columns))
	for _, col := range table.Columns {
		fields = append(fields, CQColumnToArrowField(&col))
	}
	return arrow.NewSchema(fields, nil)
}


func (c CQTypes) ToRecord(arrowSchema *arrow.Schema) arrow.Record {
	bldr := array.NewRecordBuilder(nil, arrowSchema)
	fields := bldr.Fields()
	for i := range fields {
		switch c[i].Type() {
		case TypeBool:
			bldr.Field(i).(*array.BooleanBuilder).Append(c[i].(*Bool).Bool)
		case TypeInt:
			bldr.Field(i).(*array.Int64Builder).Append(c[i].(*Int8).Int)
		case TypeFloat:
			bldr.Field(i).(*array.Float64Builder).Append(c[i].(*Float8).Float)
		case TypeString:
			bldr.Field(i).(*array.StringBuilder).Append(c[i].(*Text).Str)
		case TypeByteArray:
			bldr.Field(i).(*array.BinaryBuilder).Append(c[i].(*Bytea).Bytes)
		case TypeStringArray:
			listBldr := bldr.Field(i).(*array.ListBuilder)
			for _, str := range c[i].(*TextArray).Elements {
				listBldr.Append(true)
				bldr.Field(i).(*array.StringBuilder).Append(str.Str)
			}
		case TypeIntArray:
			listBldr := bldr.Field(i).(*array.ListBuilder)
			for _, e := range c[i].(*Int8Array).Elements {
				listBldr.Append(true)
				bldr.Field(i).(*array.Int64Builder).Append(e.Int)
			}
		case TypeTimestamp:
			bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(c[i].(*Timestamptz).Time.Unix()))
		case TypeJSON:
		case TypeUUID:
			bldr.Field(i).(*array.BinaryBuilder).Append(c[i].(*UUID).Bytes[:])
		case TypeUUIDArray:
			listBldr := bldr.Field(i).(*array.ListBuilder)
			for _, e := range c[i].(*UUIDArray).Elements {
				listBldr.Append(true)
				bldr.Field(i).(*array.BinaryBuilder).Append(e.Bytes[:])
			}
		case TypeInet:
			bldr.Field(i).(*array.BinaryBuilder).Append(c[i].(*Inet).IPNet.IP)
		case TypeInetArray:
			listBldr := bldr.Field(i).(*array.ListBuilder)
			for _, e := range c[i].(*InetArray).Elements {
				listBldr.Append(true)
				bldr.Field(i).(*array.BinaryBuilder).Append(e.IPNet.IP)
			}
		case TypeCIDR:
			bldr.Field(i).(*array.BinaryBuilder).Append(c[i].(*CIDR).IPNet.IP)
		case TypeCIDRArray:
			listBldr := bldr.Field(i).(*array.ListBuilder)
			for _, e := range c[i].(*CIDRArray).Elements {
				listBldr.Append(true)
				bldr.Field(i).(*array.BinaryBuilder).Append(e.IPNet.IP)
			}
		case TypeMacAddr:
			bldr.Field(i).(*array.BinaryBuilder).Append(c[i].(*Macaddr).Addr)
		case TypeMacAddrArray:
			listBldr := bldr.Field(i).(*array.ListBuilder)
			for _, e := range c[i].(*MacaddrArray).Elements {
				listBldr.Append(true)
				bldr.Field(i).(*array.BinaryBuilder).Append(e.Addr)
			}
		}
	}
	return bldr.NewRecord()
}