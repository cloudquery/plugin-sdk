package scalar

import (
	"encoding/hex"
	"fmt"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
)

type UUID struct {
	Valid bool
	Value uuid.UUID
}

func (s *UUID) IsValid() bool {
	return s.Valid
}

func (*UUID) DataType() arrow.DataType {
	return types.ExtensionTypes.UUID
}

func (s *UUID) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.Value.String()
}

func (s *UUID) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*UUID)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value == r.Value
}

func (s *UUID) Get() any {
	return s.Value
}

func (s *UUID) Set(src any) error {
	if src == nil {
		return nil
	}

	if sc, ok := src.(Scalar); ok {
		if !sc.IsValid() {
			s.Valid = false
			return nil
		}
		return s.Set(sc.Get())
	}

	switch value := src.(type) {
	case fmt.Stringer:
		value2 := value.String()
		return s.Set(value2)
	case [16]byte:
		s.Value = uuid.UUID(value)
	case *[]byte:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case []byte:
		if value == nil {
			s.Valid = false
			return nil
		}
		if len(value) != 16 {
			return &ValidationError{Type: types.ExtensionTypes.UUID, Msg: "[]byte must be 16 bytes to convert to UUID", Value: value}
		}
		copy(s.Value[:], value)
	case string:
		uuidVal, err := parseUUID(value)
		if err != nil {
			return err
		}
		s.Value = uuidVal
	case *string:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingUUIDType(src); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: types.ExtensionTypes.UUID, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

func (s *UUID) ByteSize() int64 { return int64(len(s.Value)) }

var (
	_ Scalar = (*UUID)(nil)
)

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		// assume invalid.
		return dst, &ValidationError{Type: types.ExtensionTypes.UUID, Msg: fmt.Sprintf("invalid %d UUID length", len(src)), Value: src}
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}
