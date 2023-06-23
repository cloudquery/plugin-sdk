package writers

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
	migrateTables []*message.MigrateTable
	inserts       []*message.Insert
	deleteStales  []*message.DeleteStale
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

func (c *testBatchClient) MigrateTables(_ context.Context, msgs []*message.MigrateTable) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.migrateTables = append(c.migrateTables, msgs...)
	return nil
}

func (c *testBatchClient) WriteTableBatch(_ context.Context, _ string, msgs []*message.Insert) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inserts = append(c.inserts, msgs...)
	return nil
}
func (c *testBatchClient) DeleteStale(_ context.Context, msgs []*message.DeleteStale) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.deleteStales = append(c.deleteStales, msgs...)
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

// TestBatchFlushDifferentMessages tests that if writer receives a message of a new type all other pending
// batches are flushed.
func TestBatchFlushDifferentMessages(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := NewBatchWriter(testClient)
	if err != nil {
		t.Fatal(err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, batchTestTables[0].ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()
	if err := wr.writeAll(ctx, []message.Message{&message.MigrateTable{Table: batchTestTables[0]}}); err != nil {
		t.Fatal(err)
	}

	if testClient.MigrateTablesLen() != 0 {
		t.Fatalf("expected 0 create table messages, got %d", testClient.MigrateTablesLen())
	}

	if err := wr.writeAll(ctx, []message.Message{&message.Insert{Record: record}}); err != nil {
		t.Fatal(err)
	}

	if testClient.MigrateTablesLen() != 1 {
		t.Fatalf("expected 1 migrate table messages, got %d", testClient.MigrateTablesLen())
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	if err := wr.writeAll(ctx, []message.Message{&message.MigrateTable{Table: batchTestTables[0]}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 1 {
		t.Fatalf("expected 1 insert messages, got %d", testClient.InsertsLen())
	}

}

func TestBatchSize(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := NewBatchWriter(testClient, WithBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	if err := wr.writeAll(ctx, []message.Message{&message.Insert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	if err := wr.writeAll(ctx, []message.Message{&message.Insert{
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

func TestBatchTimeout(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := NewBatchWriter(testClient, WithBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	if err := wr.writeAll(ctx, []message.Message{&message.Insert{
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
		t.Fatalf("expected 1 insert messages, got %d", testClient.InsertsLen())
	}
}

func TestBatchUpserts(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := NewBatchWriter(testClient, WithBatchSize(2), WithBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true}}}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	if err := wr.writeAll(ctx, []message.Message{&message.Insert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if testClient.InsertsLen() != 0 {
		t.Fatalf("expected 0 insert messages, got %d", testClient.InsertsLen())
	}

	if err := wr.writeAll(ctx, []message.Message{&message.Insert{
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
