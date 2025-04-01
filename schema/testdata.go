package schema

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
)

// TestSourceOptions controls which types are included by TestSourceColumns.
type TestSourceOptions struct {
	SkipDates      bool
	SkipDurations  bool
	SkipIntervals  bool
	SkipLargeTypes bool // e.g. large binary, large string
	SkipLists      bool // lists of all primitive types. Lists that were supported by CQTypes are always included.
	SkipMaps       bool
	SkipStructs    bool
	SkipTimes      bool // time of day types
	SkipTimestamps bool // timestamp types. Microsecond timestamp is always be included, regardless of this setting.
	TimePrecision  time.Duration
	SkipDecimals   bool
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
	columns := make([]Column, len(baseColumns)*2)
	for i := 0; i < len(columns); i += 2 {
		// we focus on string and int keys for now
		n := i / 2
		columns[i] = Column{Name: "string_" + baseColumns[n].Name + "_map", Type: arrow.MapOf(arrow.BinaryTypes.String, baseColumns[n].Type)}
		columns[i+1] = Column{Name: "int_" + baseColumns[n].Name + "_map", Type: arrow.MapOf(arrow.PrimitiveTypes.Int64, baseColumns[n].Type)}
	}
	return columns
}

func columnsToFields(columns []Column) []arrow.Field {
	fields := make([]arrow.Field, len(columns))
	for i := range columns {
		fields[i] = arrow.Field{
			Name: columns[i].Name,
			Type: columns[i].Type,
		}
	}
	return fields
}

