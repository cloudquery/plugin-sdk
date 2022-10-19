package schema

import (
	"fmt"
)

type Int64 struct {
	Int64 int64
	Valid bool
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
		return fmt.Errorf("cannot scan %T", src)
	}

	*dst = Int64{Int64: n, Valid: true}

	return nil
}