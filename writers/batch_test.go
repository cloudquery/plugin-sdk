package writers

import (
	"context"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testBatchClient struct {
}

func (c *testBatchClient) WriteTableBatch(context.Context, *schema.Table, []arrow.Record) error {
	return nil
}

func TestBatchWriter(t *testing.T) {
	ctx := context.Background()
	tables := schema.Tables{
		{
			Name: "table1",
			Columns: []schema.Column{
				{
					Name: "id",
					Type: arrow.PrimitiveTypes.Int64,
				},
			},
		},
		{
			Name: "table2",
			Columns: []schema.Column{
				{
					Name: "id",
					Type: arrow.PrimitiveTypes.Int64,
				},
			},
		},
	}

	wr, err := NewBatchWriter(tables, &testBatchClient{})
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan arrow.Record, 1)

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, tables[0].ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	ch <- bldr.NewRecord()
	close(ch)
	if err := wr.Write(ctx, ch); err != nil {
		t.Fatal(err)
	}
}
