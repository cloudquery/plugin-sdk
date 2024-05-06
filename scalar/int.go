package scalar

import (
	"math"
	"strconv"

	"github.com/apache/arrow/go/v16/arrow"
)

type Int struct {
	Valid    bool
	Value    int64
	BitWidth uint8 // defaults to 64
}

func (s *Int) IsValid() bool {
	return s.Valid
}

func (s *Int) DataType() arrow.DataType {
	switch s.getBitWidth() {
	case 64:
		return arrow.PrimitiveTypes.Int64
	case 32:
		return arrow.PrimitiveTypes.Int32
	case 16:
		return arrow.PrimitiveTypes.Int16
	case 8:
		return arrow.PrimitiveTypes.Int8
	default:
		panic("invalid bit width")
	}
}

func (s *Int) String() string {
	if !s.Valid {
		return nullValueStr
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
	return s.getBitWidth() == r.getBitWidth() && s.Valid == r.Valid && s.Value == r.Value
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
		if err := s.validateValue(value); err != nil {
			return err
		}
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
		return s.Set(int64(value))
	case uint:
		return s.Set(int64(value))
	case float32:
		return s.Set(int64(value))
	case float64:
		return s.Set(int64(value))
	case string:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return &ValidationError{Type: s.DataType(), Msg: "invalid string", Value: value}
		}
		return s.Set(v)
	case *int8:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *int16:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *int32:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *int64:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *int:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *uint8:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *uint16:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *uint32:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *uint64:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *uint:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *float32:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *float64:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *string:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Int) validateValue(value int64) error {
	switch s.getBitWidth() {
	case 8:
		if value > math.MaxInt8 {
			return &ValidationError{Type: s.DataType(), Msg: "value greater than MaxInt8", Value: value}
		}
		if value < math.MinInt8 {
			return &ValidationError{Type: s.DataType(), Msg: "value less than MinInt8", Value: value}
		}
	case 16:
		if value > math.MaxInt16 {
			return &ValidationError{Type: s.DataType(), Msg: "value greater than MaxInt16", Value: value}
		}
		if value < math.MinInt16 {
			return &ValidationError{Type: s.DataType(), Msg: "value less than MinInt16", Value: value}
		}
	case 32:
		if value > math.MaxInt32 {
			return &ValidationError{Type: s.DataType(), Msg: "value greater than MaxInt32", Value: value}
		}
		if value < math.MinInt32 {
			return &ValidationError{Type: s.DataType(), Msg: "value less than MinInt32", Value: value}
		}
	}
	return nil
}

func (s *Int) getBitWidth() uint8 {
	if s.BitWidth == 0 {
		return 64 // default
	}
	return s.BitWidth
}
