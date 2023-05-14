package scalar

import "testing"

func TestFloat64Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Float64
	}{
		{source: float32(1), result: Float64{Value: 1, Valid: true}},
		{source: float64(1), result: Float64{Value: 1, Valid: true}},
		{source: int8(1), result: Float64{Value: 1, Valid: true}},
		{source: int16(1), result: Float64{Value: 1, Valid: true}},
		{source: int32(1), result: Float64{Value: 1, Valid: true}},
		{source: int64(1), result: Float64{Value: 1, Valid: true}},
		{source: int8(-1), result: Float64{Value: -1, Valid: true}},
		{source: int16(-1), result: Float64{Value: -1, Valid: true}},
		{source: int32(-1), result: Float64{Value: -1, Valid: true}},
		{source: int64(-1), result: Float64{Value: -1, Valid: true}},
		{source: uint8(1), result: Float64{Value: 1, Valid: true}},
		{source: uint16(1), result: Float64{Value: 1, Valid: true}},
		{source: uint32(1), result: Float64{Value: 1, Valid: true}},
		{source: uint64(1), result: Float64{Value: 1, Valid: true}},
		{source: "1", result: Float64{Value: 1, Valid: true}},
		{source: _int8(1), result: Float64{Value: 1, Valid: true}},
	}

	for i, tt := range successfulTests {
		var r Float64
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}