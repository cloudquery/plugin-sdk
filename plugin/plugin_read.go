package plugin

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// readAll is used in tests to read all records from a table.
func (p *Plugin) readAll(ctx context.Context, table *schema.Table) ([]arrow.RecordBatch, error) {
	var err error
	ch := make(chan arrow.RecordBatch)
	go func() {
		defer close(ch)
		err = p.client.Read(ctx, table, ch)
	}()
	// nolint:prealloc
	var records []arrow.RecordBatch
	for record := range ch {
		records = append(records, sliceToSingleRowRecord(record)...)
	}

	return records, err
}

func sliceToSingleRowRecord(record arrow.RecordBatch) []arrow.RecordBatch {
	result := make([]arrow.RecordBatch, record.NumRows())
	for i := int64(0); i < record.NumRows(); i++ {
		result[i] = record.NewSlice(i, i+1)
	}
	return result
}
