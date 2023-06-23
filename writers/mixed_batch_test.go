package writers

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testMixedBatchClient struct {
	receivedBatches [][]message.Message
}

func (c *testMixedBatchClient) MigrateTableBatch(ctx context.Context, msgs []*message.MigrateTable, options plugin.WriteOptions) error {
	m := make([]message.Message, len(msgs))
	for i, msg := range msgs {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

func (c *testMixedBatchClient) InsertBatch(ctx context.Context, msgs []*message.Insert, options plugin.WriteOptions) error {
	m := make([]message.Message, len(msgs))
	for i, msg := range msgs {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

func (c *testMixedBatchClient) DeleteStaleBatch(ctx context.Context, msgs []*message.DeleteStale, options plugin.WriteOptions) error {
	m := make([]message.Message, len(msgs))
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
	msgMigrateTable1 := &message.MigrateTable{
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
	msgMigrateTable2 := &message.MigrateTable{
		Table: table2,
	}

	// message to insert into table1
	bldr1 := array.NewRecordBuilder(memory.DefaultAllocator, table1.ToArrowSchema())
	bldr1.Field(0).(*array.Int64Builder).Append(1)
	rec1 := bldr1.NewRecord()
	msgInsertTable1 := &message.Insert{
		Record: rec1,
	}

	// message to insert into table2
	bldr2 := array.NewRecordBuilder(memory.DefaultAllocator, table1.ToArrowSchema())
	bldr2.Field(0).(*array.Int64Builder).Append(1)
	rec2 := bldr2.NewRecord()
	msgInsertTable2 := &message.Insert{
		Record: rec2,
	}

	// message to delete stale from table1
	msgDeleteStale1 := &message.DeleteStale{
		Table:      table1,
		SourceName: "my-source",
		SyncTime:   time.Now(),
	}
	msgDeleteStale2 := &message.DeleteStale{
		Table:      table1,
		SourceName: "my-source",
		SyncTime:   time.Now(),
	}

	testCases := []struct {
		name        string
		messages    []message.Message
		wantBatches [][]message.Message
	}{
		{
			name: "create table, insert, delete stale",
			messages: []message.Message{
				msgMigrateTable1,
				msgMigrateTable2,
				msgInsertTable1,
				msgInsertTable2,
				msgDeleteStale1,
				msgDeleteStale2,
			},
			wantBatches: [][]message.Message{
				{msgMigrateTable1, msgMigrateTable2},
				{msgInsertTable1, msgInsertTable2},
				{msgDeleteStale1, msgDeleteStale2},
			},
		},
		{
			name: "interleaved messages",
			messages: []message.Message{
				msgMigrateTable1,
				msgInsertTable1,
				msgDeleteStale1,
				msgMigrateTable2,
				msgInsertTable2,
				msgDeleteStale2,
			},
			wantBatches: [][]message.Message{
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
			messages: []message.Message{
				msgMigrateTable1,
				msgMigrateTable2,
				msgInsertTable1,
				msgDeleteStale2,
				msgInsertTable2,
				msgDeleteStale1,
			},
			wantBatches: [][]message.Message{
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
				receivedBatches: make([][]message.Message, 0),
			}
			wr, err := NewMixedBatchWriter(client)
			if err != nil {
				t.Fatal(err)
			}
			ch := make(chan message.Message, len(tc.messages))
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
