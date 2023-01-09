//nolint:revive
package schema

import (
	"database/sql/driver"
	"fmt"
)

type TextTransformer interface {
	TransformText(*Text) any
}

type Text struct {
	Str    string
	Status Status
}

func (dst *Text) GetStatus() Status {
	return dst.Status
}

func (*Text) Type() ValueType {
	return TypeString
}

func (dst *Text) Size() int {
	return len(dst.Str)
}

func (dst *Text) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Text)
	if !ok {
		return false
	}
	return dst.Status == s.Status && dst.Str == s.Str
}

func (dst *Text) String() string {
	if dst.Status == Present {
		return dst.Str
	} else {
		return ""
	}
}

func (dst *Text) Set(src any) error {
	if src == nil {
		*dst = Text{Status: Null}
		return nil
	}

	if value, ok := src.(CQType); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		*dst = Text{Str: value, Status: Present}
	case *string:
		if value == nil {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{Str: *value, Status: Present}
		}
	case []byte:
		if value == nil {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{Str: string(value), Status: Present}
		}
	case fmt.Stringer:
		if value == fmt.Stringer(nil) {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{Str: value.String(), Status: Present}
		}
	default:
		// Cannot be part of the switch: If Value() returns nil on
		// non-string, we should still try to checks the underlying type
		// using reflection.
		//
		// For example the struct might implement driver.Valuer with
		// pointer receiver and fmt.Stringer with value receiver.
		if value, ok := src.(driver.Valuer); ok {
			if value == driver.Valuer(nil) {
				*dst = Text{Status: Null}
				return nil
			} else {
				v, err := value.Value()
				if err != nil {
					return &ValidationError{Type: TypeString, Msg: "driver.Valuer Value() method failed", Err: err, Value: src}
				}

				// Handles also v == nil case.
				if s, ok := v.(string); ok {
					*dst = Text{Str: s, Status: Present}
					return nil
				}
			}
		}

		if originalSrc, ok := underlyingStringType(src); ok {
			return dst.Set(originalSrc)
		}
		return &ValidationError{Type: TypeString, Msg: noConversion, Value: src}
	}

	return nil
}

func (dst Text) Get() any {
	switch dst.Status {
	case Present:
		return dst.Str
	case Null:
		return nil
	default:
		return dst.Status
	}
}
