package scalar

import (
	"strings"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
)

type Duration struct {
	Int
	Unit arrow.TimeUnit
}

func (s *Duration) DataType() arrow.DataType {
	return &arrow.DurationType{Unit: s.Unit}
}

func (s *Duration) String() string {
	if !s.Int.IsValid() {
		return nullValueStr
	}

	return s.Int.String() + s.Unit.String()
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
	switch v := value.(type) {
	case string:
		stripped := strings.TrimSuffix(v, s.Unit.String())
		return s.Int.Set(stripped)
	case *string:
		if v == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*v)
	case time.Duration:
		return s.Int.Set(v / s.Unit.Multiplier())
	case *time.Duration:
		if v == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*v)
	}
	return s.Int.Set(value)
}
