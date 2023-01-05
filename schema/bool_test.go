package schema

import (
	"testing"
)

func TestBoolSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result Bool
	}{
		{source: true, result: Bool{Bool: true, Status: Present}},
		{source: false, result: Bool{Bool: false, Status: Present}},
		{source: "true", result: Bool{Bool: true, Status: Present}},
		{source: "false", result: Bool{Bool: false, Status: Present}},
		{source: "t", result: Bool{Bool: true, Status: Present}},
		{source: "f", result: Bool{Bool: false, Status: Present}},
		{source: _bool(true), result: Bool{Bool: true, Status: Present}},
		{source: _bool(false), result: Bool{Bool: false, Status: Present}},
		{source: nil, result: Bool{Status: Null}},
	}

	for i, tt := range successfulTests {
		var r Bool
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestBool_Size(t *testing.T) {
	tests := []struct {
		name string
		b    Bool
		want int
	}{
		{
			name: "present",
			b:    Bool{Bool: true, Status: Present},
			want: 1,
		},
		{
			name: "null",
			b:    Bool{Status: Null},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Size(); got != tt.want {
				t.Errorf("Bool.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
