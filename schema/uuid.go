package schema

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

type UUIDScanner interface {
	ScanUUID(v UUID) error
}

type UUIDValuer interface {
	UUIDValue() (UUID, error)
}

type UUID struct {
	Bytes [16]byte
	Valid bool
}

func NewMustUUID(src any) *UUID {
	r := &UUID{}
	if err := r.Scan(src); err != nil {
		panic(err)
	}
	return r
}

func (*UUID) Type() ValueType {
	return TypeUUID
}

func (b *UUID) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*UUID); ok {
		return b.Valid == other.Valid && b.Bytes == other.Bytes
	}
	return false
}

func (b UUID) UUIDValue() (UUID, error) {
	return b, nil
}

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		// assume invalid.
		return dst, fmt.Errorf("cannot parse UUID %v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

// Scan implements the database/sql Scanner interface.
func (dst *UUID) Scan(src any) error {
	if src == nil {
		*dst = UUID{}
		return nil
	}

	switch src := src.(type) {
	case string:
		buf, err := parseUUID(src)
		if err != nil {
			return err
		}
		*dst = UUID{Bytes: buf, Valid: true}
	case []byte:
		if len(src) != 16 {
			return fmt.Errorf("cannot scan %T to UUID. incorrect length %d", src, len(src))
		}
		*dst = UUID{Bytes: [16]byte{}, Valid: true}
		copy(dst.Bytes[:], src)
	case [16]byte:
		*dst = UUID{Bytes: [16]byte{}, Valid: true}
		copy(dst.Bytes[:], src[:])
	default:
		return fmt.Errorf("cannot scan %T to UUID", src)
	}

	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (src UUID) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}

	return encodeUUID(src.Bytes), nil
}

