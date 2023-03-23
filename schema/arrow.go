package schema

import (
	"github.com/goccy/go-json"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/cloudquery/plugin-sdk/types"
)

const (
	MetadataPrimaryKey     = "cq:extension:primary_key"
	MetadataPrimaryKeyTrue = "true"
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
		typ = types.NewUUIDType()
	case TypeString:
		typ = arrow.BinaryTypes.String
	case TypeByteArray:
		typ = arrow.BinaryTypes.Binary
	case TypeStringArray:
		typ = arrow.ListOf(arrow.BinaryTypes.String)
	case TypeIntArray:
		typ = arrow.ListOf(arrow.PrimitiveTypes.Int64)
	case TypeTimestamp:
		typ = arrow.FixedWidthTypes.Timestamp_us
	case TypeJSON:
		typ = types.NewJSONType()
	case TypeUUIDArray:
		typ = arrow.ListOf(types.NewUUIDType())
	case TypeInet:
		typ = types.NewInetType()
	case TypeInetArray:
		typ = arrow.ListOf(types.NewInetType())
	case TypeCIDR:
		typ = types.NewInetType()
	case TypeCIDRArray:
		typ = arrow.ListOf(types.NewInetType())
	case TypeMacAddr:
		typ = types.NewMacType()
	case TypeMacAddrArray:
		typ = arrow.ListOf(types.NewMacType())
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

func CQTypesOneToRecord(mem memory.Allocator, c CQTypes, arrowSchema *arrow.Schema) arrow.Record {
	return CQTypesToRecord(mem, []CQTypes{c}, arrowSchema)
}

func CQTypesToRecord(mem memory.Allocator, c []CQTypes, arrowSchema *arrow.Schema) arrow.Record {
	bldr := array.NewRecordBuilder(mem, arrowSchema)
	fields := bldr.Fields()
	for i := range fields {
		for j := range c {
			switch c[j][i].Type() {
			case TypeBool:
				if c[j][i].(*Bool).Status == Present {
					bldr.Field(i).(*array.BooleanBuilder).Append(c[j][i].(*Bool).Bool)
				} else {
					bldr.Field(i).(*array.BooleanBuilder).AppendNull()
				}
			case TypeInt:
				if c[j][i].(*Int8).Status == Present {
					bldr.Field(i).(*array.Int64Builder).Append(c[j][i].(*Int8).Int)
				} else {
					bldr.Field(i).(*array.Int64Builder).AppendNull()
				}
			case TypeFloat:
				if c[j][i].(*Float8).Status == Present {
					bldr.Field(i).(*array.Float64Builder).Append(c[j][i].(*Float8).Float)
				} else {
					bldr.Field(i).(*array.Float64Builder).AppendNull()
				}
			case TypeString:
				if c[j][i].(*Text).Status == Present {
					bldr.Field(i).(*array.StringBuilder).Append(c[j][i].(*Text).Str)
				} else {
					bldr.Field(i).(*array.StringBuilder).AppendNull()
				}
			case TypeByteArray:
				if c[j][i].(*Bytea).Status == Present {
					bldr.Field(i).(*array.BinaryBuilder).Append(c[j][i].(*Bytea).Bytes)
				} else {
					bldr.Field(i).(*array.BinaryBuilder).AppendNull()
				}
			case TypeStringArray:
				if c[j][i].(*TextArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, str := range c[j][i].(*TextArray).Elements {
						listBldr.ValueBuilder().(*array.StringBuilder).Append(str.Str)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case TypeIntArray:
				if c[j][i].(*Int8Array).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*Int8Array).Elements {
						listBldr.ValueBuilder().(*array.Int64Builder).Append(e.Int)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case TypeTimestamp:
				if c[j][i].(*Timestamptz).Status == Present {
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(c[j][i].(*Timestamptz).Time.UnixMicro()))
				} else {
					bldr.Field(i).(*array.TimestampBuilder).AppendNull()
				}
			case TypeJSON:
				if c[j][i].(*JSON).Status == Present {
					var d any
					if err := json.Unmarshal(c[j][i].(*JSON).Bytes, &d); err != nil {
						panic(err)
					}
					bldr.Field(i).(*types.JSONBuilder).Append(d)
				} else {
					bldr.Field(i).(*types.JSONBuilder).AppendNull()
				}
			case TypeUUID:
				if c[j][i].(*UUID).Status == Present {
					bldr.Field(i).(*types.UUIDBuilder).Append(c[j][i].(*UUID).Bytes)
				} else {
					bldr.Field(i).(*types.UUIDBuilder).AppendNull()
				}
			case TypeUUIDArray:
				if c[j][i].(*UUIDArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*UUIDArray).Elements {
						listBldr.ValueBuilder().(*types.UUIDBuilder).Append(e.Bytes)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case TypeInet:
				if c[j][i].(*Inet).Status == Present {
					bldr.Field(i).(*types.InetBuilder).Append(*c[j][i].(*Inet).IPNet)
				} else {
					bldr.Field(i).(*types.InetBuilder).AppendNull()
				}
			case TypeInetArray:
				if c[j][i].(*InetArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*InetArray).Elements {
						listBldr.ValueBuilder().(*types.InetBuilder).Append(*e.IPNet)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case TypeCIDR:
				if c[j][i].(*CIDR).Status == Present {
					bldr.Field(i).(*types.InetBuilder).Append(*c[j][i].(*CIDR).IPNet)
				} else {
					bldr.Field(i).(*types.InetBuilder).AppendNull()
				}
			case TypeCIDRArray:
				if c[j][i].(*CIDRArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*CIDRArray).Elements {
						listBldr.ValueBuilder().(*types.InetBuilder).Append(*e.IPNet)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case TypeMacAddr:
				if c[j][i].(*Macaddr).Status == Present {
					bldr.Field(i).(*types.MacBuilder).Append(c[j][i].(*Macaddr).Addr)
				} else {
					bldr.Field(i).(*types.MacBuilder).AppendNull()
				}
			case TypeMacAddrArray:
				if c[j][i].(*MacaddrArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*MacaddrArray).Elements {
						listBldr.ValueBuilder().(*types.MacBuilder).Append(e.Addr)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			}
		}
	}

	return bldr.NewRecord()
}