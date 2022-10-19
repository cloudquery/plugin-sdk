package schema

import "fmt"

type ByteArray struct {
	ByteArray []byte
	Valid  bool
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