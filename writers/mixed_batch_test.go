package writers

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testMixedBatchClient struct {
	receivedBatches [][]plugin.Message
}

func (c *testMixedBatchClient) MigrateTableBatch(ctx context.Context, msgs []*plugin.MessageMigrateTable, options plugin.WriteOptions) error {
	m := make([]plugin.Message, len(msgs))
	for i, msg := range msgs {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

func (c *testMixedBatchClient) InsertBatch(ctx context.Context, msgs []*plugin.MessageInsert, options plugin.WriteOptions) error {
	m := make([]plugin.Message, len(msgs))
	for i, msg := range msgs {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

func (c *testMixedBatchClient) DeleteStaleBatch(ctx context.Context, msgs []*plugin.MessageDeleteStale, options plugin.WriteOptions) error {
	m := make([]plugin.Message, len(msgs))
	for i, msg := range msgs {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

var _ MixedBatchClient = (*testMixedBatchClient)(nil)

func TestMixedBatchWriter(t *testing.T) {
	ctx := context.Background()

	// message to create table1
	table1 := &schema.Table{
		Name: "table1",
		Columns: []schema.Column{
			{
				Name: "id",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}
	msgMigrateTable1 := &plugin.MessageMigrateTable{
		Table: table1,
	}

	// message to create table2
	table2 := &schema.Table{
		Name: "table2",
		Columns: []schema.Column{
			{
				Name: "id",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}
	msgMigrateTable2 := &plugin.MessageMigrateTable{
		Table: table2,
	}

	// message to insert into table1
	bldr1 := array.NewRecordBuilder(memory.DefaultAllocator, table1.ToArrowSchema())
	bldr1.Field(0).(*array.Int64Builder).Append(1)
	rec1 := bldr1.NewRecord()
	msgInsertTable1 := &plugin.MessageInsert{
		Record: rec1,
	}

	// message to insert into table2
	bldr2 := array.NewRecordBuilder(memory.DefaultAllocator, table1.ToArrowSchema())
	bldr2.Field(0).(*array.Int64Builder).Append(1)
	rec2 := bldr2.NewRecord()
	msgInsertTable2 := &plugin.MessageInsert{
		Record: rec2,
		Upsert: false,
	}

	// message to delete stale from table1
	msgDeleteStale1 := &plugin.MessageDeleteStale{
		Table:      table1,
		SourceName: "my-source",
		SyncTime:   time.Now(),
	}
	msgDeleteStale2 := &plugin.MessageDeleteStale{
		Table:      table1,
		SourceName: "my-source",
		SyncTime:   time.Now(),
	}

	testCases := []struct {
		name        string
		messages    []plugin.Message
		wantBatches [][]plugin.Message
	}{
		{
			name: "create table, insert, delete stale",
			messages: []plugin.Message{
				msgMigrateTable1,
				msgMigrateTable2,
				msgInsertTable1,
				msgInsertTable2,
				msgDeleteStale1,
				msgDeleteStale2,
			},
			wantBatches: [][]plugin.Message{
				{msgMigrateTable1, msgMigrateTable2},
				{msgInsertTable1, msgInsertTable2},
				{msgDeleteStale1, msgDeleteStale2},
			},
		},
		{
			name: "interleaved messages",
			messages: []plugin.Message{
				msgMigrateTable1,
				msgInsertTable1,
				msgDeleteStale1,
				msgMigrateTable2,
				msgInsertTable2,
				msgDeleteStale2,
			},
			wantBatches: [][]plugin.Message{
				{msgMigrateTable1},
				{msgInsertTable1},
				{msgDeleteStale1},
				{msgMigrateTable2},
				{msgInsertTable2},
				{msgDeleteStale2},
			},
		},
		{
			name: "interleaved messages",
			messages: []plugin.Message{
				msgMigrateTable1,
				msgMigrateTable2,
				msgInsertTable1,
				msgDeleteStale2,
				msgInsertTable2,
				msgDeleteStale1,
			},
			wantBatches: [][]plugin.Message{
				{msgMigrateTable1, msgMigrateTable2},
				{msgInsertTable1},
				{msgDeleteStale2},
				{msgInsertTable2},
				{msgDeleteStale1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &testMixedBatchClient{
				receivedBatches: make([][]plugin.Message, 0),
			}
			wr, err := NewMixedBatchWriter(client)
			if err != nil {
				t.Fatal(err)
			}
			ch := make(chan plugin.Message, len(tc.messages))
			for _, msg := range tc.messages {
				ch <- msg
			}
			close(ch)
			if err := wr.Write(ctx, plugin.WriteOptions{}, ch); err != nil {
				t.Fatal(err)
			}
			if len(client.receivedBatches) != len(tc.wantBatches) {
				t.Fatalf("got %d batches, want %d", len(client.receivedBatches), len(tc.wantBatches))
			}
			for i, wantBatch := range tc.wantBatches {
				if len(client.receivedBatches[i]) != len(wantBatch) {
					t.Fatalf("got %d messages in batch %d, want %d", len(client.receivedBatches[i]), i, len(wantBatch))
				}
			}
		})
	}
}
