package scalar

import (
	"encoding"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
)

type Date64 struct {
	Valid bool
	Value arrow.Date64
}

func (s *Date64) IsValid() bool {
	return s.Valid
}

func (*Date64) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.Date64
}

func (s *Date64) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Date64)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Date64) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.Value.FormattedString()
}

func (s *Date64) Get() any {
	return s.Value
}

func (s *Date64) Set(val any) error {
	if val == nil {
		s.Valid = false
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
	case arrow.Date64:
		s.Value = value
	case time.Time:
		return s.Set(arrow.Date64FromTime(value))
	case *time.Time:
		if value == nil {
			s.Valid = false
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
			s.Valid = false
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
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}
