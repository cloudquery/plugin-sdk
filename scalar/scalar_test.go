package scalar

import (
	"strconv"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewScalar(t *testing.T) {
	tl := []struct {
		dt  arrow.DataType
		str string
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
		{dt: arrow.FixedWidthTypes.Date32, str: "2006-01-02"},
		{dt: arrow.FixedWidthTypes.Date64, str: "2006-01-02"},
		{dt: arrow.FixedWidthTypes.Time32s, str: "21:14:00"},
		{dt: arrow.FixedWidthTypes.Time32ms, str: "21:14:00.709"},
		{dt: arrow.FixedWidthTypes.Time64us, str: "21:14:00.709229"},
		{dt: arrow.FixedWidthTypes.Time64ns, str: "21:14:00.709227000"},
		{dt: arrow.FixedWidthTypes.Timestamp_ns, str: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Timestamp_us, str: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Timestamp_ms, str: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Timestamp_s, str: "2006-01-02 15:04:05.999999999"},
		{dt: arrow.FixedWidthTypes.Duration_ns},
		{dt: arrow.FixedWidthTypes.Duration_ns, str: "1ns"},
		{dt: arrow.FixedWidthTypes.Duration_us},
		{dt: arrow.FixedWidthTypes.Duration_us, str: "1us"},
		{dt: arrow.FixedWidthTypes.Duration_ms},
		{dt: arrow.FixedWidthTypes.Duration_ms, str: "1ms"},
		{dt: arrow.FixedWidthTypes.Duration_s},
		{dt: arrow.FixedWidthTypes.Duration_s, str: "1s"},
		{dt: arrow.FixedWidthTypes.Float16},

		{dt: arrow.FixedWidthTypes.DayTimeInterval, str: `{"days":1,"milliseconds":2}`},
		{dt: arrow.FixedWidthTypes.MonthDayNanoInterval, str: `{"months":1,"days":2,"nanoseconds":3}`},
		{dt: arrow.FixedWidthTypes.MonthInterval},
	}

	for idx, tc := range tl {
		tc := tc
		if tc.str == "" {
			tc.str = "1"
		}

		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			t.Parallel()

			bldr := array.NewBuilder(memory.DefaultAllocator, tc.dt)
			defer bldr.Release()

			s := NewScalar(tc.dt)
			if s.DataType() != tc.dt {
				t.Fatalf("expected %v, got %v", tc.dt, s.DataType())
			}

			assert.NoErrorf(t, s.Set(tc.str), "failed with DataType %s", tc.dt.String())
			if t.Failed() {
				return
			}

			assert.Truef(t, s.IsValid(), "failed with DataType %s", tc.dt.String())
			AppendToBuilder(bldr, s)
		})
	}
}