func TestTable(name string, opts TestSourceOptions) *Table {
	t := &Table{
		Name:    name,
		Columns: make(ColumnList, 0),
	}
	var columns ColumnList
	columns = append(columns, ColumnList{
		// id is to be used as an auto-incrementing column that can be used to order the rows in tests
		{Name: "id", Type: arrow.PrimitiveTypes.Int64, NotNull: true},

		// primitive columns
		{Name: "int8", Type: arrow.PrimitiveTypes.Int8},
		{Name: "int16", Type: arrow.PrimitiveTypes.Int16},
		{Name: "int32", Type: arrow.PrimitiveTypes.Int32},
		{Name: "int64", Type: arrow.PrimitiveTypes.Int64},
		{Name: "uint8", Type: arrow.PrimitiveTypes.Uint8},
		{Name: "uint16", Type: arrow.PrimitiveTypes.Uint16},
		{Name: "uint32", Type: arrow.PrimitiveTypes.Uint32},
		{Name: "uint64", Type: arrow.PrimitiveTypes.Uint64},
		{Name: "float32", Type: arrow.PrimitiveTypes.Float32},
		{Name: "float64", Type: arrow.PrimitiveTypes.Float64},

		// basic columns
		{Name: "binary", Type: arrow.BinaryTypes.Binary},
		{Name: "string", Type: arrow.BinaryTypes.String},
		{Name: "boolean", Type: arrow.FixedWidthTypes.Boolean},

		// extension types
		{Name: "uuid", Type: types.ExtensionTypes.UUID},
		{Name: "inet", Type: types.ExtensionTypes.Inet},
		{Name: "mac", Type: types.ExtensionTypes.MAC},
		{Name: "json", Type: types.ExtensionTypes.JSON},
	}...)
	if !opts.SkipDates {
		columns = append(columns, ColumnList{
			{Name: "date32", Type: arrow.FixedWidthTypes.Date32},
			{Name: "date64", Type: arrow.FixedWidthTypes.Date64},
		}...)
	}
	if !opts.SkipDurations {
		columns = append(columns, ColumnList{
			{Name: "duration_s", Type: arrow.FixedWidthTypes.Duration_s},
			{Name: "duration_ms", Type: arrow.FixedWidthTypes.Duration_ms},
			{Name: "duration_us", Type: arrow.FixedWidthTypes.Duration_us},
			{Name: "duration_ns", Type: arrow.FixedWidthTypes.Duration_ns},
		}...)
	}

	if !opts.SkipIntervals {
		columns = append(columns, ColumnList{
			{Name: "interval_month", Type: arrow.FixedWidthTypes.MonthInterval},
			{Name: "interval_day_time", Type: arrow.FixedWidthTypes.DayTimeInterval},
			{Name: "interval_month_day_nano", Type: arrow.FixedWidthTypes.MonthDayNanoInterval},
		}...)
	}

	if !opts.SkipLargeTypes {
		columns = append(columns, ColumnList{
			{Name: "large_binary", Type: arrow.BinaryTypes.LargeBinary},
			{Name: "large_string", Type: arrow.BinaryTypes.LargeString},
		}...)
	}

	if !opts.SkipTimes {
		columns = append(columns, ColumnList{
			{Name: "time32_s", Type: arrow.FixedWidthTypes.Time32s},
			{Name: "time32_ms", Type: arrow.FixedWidthTypes.Time32ms},
			{Name: "time64_us", Type: arrow.FixedWidthTypes.Time64us},
			{Name: "time64_ns", Type: arrow.FixedWidthTypes.Time64ns},
		}...)
	}

	if !opts.SkipTimestamps {
		columns = append(columns, ColumnList{
			{Name: "timestamp_s", Type: arrow.FixedWidthTypes.Timestamp_s},
			{Name: "timestamp_ms", Type: arrow.FixedWidthTypes.Timestamp_ms},
			{Name: "timestamp_us", Type: arrow.FixedWidthTypes.Timestamp_us},
			{Name: "timestamp_ns", Type: arrow.FixedWidthTypes.Timestamp_ns},
		}...)
	}

	if !opts.SkipDecimals {
		columns = append(columns, ColumnList{
			{Name: "decimal128", Type: &arrow.Decimal128Type{Precision: 19, Scale: 10}},
			// {Name: "decimal256", Type: &arrow.Decimal256Type{Precision: 40, Scale: 10}},
		}...)
	}

	if !opts.SkipStructs {
		columns = append(columns, Column{Name: "struct", Type: arrow.StructOf(columnsToFields(columns)...)})

		// struct with nested struct
		// columns = append(columns, Column{Name: "nested_struct", Type: arrow.StructOf(arrow.Field{Name: "inner", Type: arrow.StructOf(columnsToFields(basicColumns...)...)})})
	}

	if !opts.SkipLists {
		cols := excludeType(columns, types.ExtensionTypes.JSON)
		columns = append(columns, listOfColumns(cols)...)
	}

	if !opts.SkipMaps {
		columns = append(columns, mapOfColumns(columns)...)
	}

	t.Columns = append(t.Columns, columns...)
	return t
}

func excludeType(columns ColumnList, typ arrow.DataType) ColumnList {
	var cols ColumnList
	for _, c := range columns {
		if !arrow.TypeEqual(c.Type, typ) {
			cols = append(cols, c)
		}
	}
	return cols
}

// GenTestDataOptions are options for generating test data
type GenTestDataOptions struct {
	// SourceName is the name of the source to set in the source_name column.
	SourceName string
	// SyncTime is the time to set in the sync_time column.
	SyncTime time.Time
	// MaxRows is the number of rows to generate.
	MaxRows int
	// StableTime is the time to use for all rows other than sync time. If set to time.Time{}, a new time will be generated
	StableTime time.Time
	// TimePrecision is the precision to use for time columns.
	TimePrecision time.Duration
	// NullRows indicates whether to generate rows with all null values.
	NullRows bool
	// UseHomogeneousType indicates whether to use a single type for JSON arrays.
	UseHomogeneousType bool
}

type TestDataGenerator struct {
	counter  int
	seed     uint64
	colToRnd map[string]*rand.Rand
}

func NewTestDataGenerator(randomSeed uint64) *TestDataGenerator {
	return &TestDataGenerator{
		counter:  int(randomSeed),
		seed:     randomSeed,
		colToRnd: map[string]*rand.Rand{},
	}
}

