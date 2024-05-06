package scalar

import (
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/decimal128"
	"github.com/apache/arrow/go/v16/arrow/decimal256"
	"github.com/stretchr/testify/require"
)

// nolint:dupl
func TestDecimal128Set(t *testing.T) {
	str := "100.32"
	decimalType := &arrow.Decimal128Type{Precision: 5, Scale: 2}
	strDecimal, _ := decimal128.FromString(str, decimalType.Precision, decimalType.Scale)

	intVal := int(1)
	int8Val := int8(1)
	int16Val := int16(1)
	int32Val := int32(1)
	int64Val := int64(1)
	uintVal := uint(1)
	uint8Val := uint8(1)
	uint16Val := uint16(1)
	uint32Val := uint32(1)
	uint64Val := uint64(1)

	intValPointer := &intVal

	successfulTests := []struct {
		source      any
		decimalType *arrow.Decimal128Type
		expect      Decimal128
	}{
		{source: str, expect: Decimal128{Value: strDecimal, Valid: true, Type: decimalType}, decimalType: decimalType},
		{source: &str, expect: Decimal128{Value: strDecimal, Valid: true, Type: decimalType}, decimalType: decimalType},
		{source: intVal, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: int8Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: int16Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: int32Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: int64Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: uintVal, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: uint8Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: uint16Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: uint32Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: uint64Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: &intVal, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: &int8Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: &int16Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: &int32Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: &int64Val, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
		{source: &uintVal, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: &uint8Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: &uint16Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: &uint32Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: &uint64Val, expect: Decimal128{Value: decimal128.FromU64(1), Valid: true}},
		{source: &intValPointer, expect: Decimal128{Value: decimal128.FromI64(1), Valid: true}},
	}

	for i, tt := range successfulTests {
		r := Decimal128{}
		r.Type = tt.decimalType
		err := r.Set(tt.source)
		require.NoError(t, err, "No error expected for test %d", i)
		require.Equal(t, tt.expect, r, "Unexpected result for test %d", i)
	}
}

// nolint:dupl
func TestDecimal256Set(t *testing.T) {
	str := "100.32"
	decimalType := &arrow.Decimal256Type{Precision: 5, Scale: 2}
	strDecimal, _ := decimal256.FromString(str, decimalType.Precision, decimalType.Scale)

	intVal := int(1)
	int8Val := int8(1)
	int16Val := int16(1)
	int32Val := int32(1)
	int64Val := int64(1)
	uintVal := uint(1)
	uint8Val := uint8(1)
	uint16Val := uint16(1)
	uint32Val := uint32(1)
	uint64Val := uint64(1)

	intValPointer := &intVal

	successfulTests := []struct {
		source      any
		decimalType *arrow.Decimal256Type
		expect      Decimal256
	}{
		{source: str, expect: Decimal256{Value: strDecimal, Valid: true, Type: decimalType}, decimalType: decimalType},
		{source: &str, expect: Decimal256{Value: strDecimal, Valid: true, Type: decimalType}, decimalType: decimalType},
		{source: intVal, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: int8Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: int16Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: int32Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: int64Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: uintVal, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: uint8Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: uint16Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: uint32Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: uint64Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: &intVal, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: &int8Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: &int16Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: &int32Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: &int64Val, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
		{source: &uintVal, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: &uint8Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: &uint16Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: &uint32Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: &uint64Val, expect: Decimal256{Value: decimal256.FromU64(1), Valid: true}},
		{source: &intValPointer, expect: Decimal256{Value: decimal256.FromI64(1), Valid: true}},
	}

	for i, tt := range successfulTests {
		r := Decimal256{}
		r.Type = tt.decimalType
		err := r.Set(tt.source)
		require.NoError(t, err, "No error expected for test %d", i)
		require.Equal(t, tt.expect, r, "Unexpected result for test %d", i)
	}
}
