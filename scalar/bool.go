package scalar

import (
	"strconv"

	"github.com/apache/arrow/go/v13/arrow"
)

type Bool struct {
	Valid bool
	Value bool
}

func (s *Bool) IsValid() bool {
	return s.Valid
}

func (*Bool) DataType() arrow.DataType {
	return arrow.FixedWidthTypes.Boolean
}

func (s *Bool) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Bool)
	if !ok {
		return false
	}
	return s.Value == r.Value && s.Valid == r.Valid
}

func (s *Bool) String() string {
	if !s.Valid {
		return "(null)"
	}
	return strconv.FormatBool(s.Value)
}

func (s *Bool) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}
	switch value := val.(type) {
	case bool:
		s.Value = value
	case string:
		bb, err := strconv.ParseBool(value)
		if err != nil {
			return &ValidationError{Type: arrow.FixedWidthTypes.Boolean, Msg: "failed to ParseBool", Value: value, Err: err}
		}
		s.Value = bb
	case *bool:
		return s.Set(*value)
	case *string:
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingBoolType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.FixedWidthTypes.Boolean, Msg: noConversion, Value: val}
	}
	s.Valid = true
	return nil
}
