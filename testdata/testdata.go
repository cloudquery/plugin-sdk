package testdata

import (
	"net"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/types"
	"github.com/google/uuid"
)

func TestSourceTable(name string) *schema.Table {
	return &schema.Table{
		Name:        name,
		Description: "Test table",
		Columns: schema.ColumnList{
			schema.CqIDColumn,
			schema.CqParentIDColumn,
			{
				Name: "bool",
				Type: schema.TypeBool,
			},
			{
				Name: "int",
				Type: schema.TypeInt,
			},
			{
				Name: "float",
				Type: schema.TypeFloat,
			},
			{
				Name:            "uuid",
				Type:            schema.TypeUUID,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name: "text",
				Type: schema.TypeString,
			},
			{
				Name: "text_with_null",
				Type: schema.TypeString,
			},
			{
				Name: "bytea",
				Type: schema.TypeByteArray,
			},
			{
				Name: "text_array",
				Type: schema.TypeStringArray,
			},
			{
				Name: "text_array_with_null",
				Type: schema.TypeStringArray,
			},
			{
				Name: "int_array",
				Type: schema.TypeIntArray,
			},
			{
				Name: "timestamp",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "json",
				Type: schema.TypeJSON,
			},
			{
				Name: "uuid_array",
				Type: schema.TypeUUIDArray,
			},
			{
				Name: "inet",
				Type: schema.TypeInet,
			},
			{
				Name: "inet_array",
				Type: schema.TypeInetArray,
			},
			{
				Name: "cidr",
				Type: schema.TypeCIDR,
			},
			{
				Name: "cidr_array",
				Type: schema.TypeCIDRArray,
			},
			{
				Name: "macaddr",
				Type: schema.TypeMacAddr,
			},
			{
				Name: "macaddr_array",
				Type: schema.TypeMacAddrArray,
			},
		},
	}
}

// TestTable returns a table with columns of all type. useful for destination testing purposes
func TestTable(name string) *schema.Table {
	sourceTable := TestSourceTable(name)
	sourceTable.Columns = append(schema.ColumnList{
		schema.CqSourceNameColumn,
		schema.CqSyncTimeColumn,
	}, sourceTable.Columns...)
	return sourceTable
}

type GenTestDataOptions struct {
	// SourceName is the name of the source to set in the source_name column.
	SourceName string
	// SyncTime is the time to set in the sync_time column.
	SyncTime time.Time
	// MaxRows is the number of rows to generate.
	// Rows alternate between not containing null values and containing only null values.
	// (Only columns that are nullable according to the schema will be null)
	MaxRows int
	// StableUUID is the UUID to use for all rows. If set to uuid.Nil, a new UUID will be generated
	StableUUID uuid.UUID
}

func GenTestData(mem memory.Allocator, sc *arrow.Schema, opts GenTestDataOptions) []arrow.Record {
	var records []arrow.Record
	for j := 0; j < opts.MaxRows; j++ {
		u := uuid.New()
		if opts.StableUUID != uuid.Nil {
			u = opts.StableUUID
		}
		nullRow := j%2 == 1
		bldr := array.NewRecordBuilder(mem, sc)
		for i, c := range sc.Fields() {
			if nullRow && c.Nullable && c.Name != schema.CqSourceNameColumn.Name && c.Name != schema.CqSyncTimeColumn.Name {
				bldr.Field(i).AppendNull()
				continue
			}
			if arrow.TypeEqual(c.Type, arrow.FixedWidthTypes.Boolean) {
				bldr.Field(i).(*array.BooleanBuilder).Append(true)
			} else if arrow.TypeEqual(c.Type, arrow.PrimitiveTypes.Int64) {
				bldr.Field(i).(*array.Int64Builder).Append(1)
			} else if arrow.TypeEqual(c.Type, arrow.PrimitiveTypes.Float64) {
				bldr.Field(i).(*array.Float64Builder).Append(1.1)
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.UUID) {
				bldr.Field(i).(*types.UUIDBuilder).Append(u)
			} else if arrow.TypeEqual(c.Type, arrow.BinaryTypes.String) {
				if c.Name == schema.CqSourceNameColumn.Name {
					bldr.Field(i).(*array.StringBuilder).AppendString(opts.SourceName)
				} else if c.Name == "text_with_null" {
					bldr.Field(i).(*array.StringBuilder).AppendString("AStringWith\x00NullBytes")
				} else {
					bldr.Field(i).(*array.StringBuilder).AppendString("AString")
				}
			} else if arrow.TypeEqual(c.Type, arrow.BinaryTypes.Binary) {
				bldr.Field(i).(*array.BinaryBuilder).Append([]byte{1, 2, 3})
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(arrow.BinaryTypes.String)) {
				if c.Name == "text_array_with_null" {
					bldr.Field(i).(*array.ListBuilder).Append(true)
					bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*array.StringBuilder).AppendString("test1")
					bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*array.StringBuilder).AppendString("test2\x00WithNull")
				} else {
					bldr.Field(i).(*array.ListBuilder).Append(true)
					bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*array.StringBuilder).AppendString("test1")
					bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*array.StringBuilder).AppendString("test2")
				}
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(arrow.PrimitiveTypes.Int64)) {
				bldr.Field(i).(*array.ListBuilder).Append(true)
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*array.Int64Builder).Append(1)
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*array.Int64Builder).Append(2)
			} else if arrow.TypeEqual(c.Type, arrow.FixedWidthTypes.Timestamp_us) {
				if c.Name == schema.CqSyncTimeColumn.Name {
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(opts.SyncTime.UTC().UnixMicro()))
				} else {
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(time.Now().UTC().UnixMicro()))
				}
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.JSON) {
				bldr.Field(i).(*types.JSONBuilder).Append(map[string]interface{}{"test": "test"})
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(types.ExtensionTypes.UUID)) {
				bldr.Field(i).(*array.ListBuilder).Append(true)
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.UUIDBuilder).Append(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.UUIDBuilder).Append(uuid.MustParse("00000000-0000-0000-0000-000000000002"))
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.Inet) {
				_, ipnet, err := net.ParseCIDR("192.0.2.0/24")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*types.InetBuilder).Append(*ipnet)
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(types.ExtensionTypes.Inet)) {
				bldr.Field(i).(*array.ListBuilder).Append(true)
				_, ipnet, err := net.ParseCIDR("192.0.2.1/24")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.InetBuilder).Append(*ipnet)
				_, ipnet, err = net.ParseCIDR("192.0.2.1/24")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.InetBuilder).Append(*ipnet)
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.Mac) {
				mac, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*types.MacBuilder).Append(mac)
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(types.ExtensionTypes.Mac)) {
				mac, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*array.ListBuilder).Append(true)
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.MacBuilder).Append(mac)
				mac, err = net.ParseMAC("11:22:33:44:55:66")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.MacBuilder).Append(mac)
			} else {
				panic("unknown type: " + c.Type.String() + " column: " + c.Name)
			}
		}
		records = append(records, bldr.NewRecord())
		bldr.Release()
	}
	return records
}
