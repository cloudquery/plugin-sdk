package scalar

import (
	"github.com/apache/arrow/go/v16/arrow"
)

type Time struct {
	Int
	Unit arrow.TimeUnit
}

func (s *Time) DataType() arrow.DataType {
	switch s.getBitWidth() {
	case 64:
		return &arrow.Time64Type{Unit: s.Unit}
	case 32:
		return &arrow.Time32Type{Unit: s.Unit}
	default:
		panic("unsupported bit width")
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

func (s *Time) String() string {
	if !s.Valid {
		return nullValueStr
	}

	switch s.getBitWidth() {
	case 64:
		return arrow.Time64(s.Int.Value).FormattedString(s.Unit)
	case 32:
		return arrow.Time32(s.Int.Value).FormattedString(s.Unit)
	default:
		panic("unsupported bit width")
	}
}

func (s *Time) Get() any {
	if !s.Valid {
		return nil
	}
	switch s.getBitWidth() {
	case 64:
		return arrow.Time64(s.Int.Get().(int64))
	case 32:
		return arrow.Time32(s.Int.Get().(int64))
	default:
		panic("unknown bit width")
	}
}

func (s *Time) Set(value any) error {
	switch v := value.(type) {
	case arrow.Time32:
		return s.Int.Set(int64(v))
	case arrow.Time64:
		return s.Int.Set(int64(v))

	case string:
		switch s.getBitWidth() {
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
