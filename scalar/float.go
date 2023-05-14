package scalar

import (
	"math"
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Float32 struct {
	Valid bool
	Value float32
}

func (s *Float32) IsValid() bool {
	return s.Valid
}

func (s *Float32) DataType() arrow.DataType {
	return arrow.PrimitiveTypes.Float32
}

func (s *Float32) String() string {
	if !s.Valid {
		return "(null)"
	}
	return strconv.FormatFloat(float64(s.Value), 'f', -1, 32)
}

func (s *Float32) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Float32)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Float32) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	switch value := val.(type) {
	case int8:
		s.Value = float32(value)
	case int16:
		if value > math.MaxInt8 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Float32, Msg: "int16 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case int32:
		if value > math.MaxInt8 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Float32, Msg: "int32 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case int64:
		if value > math.MaxInt32 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Float32, Msg: "int64 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case uint8:
		if value > math.MaxInt8 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "uint8 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case uint16:
		if value > math.MaxInt8 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "uint16 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case uint32:
		if value > math.MaxInt32 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "uint32 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case uint64:
		if value > math.MaxInt32 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "uint64 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case float32:
		s.Value = value
	case float64:
		if value > math.MaxInt32 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "float64 bigger than MaxInt8", Value: value}
		}
		s.Value = float32(value)
	case string:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "invalid string", Value: value}
		}
		s.Value = float32(num)
	case *int8:
		return s.Set(*value)
	case *int16:
		return s.Set(*value)
	case *int32:
		return s.Set(*value)
	case *int64:
		return s.Set(*value)
	case *uint8:
		return s.Set(*value)
	case *uint16:
		return s.Set(*value)
	case *uint32:
		return s.Set(*value)
	case *uint64:
		return s.Set(*value)
	case *float32:
		return s.Set(*value)
	case *float64:
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

type Float64 struct {
	Valid bool
	Value float64
}
func (s *Float64) IsValid() bool {
	return s.Valid
}

func (s *Float64) DataType() arrow.DataType {
	return arrow.PrimitiveTypes.Float64
}

func (s *Float64) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Float64)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Float64) String() string {
	if !s.Valid {
		return "(null)"
	}
	return strconv.FormatFloat(s.Value, 'f', -1, 64)
}

func (s *Float64) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	switch value := val.(type) {
	case int8:
		s.Value = float64(value)
	case int16:
		s.Value = float64(value)
	case int32:
		s.Value = float64(value)
	case int64:
		s.Value = float64(value)
	case uint8:
		s.Value = float64(value)
	case uint16:
		s.Value = float64(value)
	case uint32:
		s.Value = float64(value)
	case uint64:
		s.Value = float64(value)
	case float32:
		s.Value = float64(value)
	case float64:
		s.Value = value
	case string:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "invalid string", Value: value}
		}
		s.Value = num
	case *int8:
		return s.Set(*value)
	case *int16:
		return s.Set(*value)
	case *int32:
		return s.Set(*value)
	case *int64:
		return s.Set(*value)
	case *uint8:
		return s.Set(*value)
	case *uint16:
		return s.Set(*value)
	case *uint32:
		return s.Set(*value)
	case *uint64:
		return s.Set(*value)
	case *float32:
		return s.Set(*value)
	case *float64:
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}