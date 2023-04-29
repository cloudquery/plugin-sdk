package schema

import (
	"strings"

	"github.com/goccy/go-json"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v2/types"
)

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
					// In the new type system we wont allow null string as they are not valid utf-8
					// https://github.com/apache/arrow/pull/35161#discussion_r1170516104
					bldr.Field(i).(*array.StringBuilder).Append(strings.ReplaceAll(c[j][i].(*Text).Str, "\x00", ""))
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
						listBldr.ValueBuilder().(*array.StringBuilder).Append(strings.ReplaceAll(str.Str, "\x00", ""))
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
					bldr.Field(i).(*types.InetBuilder).Append(c[j][i].(*Inet).IPNet)
				} else {
					bldr.Field(i).(*types.InetBuilder).AppendNull()
				}
			case TypeInetArray:
				if c[j][i].(*InetArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*InetArray).Elements {
						listBldr.ValueBuilder().(*types.InetBuilder).Append(e.IPNet)
					}
				} else {
					bldr.Field(i).(*array.ListBuilder).AppendNull()
				}
			case TypeCIDR:
				if c[j][i].(*CIDR).Status == Present {
					bldr.Field(i).(*types.InetBuilder).Append(c[j][i].(*CIDR).IPNet)
				} else {
					bldr.Field(i).(*types.InetBuilder).AppendNull()
				}
			case TypeCIDRArray:
				if c[j][i].(*CIDRArray).Status == Present {
					listBldr := bldr.Field(i).(*array.ListBuilder)
					listBldr.Append(true)
					for _, e := range c[j][i].(*CIDRArray).Elements {
						listBldr.ValueBuilder().(*types.InetBuilder).Append(e.IPNet)
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
