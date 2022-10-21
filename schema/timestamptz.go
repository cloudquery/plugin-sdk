package schema

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Timestamptz represents the PostgreSQL timestamptz type.
type Timestamptz struct {
	Time time.Time
	// InfinityModifier InfinityModifier
	Valid bool
}

func (*Timestamptz) Type() ValueType {
	return TypeTimestamp
}

func (dst *Timestamptz) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*Timestamptz); ok {
		return dst.Valid == other.Valid && dst.Time.Equal(other.Time)
	}
	return false
}

// Scan implements the database/sql Scanner interface.
func (dst *Timestamptz) Scan(src any) error {
	if src == nil {
		*dst = Timestamptz{}
		return nil
	}

	switch src := src.(type) {
	case time.Time:
		*dst = Timestamptz{Time: src, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T into Timestamptz", src)
}

// Value implements the database/sql/driver Valuer interface.
func (dst Timestamptz) Value() (driver.Value, error) {
	if !dst.Valid {
		return nil, nil
	}

	return dst.Time, nil
}
