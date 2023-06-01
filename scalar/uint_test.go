package scalar

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64Set(t *testing.T) {
	successfulTests := []struct {
		source any
		expect Uint
	}{
		{source: int8(1), expect: Uint{Value: 1, Valid: true}},
		{source: int16(1), expect: Uint{Value: 1, Valid: true}},
		{source: int32(1), expect: Uint{Value: 1, Valid: true}},
		{source: int64(1), expect: Uint{Value: 1, Valid: true}},
		{source: uint8(1), expect: Uint{Value: 1, Valid: true}},
		{source: uint16(1), expect: Uint{Value: 1, Valid: true}},
		{source: uint32(1), expect: Uint{Value: 1, Valid: true}},
		{source: uint64(1), expect: Uint{Value: 1, Valid: true}},
		{source: float32(1), expect: Uint{Value: 1, Valid: true}},
		{source: float64(1), expect: Uint{Value: 1, Valid: true}},
		{source: "1", expect: Uint{Value: 1, Valid: true}},
		{source: _int8(1), expect: Uint{Value: 1, Valid: true}},
		{source: &Uint{Value: 1, Valid: true}, expect: Uint{Value: 1, Valid: true}},
	}

	for _, bitWidth := range []uint8{8, 16, 32, 64} {
		bitWidth := bitWidth
		t.Run(strconv.Itoa(int(bitWidth)), func(t *testing.T) {
			t.Parallel()

			for i, tt := range successfulTests {
				r := Uint{BitWidth: bitWidth}
				err := r.Set(tt.source)
				if err != nil {
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

func TestUintOverflows(t *testing.T) {
	cases := []struct {
		source      any
		bitWidth    uint8
		expectError bool
	}{
		{source: int16(math.MaxUint8), bitWidth: 8, expectError: false},
		{source: int16(math.MaxUint8 + 1), bitWidth: 8, expectError: true},
		{source: uint16(math.MaxUint8), bitWidth: 8, expectError: false},
		{source: float32(math.MaxUint8), bitWidth: 8, expectError: false},
		{source: float32(math.MaxUint8 + 1), bitWidth: 8, expectError: true},
		{source: &Uint{Value: math.MaxUint8 + 1, Valid: true}, bitWidth: 8, expectError: true},

		{source: int32(math.MaxUint16), bitWidth: 16, expectError: false},
		{source: int32(math.MaxUint16 + 1), bitWidth: 16, expectError: true},
		{source: uint16(math.MaxUint16), bitWidth: 16, expectError: false},
		{source: float32(math.MaxUint16), bitWidth: 16, expectError: false},
		{source: float32(math.MaxUint16 + 1), bitWidth: 16, expectError: true},

		{source: int64(math.MaxUint32), bitWidth: 32, expectError: false},
		{source: int64(math.MaxUint32 + 1), bitWidth: 32, expectError: true},
		{source: uint32(math.MaxUint32), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxUint32), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxUint32 + 1), bitWidth: 32, expectError: true},
		{source: float64(math.MaxUint32), bitWidth: 32, expectError: false},
		{source: float64(math.MaxUint32 + 1), bitWidth: 32, expectError: true},

		{source: float32(math.MaxFloat32), bitWidth: 64, expectError: true},
		{source: float64(math.MaxFloat64), bitWidth: 64, expectError: true},
	}
	for idx, tc := range cases {
		tc := tc
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			t.Parallel()

			r := Uint{BitWidth: tc.bitWidth}
			err := r.Set(tc.source)
			if tc.expectError {
				assert.Errorf(t, err, "with %T %#v", tc.source, tc.source)
			} else {
				assert.NoErrorf(t, err, "with %T %#v", tc.source, tc.source)
			}
		})
	}
}
