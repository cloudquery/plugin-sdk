//nolint:revive
package schema

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type ByteaTransformer interface {
	TransformBytea(*Bytea) interface{}
}

type Bytea struct {
	Bytes  []byte
	Status Status
}

func (*Bytea) Type() ValueType {
	return TypeByteArray
}

func (dst *Bytea) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Bytea)
	if !ok {
		return false
	}

	return dst.Status == s.Status && bytes.Equal(dst.Bytes, s.Bytes)
}

func (dst *Bytea) String() string {
	if dst.Status == Present {
		return hex.EncodeToString(dst.Bytes)
	} else {
		return ""
	}
}

func (dst *Bytea) Set(src interface{}) error {
	if src == nil {
		*dst = Bytea{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case []byte:
		if value != nil {
			*dst = Bytea{Bytes: value, Status: Present}
		} else {
			*dst = Bytea{Status: Null}
		}
	case string:
		if value != "" {
			b := make([]byte, hex.DecodedLen(len(value)))
			_, err := hex.Decode(b, []byte(value))
			if err != nil {
				return fmt.Errorf("cannot decode hex string to bytea: %w", err)
			}
			*dst = Bytea{Status: Present, Bytes: b}
		} else {
			*dst = Bytea{Status: Null}
		}
	default:
		if originalSrc, ok := underlyingBytesType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Bytea", value)
	}

	return nil
}

func (dst Bytea) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Bytes
	case Null:
		return nil
	default:
		return dst.Status
	}
}
