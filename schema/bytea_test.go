package schema

import (
	"testing"
)

func TestByteaSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result Bytea
	}{
		{source: []byte{1, 2, 3}, result: Bytea{Bytes: []byte{1, 2, 3}, Status: Present}},
		{source: []byte{}, result: Bytea{Bytes: []byte{}, Status: Present}},
		{source: []byte(nil), result: Bytea{Status: Null}},
		{source: _byteSlice{1, 2, 3}, result: Bytea{Bytes: []byte{1, 2, 3}, Status: Present}},
		{source: _byteSlice(nil), result: Bytea{Status: Null}},
	}

	for i, tt := range successfulTests {
		var r Bytea
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
