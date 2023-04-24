package testdata

import (
	"fmt"
	"sort"
	"strings"
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

func TestTableIncremental(name string) *schema.Table {
	t := TestTable(name)
	t.IsIncremental = true
	return t
}

// TestTable returns a table with columns of all CQ types. Useful for destination testing purposes
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

func GenTestData(sc *arrow.Schema, opts GenTestDataOptions) []arrow.Record {
	var records []arrow.Record
	for j := 0; j < opts.MaxRows; j++ {
		nullRow := j%2 == 1
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		for i, c := range sc.Fields() {
			if nullRow && c.Nullable && !schema.IsPk(c) &&
				c.Name != schema.CqSourceNameColumn.Name &&
				c.Name != schema.CqSyncTimeColumn.Name &&
				c.Name != schema.CqIDField.Name &&
				c.Name != schema.CqParentIDColumn.Name {
				bldr.Field(i).AppendNull()
				continue
			}
			example := getExampleJSON(c.Name, c.Type, opts)
			l := `[` + example + `]`
			err := bldr.Field(i).UnmarshalJSON([]byte(l))
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal json `%v` for column %v: %v", l, c.Name, err))
			}
		}
		records = append(records, bldr.NewRecord())
		bldr.Release()
	}
	if indices := sc.FieldIndices(schema.CqIDColumn.Name); len(indices) > 0 {
		cqIDIndex := indices[0]
		sort.Slice(records, func(i, j int) bool {
			firstUUID := records[i].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			secondUUID := records[j].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			return strings.Compare(firstUUID, secondUUID) < 0
		})
	}
	return records
}

func getExampleJSON(colName string, dataType arrow.DataType, opts GenTestDataOptions) string {
	// handle lists (including maps)
	if arrow.IsListLike(dataType.ID()) {
		if dataType.ID() == arrow.MAP {
			k := getExampleJSON(colName, dataType.(*arrow.MapType).KeyType(), opts)
			v := getExampleJSON(colName, dataType.(*arrow.MapType).ValueType().Field(1).Type, opts)
			return fmt.Sprintf(`[{"key": %s,"value": %s}]`, k, v)
		}
		inner := dataType.(*arrow.ListType).Elem()
		return `[` + getExampleJSON(colName, inner, opts) + `]`
	}
	// handle extension types
	if arrow.TypeEqual(dataType, types.ExtensionTypes.UUID) {
		u := uuid.New()
		if opts.StableUUID != uuid.Nil {
			u = opts.StableUUID
		}
		return `"` + u.String() + `"`
	}
	if arrow.TypeEqual(dataType, types.ExtensionTypes.JSON) {
		return `"{\"test\": \"test\"}"`
	}
	if arrow.TypeEqual(dataType, types.ExtensionTypes.Inet) {
		return `"192.0.2.0/24"`
	}
	if arrow.TypeEqual(dataType, types.ExtensionTypes.Mac) {
		return `"aa:bb:cc:dd:ee:ff"`
	}

	// handle integers
	if arrow.IsInteger(dataType.ID()) {
		return "-1"
	}

	// handle unsigned integers
	if arrow.IsUnsignedInteger(dataType.ID()) {
		return "1"
	}

	// handle floats
	if arrow.IsFloating(dataType.ID()) {
		return "1.1"
	}

	// handle decimals
	if arrow.IsDecimal(dataType.ID()) {
		return "1.1"
	}

	// handle booleans
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Boolean) {
		return "true"
	}

	// handle strings
	stringTypes := []arrow.DataType{
		arrow.BinaryTypes.String,
		arrow.BinaryTypes.LargeString,
	}
	for _, stringType := range stringTypes {
		if arrow.TypeEqual(dataType, stringType) {
			if colName == schema.CqSourceNameColumn.Name {
				return `"` + opts.SourceName + `"`
			} else if colName == "text_with_null" {
				return `"AStringWith` + "\x00" + `NullBytes"`
			}
			return `"AString"`
		}
	}

	// handle binary types
	binaryTypes := []arrow.DataType{
		arrow.BinaryTypes.Binary,
		arrow.BinaryTypes.LargeBinary,
	}
	for _, binaryType := range binaryTypes {
		if arrow.TypeEqual(dataType, binaryType) {
			return `"AQIDBA=="` // base64 encoded 0x01, 0x02, 0x03, 0x04
		}
	}

	// handle structs
	if dataType.ID() == arrow.STRUCT {
		var fields []string
		for _, field := range dataType.(*arrow.StructType).Fields() {
			v := getExampleJSON(field.Name, field.Type, opts)
			fields = append(fields, fmt.Sprintf(`"%s": %v`, field.Name, v))
		}
		return `{` + strings.Join(fields, ",") + `}`
	}

	// handle timestamp types
	timestampTypes := []arrow.DataType{
		arrow.FixedWidthTypes.Timestamp_s,
		arrow.FixedWidthTypes.Timestamp_ms,
		arrow.FixedWidthTypes.Timestamp_us,
		arrow.FixedWidthTypes.Timestamp_ns,
		arrow.FixedWidthTypes.Time32s,
		arrow.FixedWidthTypes.Time32ms,
		arrow.FixedWidthTypes.Time64us,
		arrow.FixedWidthTypes.Time64ns,
	}
	for _, timestampType := range timestampTypes {
		if arrow.TypeEqual(dataType, timestampType) {
			t := time.Now()
			if colName == schema.CqSyncTimeColumn.Name {
				t = opts.SyncTime.UTC()
			} else if !opts.StableTime.IsZero() {
				t = opts.StableTime
			}
			switch timestampType {
			case arrow.FixedWidthTypes.Timestamp_s:
				return fmt.Sprintf("%d", t.Unix())
			case arrow.FixedWidthTypes.Timestamp_ms:
				return fmt.Sprintf("%d", t.UnixMilli())
			case arrow.FixedWidthTypes.Timestamp_us:
				return fmt.Sprintf("%d", t.UnixMicro())
			case arrow.FixedWidthTypes.Timestamp_ns:
				return fmt.Sprintf("%d", t.UnixNano())
			case arrow.FixedWidthTypes.Time32s:
				return fmt.Sprintf("%d", t.Unix())
			case arrow.FixedWidthTypes.Time32ms:
				return fmt.Sprintf("%d", t.UnixMilli())
			case arrow.FixedWidthTypes.Time64us:
				return fmt.Sprintf("%d", t.UnixMicro())
			case arrow.FixedWidthTypes.Time64ns:
				return fmt.Sprintf("%d", t.UnixNano())
			default:
				panic("unhandled timestamp type: " + timestampType.Name())
			}
		}
	}

	// handle date types
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Date32) {
		return `1682088351`
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Date64) {
		return `1682088344338`
	}

	// handle duration and interval types
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.DayTimeInterval) {
		return `{"days": 1, "milliseconds": 1}`
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.MonthInterval) {
		return `{"months": 1}`
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.MonthDayNanoInterval) {
		return `{"months": 1, "days": 1, "nanoseconds": 1}`
	}
	durationTypes := []arrow.DataType{
		arrow.FixedWidthTypes.Duration_s,
		arrow.FixedWidthTypes.Duration_ms,
		arrow.FixedWidthTypes.Duration_us,
		arrow.FixedWidthTypes.Duration_ns,
	}
	for _, durationType := range durationTypes {
		if arrow.TypeEqual(dataType, durationType) {
			return `123456789`
		}
	}

	panic("unknown type: " + dataType.String() + " column: " + colName)
}

// GenTestDataV1 does approximately the same job as GenTestData, however, it's intended for simpler use-cases.
// Deprecated: Will be removed in future release.
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
