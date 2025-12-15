package plugin

import (
	"sort"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// sort records by specific column. This is intended for testing purposes only.
// Because "id" is auto-incrementing in the test  data generator, if passed "id"
// this should result in records being returned in insertion order.
// nolint:unparam
func sortRecords(table *schema.Table, records []arrow.RecordBatch, columnName string) {
	sch := table.ToArrowSchema()
	if !sch.HasField(columnName) {
		panic("table has no '" + columnName + "' column to sort on")
	}
	colIndex := sch.FieldIndices(columnName)[0]
	sort.Slice(records, func(i, j int) bool {
		switch records[i].Column(colIndex).DataType().(type) {
		case *arrow.Int64Type:
			v1 := records[i].Column(colIndex).(*array.Int64).Value(0)
			v2 := records[j].Column(colIndex).(*array.Int64).Value(0)
			return v1 < v2
		case *arrow.StringType:
			v1 := records[i].Column(colIndex).(*array.String).Value(0)
			v2 := records[j].Column(colIndex).(*array.String).Value(0)
			return v1 < v2
		default:
			panic("unsupported type for sorting")
		}
	})
}
