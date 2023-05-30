package scalar

import (
	"strconv"
	"testing"
)

func TestFloat32Set(t *testing.T) {
	successfulTests := []struct {
		source any
		expect Float
	}{
		{source: float32(1), expect: Float{Value: 1, Valid: true}},
		{source: float64(1), expect: Float{Value: 1, Valid: true}},
		{source: int8(1), expect: Float{Value: 1, Valid: true}},
		{source: int16(1), expect: Float{Value: 1, Valid: true}},
		{source: int32(1), expect: Float{Value: 1, Valid: true}},
		{source: int64(1), expect: Float{Value: 1, Valid: true}},
		{source: int8(-1), expect: Float{Value: -1, Valid: true}},
		{source: int16(-1), expect: Float{Value: -1, Valid: true}},
		{source: int32(-1), expect: Float{Value: -1, Valid: true}},
		{source: int64(-1), expect: Float{Value: -1, Valid: true}},
		{source: uint8(1), expect: Float{Value: 1, Valid: true}},
		{source: uint16(1), expect: Float{Value: 1, Valid: true}},
		{source: uint32(1), expect: Float{Value: 1, Valid: true}},
		{source: uint64(1), expect: Float{Value: 1, Valid: true}},
		{source: "1", expect: Float{Value: 1, Valid: true}},
		{source: _int8(1), expect: Float{Value: 1, Valid: true}},
		{source: &Float{Value: 1, Valid: true, BitWidth: 32}, expect: Float{Value: 1, Valid: true}},
	}

	for _, bitWidth := range []uint8{8, 16, 32, 64} {
		bitWidth := bitWidth
		t.Run(strconv.Itoa(int(bitWidth)), func(t *testing.T) {
			t.Parallel()

			for i, tt := range successfulTests {
				r := Float{BitWidth: bitWidth}
				if err := r.Set(tt.source); err != nil {
					t.Errorf("%d: %v", i, err)
				}

				tt.expect.BitWidth = bitWidth
				if !r.Equal(&tt.expect) {
					t.Errorf("%d: %v != %v", i, r, tt.expect)
				}
			}
		})
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
