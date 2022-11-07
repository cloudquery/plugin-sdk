//nolint:revive
package schema

import (
	"fmt"
	"math"
	"strconv"
)

type Int8Transformer interface {
	TransformInt8(*Int8) interface{}
}

type Int8 struct {
	Int    int64
	Status Status
}

func (*Int8) Type() ValueType {
	return TypeInt
}

func (dst *Int8) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Int8)
	if !ok {
		return false
	}
	return dst.Status == s.Status && dst.Int == s.Int
}

func (dst *Int8) String() string {
	if dst.Status == Present {
		return strconv.FormatInt(dst.Int, 10)
	} else {
		return ""
	}
}

func (dst *Int8) Set(src interface{}) error {
	if src == nil {
		*dst = Int8{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case int8:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint8:
		*dst = Int8{Int: int64(value), Status: Present}
	case int16:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint16:
		*dst = Int8{Int: int64(value), Status: Present}
	case int32:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint32:
		*dst = Int8{Int: int64(value), Status: Present}
	case int64:
		*dst = Int8{Int: value, Status: Present}
	case uint64:
		if value > math.MaxInt64 {
			return fmt.Errorf("%d is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case int:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint:
		if uint64(value) > math.MaxInt64 {
			return fmt.Errorf("%d is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case string:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*dst = Int8{Int: num, Status: Present}
	case float32:
		if value > math.MaxInt64 {
			return fmt.Errorf("%f is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case float64:
		if value > math.MaxInt64 {
			return fmt.Errorf("%f is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case *int8:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint8:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int16:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint16:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int32:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint32:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int64:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint64:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float32:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float64:
		if value == nil {
			*dst = Int8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Int8", value)
	}

	return nil
}

func (dst Int8) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Int
	case Null:
		return nil
	default:
		return dst.Status
	}
}
