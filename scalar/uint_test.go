package scalar

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
)

func TestUint64Set(t *testing.T) {
	successfulTests := []struct {
		source any
		result Uint
	}{
		{source: int8(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: int16(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: int32(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: int64(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: uint8(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: uint16(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: uint32(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: uint64(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: float32(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: float64(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: "1", result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: _int8(1), result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
		{source: &Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}, result: Uint{Value: 1, Valid: true, Type: arrow.PrimitiveTypes.Uint64}},
	}

	for i, tt := range successfulTests {
		var r Uint
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
