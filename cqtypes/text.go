package cqtypes

import (
	"database/sql/driver"
	"fmt"
)

type Text struct {
	String string
	Status Status
}

func (dst *Text) Set(src interface{}) error {
	if src == nil {
		*dst = Text{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		*dst = Text{String: value, Status: Present}
	case *string:
		if value == nil {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{String: *value, Status: Present}
		}
	case []byte:
		if value == nil {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{String: string(value), Status: Present}
		}
	case fmt.Stringer:
		if value == fmt.Stringer(nil) {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{String: value.String(), Status: Present}
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
					return fmt.Errorf("driver.Valuer Value() method failed: %w", err)
				}

				// Handles also v == nil case.
				if s, ok := v.(string); ok {
					*dst = Text{String: s, Status: Present}
					return nil
				}
			}
		}

		if originalSrc, ok := underlyingStringType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Text", value)
	}

	return nil
}

func (dst Text) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.String
	case Null:
		return nil
	default:
		return dst.Status
	}
}
