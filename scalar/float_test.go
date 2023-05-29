package scalar

import (
	"testing"
)

func TestFloat32Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Float
	}{
		{source: float32(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: float64(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: int8(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: int16(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: int32(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: int64(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: int8(-1), result: Float{Value: -1, Valid: true, BitWidth: 32}},
		{source: int16(-1), result: Float{Value: -1, Valid: true, BitWidth: 32}},
		{source: int32(-1), result: Float{Value: -1, Valid: true, BitWidth: 32}},
		{source: int64(-1), result: Float{Value: -1, Valid: true, BitWidth: 32}},
		{source: uint8(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: uint16(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: uint32(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: uint64(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: "1", result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: _int8(1), result: Float{Value: 1, Valid: true, BitWidth: 32}},
		{source: &Float{Value: 1, Valid: true, BitWidth: 32}, result: Float{Value: 1, Valid: true, BitWidth: 32}},
	}

	for i, tt := range successfulTests {
		r := Float{BitWidth: 32}
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestFloat64Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Float
	}{
		{source: float32(1), result: Float{Value: 1, Valid: true}},
		{source: float64(1), result: Float{Value: 1, Valid: true}},
		{source: int8(1), result: Float{Value: 1, Valid: true}},
		{source: int16(1), result: Float{Value: 1, Valid: true}},
		{source: int32(1), result: Float{Value: 1, Valid: true}},
		{source: int64(1), result: Float{Value: 1, Valid: true}},
		{source: int8(-1), result: Float{Value: -1, Valid: true}},
		{source: int16(-1), result: Float{Value: -1, Valid: true}},
		{source: int32(-1), result: Float{Value: -1, Valid: true}},
		{source: int64(-1), result: Float{Value: -1, Valid: true}},
		{source: uint8(1), result: Float{Value: 1, Valid: true}},
		{source: uint16(1), result: Float{Value: 1, Valid: true}},
		{source: uint32(1), result: Float{Value: 1, Valid: true}},
		{source: uint64(1), result: Float{Value: 1, Valid: true}},
		{source: "1", result: Float{Value: 1, Valid: true}},
		{source: _int8(1), result: Float{Value: 1, Valid: true}},
		{source: &Float{Value: 1, Valid: true}, result: Float{Value: 1, Valid: true}},
	}

	for i, tt := range successfulTests {
		var r Float
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
