package scalar

import (
	"encoding"
	"fmt"
	"math"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
)

type Date64 struct {
	Valid bool
	Value int64 // int64 milliseconds since the UNIX epoch
}

func (s *Date64) IsValid() bool {
	return s.Valid
}

func (s *Date64) DataType() arrow.DataType {
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
		return "(null)"
	}
	return time.UnixMilli(s.Value).UTC().Format(arrowStringFormat)
}

func (s *Date64) Get() any {
	return s.Value
}

func (s *Date64) Set(val any) error {
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
	case arrow.Date64:
		s.Value = int64(value)
	case int:
		s.Value = int64(value)
	case int64:
		s.Value = value
	case uint64:
		if value > math.MaxInt64 {
			return &ValidationError{Type: s.DataType(), Msg: "uint64 bigger than MaxInt64", Value: value}
		}
		s.Value = int64(value)
	case time.Time:
		return s.Set(arrow.Date64FromTime(value))
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

		p, err := time.Parse(arrowStringFormat, value)
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
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}
