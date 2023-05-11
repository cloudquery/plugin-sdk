package schema

import (
	"net"
	"sort"
	"strings"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v3/types"
	"github.com/google/uuid"
)

func TestSourceTable(name string) *Table {
	return &Table{
		Name:        name,
		Description: "Test table",
		Columns: ColumnList{
			CqIDColumn,
			CqParentIDColumn,
			{
				Name:            "uuid_pk",
				Type:            types.ExtensionTypes.UUID,
				CreationOptions: ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name:            "string_pk",
				Type:            arrow.BinaryTypes.String,
				CreationOptions: ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name: "bool",
				Type: arrow.FixedWidthTypes.Boolean,
			},
			{
				Name: "int",
				Type: arrow.PrimitiveTypes.Int64,
			},
			{
				Name: "float",
				Type: arrow.PrimitiveTypes.Float64,
			},
			{
				Name: "uuid",
				Type: types.ExtensionTypes.UUID,
			},
			{
				Name: "text",
				Type: arrow.BinaryTypes.String,
			},
			{
				Name: "text_with_null",
				Type: arrow.BinaryTypes.String,
			},
			{
				Name: "bytea",
				Type: arrow.BinaryTypes.Binary,
			},
			{
				Name: "text_array",
				Type: arrow.ListOf(arrow.BinaryTypes.String),
			},
			{
				Name: "text_array_with_null",
				Type: arrow.ListOf(arrow.BinaryTypes.String),
			},
			{
				Name: "int_array",
				Type: arrow.ListOf(arrow.PrimitiveTypes.Int64),
			},
			{
				Name: "timestamp",
				Type: arrow.FixedWidthTypes.Timestamp_us,
			},
			{
				Name: "json",
				Type: types.ExtensionTypes.JSON,
			},
			{
				Name: "uuid_array",
				Type: arrow.ListOf(types.ExtensionTypes.UUID),
			},
			{
				Name: "inet",
				Type: types.ExtensionTypes.Inet,
			},
			{
				Name: "inet_array",
				Type: arrow.ListOf(types.ExtensionTypes.Inet),
			},
			{
				Name: "cidr",
				Type: types.ExtensionTypes.Inet,
			},
			{
				Name: "cidr_array",
				Type: arrow.ListOf(types.ExtensionTypes.Inet),
			},
			{
				Name: "macaddr",
				Type: types.ExtensionTypes.Mac,
			},
			{
				Name: "macaddr_array",
				Type: arrow.ListOf(types.ExtensionTypes.Mac),
			},
			// This column is added to be able to test the Postgresql destination handling of reserved keywords in PKs
			{
				Name:            "user",
				Type:            schema.TypeString,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
		},
	}
}

func TestTableIncremental(name string) *Table {
	t := TestTable(name)
	t.IsIncremental = true
	return t
}

// TestTable returns a table with columns of all type. useful for destination testing purposes
func TestTable(name string) *Table {
	sourceTable := TestSourceTable(name)
	sourceTable.Columns = append(ColumnList{
		CqSourceNameColumn,
		CqSyncTimeColumn,
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

func GenTestData(table *Table, opts GenTestDataOptions) []arrow.Record {
	var records []arrow.Record
	sc := table.ToArrowSchema()
	for j := 0; j < opts.MaxRows; j++ {
		u := uuid.New()
		if opts.StableUUID != uuid.Nil {
			u = opts.StableUUID
		}
		nullRow := j%2 == 1
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		for i, c := range table.Columns {
			if nullRow && !c.CreationOptions.NotNull && !c.CreationOptions.PrimaryKey &&
				c.Name != CqSourceNameColumn.Name &&
				c.Name != CqSyncTimeColumn.Name &&
				c.Name != CqIDColumn.Name &&
				c.Name != CqParentIDColumn.Name {
				bldr.Field(i).AppendNull()
				continue
			}
			// nolint:gocritic
			if arrow.TypeEqual(c.Type, arrow.FixedWidthTypes.Boolean) {
				bldr.Field(i).(*array.BooleanBuilder).Append(true)
			} else if arrow.TypeEqual(c.Type, arrow.PrimitiveTypes.Int64) {
				bldr.Field(i).(*array.Int64Builder).Append(1)
			} else if arrow.TypeEqual(c.Type, arrow.PrimitiveTypes.Float64) {
				bldr.Field(i).(*array.Float64Builder).Append(1.1)
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.UUID) {
				bldr.Field(i).(*types.UUIDBuilder).Append(u)
			} else if arrow.TypeEqual(c.Type, arrow.BinaryTypes.String) {
				switch c.Name {
				case CqSourceNameColumn.Name:
					bldr.Field(i).(*array.StringBuilder).AppendString(opts.SourceName)
				case "text_with_null":
					bldr.Field(i).(*array.StringBuilder).AppendString("AStringWith\x00NullBytes")
				default:
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
				if c.Name == CqSyncTimeColumn.Name {
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(opts.SyncTime.UTC().Truncate(time.Millisecond).UnixMicro()))
				} else {
					t := time.Now()
					if !opts.StableTime.IsZero() {
						t = opts.StableTime
					}
					bldr.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(t.UTC().Truncate(time.Millisecond).UnixMicro()))
				}
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.JSON) {
				bldr.Field(i).(*types.JSONBuilder).Append(map[string]any{"test": "test"})
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(types.ExtensionTypes.UUID)) {
				bldr.Field(i).(*array.ListBuilder).Append(true)
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.UUIDBuilder).Append(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.UUIDBuilder).Append(uuid.MustParse("00000000-0000-0000-0000-000000000002"))
			} else if arrow.TypeEqual(c.Type, types.ExtensionTypes.Inet) {
				_, ipnet, err := net.ParseCIDR("192.0.2.0/24")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*types.InetBuilder).Append(ipnet)
			} else if arrow.TypeEqual(c.Type, arrow.ListOf(types.ExtensionTypes.Inet)) {
				bldr.Field(i).(*array.ListBuilder).Append(true)
				_, ipnet, err := net.ParseCIDR("192.0.2.1/24")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.InetBuilder).Append(ipnet)
				_, ipnet, err = net.ParseCIDR("192.0.2.1/24")
				if err != nil {
					panic(err)
				}
				bldr.Field(i).(*array.ListBuilder).ValueBuilder().(*types.InetBuilder).Append(ipnet)
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
	if indices := sc.FieldIndices(CqIDColumn.Name); len(indices) > 0 {
		cqIDIndex := indices[0]
		sort.Slice(records, func(i, j int) bool {
			firstUUID := records[i].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			secondUUID := records[j].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			return strings.Compare(firstUUID, secondUUID) < 0
		})
	}
	return records
}
