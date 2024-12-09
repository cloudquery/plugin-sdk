package mixedbatchwriter

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"golang.org/x/sync/errgroup"
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

func (c *testMixedBatchClient) DeleteRecordsBatch(_ context.Context, messages message.WriteDeleteRecords) error {
	m := make([]message.WriteMessage, len(messages))
	for i, msg := range messages {
		m[i] = msg
	}
	c.receivedBatches = append(c.receivedBatches, m)
	return nil
}

var _ Client = (*testMixedBatchClient)(nil)

type testMessages struct {
	migrateTable1 *message.WriteMigrateTable
	migrateTable2 *message.WriteMigrateTable
	insert1       *message.WriteInsert
	insert2       *message.WriteInsert
	deleteStale1  *message.WriteDeleteStale
	deleteStale2  *message.WriteDeleteStale
}

func getTestMessages() testMessages {
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

	return testMessages{
		migrateTable1: msgMigrateTable1,
		migrateTable2: msgMigrateTable2,
		insert1:       msgInsertTable1,
		insert2:       msgInsertTable2,
		deleteStale1:  msgDeleteStale1,
		deleteStale2:  msgDeleteStale2,
	}
}

func TestMixedBatchWriter(t *testing.T) {
	tm := getTestMessages()
	testCases := []struct {
		name        string
		messages    []message.WriteMessage
		wantBatches [][]message.WriteMessage
	}{
		{
			name: "create table, insert, delete stale",
			messages: []message.WriteMessage{
				tm.migrateTable1,
				tm.migrateTable2,
				tm.insert1,
				tm.insert2,
				tm.deleteStale1,
				tm.deleteStale2,
			},
			wantBatches: [][]message.WriteMessage{
				{tm.migrateTable1, tm.migrateTable2},
				{tm.insert1, tm.insert2},
				{tm.deleteStale1, tm.deleteStale2},
			},
		},
		{
			name: "interleaved messages",
			messages: []message.WriteMessage{
				tm.migrateTable1,
				tm.insert1,
				tm.deleteStale1,
				tm.migrateTable2,
				tm.insert2,
				tm.deleteStale2,
			},
			wantBatches: [][]message.WriteMessage{
				{tm.migrateTable1},
				{tm.insert1},
				{tm.deleteStale1},
				{tm.migrateTable2},
				{tm.insert2},
				{tm.deleteStale2},
			},
		},
		{
			name: "interleaved messages",
			messages: []message.WriteMessage{
				tm.migrateTable1,
				tm.migrateTable2,
				tm.insert1,
				tm.deleteStale2,
				tm.insert2,
				tm.deleteStale1,
			},
			wantBatches: [][]message.WriteMessage{
				{tm.migrateTable1, tm.migrateTable2},
				{tm.insert1},
				{tm.deleteStale2},
				{tm.insert2},
				{tm.deleteStale1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client := &testMixedBatchClient{
				receivedBatches: make([][]message.WriteMessage, 0),
			}
			wr, err := New(client)
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

type mockTicker struct {
	C       chan time.Time
	trigger chan struct{}
}

func (m *mockTicker) Chan() <-chan time.Time {
	return m.C
}

func (m *mockTicker) Trigger() chan<- struct{} {
	return m.trigger
}

func (m *mockTicker) Stop() {
	close(m.C)
}

func (*mockTicker) Reset(_ time.Duration) {}

func newMockTicker(trigger chan struct{}) *mockTicker {
	c := make(chan time.Time)
	go func() {
		for range trigger {
			c <- time.Now()
		}
	}()
	return &mockTicker{
		C:       c,
		trigger: trigger,
	}
}

func TestMixedBatchWriterTimeout(t *testing.T) {
	tm := getTestMessages()
	cases := []struct {
		name        string
		messages    []message.WriteMessage
		wantBatches [][]message.WriteMessage
	}{
		{
			name: "one_message_batches",
			messages: []message.WriteMessage{
				tm.insert1,
				tm.insert2,
			},
			wantBatches: [][]message.WriteMessage{
				{tm.insert1},
				{tm.insert2},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client := &testMixedBatchClient{
				receivedBatches: make([][]message.WriteMessage, 0),
			}
			triggerTimeout := make(chan struct{})
			defer close(triggerTimeout)
			wr, err := New(client,
				WithBatchSize(1000),
				WithBatchSizeBytes(1000000),
				withTickerFn(func(_ time.Duration) writers.Ticker {
					return newMockTicker(triggerTimeout)
				}),
			)
			if err != nil {
				t.Fatal(err)
			}
			ch := make(chan message.WriteMessage)

			eg := errgroup.Group{}
			eg.Go(func() error {
				return wr.Write(ctx, ch)
			})

			for _, msg := range tc.messages {
				ch <- msg
				time.Sleep(100 * time.Millisecond)
				triggerTimeout <- struct{}{}
				time.Sleep(100 * time.Millisecond)
			}
			close(ch)
			err = eg.Wait()
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
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