func (tg *TestDataGenerator) Reset() {
	tg.counter = 0
	tg.colToRnd = map[string]*rand.Rand{}
}

// Generate will produce a single arrow.Record with the given schema and options.
func (tg *TestDataGenerator) Generate(table *Table, opts GenTestDataOptions) arrow.Record {
	sc := table.ToArrowSchema()
	if opts.MaxRows == 0 {
		// We generate an empty record
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		defer bldr.Release()
		return bldr.NewRecord()
	}

	var records []arrow.Record
	for j := 0; j < opts.MaxRows; j++ {
		tg.counter++
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		for i, c := range table.Columns {
			if opts.NullRows && !c.NotNull && !c.PrimaryKey &&
				c.Name != CqSourceNameColumn.Name &&
				c.Name != CqSyncTimeColumn.Name &&
				c.Name != CqIDColumn.Name {
				bldr.Field(i).AppendNull()
				continue
			}
			example := tg.getExampleJSON(c.Name, c.Type, opts)
			l := `[` + example + `]`
			err := bldr.Field(i).UnmarshalJSON([]byte(l))
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal json `%v` for column %v: %v", l, c.Name, err))
			}
		}
		records = append(records, bldr.NewRecord())
		bldr.Release()
	}

	arrowTable := array.NewTableFromRecords(sc, records)
	columns := make([]arrow.Array, sc.NumFields())
	for n := 0; n < sc.NumFields(); n++ {
		concatenated, err := array.Concatenate(arrowTable.Column(n).Data().Chunks(), memory.DefaultAllocator)
		if err != nil {
			panic(fmt.Sprintf("failed to concatenate arrays: %v", err))
		}
		columns[n] = concatenated
	}

	return array.NewRecord(sc, columns, -1)
}

