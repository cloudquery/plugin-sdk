//nolint:revive
package schema

import (
	"strconv"
)

type BoolTransformer interface {
	TransformBool(*Bool) any
}

type Bool struct {
	Bool   bool
	Status Status
}

func (dst *Bool) GetStatus() Status {
	return dst.Status
}

func (*Bool) Type() ValueType {
	return TypeBool
}

func (*Bool) Size() int {
	return 1
}

func (dst *Bool) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Bool)
	if !ok {
		return false
	}
	return dst.Status == s.Status && dst.Bool == s.Bool
}

func (dst *Bool) String() string {
	if dst.Status == Present {
		if dst.Bool {
			return "true"
		} else {
			return "false"
		}
	} else {
		return ""
	}
}

func (dst *Bool) Set(src any) error {
	if src == nil {
		*dst = Bool{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() any }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case bool:
		*dst = Bool{Bool: value, Status: Present}
	case string:
		bb, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*dst = Bool{Bool: bb, Status: Present}
	case *bool:
		if value == nil {
			*dst = Bool{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Bool{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingBoolType(src); ok {
			return dst.Set(originalSrc)
		}
		return &ValidationError{Type: TypeBool, Msg: noConversion, Value: src}
	}

	return nil
}

func (dst Bool) Get() any {
	switch dst.Status {
	case Present:
		return dst.Bool
	case Null:
		return nil
	default:
		return dst.Status
	}
}
