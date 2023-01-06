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

func TestBytea_Size(t *testing.T) {
	tests := []struct {
		name string
		b    Bytea
		want int
	}{
		{
			name: "present",
			b:    Bytea{Bytes: []byte{1, 2, 3}, Status: Present},
			want: 3,
		},
		{
			name: "null",
			b:    Bytea{Status: Null},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Size(); got != tt.want {
				t.Errorf("Bytea.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
