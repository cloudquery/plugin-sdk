package scalar

import (
	"bytes"

	"github.com/apache/arrow/go/v15/arrow"
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
		return nullValueStr
	}
	return string(s.Value)
}

func (s *Binary) Get() any {
	return s.Value
}

func (s *Binary) Set(val any) error {
	// this will check for typed nils as well, so no need to check below
	if IsNil(val) {
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
	case *[]byte:
		return s.Set(*value)
	case []byte:
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

func (*Binary) DataType() arrow.DataType {
	return arrow.BinaryTypes.Binary
}

type LargeBinary struct {
	Binary
}

func (*LargeBinary) DataType() arrow.DataType {
	return arrow.BinaryTypes.LargeBinary
}
