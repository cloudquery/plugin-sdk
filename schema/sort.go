package schema

import (
	"sort"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
)

// SortTable will sort rows by specific column. This is intended for testing purposes only.
// Because "id" is auto-incrementing in the test  data generator, if passed "id"
// this should result in records being returned in insertion order.
func SortTable(table arrow.Table, name string) arrow.Table {
	if !table.Schema().HasField(name) {
		panic("table has no '" + name + "' column to sort on")
	}
	rows := slice(table)
	colIndex := table.Schema().FieldIndices(name)[0]
	sort.Slice(rows, func(i, j int) bool {
		switch rows[i].Column(colIndex).DataType().(type) {
		case *arrow.Int64Type:
			v1 := rows[i].Column(colIndex).(*array.Int64).Value(0)
			v2 := rows[j].Column(colIndex).(*array.Int64).Value(0)
			return v1 < v2
		case *arrow.StringType:
			v1 := rows[i].Column(colIndex).(*array.String).Value(0)
			v2 := rows[j].Column(colIndex).(*array.String).Value(0)
			return v1 < v2
		default:
			panic("unsupported type for sorting")
		}
	})
	return array.NewTableFromRecords(table.Schema(), rows)
}

func slice(table arrow.Table) []arrow.Record {
	res := make([]arrow.Record, 0, table.NumRows())
	reader := array.NewTableReader(table, -1)
	for reader.Next() {
		record := reader.Record()
		for i := int64(0); i < record.NumRows(); i++ {
			res = append(res, record.NewSlice(i, i+1))
		}
	}
	return res
}
