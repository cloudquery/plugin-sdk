package cqtypes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

		if diff := cmp.Diff(r, tt.result); diff != "" {
			t.Errorf("%d: got diff:\n%s", i, diff)
		}
	}
}
