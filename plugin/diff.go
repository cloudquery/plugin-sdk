package plugin

import (
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

func RecordsDiff(sc *arrow.Schema, have, want []arrow.Record) string {
	return TableDiff(array.NewTableFromRecords(sc, have), array.NewTableFromRecords(sc, want))
}

func getUnifiedDiff(edits array.Edits, wantCol, haveCol arrow.Array) string {
	defer func() {
		if r := recover(); r != nil {
			wantDataType := wantCol.DataType()
			wantData := make([]string, wantCol.Len())
			for i := 0; i < wantCol.Len(); i++ {
				wantData[i] = wantCol.ValueStr(i)
			}
			haveDataType := haveCol.DataType()
			haveData := make([]string, haveCol.Len())
			for i := 0; i < haveCol.Len(); i++ {
				haveData[i] = haveCol.ValueStr(i)
			}
			panic(fmt.Errorf("panic in getUnifiedDiff: %s, want: [%s], have: [%s], want type: %s, have type: %s", r, strings.Join(wantData, ", "), strings.Join(haveData, ", "), wantDataType, haveDataType))
		}
	}()
	return edits.UnifiedDiff(wantCol, haveCol)
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
		diff := getUnifiedDiff(edits, wantCol, haveCol)
		if diff != "" {
			sb.WriteString(have.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
