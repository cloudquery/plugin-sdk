package scalar

import (
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
	switch s.BitWidth {
	case 0, 64:
		return arrow.PrimitiveTypes.Float64
	case 32:
		return arrow.PrimitiveTypes.Float32
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
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Float) String() string {
	if !s.Valid {
		return "(null)"
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
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}
