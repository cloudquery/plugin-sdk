package scalar

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var nilPointerInt8 *int8
var nilPointerInt16 *int16
var nilPointerInt32 *int32
var nilPointerInt64 *int64
var nilPointerUint8 *uint8
var nilPointerUint16 *uint16
var nilPointerUint32 *uint32
var nilPointerUint64 *uint64
var nilPointerFloat32 *float32
var nilPointerFloat64 *float64

func TestFloatAllWidthsSet(t *testing.T) {
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
		{source: nilPointerInt8, expect: Float{Valid: false}},
		{source: nilPointerInt16, expect: Float{Valid: false}},
		{source: nilPointerInt32, expect: Float{Valid: false}},
		{source: nilPointerInt64, expect: Float{Valid: false}},
		{source: nilPointerUint8, expect: Float{Valid: false}},
		{source: nilPointerUint16, expect: Float{Valid: false}},
		{source: nilPointerUint32, expect: Float{Valid: false}},
		{source: nilPointerUint64, expect: Float{Valid: false}},
		{source: nilPointerFloat32, expect: Float{Valid: false}},
		{source: nilPointerFloat64, expect: Float{Valid: false}},
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
		{source: nilPointerInt8, result: Float{Valid: false}},
		{source: nilPointerInt16, result: Float{Valid: false}},
		{source: nilPointerInt32, result: Float{Valid: false}},
		{source: nilPointerInt64, result: Float{Valid: false}},
		{source: nilPointerUint8, result: Float{Valid: false}},
		{source: nilPointerUint16, result: Float{Valid: false}},
		{source: nilPointerUint32, result: Float{Valid: false}},
		{source: nilPointerUint64, result: Float{Valid: false}},
		{source: nilPointerFloat32, result: Float{Valid: false}},
		{source: nilPointerFloat64, result: Float{Valid: false}},
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

func TestFloatOverflows(t *testing.T) {
	cases := []struct {
		source      any
		bitWidth    uint8
		expectError bool
	}{
		{source: uint64(math.MaxInt32), bitWidth: 64, expectError: false},
		{source: uint64(math.MaxInt32 + 1), bitWidth: 64, expectError: false},
		{source: uint64(math.MaxInt64), bitWidth: 64, expectError: true},
		{source: uint64(math.MaxUint32), bitWidth: 64, expectError: false},
		{source: uint64(math.MaxUint64), bitWidth: 64, expectError: true},
		{source: int64(math.MinInt64), bitWidth: 64, expectError: true},
		{source: int64(math.MinInt32), bitWidth: 64, expectError: false},
		{source: int64(math.MinInt32 - 1), bitWidth: 64, expectError: false},

		{source: uint64(math.MaxInt16), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxInt16 + 1), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxInt32), bitWidth: 32, expectError: true},
		{source: uint64(math.MaxUint16), bitWidth: 32, expectError: false},
		{source: uint64(math.MaxUint32), bitWidth: 32, expectError: true},
		{source: int64(math.MinInt32), bitWidth: 32, expectError: true},
		{source: int64(math.MinInt16), bitWidth: 32, expectError: false},
		{source: int64(math.MinInt16 - 1), bitWidth: 32, expectError: false},
	}
	for idx, tc := range cases {
		tc := tc
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			t.Parallel()

			r := Float{BitWidth: tc.bitWidth}
			err := r.Set(tc.source)
			if tc.expectError {
				assert.Errorf(t, err, "with %T %#v", tc.source, tc.source)
			} else {
				assert.NoErrorf(t, err, "with %T %#v", tc.source, tc.source)
			}
		})
	}
}
