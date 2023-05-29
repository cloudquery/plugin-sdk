package scalar

import (
	"github.com/apache/arrow/go/v13/arrow"
)

type Duration struct {
	Int
	Unit arrow.TimeUnit
}

func (s *Duration) DataType() arrow.DataType {
	switch s.Unit {
	case arrow.Second:
		return arrow.FixedWidthTypes.Duration_s
	case arrow.Millisecond:
		return arrow.FixedWidthTypes.Duration_ms
	case arrow.Nanosecond:
		return arrow.FixedWidthTypes.Duration_ns
	case arrow.Microsecond:
		return arrow.FixedWidthTypes.Duration_us
	default:
		panic("unknown duration unit")
	}
}

func (s *Duration) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Duration)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Unit == r.Unit && s.Value == r.Value
}

func (s *Duration) Set(value any) error {
	if dur, ok := value.(arrow.Duration); ok {
		return s.Int.Set(int64(dur))
	}
	return s.Int.Set(value)
}
