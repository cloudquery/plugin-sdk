package plugin

import (
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
)

func RecordDiff(l, r arrow.Record) string {
	if array.RecordApproxEqual(l, r, array.WithUnorderedMapKeys(true)) {
		return ""
	}
	var sb strings.Builder
	if l.NumCols() != r.NumCols() {
		return fmt.Sprintf("different number of columns: %d vs %d", l.NumCols(), r.NumCols())
	}
	if l.NumRows() != r.NumRows() {
		return fmt.Sprintf("different number of rows: %d vs %d", l.NumRows(), r.NumRows())
	}
	for i := 0; i < int(l.NumCols()); i++ {
		edits, err := array.Diff(l.Column(i), r.Column(i))
		if err != nil {
			panic(fmt.Sprintf("left: %v, right: %v, error: %v", l.Column(i).DataType(), r.Column(i).DataType(), err))
		}
		diff := edits.UnifiedDiff(l.Column(i), r.Column(i))
		if diff != "" {
			sb.WriteString(l.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