func (tg TestDataGenerator) getExampleJSON(colName string, dataType arrow.DataType, opts GenTestDataOptions) string {
	var rnd, found = tg.colToRnd[colName]
	if !found {
		tg.colToRnd[colName] = rand.New(rand.NewSource(int64(tg.seed)))
		rnd = tg.colToRnd[colName]
	}

	// special case for auto-incrementing id column, used for to determine ordering in tests
	if arrow.IsInteger(dataType.ID()) && colName == "id" {
		return `` + strconv.Itoa(tg.counter) + ``
	}

	// handle lists (including maps)
	if arrow.IsListLike(dataType.ID()) {
		if dataType.ID() == arrow.MAP {
			k := tg.getExampleJSON(colName, dataType.(*arrow.MapType).KeyType(), opts)
			v := tg.getExampleJSON(colName, dataType.(*arrow.MapType).ItemType(), opts)
			k2 := tg.getExampleJSON(colName, dataType.(*arrow.MapType).KeyType(), opts)
			v2 := tg.getExampleJSON(colName, dataType.(*arrow.MapType).ItemType(), opts)
			return fmt.Sprintf(`[{"key": %s,"value": %s},{"key": %s,"value": %s}]`, k, v, k2, v2)
		}
		inner := dataType.(*arrow.ListType).Elem()
		return `[` + tg.getExampleJSON(colName, inner, opts) + `,null,` + tg.getExampleJSON(colName, inner, opts) + `]`
	}
	// handle extension types
	if arrow.TypeEqual(dataType, types.ExtensionTypes.UUID) {
		// This will make UUIDs deterministic like all other types
		hash := sha256.New()
		hash.Write([]byte(fmt.Sprintf(`"AString%d"`, rnd.Intn(100000))))
		u := uuid.NewSHA1(uuid.UUID{}, hash.Sum(nil))
		return `"` + u.String() + `"`
	}
	if arrow.TypeEqual(dataType, types.ExtensionTypes.JSON) {
		if strings.HasSuffix(colName, "_array") {
			if opts.UseHomogeneousType {
				return `[{"test1":"test1"},{"test2":"test2"},{"test3":"test3"}]`
			}
			return `[{"test":"test"},123,{"test_number":456}]`
		}
		if opts.UseHomogeneousType {
			return `{"test":["a", "b", "c"]}`
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
		switch dataType {
		case arrow.PrimitiveTypes.Int8:
			return fmt.Sprintf("-%d", rnd.Intn(int(^uint8(0)>>1)))
		case arrow.PrimitiveTypes.Int16:
			return fmt.Sprintf("-%d", rnd.Intn(int(^uint16(0)>>1)))
		case arrow.PrimitiveTypes.Int32:
			return fmt.Sprintf("-%d", rnd.Intn(int(^uint32(0)>>1)))
		case arrow.PrimitiveTypes.Int64:
			return fmt.Sprintf("-%d", rnd.Int63n(int64(^uint64(0)>>1)))
		}
	}

	// handle unsigned integers
	if arrow.IsUnsignedInteger(dataType.ID()) {
		switch dataType {
		case arrow.PrimitiveTypes.Uint8:
			return fmt.Sprintf("%d", rnd.Uint64()%(uint64(^uint8(0))))
		case arrow.PrimitiveTypes.Uint16:
			return fmt.Sprintf("%d", rnd.Uint64()%(uint64(^uint16(0))))
		case arrow.PrimitiveTypes.Uint32:
			return fmt.Sprintf("%d", rnd.Uint64()%(uint64(^uint32(0))))
		case arrow.PrimitiveTypes.Uint64:
			return fmt.Sprintf("%d", rnd.Uint64())
		}
	}

	// handle floats
	if arrow.IsFloating(dataType.ID()) {
		return fmt.Sprintf("%d.%d", rnd.Intn(1e3), rnd.Intn(1e3))
	}

	// handle decimals
	if arrow.IsDecimal(dataType.ID()) {
		return fmt.Sprintf("%d.%d", rnd.Int63n(1e9), rnd.Int63n(1e10))
	}

	// handle booleans
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Boolean) {
		if rnd.Intn(2) == 0 {
			return "false"
		}
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
			n := rnd.Intn(100000)
			return fmt.Sprintf(`"AString%d"`, n)
		}
	}

	// handle binary types
	binaryTypes := []arrow.DataType{
		arrow.BinaryTypes.Binary,
		arrow.BinaryTypes.LargeBinary,
	}
	for _, binaryType := range binaryTypes {
		if arrow.TypeEqual(dataType, binaryType) {
			bytes := make([]byte, 4)
			_, _ = rnd.Read(bytes)
			return `"` + base64.StdEncoding.EncodeToString(bytes) + `"`
		}
	}

	// handle structs
	if dataType.ID() == arrow.STRUCT {
		var columns []string
		for _, field := range dataType.(*arrow.StructType).Fields() {
			v := tg.getExampleJSON(field.Name, field.Type, opts)
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
		return fmt.Sprintf("%d", 19471+rnd.Intn(100))
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.Date64) {
		ms := (19471 + rnd.Intn(100)) * 86400000
		return fmt.Sprintf("%d", ms)
	}

	// handle duration and interval types
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.DayTimeInterval) {
		n := rnd.Intn(10000)
		return fmt.Sprintf(`{"days": %[1]d, "milliseconds": %[1]d}`, n)
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.MonthInterval) {
		return `{"months": 1}`
	}
	if arrow.TypeEqual(dataType, arrow.FixedWidthTypes.MonthDayNanoInterval) {
		n := rnd.Intn(10000)
		return fmt.Sprintf(`{"months": %[1]d, "days": %[1]d, "nanoseconds": %[1]d}`, n)
	}
	durationTypes := []arrow.DataType{
		arrow.FixedWidthTypes.Duration_s,
		arrow.FixedWidthTypes.Duration_ms,
		arrow.FixedWidthTypes.Duration_us,
		arrow.FixedWidthTypes.Duration_ns,
	}
	for _, durationType := range durationTypes {
		if arrow.TypeEqual(dataType, durationType) {
			n := rnd.Intn(10000000)
			return fmt.Sprintf("%d", n)
		}
	}

	panic("unknown type: " + dataType.String() + " column: " + colName)
}
