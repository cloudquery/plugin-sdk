package schema

import (
	"testing"
)

func TestByteaSet(t *testing.T) {
	successfulTests := []struct {
		source any
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

func TestBytea_LessThan(t *testing.T) {
	cases := []struct {
		a    Bytea
		b    Bytea
		want bool
	}{
		{a: Bytea{Bytes: []byte{1, 2, 3}, Status: Present}, b: Bytea{Bytes: []byte{1, 2, 4}, Status: Present}, want: true},
		{a: Bytea{Bytes: []byte{1, 2, 4}, Status: Present}, b: Bytea{Bytes: []byte{1, 2, 3}, Status: Present}, want: false},
		{a: Bytea{Bytes: []byte{1, 2, 3}, Status: Undefined}, b: Bytea{Bytes: []byte{1, 2, 3}, Status: Present}, want: true},
		{a: Bytea{Bytes: []byte{1, 2, 3}, Status: Present}, b: Bytea{Bytes: []byte{1, 2, 3}, Status: Undefined}, want: false},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
