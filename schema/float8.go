//nolint:revive
package schema

import (
	"fmt"
	"math"
	"strconv"
)

type Float8Transformer interface {
	TransformFloat8(*Float8) interface{}
}

type Float8 struct {
	Float  float64
	Status Status
}

const float64EqualityThreshold = 1e-9

func (*Float8) Type() ValueType {
	return TypeFloat
}

func (dst *Float8) String() string {
	if dst.Status == Present {
		return strconv.FormatFloat(dst.Float, 'f', -1, 64)
	} else {
		return ""
	}
}

func (dst *Float8) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Float8)
	if !ok {
		return false
	}
	if dst.Status != s.Status {
		return false
	}
	return math.Abs(dst.Float-s.Float) <= float64EqualityThreshold
}

func (dst *Float8) Set(src interface{}) error {
	if src == nil {
		*dst = Float8{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case float32:
		*dst = Float8{Float: float64(value), Status: Present}
	case float64:
		*dst = Float8{Float: value, Status: Present}
	case int8:
		*dst = Float8{Float: float64(value), Status: Present}
	case uint8:
		*dst = Float8{Float: float64(value), Status: Present}
	case int16:
		*dst = Float8{Float: float64(value), Status: Present}
	case uint16:
		*dst = Float8{Float: float64(value), Status: Present}
	case int32:
		*dst = Float8{Float: float64(value), Status: Present}
	case uint32:
		*dst = Float8{Float: float64(value), Status: Present}
	case int64:
		f64 := float64(value)
		if int64(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case uint64:
		f64 := float64(value)
		if uint64(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case int:
		f64 := float64(value)
		if int(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case uint:
		f64 := float64(value)
		if uint(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case string:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		*dst = Float8{Float: num, Status: Present}
	case *float64:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float32:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int8:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint8:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int16:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint16:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int32:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint32:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int64:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint64:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Float8", value)
	}

	return nil
}

func (dst Float8) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Float
	case Null:
		return nil
	default:
		return dst.Status
	}
}
