package plugin

import (
	"context"
	"sort"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// readAll is used in tests to read all records from a table.
func (p *Plugin) readAll(ctx context.Context, table *schema.Table) ([]arrow.Record, error) {
	// TODO (Go 1.21+) use testing.Testing() to guard against use of this function outside of tests
	// https://stackoverflow.com/a/75787351

	var err error
	ch := make(chan arrow.Record)
	go func() {
		defer close(ch)
		err = p.client.Read(ctx, table, ch)
	}()
	// nolint:prealloc
	var records []arrow.Record
	for record := range ch {
		records = append(records, record)
	}

	// sort records by "id" column, if present. Because "id" is auto-incrementing in the test
	// data generator, this should result in records being returned in insertion order.
	sch := table.ToArrowSchema()
	if sch.HasField("id") {
		idColIndex := sch.FieldIndices("id")[0]
		sort.Slice(records, func(i, j int) bool {
			v1 := records[i].Column(idColIndex).(*array.Int64).Value(0)
			v2 := records[j].Column(idColIndex).(*array.Int64).Value(0)
			return v1 < v2
		})
	} else if len(records) > 1 {
		panic("table has no 'id' column to sort on")
	}

	return records, err
}
