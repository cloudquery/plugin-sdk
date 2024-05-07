package scalar

import (
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewScalar(t *testing.T) {
	tl := []struct {
		dt    arrow.DataType
		input any
	}{
		{dt: arrow.PrimitiveTypes.Uint8},
		{dt: arrow.PrimitiveTypes.Uint16},
		{dt: arrow.PrimitiveTypes.Uint32},
		{dt: arrow.PrimitiveTypes.Uint64},
		{dt: arrow.PrimitiveTypes.Int8},
		{dt: arrow.PrimitiveTypes.Int16},
		{dt: arrow.PrimitiveTypes.Int32},
		{dt: arrow.PrimitiveTypes.Int64},
		{dt: arrow.PrimitiveTypes.Float32},
		{dt: arrow.PrimitiveTypes.Float64},

		{dt: arrow.BinaryTypes.String},
		{dt: arrow.BinaryTypes.Binary},
		{dt: arrow.BinaryTypes.LargeString},
		{dt: arrow.BinaryTypes.LargeBinary},

		{dt: arrow.FixedWidthTypes.Boolean},
		{dt: arrow.FixedWidthTypes.Date32, input: "2006-01-02"},
		{dt: arrow.FixedWidthTypes.Date64, input: "2006-01-02"},
		{dt: arrow.FixedWidthTypes.Time32s, input: "21:14:00"},
		{dt: arrow.FixedWidthTypes.Time32ms, input: "21:14:00.709"},
		{dt: arrow.FixedWidthTypes.Time64us, input: "21:14:00.709229"},
		{dt: arrow.FixedWidthTypes.Time64ns, input: "21:14:00.709227000"},
		{dt: arrow.FixedWidthTypes.Timestamp_ns, input: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Timestamp_us, input: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Timestamp_ms, input: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Timestamp_s, input: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Duration_ns},
		{dt: arrow.FixedWidthTypes.Duration_ns, input: "1ns"},
		{dt: arrow.FixedWidthTypes.Duration_us},
		{dt: arrow.FixedWidthTypes.Duration_us, input: "1us"},
		{dt: arrow.FixedWidthTypes.Duration_ms},
		{dt: arrow.FixedWidthTypes.Duration_ms, input: "1ms"},
		{dt: arrow.FixedWidthTypes.Duration_s},
		{dt: arrow.FixedWidthTypes.Duration_s, input: "1s"},
		{dt: arrow.FixedWidthTypes.Float16},

		{dt: arrow.FixedWidthTypes.DayTimeInterval, input: map[string]any{"days": 1, "milliseconds": 2}},
		{dt: arrow.FixedWidthTypes.MonthDayNanoInterval, input: map[string]any{"months": 1, "days": 2, "nanoseconds": 3}},
		{dt: arrow.FixedWidthTypes.MonthInterval, input: map[string]any{"months": 1}},

		{dt: arrow.StructOf(arrow.Field{Name: "i64", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "s", Type: arrow.BinaryTypes.String}), input: `{"i64": 1, "s": "foo"}`},
		{dt: &arrow.Decimal128Type{Precision: 10, Scale: 5}},
		{dt: &arrow.Decimal256Type{Precision: 10, Scale: 5}},
	}

	for _, tc := range tl {
		tc := tc
		if tc.input == nil {
			tc.input = "1"
		}

		t.Run("create_append:"+tc.dt.String(), func(t *testing.T) {
			t.Parallel()

			bldr := array.NewBuilder(memory.DefaultAllocator, tc.dt)
			defer bldr.Release()

			s := NewScalar(tc.dt)
			if !arrow.TypeEqual(s.DataType(), tc.dt) {
				t.Fatalf("expected %v, got %v", tc.dt, s.DataType())
			}

			assert.NoError(t, s.Set(tc.input))
			if t.Failed() {
				return
			}

			assert.True(t, s.IsValid())
			AppendToBuilder(bldr, s)

			t.Run("double_set_nil", genDoubleSetTest(tc.dt, tc.input, nil))

			var str *string
			t.Run("double_set_typed_nil_string", genDoubleSetTest(tc.dt, tc.input, str))

			switch {
			case
				tc.dt.ID() == arrow.BOOL,
				tc.dt.ID() == arrow.DATE32,
				tc.dt.ID() == arrow.DATE64,
				tc.dt.ID() == arrow.DURATION,
				tc.dt.ID() == arrow.LARGE_STRING,
				tc.dt.ID() == arrow.STRING,
				tc.dt.ID() == arrow.TIME32,
				tc.dt.ID() == arrow.TIME64,
				tc.dt.ID() == arrow.TIMESTAMP,
				arrow.IsDecimal(tc.dt.ID()),
				arrow.IsNested(tc.dt.ID()):

			case arrow.IsInteger(tc.dt.ID()), arrow.IsFloating(tc.dt.ID()):
				var i8 *int8
				t.Run("double_set_typed_nil_int8", genDoubleSetTest(tc.dt, tc.input, i8))

				var i16 *int16
				t.Run("double_set_typed_nil_int16", genDoubleSetTest(tc.dt, tc.input, i16))

				var i32 *int32
				t.Run("double_set_typed_nil_int32", genDoubleSetTest(tc.dt, tc.input, i32))

				var i64 *int64
				t.Run("double_set_typed_nil_int64", genDoubleSetTest(tc.dt, tc.input, i64))

				var f32 *float32
				t.Run("double_set_typed_nil_f32", genDoubleSetTest(tc.dt, tc.input, f32))

				var f64 *float64
				t.Run("double_set_typed_nil_f64", genDoubleSetTest(tc.dt, tc.input, f64))

			default:
				var val []byte
				t.Run("double_set_typed_nil_byteslice", genDoubleSetTest(tc.dt, tc.input, val))
			}
		})
	}
}

func genDoubleSetTest(dt arrow.DataType, input any, setToNil any) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		bldr := array.NewBuilder(memory.DefaultAllocator, dt)
		defer bldr.Release()

		s := NewScalar(dt)
		if !arrow.TypeEqual(s.DataType(), dt) {
			t.Fatalf("expected %v, got %v", dt, s.DataType())
		}

		assert.NoError(t, s.Set(input))
		if t.Failed() {
			return
		}

		assert.NoError(t, s.Set(setToNil))
		if t.Failed() {
			return
		}

		assert.False(t, s.IsValid())
	}
}
