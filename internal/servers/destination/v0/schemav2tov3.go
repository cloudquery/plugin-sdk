package destination

import (
	"encoding/json"
	"strings"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/apache/arrow/go/v17/arrow/memory"
	schemav2 "github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

func TablesV2ToV3(tables schemav2.Tables) schema.Tables {
	res := make(schema.Tables, len(tables))
	for i, t := range tables {
		res[i] = TableV2ToV3(t)
	}
	return res
}

func TableV2ToV3(table *schemav2.Table) *schema.Table {
	newTable := &schema.Table{
		Name:          table.Name,
		Description:   table.Description,
		Columns:       ColumnsV2ToV3(table.Columns),
		IgnoreInTests: table.IgnoreInTests,
		IsIncremental: table.IsIncremental,
	}
	if len(table.Relations) > 0 {
		newTable.Relations = TablesV2ToV3(table.Relations)
	}
	return newTable
}

func ColumnsV2ToV3(columns []schemav2.Column) []schema.Column {
	res := make([]schema.Column, len(columns))
	for i, c := range columns {
		res[i] = ColumnV2ToV3(c)
	}
	return res
}

func ColumnV2ToV3(column schemav2.Column) schema.Column {
	return schema.Column{
		Name:           column.Name,
		Description:    column.Description,
		Type:           TypeV2ToV3(column.Type),
		NotNull:        column.CreationOptions.NotNull,
		Unique:         column.CreationOptions.Unique,
		PrimaryKey:     column.CreationOptions.PrimaryKey,
		IncrementalKey: column.CreationOptions.IncrementalKey,
		IgnoreInTests:  column.IgnoreInTests,
	}
}

func TypeV2ToV3(dataType schemav2.ValueType) arrow.DataType {
	var typ arrow.DataType

	switch dataType {
	case schemav2.TypeBool:
		return arrow.FixedWidthTypes.Boolean
	case schemav2.TypeInt:
		return arrow.PrimitiveTypes.Int64
	case schemav2.TypeFloat:
		return arrow.PrimitiveTypes.Float64
	case schemav2.TypeUUID:
		return types.ExtensionTypes.UUID
	case schemav2.TypeString:
		return arrow.BinaryTypes.String
	case schemav2.TypeByteArray:
		return arrow.BinaryTypes.Binary
	case schemav2.TypeStringArray:
		return arrow.ListOf(arrow.BinaryTypes.String)
	case schemav2.TypeIntArray:
		return arrow.ListOf(arrow.PrimitiveTypes.Int64)
	case schemav2.TypeTimestamp:
		return arrow.FixedWidthTypes.Timestamp_us
	case schemav2.TypeJSON:
		return types.ExtensionTypes.JSON
	case schemav2.TypeUUIDArray:
		return arrow.ListOf(types.ExtensionTypes.UUID)
	case schemav2.TypeInet:
		return types.ExtensionTypes.Inet
	case schemav2.TypeInetArray:
		return arrow.ListOf(types.ExtensionTypes.Inet)
	case schemav2.TypeCIDR:
		return types.ExtensionTypes.Inet
	case schemav2.TypeCIDRArray:
		return arrow.ListOf(types.ExtensionTypes.Inet)
	case schemav2.TypeMacAddr:
		return types.ExtensionTypes.MAC
	case schemav2.TypeMacAddrArray:
		return arrow.ListOf(types.ExtensionTypes.MAC)
	default:
		panic("unknown type " + typ.Name())
	}
}

func CQTypesOneToRecord(mem memory.Allocator, c schemav2.CQTypes, arrowSchema *arrow.Schema) arrow.Record {
	return CQTypesToRecord(mem, []schemav2.CQTypes{c}, arrowSchema)
}

func CQTypesToRecord(mem memory.Allocator, c []schemav2.CQTypes, arrowSchema *arrow.Schema) arrow.Record {
	bldr := array.NewRecordBuilder(mem, arrowSchema)
	fields := bldr.Fields()
	for i := range fields {
		for j := range c {
			switch c[j][i].Type() {
			case schemav2.TypeBool:
				if c[j][i].(*schemav2.Bool).Status == schemav2.Present {
					bldr.Field(i).(*array.BooleanBuilder).Append(c[j][i].(*schemav2.Bool).Bool)
				} else {
					bldr.Field(i).(*array.BooleanBuilder).AppendNull()
				}
			case schemav2.TypeInt:
				if c[j][i].(*schemav2.Int8).Status == schemav2.Present {
					bldr.Field(i).(*array.Int64Builder).Append(c[j][i].(*schemav2.Int8).Int)
				} else {
					bldr.Field(i).(*array.Int64Builder).AppendNull()
				}
			case schemav2.TypeFloat:
				if c[j][i].(*schemav2.Float8).Status == schemav2.Present {
					bldr.Field(i).(*array.Float64Builder).Append(c[j][i].(*schemav2.Float8).Float)
				} else {
					bldr.Field(i).(*array.Float64Builder).AppendNull()
				}
			case schemav2.TypeString:
				if c[j][i].(*schemav2.Text).Status == schemav2.Present {
					// In the new type system we wont allow null string as they are not valid utf-8
					// https://github.com/apache/arrow/pull/35161#discussion_r1170516104
					bldr.Field(i).(*array.StringBuilder).Append(strings.ReplaceAll(c[j][i].(*schemav2.Text).Str, "\x00", ""))
				} else {
					bldr.Field(i).(*array.StringBuilder).AppendNull()
				}
			case schemav2.TypeByteArray:
				if c[j][i].(*schemav2.Bytea).Status == schemav2.Present {
					bldr.Field(i).(*array.BinaryBuilder).Append(c[j][i].(*schemav2.Bytea).Bytes)
				} else {
					bldr.Field(i).(*array.BinaryBuilder).AppendNull()
				}
			case schemav2.TypeStringArray:
				if c[j][i].(*schemav2.TextArray).Status == schemav2.Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, str := range c[j][i].(*schemav2.TextArray).Elements {
						listBldr.ValueBuilder().(*array.StringBuilder).Append(strings.ReplaceAll(str.Str, "\x00", ""))
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case schemav2.TypeIntArray:
				if c[j][i].(*schemav2.Int8Array).Status == schemav2.Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*schemav2.Int8Array).Elements {
						listBldr.ValueBuilder().(*array.Int64Builder).Append(e.Int)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case schemav2.TypeTimestamp:
				if c[j][i].(*schemav2.Timestamptz).Status == schemav2.Present {
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(c[j][i].(*schemav2.Timestamptz).Time.UnixMicro()))
				} else {
					bldr.Field(i).(*array.TimestampBuilder).AppendNull()
				}
			case schemav2.TypeJSON:
				if c[j][i].(*schemav2.JSON).Status == schemav2.Present {
					var d any
					if err := json.Unmarshal(c[j][i].(*schemav2.JSON).Bytes, &d); err != nil {
						panic(err)
					}
					bldr.Field(i).(*types.JSONBuilder).Append(d)
				} else {
					bldr.Field(i).(*types.JSONBuilder).AppendNull()
				}
			case schemav2.TypeUUID:
				if c[j][i].(*schemav2.UUID).Status == schemav2.Present {
					bldr.Field(i).(*types.UUIDBuilder).Append(c[j][i].(*schemav2.UUID).Bytes)
				} else {
					bldr.Field(i).(*types.UUIDBuilder).AppendNull()
				}
			case schemav2.TypeUUIDArray:
				if c[j][i].(*schemav2.UUIDArray).Status == schemav2.Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*schemav2.UUIDArray).Elements {
						listBldr.ValueBuilder().(*types.UUIDBuilder).Append(e.Bytes)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case schemav2.TypeInet:
				if c[j][i].(*schemav2.Inet).Status == schemav2.Present {
					bldr.Field(i).(*types.InetBuilder).Append(c[j][i].(*schemav2.Inet).IPNet)
				} else {
					bldr.Field(i).(*types.InetBuilder).AppendNull()
				}
			case schemav2.TypeInetArray:
				if c[j][i].(*schemav2.InetArray).Status == schemav2.Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*schemav2.InetArray).Elements {
						listBldr.ValueBuilder().(*types.InetBuilder).Append(e.IPNet)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case schemav2.TypeCIDR:
				if c[j][i].(*schemav2.CIDR).Status == schemav2.Present {
					bldr.Field(i).(*types.InetBuilder).Append(c[j][i].(*schemav2.CIDR).IPNet)
				} else {
					bldr.Field(i).(*types.InetBuilder).AppendNull()
				}
			case schemav2.TypeCIDRArray:
				if c[j][i].(*schemav2.CIDRArray).Status == schemav2.Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*schemav2.CIDRArray).Elements {
						listBldr.ValueBuilder().(*types.InetBuilder).Append(e.IPNet)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case schemav2.TypeMacAddr:
				if c[j][i].(*schemav2.Macaddr).Status == schemav2.Present {
					bldr.Field(i).(*types.MACBuilder).Append(c[j][i].(*schemav2.Macaddr).Addr)
				} else {
					bldr.Field(i).(*types.MACBuilder).AppendNull()
				}
			case schemav2.TypeMacAddrArray:
				if c[j][i].(*schemav2.MacaddrArray).Status == schemav2.Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*schemav2.MacaddrArray).Elements {
						listBldr.ValueBuilder().(*types.MACBuilder).Append(e.Addr)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			}
		}
	}

	return bldr.NewRecord()
}
