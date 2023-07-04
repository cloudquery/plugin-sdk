package mixedbatchwriter_test

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers/mixedbatchwriter"
)

type testMixedBatchClient struct {
	receivedBatches [][]message.WriteMessage
}

func (c *testMixedBatchClient) MigrateTableBatch(_ context.Context, messages message.WriteMigrateTables) error {
	m := make([]message.WriteMessage, len(messages))
	for i, msg := range messages {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

func (c *testMixedBatchClient) InsertBatch(_ context.Context, messages message.WriteInserts) error {
	m := make([]message.WriteMessage, len(messages))
	for i, msg := range messages {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

func (c *testMixedBatchClient) DeleteStaleBatch(_ context.Context, messages message.WriteDeleteStales) error {
	m := make([]message.WriteMessage, len(messages))
	for i, msg := range messages {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

var _ mixedbatchwriter.Client = (*testMixedBatchClient)(nil)

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
	msgMigrateTable1 := &message.WriteMigrateTable{
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
	msgMigrateTable2 := &message.WriteMigrateTable{
		Table: table2,
	}

	// message to insert into table1
	bldr1 := array.NewRecordBuilder(memory.DefaultAllocator, table1.ToArrowSchema())
	bldr1.Field(0).(*array.Int64Builder).Append(1)
	rec1 := bldr1.NewRecord()
	msgInsertTable1 := &message.WriteInsert{
		Record: rec1,
	}

	// message to insert into table2
	bldr2 := array.NewRecordBuilder(memory.DefaultAllocator, table1.ToArrowSchema())
	bldr2.Field(0).(*array.Int64Builder).Append(1)
	rec2 := bldr2.NewRecord()
	msgInsertTable2 := &message.WriteInsert{
		Record: rec2,
	}

	// message to delete stale from table1
	msgDeleteStale1 := &message.WriteDeleteStale{
		TableName:  table1.Name,
		SourceName: "my-source",
		SyncTime:   time.Now(),
	}
	msgDeleteStale2 := &message.WriteDeleteStale{
		TableName:  table1.Name,
		SourceName: "my-source",
		SyncTime:   time.Now(),
	}

	testCases := []struct {
		name        string
		messages    []message.WriteMessage
		wantBatches [][]message.WriteMessage
	}{
		{
			name: "create table, insert, delete stale",
			messages: []message.WriteMessage{
				msgMigrateTable1,
				msgMigrateTable2,
				msgInsertTable1,
				msgInsertTable2,
				msgDeleteStale1,
				msgDeleteStale2,
			},
			wantBatches: [][]message.WriteMessage{
				{msgMigrateTable1, msgMigrateTable2},
				{msgInsertTable1, msgInsertTable2},
				{msgDeleteStale1, msgDeleteStale2},
			},
		},
		{
			name: "interleaved messages",
			messages: []message.WriteMessage{
				msgMigrateTable1,
				msgInsertTable1,
				msgDeleteStale1,
				msgMigrateTable2,
				msgInsertTable2,
				msgDeleteStale2,
			},
			wantBatches: [][]message.WriteMessage{
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
			messages: []message.WriteMessage{
				msgMigrateTable1,
				msgMigrateTable2,
				msgInsertTable1,
				msgDeleteStale2,
				msgInsertTable2,
				msgDeleteStale1,
			},
			wantBatches: [][]message.WriteMessage{
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
				receivedBatches: make([][]message.WriteMessage, 0),
			}
			wr, err := mixedbatchwriter.New(client)
			if err != nil {
				t.Fatal(err)
			}
			ch := make(chan message.WriteMessage, len(tc.messages))
			for _, msg := range tc.messages {
				ch <- msg
			}
			close(ch)
			if err := wr.Write(ctx, ch); err != nil {
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
