package scalar

import (
	"math"
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Int64 struct {
	Valid bool
	Value int64
}

func (s *Int64) IsValid() bool {
	return s.Valid
}

func (s *Int64) DataType() arrow.DataType {
	return arrow.PrimitiveTypes.Int64
}

func (s *Int64) String() string {
	if !s.Valid {
		return "(null)"
	}
	return strconv.FormatInt(int64(s.Value), 10)
}

func (s *Int64) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Int64)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Int64) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	switch value := val.(type) {
	case int8:
		s.Value = int64(value)
	case int16:
		s.Value = int64(value)
	case int32:
		s.Value = int64(value)
	case int64:
		s.Value = value
	case int:
		s.Value = int64(value)
	case uint8:
		s.Value = int64(value)
	case uint16:
		s.Value = int64(value)
	case uint32:
		s.Value = int64(value)
	case uint64:
		if value > math.MaxInt64 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int64, Msg: "uint64 bigger than MaxInt64", Value: value}
		}
		s.Value = int64(value)
	case uint:
		if value > math.MaxInt64 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int64, Msg: "uint bigger than MaxInt64", Value: value}
		}
		s.Value = int64(value)
	case float32:
		s.Value = int64(value)
	case float64:
		s.Value = int64(value)
	case string:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int64, Msg: "invalid string", Value: value}
		}
		s.Value = num
	case *int8:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *int16:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *int32:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *int64:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *int:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *uint8:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *uint16:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *uint32:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *uint64:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *uint:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *float32:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *float64:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	case *string:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.PrimitiveTypes.Int64, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}