package plugin

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// readAll is used in tests to read all records from a table.
func (p *Plugin) readAll(ctx context.Context, table *schema.Table) (arrow.Table, error) {
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

	return array.NewTableFromRecords(table.ToArrowSchema(), records), err
}
