package scalar

import (
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/float16"
	"github.com/cloudquery/plugin-sdk/v3/types"
)

// Scalar represents a single value of a specific DataType as opposed to
// an array.
//
// Scalars are useful for passing single value inputs to compute functions
// (not yet implemented) or for representing individual array elements,
// (with a non-trivial cost though).
type Scalar interface {
	fmt.Stringer
	// IsValid returns true if the value is non-null, otherwise false.
	IsValid() bool
	// The datatype of the value in this scalar
	DataType() arrow.DataType
	// Performs cheap validation checks, returns nil if successful
	// Validate() error
	// tries to set the value of the scalar to the given value
	Set(val any) error

	Get() any
	Equal(other Scalar) bool
}

type Vector []Scalar

func (v Vector) Equal(r Vector) bool {
	if len(v) != len(r) {
		return false
	}
	for i := range v {
		if !v[i].Equal(r[i]) {
			return false
		}
	}
	return true
}

func NewScalar(dt arrow.DataType) Scalar {
	switch dt.ID() {
	case arrow.TIMESTAMP:
		switch {
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Timestamp_ms):
			return &Timestamp{Unit: arrow.Millisecond}
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Timestamp_ns):
			return &Timestamp{Unit: arrow.Nanosecond}
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Timestamp_s):
			return &Timestamp{Unit: arrow.Second}
		default:
			return &Timestamp{Unit: arrow.Microsecond}
		}
	case arrow.BINARY:
		return &Binary{}
	case arrow.STRING:
		return &String{}
	case arrow.LARGE_BINARY:
		return &LargeBinary{}
	case arrow.LARGE_STRING:
		return &LargeString{}
	case arrow.INT64:
		return &Int{}
	case arrow.INT32:
		return &Int{BitWidth: 32}
	case arrow.INT16:
		return &Int{BitWidth: 16}
	case arrow.INT8:
		return &Int{BitWidth: 8}
	case arrow.UINT64:
		return &Uint{}
	case arrow.UINT32:
		return &Uint{BitWidth: 32}
	case arrow.UINT16:
		return &Uint{BitWidth: 16}
	case arrow.UINT8:
		return &Uint{BitWidth: 8}
	case arrow.FLOAT64:
		return &Float{}
	case arrow.FLOAT32:
		return &Float{BitWidth: 32}
	case arrow.FLOAT16:
		return &Float{BitWidth: 16}
	case arrow.BOOL:
		return &Bool{}
	case arrow.EXTENSION:
		switch {
		case arrow.TypeEqual(dt, types.ExtensionTypes.UUID):
			return &UUID{}
		case arrow.TypeEqual(dt, types.ExtensionTypes.JSON):
			return &JSON{}
		case arrow.TypeEqual(dt, types.ExtensionTypes.MAC):
			return &Mac{}
		case arrow.TypeEqual(dt, types.ExtensionTypes.Inet):
			return &Inet{}
		default:
			panic("not implemented extension: " + dt.Name())
		}
	case arrow.LIST:
		return &List{
			Type: dt,
		}
	case arrow.DATE64:
		return &Date64{}
	case arrow.DATE32:
		return &Date32{}
	case arrow.DURATION:
		switch {
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Duration_ms):
			return &Duration{
				Int:  Int{BitWidth: 64},
				Unit: arrow.Millisecond,
			}
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Duration_ns):
			return &Duration{
				Int:  Int{BitWidth: 64},
				Unit: arrow.Nanosecond,
			}
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Duration_s):
			return &Duration{
				Int:  Int{BitWidth: 64},
				Unit: arrow.Second,
			}
		default:
			return &Duration{
				Int:  Int{BitWidth: 64},
				Unit: arrow.Microsecond,
			}
		}
	case arrow.TIME32, arrow.TIME64:
		switch {
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Time32s):
			return &Time{
				Int:      Int{BitWidth: 64},
				BitWidth: 32,
				Unit:     arrow.Second,
			}
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Time32ms):
			return &Time{
				Int:      Int{BitWidth: 64},
				BitWidth: 32,
				Unit:     arrow.Millisecond,
			}
		case arrow.TypeEqual(dt, arrow.FixedWidthTypes.Time64us):
			return &Time{
				Int:      Int{BitWidth: 64},
				BitWidth: 64,
				Unit:     arrow.Microsecond,
			}
		default:
			return &Time{
				Int:      Int{BitWidth: 64},
				BitWidth: 64,
				Unit:     arrow.Nanosecond,
			}
		}

	case arrow.INTERVAL_MONTHS:
		return &MonthInterval{Int{BitWidth: 32}}
	case arrow.INTERVAL_DAY_TIME:
		return &DayTimeInterval{}
	case arrow.INTERVAL_MONTH_DAY_NANO:
		return &MonthDayNanoInterval{}

	case arrow.STRUCT:
		return &Struct{Type: dt}

	default:
		panic("not implemented: " + dt.Name())
	}
}

