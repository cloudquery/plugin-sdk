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
				Name:            "uuid_pk",
				Type:            schema.TypeUUID,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name:            "string_pk",
				Type:            schema.TypeString,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
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
				Name: "uuid",
				Type: schema.TypeUUID,
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
	// StableTime is the time to use for all rows other than sync time. If set to time.Time{}, a new time will be generated
	StableTime time.Time
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
			if nullRow && c.Nullable && !schema.IsPk(c) &&
				c.Name != schema.CqSourceNameColumn.Name &&
				c.Name != schema.CqSyncTimeColumn.Name &&
				c.Name != schema.CqIDField.Name &&
				c.Name != schema.CqParentIDColumn.Name {
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
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(opts.SyncTime.UTC().Truncate(time.Millisecond).UnixMicro()))
				} else {
					t := time.Now()
					if !opts.StableTime.IsZero() {
						t = opts.StableTime
					}
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(t.UTC().Truncate(time.Millisecond).UnixMicro()))
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

// GenTestDataV1 does approximately the same job as GenTestData, however, it's intended for simpler use-cases.
// Deprecated. Will be removed in future release.
func GenTestDataV1(table *schema.Table) schema.CQTypes {
	data := make(schema.CQTypes, len(table.Columns))
	for i, c := range table.Columns {
		switch c.Type {
		case schema.TypeBool:
			data[i] = &schema.Bool{
				Bool:   true,
				Status: schema.Present,
			}
		case schema.TypeInt:
			data[i] = &schema.Int8{
				Int:    1,
				Status: schema.Present,
			}
		case schema.TypeFloat:
			data[i] = &schema.Float8{
				Float:  1.1,
				Status: schema.Present,
			}
		case schema.TypeUUID:
			uuidColumn := &schema.UUID{}
			if err := uuidColumn.Set(uuid.NewString()); err != nil {
				panic(err)
			}
			data[i] = uuidColumn
		case schema.TypeString:
			if c.Name == "text_with_null" {
				data[i] = &schema.Text{
					Str:    "AStringWith\x00NullBytes",
					Status: schema.Present,
				}
			} else {
				data[i] = &schema.Text{
					Str:    "test",
					Status: schema.Present,
				}
			}
		case schema.TypeByteArray:
			data[i] = &schema.Bytea{
				Bytes:  []byte{1, 2, 3},
				Status: schema.Present,
			}
		case schema.TypeStringArray:
			if c.Name == "text_array_with_null" {
				data[i] = &schema.TextArray{
					Elements: []schema.Text{{Str: "test1", Status: schema.Present}, {Str: "test2\x00WithNull", Status: schema.Present}},
					Status:   schema.Present,
				}
			} else {
				data[i] = &schema.TextArray{
					Elements: []schema.Text{{Str: "test1", Status: schema.Present}, {Str: "test2", Status: schema.Present}},
					Status:   schema.Present,
				}
			}

		case schema.TypeIntArray:
			data[i] = &schema.Int8Array{
				Elements: []schema.Int8{{Int: 1, Status: schema.Present}, {Int: 2, Status: schema.Present}},
				Status:   schema.Present,
			}
		case schema.TypeTimestamp:
			data[i] = &schema.Timestamptz{
				Time:   time.Now().UTC().Round(time.Second),
				Status: schema.Present,
			}
		case schema.TypeJSON:
			data[i] = &schema.JSON{
				Bytes:  []byte(`{"test": "test"}`),
				Status: schema.Present,
			}
		case schema.TypeUUIDArray:
			uuidArrayColumn := &schema.UUIDArray{}
			if err := uuidArrayColumn.Set([]string{"00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}); err != nil {
				panic(err)
			}
			data[i] = uuidArrayColumn
		case schema.TypeInet:
			inetColumn := &schema.Inet{}
			if err := inetColumn.Set("192.0.2.0/24"); err != nil {
				panic(err)
			}
			data[i] = inetColumn
		case schema.TypeInetArray:
			inetArrayColumn := &schema.InetArray{}
			if err := inetArrayColumn.Set([]string{"192.0.2.1/24", "192.0.2.1/24"}); err != nil {
				panic(err)
			}
			data[i] = inetArrayColumn
		case schema.TypeCIDR:
			cidrColumn := &schema.CIDR{}
			if err := cidrColumn.Set("192.0.2.1"); err != nil {
				panic(err)
			}
			data[i] = cidrColumn
		case schema.TypeCIDRArray:
			cidrArrayColumn := &schema.CIDRArray{}
			if err := cidrArrayColumn.Set([]string{"192.0.2.1", "192.0.2.1"}); err != nil {
				panic(err)
			}
			data[i] = cidrArrayColumn
		case schema.TypeMacAddr:
			macaddrColumn := &schema.Macaddr{}
			if err := macaddrColumn.Set("aa:bb:cc:dd:ee:ff"); err != nil {
				panic(err)
			}
			data[i] = macaddrColumn
		case schema.TypeMacAddrArray:
			macaddrArrayColumn := &schema.MacaddrArray{}
			if err := macaddrArrayColumn.Set([]string{"aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"}); err != nil {
				panic(err)
			}
			data[i] = macaddrArrayColumn
		default:
			panic("unknown type" + c.Type.String())
		}
	}
	return data
}
