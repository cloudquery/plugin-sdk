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

func TestFloat8_Size(t *testing.T) {
	tests := []struct {
		name string
		f    Float8
		want int
	}{
		{
			name: "present",
			f:    Float8{Float: 1, Status: Present},
			want: 8,
		},
		{
			name: "null",
			f:    Float8{Status: Null},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Size(); got != tt.want {
				t.Errorf("Float8.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
