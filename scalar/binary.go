package scalar

import (
	"bytes"

	"github.com/apache/arrow/go/v13/arrow"
)

type Binary struct {
	Valid bool
	Value []byte
}

func (s *Binary) IsValid() bool {
	return s.Valid
}

func (s *Binary) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Binary)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && bytes.Equal(s.Value, r.Value)
}

func (s *Binary) String() string {
	if !s.Valid {
		return "(null)"
	}
	return string(s.Value)
}

func (s *Binary) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	switch value := val.(type) {
	case []byte:
		if value == nil {
			return nil
		}
		s.Value = value
	case string:
		s.Value = []byte(value)
	case *string:
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingBytesType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: arrow.BinaryTypes.Binary, Msg: noConversion, Value: value}
	}

	s.Valid = true
	return nil
}

func (s *Binary) DataType() arrow.DataType {
	return arrow.BinaryTypes.Binary
}

type LargeBinary struct {
	Binary
}

func (s *LargeBinary) DataType() arrow.DataType {
	return arrow.BinaryTypes.LargeBinary
}
