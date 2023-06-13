package writers

import (
	"context"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testMixedBatchClient struct {
}

func (c *testMixedBatchClient) CreateTableBatch(ctx context.Context, resources []plugin.MessageCreateTable) error {
	return nil
}

func (c *testMixedBatchClient) InsertBatch(ctx context.Context, resources []plugin.MessageInsert) error {
	return nil
}

func (c *testMixedBatchClient) DeleteStaleBatch(ctx context.Context, resources []plugin.MessageDeleteStale) error {
	return nil
}

func TestMixedBatchWriter(t *testing.T) {
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

	wr, err := NewMixedBatchWriter(tables, &testMixedBatchClient{})
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan plugin.Message, 1)

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, tables[0].ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	rec := bldr.NewRecord()
	msg := plugin.MessageInsert{
		Record: rec,
	}
	ch <- msg
	close(ch)
	if err := wr.Write(ctx, ch); err != nil {
		t.Fatal(err)
	}
}
