package schema

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v3/types"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

// TestSourceOptions controls which types are included by TestSourceColumns.
type TestSourceOptions struct {
	SkipLists      bool // lists of all primitive types. Lists that were supported by CQTypes are always included.
	SkipTimestamps bool // timestamp types. Microsecond timestamp is always be included, regardless of this setting.
	SkipDates      bool
	SkipMaps       bool
	SkipStructs    bool
	SkipIntervals  bool
	SkipDurations  bool
	SkipTimes      bool // time of day types
	SkipLargeTypes bool // e.g. large binary, large string
	TimePrecision  time.Duration
}

// TestSourceColumns returns columns for all Arrow types and composites thereof. TestSourceOptions controls
// which types are included.
func TestSourceColumns(testOpts TestSourceOptions) []Column {
	// cq columns
	var cqColumns []Column
	cqColumns = append(cqColumns, Column{Name: CqIDColumn.Name, Type: types.NewUUIDType(), NotNull: true, Unique: true, PrimaryKey: true})
	cqColumns = append(cqColumns, Column{Name: CqParentIDColumn.Name, Type: types.NewUUIDType()})

	var basicColumns []Column
	basicColumns = append(basicColumns, primitiveColumns()...)
	basicColumns = append(basicColumns, binaryColumns()...)
	basicColumns = append(basicColumns, fixedWidthColumns()...)

	// add extensions
	basicColumns = append(basicColumns, Column{Name: "uuid", Type: types.NewUUIDType()})
	basicColumns = append(basicColumns, Column{Name: "inet", Type: types.NewInetType()})
	basicColumns = append(basicColumns, Column{Name: "mac", Type: types.NewMACType()})

	// sort and remove duplicates (e.g. date32 and date64 appear twice)
	sort.Slice(basicColumns, func(i, j int) bool {
		return basicColumns[i].Name < basicColumns[j].Name
	})
	basicColumns = removeDuplicates(basicColumns)

	// we don't support float16 right now
	basicColumns = removeColumnsByType(basicColumns, arrow.FLOAT16)

	if testOpts.SkipTimestamps {
		// for backwards-compatibility, microsecond timestamps are not removed here
		basicColumns = removeColumnsByDataType(basicColumns, &arrow.TimestampType{Unit: arrow.Second, TimeZone: "UTC"})
		basicColumns = removeColumnsByDataType(basicColumns, &arrow.TimestampType{Unit: arrow.Millisecond, TimeZone: "UTC"})
		basicColumns = removeColumnsByDataType(basicColumns, &arrow.TimestampType{Unit: arrow.Nanosecond, TimeZone: "UTC"})
	}
	if testOpts.SkipDates {
		basicColumns = removeColumnsByType(basicColumns, arrow.DATE32, arrow.DATE64)
	}
	if testOpts.SkipTimes {
		basicColumns = removeColumnsByType(basicColumns, arrow.TIME32, arrow.TIME64)
	}
	if testOpts.SkipIntervals {
		basicColumns = removeColumnsByType(basicColumns, arrow.INTERVAL_DAY_TIME, arrow.INTERVAL_MONTHS, arrow.INTERVAL_MONTH_DAY_NANO)
	}
	if testOpts.SkipDurations {
		basicColumns = removeColumnsByType(basicColumns, arrow.DURATION)
	}
	if testOpts.SkipLargeTypes {
		basicColumns = removeColumnsByType(basicColumns, arrow.LARGE_BINARY, arrow.LARGE_STRING)
	}

	var compositeColumns []Column

	// we don't need to include lists of binary or large binary right now; probably no destinations or sources need to support that
	basicColumnsWithExclusions := removeColumnsByType(basicColumns, arrow.BINARY, arrow.LARGE_BINARY)
	if testOpts.SkipLists {
		// only include lists that were originally supported by CQTypes
		cqListColumns := []Column{
			{Name: "string", Type: arrow.BinaryTypes.String},
			{Name: "uuid", Type: types.NewUUIDType()},
			{Name: "inet", Type: types.NewInetType()},
			{Name: "mac", Type: types.NewMACType()},
		}
		compositeColumns = append(compositeColumns, listOfColumns(cqListColumns)...)
	} else {
		compositeColumns = append(compositeColumns, listOfColumns(basicColumnsWithExclusions)...)
	}

	// if !opts.SkipMaps {
	// 	compositeColumns = append(compositeColumns, mapOfColumns(basicColumnsWithExclusions)...)
	// }

	// add JSON later, we don't want to include it as a list or map right now (it causes complications with JSON unmarshalling)
	basicColumns = append(basicColumns, Column{Name: "json", Type: types.NewJSONType()})
	basicColumns = append(basicColumns, Column{Name: "json_array", Type: types.NewJSONType()}) // GenTestData knows to populate this with a JSON array

	if !testOpts.SkipStructs {
		// struct with all the types
		compositeColumns = append(compositeColumns, Column{Name: "struct", Type: arrow.StructOf(columnsToFields(basicColumns...)...)})

		// struct with nested struct
		compositeColumns = append(compositeColumns, Column{Name: "nested_struct", Type: arrow.StructOf(arrow.Field{Name: "inner", Type: arrow.StructOf(columnsToFields(basicColumns...)...)})})
	}

	allColumns := append(append(cqColumns, basicColumns...), compositeColumns...)
	return allColumns
}

