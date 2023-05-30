package scalar

import (
	"math"
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Uint struct {
	Valid    bool
	Value    uint64
	BitWidth uint8 // defaults to 64
}

func (s *Uint) IsValid() bool {
	return s.Valid
}

func (s *Uint) DataType() arrow.DataType {
	switch s.getBitWidth() {
	case 64:
		return arrow.PrimitiveTypes.Uint64
	case 8:
		return arrow.PrimitiveTypes.Uint8
	case 16:
		return arrow.PrimitiveTypes.Uint16
	case 32:
		return arrow.PrimitiveTypes.Uint32
	default:
		panic("invalid bit width")
	}
}

func (s *Uint) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return strconv.FormatUint(s.Value, 10)
}

func (s *Uint) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Uint)
	if !ok {
		return false
	}
	return s.getBitWidth() == r.getBitWidth() && s.Valid == r.Valid && s.Value == r.Value
}

func (s *Uint) Get() any {
	return s.Value
}

func (s *Uint) Set(val any) error {
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
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "int8 less than 0", Value: value}
		}
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int16:
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "int16 less than 0", Value: value}
		}
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int32:
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "int32 less than 0", Value: value}
		}
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case int64:
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "int64 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case int:
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "int less than 0", Value: value}
		}
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint8:
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint16:
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint32:
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case uint64:
		s.Value = value
	case uint:
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case float32:
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "float32 less than 0", Value: value}
		}
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case float64:
		if value < 0 {
			return &ValidationError{Type: s.DataType(), Msg: "float64 less than 0", Value: value}
		}
		v := uint64(value)
		if err := s.validateValue(v); err != nil {
			return err
		}
		s.Value = v
	case string:
		v, err := strconv.ParseUint(value, 10, 64)
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
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *Uint) validateValue(value uint64) error {
	switch {
	case arrow.TypeEqual(s.DataType(), arrow.PrimitiveTypes.Uint8):
		if value > math.MaxUint8 {
			return &ValidationError{Type: s.DataType(), Msg: "value bigger than MaxUint8", Value: value}
		}
	case arrow.TypeEqual(s.DataType(), arrow.PrimitiveTypes.Uint16):
		if value > math.MaxUint16 {
			return &ValidationError{Type: s.DataType(), Msg: "value bigger than MaxUint16", Value: value}
		}
	case arrow.TypeEqual(s.DataType(), arrow.PrimitiveTypes.Uint32):
		if value > math.MaxUint32 {
			return &ValidationError{Type: s.DataType(), Msg: "value bigger than MaxUint32", Value: value}
		}
	}
	return nil
}

func (s *Uint) getBitWidth() uint8 {
	if s.BitWidth == 0 {
		return 64 // default
	}
	return s.BitWidth
}
