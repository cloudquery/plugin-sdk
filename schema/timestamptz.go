package schema

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Timestamptz represents the PostgreSQL timestamptz type.
type Timestamptz struct {
	Time             time.Time
	// InfinityModifier InfinityModifier
	Valid            bool
}

func (*Timestamptz) Type() ValueType {
	return TypeTimestamp
}

func (t *Timestamptz) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*Timestamptz); ok {
		return t.Valid == other.Valid && t.Time.Equal(other.Time)
	}
	return false
}

// Scan implements the database/sql Scanner interface.
func (tstz *Timestamptz) Scan(src any) error {
	if src == nil {
		*tstz = Timestamptz{}
		return nil
	}

	switch src := src.(type) {
	case time.Time:
		*tstz = Timestamptz{Time: src, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T into Timestamptz", src)
}

// Value implements the database/sql/driver Valuer interface.
func (tstz Timestamptz) Value() (driver.Value, error) {
	if !tstz.Valid {
		return nil, nil
	}

	return tstz.Time, nil
}
