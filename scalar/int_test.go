package scalar

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt8Set(t *testing.T) {
	successfulTests := []struct {
		source any
		expect Int
	}{
		{source: int8(1), expect: Int{Value: 1, Valid: true}},
		{source: int16(1), expect: Int{Value: 1, Valid: true}},
		{source: int32(1), expect: Int{Value: 1, Valid: true}},
		{source: int64(1), expect: Int{Value: 1, Valid: true}},
		{source: int8(-1), expect: Int{Value: -1, Valid: true}},
		{source: int16(-1), expect: Int{Value: -1, Valid: true}},
		{source: int32(-1), expect: Int{Value: -1, Valid: true}},
		{source: int64(-1), expect: Int{Value: -1, Valid: true}},
		{source: uint8(1), expect: Int{Value: 1, Valid: true}},
		{source: uint16(1), expect: Int{Value: 1, Valid: true}},
		{source: uint32(1), expect: Int{Value: 1, Valid: true}},
		{source: uint64(1), expect: Int{Value: 1, Valid: true}},
		{source: float32(1), expect: Int{Value: 1, Valid: true}},
		{source: float64(1), expect: Int{Value: 1, Valid: true}},
		{source: "1", expect: Int{Value: 1, Valid: true}},
		{source: _int8(1), expect: Int{Value: 1, Valid: true}},
		{source: &Int{Value: 1, Valid: true}, expect: Int{Value: 1, Valid: true}},
	}

	for _, bitWidth := range []uint8{8, 16, 32, 64} {
		bitWidth := bitWidth
		t.Run(strconv.Itoa(int(bitWidth)), func(t *testing.T) {
			t.Parallel()

			for i, tt := range successfulTests {
				r := Int{BitWidth: bitWidth}
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

func TestIntOverflows(t *testing.T) {
	cases := []struct {
		source      any
		bitWidth    uint8
		expectError bool
	}{
		{source: int16(math.MaxInt8), bitWidth: 8, expectError: false},
		{source: int16(math.MaxInt8 + 1), bitWidth: 8, expectError: true},
		{source: int16(math.MinInt8), bitWidth: 8, expectError: false},
		{source: int16(math.MinInt8 - 1), bitWidth: 8, expectError: true},
		{source: uint16(math.MaxInt8), bitWidth: 8, expectError: false},
		{source: float32(math.MaxInt8), bitWidth: 8, expectError: false},
		{source: float32(math.MaxInt8 + 1), bitWidth: 8, expectError: true},
		{source: float32(math.MinInt8 - 1), bitWidth: 8, expectError: true},
		{source: &Int{Value: math.MaxInt8 + 1, Valid: true}, bitWidth: 8, expectError: true},

		{source: int32(math.MaxInt16), bitWidth: 16, expectError: false},
		{source: int32(math.MaxInt16 + 1), bitWidth: 16, expectError: true},
		{source: int32(math.MinInt16), bitWidth: 16, expectError: false},
		{source: int32(math.MinInt16 - 1), bitWidth: 16, expectError: true},
		{source: uint16(math.MaxInt16), bitWidth: 16, expectError: false},
		{source: float32(math.MaxInt16), bitWidth: 16, expectError: false},
		{source: float32(math.MaxInt16 + 1), bitWidth: 16, expectError: true},
		{source: float32(math.MinInt16 - 1), bitWidth: 16, expectError: true},

		{source: int64(math.MaxInt32), bitWidth: 32, expectError: false},
		{source: int64(math.MaxInt32 + 1), bitWidth: 32, expectError: true},
		{source: int64(math.MinInt32), bitWidth: 32, expectError: false},
		{source: int64(math.MinInt32 - 1), bitWidth: 32, expectError: true},
		{source: uint32(math.MaxInt32), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxInt32), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxInt32 + 1), bitWidth: 32, expectError: true},
		{source: float64(math.MaxInt32), bitWidth: 32, expectError: false},
		{source: float64(math.MaxInt32 + 1), bitWidth: 32, expectError: true},
		{source: float64(math.MinInt32 - 1), bitWidth: 32, expectError: true},
	}
	for idx, tc := range cases {
		tc := tc
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			t.Parallel()

			r := Int{BitWidth: tc.bitWidth}
			err := r.Set(tc.source)
			if tc.expectError {
				assert.Errorf(t, err, "with %T %#v", tc.source, tc.source)
			} else {
				assert.NoErrorf(t, err, "with %T %#v", tc.source, tc.source)
			}
		})
	}
}