// primitiveColumns returns a list of primitive columns as defined by Arrow types.
func primitiveColumns() []Column {
	primitiveTypesValue := reflect.ValueOf(arrow.PrimitiveTypes)
	primitiveTypesType := reflect.TypeOf(arrow.PrimitiveTypes)
	columns := make([]Column, primitiveTypesType.NumField())
	for i := 0; i < primitiveTypesType.NumField(); i++ {
		fieldName := primitiveTypesType.Field(i).Name
		dataType := primitiveTypesValue.FieldByName(fieldName).Interface().(arrow.DataType)
		columns[i] = Column{Name: strings.ToLower(fieldName), Type: dataType}
	}
	return columns
}

// binaryColumns returns a list of binary columns as defined by Arrow types.
func binaryColumns() []Column {
	binaryTypesValue := reflect.ValueOf(arrow.BinaryTypes)
	binaryTypesType := reflect.TypeOf(arrow.BinaryTypes)
	columns := make([]Column, binaryTypesType.NumField())
	for i := 0; i < binaryTypesType.NumField(); i++ {
		fieldName := binaryTypesType.Field(i).Name
		dataType := binaryTypesValue.FieldByName(fieldName).Interface().(arrow.DataType)
		columns[i] = Column{Name: strings.ToLower(fieldName), Type: dataType}
	}
	return columns
}

// fixedWidthColumns returns a list of fixed width columns as defined by Arrow types.
func fixedWidthColumns() []Column {
	fixedWidthTypesValue := reflect.ValueOf(arrow.FixedWidthTypes)
	fixedWidthTypesType := reflect.TypeOf(arrow.FixedWidthTypes)
	columns := make([]Column, fixedWidthTypesType.NumField())
	for i := 0; i < fixedWidthTypesType.NumField(); i++ {
		fieldName := fixedWidthTypesType.Field(i).Name
		dataType := fixedWidthTypesValue.FieldByName(fieldName).Interface().(arrow.DataType)
		columns[i] = Column{Name: strings.ToLower(fieldName), Type: dataType}
	}
	return columns
}

func removeDuplicates(columns []Column) []Column {
	newColumns := make([]Column, 0, len(columns))
	seen := map[string]struct{}{}
	for _, c := range columns {
		if _, ok := seen[c.Name]; ok {
			continue
		}
		newColumns = append(newColumns, c)
		seen[c.Name] = struct{}{}
	}
	return slices.Clip(newColumns)
}

