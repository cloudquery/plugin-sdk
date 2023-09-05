package schema

import (
	"github.com/apache/arrow/go/v14/arrow"
)

func FindEmptyColumns(table *Table, records []arrow.Record) []string {
	columnsWithValues := make([]bool, len(table.Columns))
	emptyColumns := make([]string, 0)

	for _, resource := range records {
		for colIndex, arr := range resource.Columns() {
			for i := 0; i < arr.Len(); i++ {
				if arr.IsValid(i) {
					columnsWithValues[colIndex] = true
				}
			}
		}
	}

	// Make sure every column has at least one value.
	for i, hasValue := range columnsWithValues {
		col := table.Columns[i]
		emptyExpected := col.Name == "_cq_parent_id" && table.Parent == nil
		if !hasValue && !emptyExpected && !col.IgnoreInTests {
			emptyColumns = append(emptyColumns, col.Name)
		}
	}
	return emptyColumns
}
