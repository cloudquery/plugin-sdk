package batchwriter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testBatchClient struct {
	mutex         sync.Mutex
	migrateTables message.WriteMigrateTables
	inserts       message.WriteInserts
	deleteStales  message.WriteDeleteStales
	deleteRecords message.WriteDeleteRecords
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

func (c *testBatchClient) DeleteRecord(_ context.Context, messages message.WriteDeleteRecords) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.deleteRecords = append(c.deleteRecords, messages...)
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

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, batchTestTables[0].ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()
	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteMigrateTable{Table: batchTestTables[0]}}); err != nil {
		t.Fatal(err)
	}

	if testClient.MigrateTablesLen() != 0 {
		t.Fatalf("expected 0 create table messages, got %d", testClient.MigrateTablesLen())
	}

	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteInsert{Record: record}}); err != nil {
		t.Fatal(err)
	}

	if testClient.MigrateTablesLen() != 1 {
		t.Fatalf("expected 1 migrate table message, got %d", testClient.MigrateTablesLen())
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteMigrateTable{Table: batchTestTables[0]}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 1 {
		t.Fatalf("expected 1 insert message, got %d", testClient.InsertsLen())
	}
}

func TestBatchSize(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	if err := wr.writeAll(ctx, []message.WriteMessage{
		&message.WriteInsert{
			Record: record,
		},
		&message.WriteInsert{ // third message to exceed the batch size
			Record: record,
		},
	}); err != nil {
		t.Fatal(err)
	}
	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if testClient.InsertsLen() != 2 {
		t.Fatalf("expected 2 insert messages, got %d", testClient.InsertsLen())
	}
}

func TestBatchTimeout(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Millisecond * 250)

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 1)

	if testClient.InsertsLen() != 1 {
		t.Fatalf("expected 1 insert message, got %d", testClient.InsertsLen())
	}
}

func TestBatchUpserts(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchSize(2), WithBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true}}}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	if err := wr.writeAll(ctx, []message.WriteMessage{&message.WriteInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}
	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if testClient.InsertsLen() != 2 {
		t.Fatalf("expected 2 insert messages, got %d", testClient.InsertsLen())
	}
}
