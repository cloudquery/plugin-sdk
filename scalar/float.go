package scalar

import (
	"math"
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Float struct {
	Valid    bool
	Value    float64
	BitWidth uint8 // defaults to 64
}

func (s *Float) IsValid() bool {
	return s.Valid
}

func (s *Float) DataType() arrow.DataType {
	switch s.getBitWidth() {
	case 64:
		return arrow.PrimitiveTypes.Float64
	case 32:
		return arrow.PrimitiveTypes.Float32
	case 16:
		return arrow.FixedWidthTypes.Float16
	default:
		panic("invalid bit width")
	}
}

func (s *Float) Get() any {
	return s.Value
}

func (s *Float) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Float)
	if !ok {
		return false
	}
	return s.getBitWidth() == r.getBitWidth() && s.Valid == r.Valid && s.Value == r.Value
}

func (s *Float) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return strconv.FormatFloat(s.Value, 'f', -1, 64)
}

func (s *Float) Set(val any) error {
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
		s.Value = float64(value)
	case int16:
		s.Value = float64(value)
	case int32:
		v := float64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int64:
		v := float64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint8:
		s.Value = float64(value)
	case uint16:
		s.Value = float64(value)
	case uint32:
		v := float64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint64:
		v := float64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case float32:
		s.Value = float64(value)
	case float64:
		if err := s.validateValue(value); err != nil {
			return err
		}
		s.Value = value
	case string:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return &ValidationError{Type: s.DataType(), Msg: "invalid string", Value: value}
		}
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case *string:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
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
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Float) validateValue(value float64) error {
	const maxFloat16 = 65504.0

	switch s.getBitWidth() {
	case 16:
		if value > maxFloat16 {
			return &ValidationError{Type: s.DataType(), Msg: "value bigger than maxFloat16", Value: value}
		}
	case 32:
		if value > math.MaxFloat32 {
			return &ValidationError{Type: s.DataType(), Msg: "value bigger than MaxFloat32", Value: value}
		}
	}
	return nil
}

func (s *Float) getBitWidth() uint8 {
	if s.BitWidth == 0 {
		return 64 // default
	}
	return s.BitWidth
}
