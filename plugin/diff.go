package plugin

import (
	"encoding/base64"
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
			wantData := make([]byte, wantCol.Len())
			for _, buffer := range wantCol.Data().Buffers() {
				wantData = append(wantData, buffer.Bytes()...)
			}
			haveDataType := haveCol.DataType()
			haveData := make([]byte, haveCol.Len())
			for _, buffer := range haveCol.Data().Buffers() {
				haveData = append(haveData, buffer.Bytes()...)
			}

			wantBase64 := base64.StdEncoding.EncodeToString(wantData)
			haveBase64 := base64.StdEncoding.EncodeToString(haveData)

			panic(fmt.Errorf("panic in getUnifiedDiff: %s, want: (%s), have: (%s), want type: %s, have type: %s", r, wantBase64, haveBase64, wantDataType, haveDataType))
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
