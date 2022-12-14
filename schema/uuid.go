//nolint:revive
package schema

import (
	"encoding/hex"
	"fmt"
)

type UUIDTransformer interface {
	TransformUUID(*UUID) any
}

type UUID struct {
	Bytes  [16]byte
	Status Status
}

func (*UUID) Type() ValueType {
	return TypeUUID
}

func (dst *UUID) Size() int {
	return 16
}

func (dst *UUID) GetStatus() Status {
	return dst.Status
}

func (dst *UUID) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*UUID)
	if !ok {
		return false
	}

	return dst.Status == s.Status && dst.Bytes == s.Bytes
}

func (dst *UUID) String() string {
	if dst.Status == Present {
		return hex.EncodeToString(dst.Bytes[:])
	} else {
		return ""
	}
}

func (dst *UUID) Set(src any) error {
	if src == nil {
		*dst = UUID{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case interface{ Get() any }:
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	case fmt.Stringer:
		value2 := value.String()
		return dst.Set(value2)
	case [16]byte:
		*dst = UUID{Bytes: value, Status: Present}
	case []byte:
		if value != nil {
			if len(value) != 16 {
				return &ValidationError{Type: TypeUUID, Msg: "[]byte must be 16 bytes to convert to UUID", Value: value}
			}
			*dst = UUID{Status: Present}
			copy(dst.Bytes[:], value)
		} else {
			*dst = UUID{Status: Null}
		}
	case string:
		uuid, err := parseUUID(value)
		if err != nil {
			return err
		}
		*dst = UUID{Bytes: uuid, Status: Present}
	case *string:
		if value == nil {
			*dst = UUID{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingUUIDType(src); ok {
			return dst.Set(originalSrc)
		}
		return &ValidationError{Type: TypeUUID, Msg: noConversion, Value: value}
	}

	return nil
}

func (dst UUID) Get() any {
	switch dst.Status {
	case Present:
		return dst.Bytes
	case Null:
		return nil
	default:
		return dst.Status
	}
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
		return dst, &ValidationError{Type: TypeUUID, Msg: fmt.Sprintf("invalid %d UUID length", len(src)), Value: src}
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}
