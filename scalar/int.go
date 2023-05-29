package scalar

import (
	"math"
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Int struct {
	Valid bool
	Value int64
	Type  arrow.DataType
}

func (s *Int) IsValid() bool {
	return s.Valid
}

func (s *Int) DataType() arrow.DataType {
	return s.Type
}

func (s *Int) String() string {
	if !s.Valid {
		return "(null)"
	}
	return strconv.FormatInt(s.Value, 10)
}

func (s *Int) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Int)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Int) Get() any {
	return s.Value
}

func (s *Int) Set(val any) error {
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
	case int8:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int16:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int32:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int64:
		s.Value = value
	case int:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint8:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint16:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint32:
		v := int64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint64:
		if value > math.MaxInt64 {
			return &ValidationError{Type: s.Type, Msg: "uint64 bigger than MaxInt64", Value: value}
		}
		s.Value = int64(value)
	case uint:
		if value > math.MaxInt64 {
			return &ValidationError{Type: s.Type, Msg: "uint bigger than MaxInt64", Value: value}
		}
		s.Value = int64(value)
	case float32:
		s.Value = int64(value)
	case float64:
		s.Value = int64(value)
	case string:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return &ValidationError{Type: s.Type, Msg: "invalid string", Value: v}
		}
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
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
		return &ValidationError{Type: s.Type, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Int) validateValue(value int64) error {
	switch {
	case arrow.TypeEqual(s.Type, arrow.PrimitiveTypes.Int8):
		if value > math.MaxInt8 {
			return &ValidationError{Type: s.Type, Msg: "value bigger than MaxInt8", Value: value}
		}
	case arrow.TypeEqual(s.Type, arrow.PrimitiveTypes.Int16):
		if value > math.MaxInt16 {
			return &ValidationError{Type: s.Type, Msg: "value bigger than MaxInt16", Value: value}
		}
	case arrow.TypeEqual(s.Type, arrow.PrimitiveTypes.Int32):
		if value > math.MaxInt32 {
			return &ValidationError{Type: s.Type, Msg: "value bigger than MaxInt32", Value: value}
		}
	}
	return nil
}
