package scalar

import (
	"encoding"
	"fmt"
	"math"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
)

type Date32 struct {
	Valid bool
	Value int32 // days since unix epoch
}

func (s *Date32) IsValid() bool {
	return s.Valid
}

func (s *Date32) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.Date32
}

func (s *Date32) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Date32)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Date32) String() string {
	if !s.Valid {
		return "(null)"
	}
	return time.Unix(86400*int64(s.Value), 0).UTC().Format("2006-01-02")
}

func (s *Date32) Get() any {
	return s.Value
}

func (s *Date32) Set(val any) error {
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
		s.Value = int32(value)
	case int64:
		if value > math.MaxInt32 {
			return &ValidationError{Type: s.DataType(), Msg: "int64 bigger than MaxInt32", Value: value}
		}
		s.Value = int32(value)
	case uint64:
		if value > math.MaxInt32 {
			return &ValidationError{Type: s.DataType(), Msg: "uint64 bigger than MaxInt32", Value: value}
		}
		s.Value = int32(value)
	case time.Time:
		val := value.UTC().Unix() / 86400
		return s.Set(val)
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
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}
