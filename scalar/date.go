package scalar

import (
	"encoding"
	"fmt"
	"math"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
)

type Date struct {
	Valid    bool
	Value    time.Time
	BitWidth uint8 // defaults to 64
}

func (s *Date) IsValid() bool {
	return s.Valid
}

func (s *Date) DataType() arrow.DataType {
	switch s.getBitWidth() {
	case 32:
		return arrow.FixedWidthTypes.Date32
	case 64:
		return arrow.FixedWidthTypes.Date64
	default:
		panic("invalid bit width")
	}
}

func (s *Date) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Date)
	if !ok {
		return false
	}
	return s.getBitWidth() == r.getBitWidth() && s.Valid == r.Valid && s.Value.Equal(r.Value)
}

func (s *Date) String() string {
	if !s.Valid {
		return "(null)"
	}
	return s.Value.Format("2006-01-02")
}

func (s *Date) Get() any {
	return s.Value
}

func (s *Date) Set(val any) error {
	if val == nil {
		return nil
	}

	if sc, ok := val.(Scalar); ok {
		if !sc.IsValid() {
			s.Valid = false
			return nil
		}
		return s.Set(sc.Get())
	}

	switch value := val.(type) {
	case int:
		s.Value = time.Unix(int64(value), 0).UTC().Truncate(24 * time.Hour)
	case int64:
		s.Value = time.Unix(value, 0).UTC().Truncate(24 * time.Hour)
	case uint64:
		if value > math.MaxInt64 {
			return &ValidationError{Type: s.DataType(), Msg: "uint64 bigger than MaxInt64", Value: value}
		}
		s.Value = time.Unix(int64(value), 0).UTC().Truncate(24 * time.Hour)
	case time.Time:
		s.Value = value.UTC().Truncate(24 * time.Hour)
	case *time.Time:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case string:
		if value == "" {
			s.Valid = false
			return nil
		}

		p, err := time.Parse("2006-01-02", value)
		if err != nil {
			return &ValidationError{Type: s.DataType(), Msg: "cannot parse date", Value: value, Err: err}
		}
		return s.Set(p)
	case *string:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingTimeType(val); ok {
			return s.Set(originalSrc)
		}
		if value, ok := value.(encoding.TextMarshaler); ok {
			text, err := value.MarshalText()
			if err == nil {
				return s.Set(string(text))
			}
			// fall through to String() method
		}
		if value, ok := value.(fmt.Stringer); ok {
			str := value.String()
			return s.Set(str)
		}
		return &ValidationError{Type: arrow.FixedWidthTypes.Timestamp_us, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Date) getBitWidth() uint8 {
	if s.BitWidth == 0 {
		return 64 // default
	}
	return s.BitWidth
}
