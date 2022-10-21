package schema

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)


type BoolValuer interface {
	BoolValue() (Bool, error)
}

type Bool struct {
	Bool  bool
	Valid bool
}

func (*Bool) Type() ValueType {
	return TypeBool
}

func (b *Bool) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*Bool); ok {
		return b.Valid == other.Valid && b.Bool == other.Bool 
	}
	return false
}

// Scan implements the database/sql Scanner interface.
func (dst *Bool) Scan(src any) error {
	if src == nil {
		*dst = Bool{}
		return nil
	}

	switch src := src.(type) {
	case bool:
		*dst = Bool{Bool: src, Valid: true}
		return nil
	case string:
		b, err := strconv.ParseBool(src)
		if err != nil {
			return err
		}
		*dst = Bool{Bool: b, Valid: true}
		return nil
	case []byte:
		b, err := strconv.ParseBool(string(src))
		if err != nil {
			return err
		}
		*dst = Bool{Bool: b, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T into Bool", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Bool) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}

	return src.Bool, nil
}
