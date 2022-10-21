package schema

import (
	"fmt"
)

type Int64 struct {
	Int64 int64
	Valid bool
}

func (*Int64) Type() ValueType {
	return TypeInt
}

func (i *Int64) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*Int64); ok {
		return i.Valid == other.Valid && i.Int64 == other.Int64
	}
	return false
}

// ScanInt64 implements the Int64Scanner interface.
func (dst *Int64) Scan(src any) error {
	if src == nil {
		*dst = Int64{}
		return nil
	}

	var n int64

	switch src := src.(type) {
	case int64:
		n = src
	case int32:
		n = int64(src)
	case int16:
		n = int64(src)
	case int8:
		n = int64(src)
	case int:
		n = int64(src)
	case uint64:
		n = int64(src)
	case uint32:
		n = int64(src)
	case uint16:
		n = int64(src)
	case uint8:
		n = int64(src)
	case uint:
		n = int64(src)
	default:
		return fmt.Errorf("cannot scan %T into Int64", src)
	}

	*dst = Int64{Int64: n, Valid: true}

	return nil
}