package scalar

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
)

func TestFloat32Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Float
	}{
		{source: float32(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: float64(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int8(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int16(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int32(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int64(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int8(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int16(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int32(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: int64(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: uint8(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: uint16(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: uint32(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: uint64(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: "1", result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: _int8(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
		{source: &Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}, result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float32}},
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

func TestFloat64Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Float
	}{
		{source: float32(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: float64(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int8(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int16(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int32(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int64(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int8(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int16(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int32(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: int64(-1), result: Float{Value: -1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: uint8(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: uint16(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: uint32(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: uint64(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: "1", result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: _int8(1), result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
		{source: &Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}, result: Float{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Float64}},
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
