package scalar

import (
	"fmt"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
)

const nullValueStr = array.NullValueStr

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
		return nullValueStr
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

func (s *String) Get() any {
	return s.Value
}

func (s *String) Set(val any) error {
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
	case []byte:
		s.Value = string(value)
	case string:
		s.Value = (value)
	case fmt.Stringer:
		s.Value = value.String()
	case *string:
		if value == nil {
			s.Valid = false
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

func (s *String) ByteSize() int64 { return int64(len(s.Value)) }

type LargeString struct {
	String
}

func (*LargeString) DataType() arrow.DataType {
	return arrow.BinaryTypes.LargeString
}

func (s *LargeString) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*LargeString)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

var (
	_ Scalar = (*String)(nil)
	_ Scalar = (*LargeString)(nil)
)
