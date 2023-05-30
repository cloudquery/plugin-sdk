package scalar

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/decimal128"
	"github.com/apache/arrow/go/v13/arrow/decimal256"
)

type Decimal256 struct {
	Valid bool
	Value decimal256.Num
}

func (s *Decimal256) IsValid() bool {
	return s.Valid
}

func (*Decimal256) DataType() arrow.DataType {
	return &arrow.Decimal256Type{}
}

func (s *Decimal256) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Decimal256)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Decimal256) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.DataType().String() + " value"
}

func (s *Decimal256) Get() any {
	return s.Value
}

func (s *Decimal256) Set(val any) error {
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
	case decimal256.Num:
		s.Value = value
	case decimal128.Num:
		s.Value = decimal256.FromDecimal128(value)
	case *decimal256.Num:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *decimal128.Num:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	default:
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

type Decimal128 struct {
	Valid bool
	Value decimal128.Num
}

func (s *Decimal128) IsValid() bool {
	return s.Valid
}

func (*Decimal128) DataType() arrow.DataType {
	return &arrow.Decimal128Type{}
}

func (s *Decimal128) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Decimal128)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *Decimal128) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.DataType().String() + " value"
}

func (s *Decimal128) Get() any {
	return s.Value
}

func (s *Decimal128) Set(val any) error {
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
	case decimal128.Num:
		s.Value = value
	case *decimal128.Num:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	default:
		return &ValidationError{Type: s.DataType(), Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}
