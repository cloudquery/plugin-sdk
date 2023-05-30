package scalar

import (
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Uint64 struct {
	Valid bool
	Value uint64
}

func (s *Uint64) IsValid() bool {
	return s.Valid
}

func (*Uint64) DataType() arrow.DataType {
	return arrow.PrimitiveTypes.Uint64
}

func (s *Uint64) String() string {
	if !s.Valid {
		return "(null)"
	}
	return strconv.FormatUint(s.Value, 10)
}

func (s *Uint64) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Uint64)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Uint64) Get() any {
	return s.Value
}

func (s *Uint64) Set(val any) error {
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
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "int8 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case int16:
		if value < 0 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "int16 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case int32:
		if value < 0 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "int32 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case int64:
		if value < 0 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "int64 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case int:
		if value < 0 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "int less than 0", Value: value}
		}
		s.Value = uint64(value)
	case uint8:
		s.Value = uint64(value)
	case uint16:
		s.Value = uint64(value)
	case uint32:
		s.Value = uint64(value)
	case uint64:
		s.Value = value
	case uint:
		s.Value = uint64(value)
	case float32:
		if value < 0 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "float32 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case float64:
		if value < 0 {
			return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: "float64 less than 0", Value: value}
		}
		s.Value = uint64(value)
	case string:
		num, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return &ValidationError{Type: arrow.PrimitiveTypes.Int8, Msg: "invalid string", Value: value}
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
	default:
		if originalSrc, ok := underlyingNumberType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.PrimitiveTypes.Uint64, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}