func AppendToBuilder(bldr array.Builder, s Scalar) {
	if !s.IsValid() {
		bldr.AppendNull()
		return
	}
	switch s.DataType().ID() {
	case arrow.BINARY:
		bldr.(*array.BinaryBuilder).Append(s.(*Binary).Value)
	case arrow.LARGE_BINARY:
		bldr.(*array.BinaryBuilder).Append(s.(*LargeBinary).Value)
	case arrow.STRING:
		bldr.(*array.StringBuilder).Append(s.(*String).Value)
	case arrow.LARGE_STRING:
		bldr.(*array.LargeStringBuilder).Append(s.(*LargeString).s.Value)
	case arrow.INT64:
		bldr.(*array.Int64Builder).Append(s.(*Int).Value)
	case arrow.INT32:
		bldr.(*array.Int32Builder).Append(int32(s.(*Int).Value))
	case arrow.INT16:
		bldr.(*array.Int16Builder).Append(int16(s.(*Int).Value))
	case arrow.INT8:
		bldr.(*array.Int8Builder).Append(int8(s.(*Int).Value))
	case arrow.UINT64:
		bldr.(*array.Uint64Builder).Append(s.(*Uint).Value)
	case arrow.UINT32:
		bldr.(*array.Uint32Builder).Append(uint32(s.(*Uint).Value))
	case arrow.UINT16:
		bldr.(*array.Uint16Builder).Append(uint16(s.(*Uint).Value))
	case arrow.UINT8:
		bldr.(*array.Uint8Builder).Append(uint8(s.(*Uint).Value))
	case arrow.FLOAT16:
		bldr.(*array.Float16Builder).Append(float16.New(float32(s.(*Float).Value)))
	case arrow.FLOAT32:
		bldr.(*array.Float32Builder).Append(float32(s.(*Float).Value))
	case arrow.FLOAT64:
		bldr.(*array.Float64Builder).Append(s.(*Float).Value)
	case arrow.BOOL:
		bldr.(*array.BooleanBuilder).Append(s.(*Bool).Value)
	case arrow.TIMESTAMP:
		bldr.(*array.TimestampBuilder).Append(arrow.Timestamp(s.(*Timestamp).Value.UnixMicro()))
	case arrow.DURATION:
		bldr.(*array.DurationBuilder).Append(arrow.Duration(s.(*Duration).Value))
	case arrow.DATE32:
		bldr.(*array.Date32Builder).Append(arrow.Date32(s.(*Date32).Value))
	case arrow.DATE64:
		bldr.(*array.Date64Builder).Append(arrow.Date64(s.(*Date64).Value))
	case arrow.TIME32:
		bldr.(*array.Time32Builder).Append(arrow.Time32(int32(s.(*Time).Value)))
	case arrow.TIME64:
		bldr.(*array.Time64Builder).Append(arrow.Time64(s.(*Time).Value))
	case arrow.INTERVAL_MONTHS:
		bldr.(*array.MonthIntervalBuilder).Append(arrow.MonthInterval(int32(s.(*MonthInterval).Value)))
	case arrow.INTERVAL_DAY_TIME:
		bldr.(*array.DayTimeIntervalBuilder).Append(s.(*DayTimeInterval).Value)
	case arrow.INTERVAL_MONTH_DAY_NANO:
		bldr.(*array.MonthDayNanoIntervalBuilder).Append(s.(*MonthDayNanoInterval).Value)

	case arrow.STRUCT:
		sb := bldr.(*array.StructBuilder)
		sb.Append(true)

		v := s.(*Struct).Value
		m := v.(map[string]any)
		st := sb.Type().(*arrow.StructType)

		for i, f := range st.Fields() {
			sc := NewScalar(sb.FieldBuilder(i).Type())

			if sv, ok := m[f.Name]; ok {
				if err := sc.Set(sv); err != nil {
					panic(err)
				}
			}

			AppendToBuilder(sb.FieldBuilder(i), sc)
		}

	case arrow.LIST:
		lb := bldr.(*array.ListBuilder)
		lb.Append(true)
		for _, v := range s.(*List).Value {
			AppendToBuilder(lb.ValueBuilder(), v)
		}
	case arrow.EXTENSION:
		switch {
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.UUID):
			bldr.(*types.UUIDBuilder).Append(s.(*UUID).Value)
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.JSON):
			bldr.(*types.JSONBuilder).AppendBytes(s.(*JSON).Value)
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.MAC):
			bldr.(*types.MACBuilder).Append(s.(*Mac).Value)
		case arrow.TypeEqual(s.DataType(), types.ExtensionTypes.Inet):
			bldr.(*types.InetBuilder).Append(s.(*Inet).Value)
		default:
			panic("not implemented extension: " + s.DataType().Name())
		}
	default:
		panic("not implemented: " + s.DataType().String())
	}
}

func AppendToRecordBuilder(bldr *array.RecordBuilder, vector Vector) {
	for i, scalar := range vector {
		AppendToBuilder(bldr.Field(i), scalar)
	}
}
