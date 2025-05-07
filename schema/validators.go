package schema

import (
	"encoding/json"
	"slices"
	"strings"

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

func FindNotMatchingSensitiveColumns(table *Table) (nonMatchingColumns []string, nonMatchingJSONColumns []string) {
	if len(table.SensitiveColumns) == 0 {
		return []string{}, []string{}
	}

	nonMatchingColumns = make([]string, 0)
	nonMatchingJSONColumns = make([]string, 0)
	tableColumns := table.Columns.Names()
	for _, c := range table.SensitiveColumns {
		isJSONPath := false
		if strings.Contains(c, ".") {
			c = strings.Split(c, ".")[0]
			isJSONPath = true
		}
		if !slices.Contains(tableColumns, c) {
			nonMatchingColumns = append(nonMatchingColumns, c)
			continue
		}
		if !isJSONPath {
			continue
		}
		col := table.Columns.Get(c)
		if !arrow.TypeEqual(col.Type, types.ExtensionTypes.JSON) {
			nonMatchingJSONColumns = append(nonMatchingJSONColumns, c)
			continue
		}
	}
	return nonMatchingColumns, nonMatchingJSONColumns
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