func removeColumnsByType(columns []Column, t ...arrow.Type) []Column {
	var newColumns []Column
	for _, c := range columns {
		shouldRemove := false
		for _, d := range t {
			if c.Type.ID() == d {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			newColumns = append(newColumns, c)
		}
	}
	return newColumns
}

func removeColumnsByDataType(columns []Column, dt ...arrow.DataType) []Column {
	var newColumns []Column
	for _, c := range columns {
		shouldRemove := false
		for _, d := range dt {
			if arrow.TypeEqual(c.Type, d) {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			newColumns = append(newColumns, c)
		}
	}
	return newColumns
}

// listOfColumns returns a list of columns that are lists of the given columns.
func listOfColumns(baseColumns []Column) []Column {
	columns := make([]Column, len(baseColumns))
	for i := 0; i < len(baseColumns); i++ {
		columns[i] = Column{Name: baseColumns[i].Name + "_list", Type: arrow.ListOf(baseColumns[i].Type)}
	}
	return columns
}

// mapOfColumns returns a list of columns that are maps of the given columns.
// nolint:unused
func mapOfColumns(baseColumns []Column) []Column {
	columns := make([]Column, len(baseColumns))
	for i := 0; i < len(baseColumns); i++ {
		columns[i] = Column{Name: baseColumns[i].Name + "_map", Type: arrow.MapOf(baseColumns[i].Type, baseColumns[i].Type)}
	}
	return columns
}

func columnsToFields(columns ...Column) []arrow.Field {
	fields := make([]arrow.Field, len(columns))
	for i := range columns {
		fields[i] = arrow.Field{
			Name: columns[i].Name,
			Type: columns[i].Type,
		}
	}
	return fields
}

// var PKColumnNames = []string{"uuid_pk"}

// TestTable returns a table with columns of all types. Useful for destination testing purposes
func TestTable(name string, testOpts TestSourceOptions) *Table {
	var columns []Column
	// columns = append(columns, Column{Name: "uuid", Type: types.NewUUIDType()})
	// columns = append(columns, Column{Name: "string_pk", Type: arrow.BinaryTypes.String})
	columns = append(columns, Column{Name: CqSourceNameColumn.Name, Type: arrow.BinaryTypes.String})
	columns = append(columns, Column{Name: CqSyncTimeColumn.Name, Type: arrow.FixedWidthTypes.Timestamp_us})
	columns = append(columns, TestSourceColumns(testOpts)...)
	return &Table{Name: name, Columns: columns}
}

// GenTestDataOptions are options for generating test data
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
	StableTime    time.Time
	TimePrecision time.Duration
}

// GenTestData generates a slice of arrow.Records with the given schema and options.
func GenTestData(table *Table, opts GenTestDataOptions) []arrow.Record {
	var records []arrow.Record
	sc := table.ToArrowSchema()
	for j := 0; j < opts.MaxRows; j++ {
		nullRow := j%2 == 1
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		for i, c := range table.Columns {
			if nullRow && !c.NotNull && !c.PrimaryKey &&
				c.Name != CqSourceNameColumn.Name &&
				c.Name != CqSyncTimeColumn.Name &&
				c.Name != CqIDColumn.Name {
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

func getExampleJSON(colName string, dataType arrow.DataType, opts GenTestDataOptions) string {
	// handle lists (including maps)
	if arrow.IsListLike(dataType.ID()) {
		if dataType.ID() == arrow.MAP {
			k := getExampleJSON(colName, dataType.(*arrow.MapType).KeyType(), opts)
			v := getExampleJSON(colName, dataType.(*arrow.MapType).ItemType(), opts)
			return fmt.Sprintf(`[{"key": %s,"value": %s}]`, k, v)
		}
		inner := dataType.(*arrow.ListType).Elem()
		return `[` + getExampleJSON(colName, inner, opts) + `,null,` + getExampleJSON(colName, inner, opts) + `]`
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
		if strings.HasSuffix(colName, "_array") {
			return `[{"test":"test"},123,{"test_number":456}]`
		}
		return `{"test":["a","b",3]}`
	}
	if arrow.TypeEqual(dataType, types.ExtensionTypes.Inet) {
		return `"192.0.2.0/24"`
	}
	if arrow.TypeEqual(dataType, types.ExtensionTypes.MAC) {
		return `"aa:bb:cc:dd:ee:ff"`
	}

	// handle signed integers
	if arrow.IsSignedInteger(dataType.ID()) {
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
			if colName == CqSourceNameColumn.Name {
				return `"` + opts.SourceName + `"`
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
		var columns []string
		for _, field := range dataType.(*arrow.StructType).Fields() {
			v := getExampleJSON(field.Name, field.Type, opts)
			columns = append(columns, fmt.Sprintf(`"%s": %v`, field.Name, v))
		}
		return `{` + strings.Join(columns, ",") + `}`
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
			if colName == CqSyncTimeColumn.Name {
				t = opts.SyncTime.UTC()
			} else if !opts.StableTime.IsZero() {
				t = opts.StableTime
			}
			t = t.Truncate(opts.TimePrecision)

			switch timestampType {
			case arrow.FixedWidthTypes.Timestamp_s:
				return strconv.FormatInt(t.Unix(), 10)
			case arrow.FixedWidthTypes.Timestamp_ms:
				return strconv.FormatInt(t.UnixMilli(), 10)
			case arrow.FixedWidthTypes.Timestamp_us:
				return strconv.FormatInt(t.UnixMicro(), 10)
			case arrow.FixedWidthTypes.Timestamp_ns:
				// Note: We use microseconds instead of nanoseconds here because
				//       nanosecond precision is not supported by many destinations.
				//       For now, we begrudgingly accept loss of precision in these cases.
				//       See https://github.com/cloudquery/plugin-sdk/issues/830
				t = t.Truncate(time.Microsecond)
				// Use string timestamp string format here because JSON integers are
				// unmarshalled as float64, losing precision for nanosecond timestamps.
				return t.Format(`"2006-01-02 15:04:05.999999999"`)
			case arrow.FixedWidthTypes.Time32s:
				h, m, s := t.Clock()
				return strconv.Itoa(h*3600 + m*60 + s)
			case arrow.FixedWidthTypes.Time32ms:
				h, m, s := t.Clock()
				ns := t.Nanosecond()
				return strconv.Itoa(h*3600000 + m*60000 + s*1000 + ns/1000000)
			case arrow.FixedWidthTypes.Time64us:
				h, m, s := t.Clock()
				ns := t.Nanosecond()
				return strconv.Itoa(h*3600000000 + m*60000000 + s*1000000 + ns/1000)
			case arrow.FixedWidthTypes.Time64ns:
				h, m, s := t.Clock()
				ns := t.Nanosecond()
				return strconv.Itoa(h*3600000000000 + m*60000000000 + s*1000000000 + ns)
			default:
				panic("unhandled timestamp type: " + timestampType.Name())
			}
		}
	}

	// handle date types
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Date32) {
		return `19471`
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Date64) {
		ms := 19471 * 86400000
		return fmt.Sprintf("%d", ms)
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
