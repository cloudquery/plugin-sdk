package schema

import "fmt"

type String struct {
	String string
	Valid  bool
}

func (*String) Type() ValueType {
	return TypeString
}

func (s *String) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*String); ok {
		return s.Valid == other.Valid && s.String == other.String
	}
	return false
}

func (dst *String) Scan(src interface{}) error {
	if src == nil {
		*dst = String{}
		return nil
	}

	switch src := src.(type) {
	case string:
		*dst = String{String: src, Valid: true}
		return nil
	case []byte:
		*dst = String{String: string(src), Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T into String", src)
}