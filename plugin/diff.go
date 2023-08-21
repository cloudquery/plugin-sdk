package plugin

import (
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
)

func RecordsDiff(sc *arrow.Schema, l, r []arrow.Record) string {
	return TableDiff(array.NewTableFromRecords(sc, l), array.NewTableFromRecords(sc, r))
}

func TableDiff(l, r arrow.Table) string {
	if array.TableApproxEqual(l, r, array.WithUnorderedMapKeys(true)) {
		return ""
	}

	if l.NumCols() != r.NumCols() {
		return fmt.Sprintf("different number of columns: %d vs %d", l.NumCols(), r.NumCols())
	}
	if l.NumRows() != r.NumRows() {
		return fmt.Sprintf("different number of rows: %d vs %d", l.NumRows(), r.NumRows())
	}

	var sb strings.Builder
	for i := 0; i < int(l.NumCols()); i++ {
		lCol, err := array.Concatenate(l.Column(i).Data().Chunks(), memory.DefaultAllocator)
		if err != nil {
			panic(fmt.Errorf("failed to concat left columns at idx %d: %w", i, err))
		}
		rCol, err := array.Concatenate(r.Column(i).Data().Chunks(), memory.DefaultAllocator)
		if err != nil {
			panic(fmt.Errorf("failed to concat right columns at idx %d: %w", i, err))
		}
		edits, err := array.Diff(lCol, rCol)
		if err != nil {
			panic(fmt.Errorf("left: %v, right: %v, error: %w", lCol.DataType(), rCol.DataType(), err))
		}
		diff := edits.UnifiedDiff(lCol, rCol)
		if diff != "" {
			sb.WriteString(l.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
