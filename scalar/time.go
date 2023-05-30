package scalar

import (
	"github.com/apache/arrow/go/v13/arrow"
)

type Time struct {
	Int
	Unit arrow.TimeUnit
}

func (s *Time) DataType() arrow.DataType {
	switch {
	case s.Unit == arrow.Second && s.Int.BitWidth == 32:
		return arrow.FixedWidthTypes.Time32s
	case s.Unit == arrow.Millisecond && s.Int.BitWidth == 32:
		return arrow.FixedWidthTypes.Time32ms
	case s.Unit == arrow.Nanosecond && s.Int.BitWidth == 64:
		return arrow.FixedWidthTypes.Time64ns
	case s.Unit == arrow.Microsecond && s.Int.BitWidth == 64:
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
	switch v := value.(type) {
	case arrow.Time32:
		return s.Int.Set(int64(v))
	case arrow.Time64:
		return s.Int.Set(int64(v))

	case string:
		switch s.BitWidth {
		case 64:
			t64, err := arrow.Time64FromString(v, s.Unit)
			if err != nil {
				return err
			}
			return s.Set(t64)
		case 32:
			t32, err := arrow.Time32FromString(v, s.Unit)
			if err != nil {
				return err
			}
			return s.Set(t32)
		default:
			return s.Int.Set(v)
		}
	case *string:
		if v == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*v)
	default:
		return s.Int.Set(value)
	}
}
