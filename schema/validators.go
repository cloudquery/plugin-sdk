package schema

import (
	"encoding/json"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

func FindEmptyColumns(table *Table, records []arrow.Record) []string {
	columnsWithValues := make([]bool, len(table.Columns))
	emptyColumns := make([]string, 0)

	for _, resource := range records {
		for colIndex, arr := range resource.Columns() {
			for i := 0; i < arr.Len(); i++ {
				if arr.IsValid(i) {
					if arrow.TypeEqual(arr.DataType(), types.ExtensionTypes.JSON) {
						// JSON column shouldn't be empty
						val := arr.GetOneForMarshal(i).(json.RawMessage)
						if isEmptyJSON(val) {
							continue
						}
					}

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

func isEmptyJSON(msg json.RawMessage) bool {
	if len(msg) == 0 {
		return true
	}
	switch string(msg) {
	case "null", "{}", "[]":
		return true
	default:
		return false
	}
}
