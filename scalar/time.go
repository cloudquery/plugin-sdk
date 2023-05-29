package scalar

import (
	"github.com/apache/arrow/go/v13/arrow"
)

type Time struct {
	Int
	Unit     arrow.TimeUnit
	BitWidth uint8
}

func (s *Time) DataType() arrow.DataType {
	switch {
	case s.Unit == arrow.Second && s.getBitWidth() == 32:
		return arrow.FixedWidthTypes.Time32s
	case s.Unit == arrow.Millisecond && s.getBitWidth() == 32:
		return arrow.FixedWidthTypes.Time32ms
	case s.Unit == arrow.Nanosecond && s.getBitWidth() == 64:
		return arrow.FixedWidthTypes.Time64ns
	case s.Unit == arrow.Microsecond && s.getBitWidth() == 64:
		return arrow.FixedWidthTypes.Time64us
	default:
		panic("unknown time unit")
	}
}

func (s *Time) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Time)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Unit == r.Unit && s.getBitWidth() == r.getBitWidth() && s.Value == r.Value
}

func (s *Time) Set(value any) error {
	if t32, ok := value.(arrow.Time32); ok {
		return s.Int.Set(int64(t32))
	}
	return s.Int.Set(value)
}

func (s *Time) getBitWidth() uint8 {
	if s.BitWidth == 0 {
		return 32 // default is 32 because arrow.TimeUnit's zero value is Second
	}
	return s.BitWidth
}
