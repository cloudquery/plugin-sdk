package batchwriter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testBatchClient struct {
	mutex         sync.Mutex
	migrateTables message.WriteMigrateTables
	inserts       message.WriteInserts
	deleteStales  message.WriteDeleteStales
}

func (c *testBatchClient) MigrateTablesLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.migrateTables)
}

func (c *testBatchClient) InsertsLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.inserts)
}

func (c *testBatchClient) DeleteStalesLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.deleteStales)
}

func (c *testBatchClient) MigrateTables(_ context.Context, messages message.WriteMigrateTables) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.migrateTables = append(c.migrateTables, messages...)
	return nil
}

func (c *testBatchClient) WriteTableBatch(_ context.Context, _ string, messages message.WriteInserts) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inserts = append(c.inserts, messages...)
	return nil
}
func (c *testBatchClient) DeleteStale(_ context.Context, messages message.WriteDeleteStales) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.deleteStales = append(c.deleteStales, messages...)
	return nil
}

var batchTestTables = schema.Tables{
	{
		Name: "table1",
		Columns: []schema.Column{
			{
				Name: "id",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	},
}

// TestBatchFlushDifferentMessages tests that if writer receives a message of a new type all other pending
// batches are flushed.
func TestBatchFlushDifferentMessages(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := New(testClient)
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan message.WriteMessage)
	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, batchTestTables[0].ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	writeTo(ch, &message.WriteMigrateTable{Table: batchTestTables[0]})

	if testClient.MigrateTablesLen() != 0 {
		t.Fatalf("expected 0 create table messages, got %d", testClient.MigrateTablesLen())
	}

	writeTo(ch, &message.WriteInsert{Record: record})

	if testClient.MigrateTablesLen() != 1 {
		t.Fatalf("expected 1 migrate table message, got %d", testClient.MigrateTablesLen())
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	writeTo(ch, &message.WriteMigrateTable{Table: batchTestTables[0]})

	if testClient.InsertsLen() != 1 {
		t.Fatalf("expected 1 insert message, got %d", testClient.InsertsLen())
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
}

func TestBatchSize(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan message.WriteMessage)
	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	writeTo(ch, &message.WriteInsert{Record: record})

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	writeTo(ch,
		&message.WriteInsert{Record: record},
		&message.WriteInsert{Record: record}, // third message to exceed the batch size
	)
	// we need to wait for the batch to be flushed
	time.Sleep(time.Millisecond * 50)

	if testClient.InsertsLen() != 2 {
		t.Fatalf("expected 2 insert messages, got %d", testClient.InsertsLen())
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
}

func TestBatchTimeout(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	timerFn, timerExpire := newMockTimer()

	wr, err := New(testClient, WithBatchTimeout(time.Second), withTimerFn(timerFn))
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan message.WriteMessage)
	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	writeTo(ch, &message.WriteInsert{Record: record})

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Millisecond * 50)

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	close(timerExpire)
	time.Sleep(time.Millisecond * 50)

	if testClient.InsertsLen() != 1 {
		t.Fatalf("expected 1 insert message, got %d", testClient.InsertsLen())
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
}

func TestBatchUpserts(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	timerFn, timerExpire := newMockTimer()

	wr, err := New(testClient, WithBatchSize(2), WithBatchTimeout(time.Second), withTimerFn(timerFn))
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan message.WriteMessage)
	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true}}}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	writeTo(ch, &message.WriteInsert{Record: record})
	time.Sleep(time.Millisecond * 50)

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	writeTo(ch, &message.WriteInsert{Record: record})
	time.Sleep(time.Millisecond * 50)

	close(timerExpire)
	time.Sleep(time.Millisecond * 50)

	if testClient.InsertsLen() != 2 {
		t.Fatalf("expected 2 insert messages, got %d", testClient.InsertsLen())
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
}

func writeTo(ch chan message.WriteMessage, msgs ...message.WriteMessage) {
	for _, msg := range msgs {
		ch <- msg
	}
}
