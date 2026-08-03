package plugin

import (
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

func RecordsDiff(sc *arrow.Schema, have, want []arrow.RecordBatch) string {
	return TableDiff(array.NewTableFromRecords(sc, have), array.NewTableFromRecords(sc, want))
}

func containsInt64(dt arrow.DataType) bool {
	switch dt.ID() {
	case arrow.INT64, arrow.UINT64:
		return true
	}
	if nested, ok := dt.(arrow.NestedType); ok {
		for _, field := range nested.Fields() {
			if containsInt64(field.Type) {
				return true
			}
		}
	}
	return false
}

func equalAsFloat64(want, have arrow.Array) bool {
	if want.Len() != have.Len() || !arrow.TypeEqual(want.DataType(), have.DataType()) {
		return false
	}
	for i := 0; i < want.Len(); i++ {
		if want.IsNull(i) != have.IsNull(i) {
			return false
		}
	}

	switch wantCol := want.(type) {
	case *array.Int64:
		haveCol := have.(*array.Int64)
		for i := 0; i < wantCol.Len(); i++ {
			if !wantCol.IsNull(i) && float64(wantCol.Value(i)) != float64(haveCol.Value(i)) {
				return false
			}
		}
		return true
	case *array.Uint64:
		haveCol := have.(*array.Uint64)
		for i := 0; i < wantCol.Len(); i++ {
			if !wantCol.IsNull(i) && float64(wantCol.Value(i)) != float64(haveCol.Value(i)) {
				return false
			}
		}
		return true
	case *array.Struct:
		haveCol := have.(*array.Struct)
		for i := 0; i < wantCol.NumField(); i++ {
			if !equalAsFloat64(wantCol.Field(i), haveCol.Field(i)) {
				return false
			}
		}
		return true
	case *array.List:
		haveCol := have.(*array.List)
		return sameOffsets(wantCol.Offsets(), haveCol.Offsets()) && equalAsFloat64(wantCol.ListValues(), haveCol.ListValues())
	case *array.LargeList:
		haveCol := have.(*array.LargeList)
		return sameLargeOffsets(wantCol.Offsets(), haveCol.Offsets()) && equalAsFloat64(wantCol.ListValues(), haveCol.ListValues())
	case *array.Map:
		haveCol := have.(*array.Map)
		return sameOffsets(wantCol.Offsets(), haveCol.Offsets()) &&
			equalAsFloat64(wantCol.Keys(), haveCol.Keys()) &&
			equalAsFloat64(wantCol.Items(), haveCol.Items())
	}

	if containsInt64(want.DataType()) {
		return false
	}
	return array.Equal(want, have)
}

func sameOffsets(want, have []int32) bool {
	if len(want) != len(have) {
		return false
	}
	for i := range want {
		if want[i] != have[i] {
			return false
		}
	}
	return true
}

func sameLargeOffsets(want, have []int64) bool {
	if len(want) != len(have) {
		return false
	}
	for i := range want {
		if want[i] != have[i] {
			return false
		}
	}
	return true
}

func TableDiff(have, want arrow.Table) string {
	if array.TableApproxEqual(have, want, array.WithUnorderedMapKeys(true)) {
		return ""
	}

	if have.NumCols() != want.NumCols() {
		return fmt.Sprintf("different number of columns: %d vs %d", have.NumCols(), want.NumCols())
	}
	if have.NumRows() != want.NumRows() {
		return fmt.Sprintf("different number of rows: %d vs %d", have.NumRows(), want.NumRows())
	}

	var sb strings.Builder
	for i := 0; i < int(have.NumCols()); i++ {
		haveCol, err := array.Concatenate(have.Column(i).Data().Chunks(), memory.DefaultAllocator)
		if err != nil {
			panic(fmt.Errorf("failed to concat left columns at idx %d: %w", i, err))
		}
		wantCol, err := array.Concatenate(want.Column(i).Data().Chunks(), memory.DefaultAllocator)
		if err != nil {
			panic(fmt.Errorf("failed to concat right columns at idx %d: %w", i, err))
		}
		edits, err := array.Diff(wantCol, haveCol)
		if err != nil {
			panic(fmt.Errorf("want: %v, have: %v, error: %w", wantCol.DataType(), haveCol.DataType(), err))
		}
		diff := edits.UnifiedDiff(wantCol, haveCol)
		if diff != "" {
			sb.WriteString(have.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			if equalAsFloat64(wantCol, haveCol) {
				sb.WriteString("values are equal after a float64 round-trip. If this destination cannot store 64-bit integers exactly, set schema.TestSourceOptions.MaxIntegerBits to schema.Float64SafeIntegerBits\n")
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
