package schema

import (
	"bytes"
	"fmt"
)

type ByteArray struct {
	ByteArray []byte
	Valid  bool
}

func (*ByteArray) Type() ValueType {
	return TypeByteArray
}

func (b *ByteArray) Equal(other CQType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(*ByteArray); ok {
		return  b.Valid == other.Valid && bytes.Compare(b.ByteArray, other.ByteArray) == 0
	}
	return false
}

func (dst *ByteArray) Scan(src interface{}) error {
	if src == nil {
		*dst = ByteArray{}
		return nil
	}

	switch src := src.(type) {
	case []byte:
		dstBuf := make([]byte, len(src))
		copy(dstBuf, src)
		*dst = ByteArray{ByteArray: dstBuf, Valid: true}
	}

	return fmt.Errorf("cannot scan %T", src)
}