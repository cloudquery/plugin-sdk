package schema

import (
	"database/sql/driver"
	"encoding/json"
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

func (b Bool) BoolValue() (Bool, error) {
	return b, nil
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

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Bool) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}

	return src.Bool, nil
}

func (src Bool) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}

	if src.Bool {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}

func (dst *Bool) UnmarshalJSON(b []byte) error {
	var v *bool
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = Bool{}
	} else {
		*dst = Bool{Bool: *v, Valid: true}
	}

	return nil
}
