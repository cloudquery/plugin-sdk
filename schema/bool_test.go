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

func TestBool_LessThan(t *testing.T) {
	cases := []struct {
		a    Bool
		b    Bool
		want bool
	}{
		{a: Bool{Bool: true, Status: Present}, b: Bool{Bool: false, Status: Present}, want: false},
		{a: Bool{Bool: false, Status: Present}, b: Bool{Bool: true, Status: Present}, want: true},
		{a: Bool{Bool: true, Status: Undefined}, b: Bool{Bool: true, Status: Present}, want: true},
		{a: Bool{Bool: true, Status: Present}, b: Bool{Bool: true, Status: Undefined}, want: false},
		{a: Bool{Bool: true, Status: Present}, b: Bool{Bool: true, Status: Present}, want: false},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
