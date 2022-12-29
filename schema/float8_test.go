package schema

import (
	"testing"
)

func TestFloat8Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Float8
	}{
		{source: float32(1), result: Float8{Float: 1, Status: Present}},
		{source: float64(1), result: Float8{Float: 1, Status: Present}},
		{source: int8(1), result: Float8{Float: 1, Status: Present}},
		{source: int16(1), result: Float8{Float: 1, Status: Present}},
		{source: int32(1), result: Float8{Float: 1, Status: Present}},
		{source: int64(1), result: Float8{Float: 1, Status: Present}},
		{source: int8(-1), result: Float8{Float: -1, Status: Present}},
		{source: int16(-1), result: Float8{Float: -1, Status: Present}},
		{source: int32(-1), result: Float8{Float: -1, Status: Present}},
		{source: int64(-1), result: Float8{Float: -1, Status: Present}},
		{source: uint8(1), result: Float8{Float: 1, Status: Present}},
		{source: uint16(1), result: Float8{Float: 1, Status: Present}},
		{source: uint32(1), result: Float8{Float: 1, Status: Present}},
		{source: uint64(1), result: Float8{Float: 1, Status: Present}},
		{source: "1", result: Float8{Float: 1, Status: Present}},
		{source: _int8(1), result: Float8{Float: 1, Status: Present}},
	}

	for i, tt := range successfulTests {
		var r Float8
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestFloat8_LessThan(t *testing.T) {
	cases := []struct {
		a    Float8
		b    Float8
		want bool
	}{
		{a: Float8{Float: 1, Status: Present}, b: Float8{Float: 2, Status: Present}, want: true},
		{a: Float8{Float: 2, Status: Present}, b: Float8{Float: 1, Status: Present}, want: false},
		{a: Float8{Float: 1, Status: Undefined}, b: Float8{Float: 1, Status: Present}, want: true},
		{a: Float8{Float: 1, Status: Present}, b: Float8{Float: 1, Status: Undefined}, want: false},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
