package schema

import (
	"testing"
)

func TestInt8Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Int8
	}{
		{source: int8(1), result: Int8{Int: 1, Status: Present}},
		{source: int16(1), result: Int8{Int: 1, Status: Present}},
		{source: int32(1), result: Int8{Int: 1, Status: Present}},
		{source: int64(1), result: Int8{Int: 1, Status: Present}},
		{source: int8(-1), result: Int8{Int: -1, Status: Present}},
		{source: int16(-1), result: Int8{Int: -1, Status: Present}},
		{source: int32(-1), result: Int8{Int: -1, Status: Present}},
		{source: int64(-1), result: Int8{Int: -1, Status: Present}},
		{source: uint8(1), result: Int8{Int: 1, Status: Present}},
		{source: uint16(1), result: Int8{Int: 1, Status: Present}},
		{source: uint32(1), result: Int8{Int: 1, Status: Present}},
		{source: uint64(1), result: Int8{Int: 1, Status: Present}},
		{source: float32(1), result: Int8{Int: 1, Status: Present}},
		{source: float64(1), result: Int8{Int: 1, Status: Present}},
		{source: "1", result: Int8{Int: 1, Status: Present}},
		{source: _int8(1), result: Int8{Int: 1, Status: Present}},
	}

	for i, tt := range successfulTests {
		var r Int8
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestInt8_LessThan(t *testing.T) {
	cases := []struct {
		a    Int8
		b    Int8
		want bool
	}{
		{a: Int8{Int: 0, Status: Present}, b: Int8{Int: 0, Status: Present}, want: false},
		{a: Int8{Int: 1, Status: Present}, b: Int8{Int: 0, Status: Present}, want: false},
		{a: Int8{Int: 0, Status: Present}, b: Int8{Int: 1, Status: Present}, want: true},
		{a: Int8{Int: 1, Status: Undefined}, b: Int8{Int: 1, Status: Present}, want: true},
		{a: Int8{Int: 1, Status: Present}, b: Int8{Int: 1, Status: Undefined}, want: false},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
