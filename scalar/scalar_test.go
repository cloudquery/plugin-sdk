package scalar

import (
	"strconv"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
)

func TestNewScalar(t *testing.T) {
	tl := []struct {
		dt arrow.DataType
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
		{dt: arrow.PrimitiveTypes.Date32},
		{dt: arrow.PrimitiveTypes.Date64},

		{dt: arrow.BinaryTypes.String},
		{dt: arrow.BinaryTypes.Binary},
		//{dt: arrow.BinaryTypes.LargeString},
		//{dt: arrow.BinaryTypes.LargeBinary},

		{dt: arrow.FixedWidthTypes.Boolean},
		{dt: arrow.FixedWidthTypes.Date32},
		{dt: arrow.FixedWidthTypes.Date64},
		//{dt: arrow.FixedWidthTypes.Time32s},
		//{dt: arrow.FixedWidthTypes.Time32ms},
		//{dt: arrow.FixedWidthTypes.Time64us},
		//{dt: arrow.FixedWidthTypes.Time64ns},
		//{dt: arrow.FixedWidthTypes.Timestamp_ns},
		{dt: arrow.FixedWidthTypes.Timestamp_us},
		//{dt: arrow.FixedWidthTypes.Timestamp_ms},
		//{dt: arrow.FixedWidthTypes.Timestamp_s},
		//{dt: arrow.FixedWidthTypes.Duration_ns},
		//{dt: arrow.FixedWidthTypes.Duration_us},
		//{dt: arrow.FixedWidthTypes.Duration_ms},
		//{dt: arrow.FixedWidthTypes.Duration_s},
		{dt: arrow.FixedWidthTypes.Float16},
		//{dt: arrow.FixedWidthTypes.DayTimeInterval},
		//{dt: arrow.FixedWidthTypes.MonthDayNanoInterval},
		//{dt: arrow.FixedWidthTypes.MonthInterval},
	}

	for idx, tc := range tl {
		tc := tc
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			t.Parallel()
			s := NewScalar(tc.dt)
			if s.DataType() != tc.dt {
				t.Fatalf("expected %v, got %v", tc.dt, s.DataType())
			}
		})
	}
}
