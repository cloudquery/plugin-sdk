package scalar

import (
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
)

type String struct {
	Valid bool
	Value string
}

func (s *String) IsValid() bool {
	return s.Valid
}

func (*String) DataType() arrow.DataType {
	return arrow.BinaryTypes.String
}

func (s *String) String() string {
	if !s.Valid {
		return "(null)"
	}
	return s.Value
}

func (s *String) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*String)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *String) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	switch value := val.(type) {
	case []byte:
		s.Value = string(value)
	case string:
		s.Value = (value)
	case fmt.Stringer:
		s.Value = value.String()
	case *string:
		if value == nil {
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingStringType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.BinaryTypes.String, Msg: noConversion, Value: value}
	}

	s.Valid = true
	return nil
}

type _smallString String

type LargeString struct {
	_smallString
}

func (*LargeString) DataType() arrow.DataType {
	return arrow.BinaryTypes.LargeString
}
