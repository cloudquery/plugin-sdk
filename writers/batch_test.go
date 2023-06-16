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

type testBatchClient struct {
	migrateTables []*plugin.MessageMigrateTable
	inserts       []*plugin.MessageInsert
	deleteStales  []*plugin.MessageDeleteStale
}

func (c *testBatchClient) MigrateTables(_ context.Context, msgs []*plugin.MessageMigrateTable) error {
	c.migrateTables = append(c.migrateTables, msgs...)
	return nil
}

func (c *testBatchClient) WriteTableBatch(_ context.Context, _ string, _ bool, msgs []*plugin.MessageInsert) error {
	c.inserts = append(c.inserts, msgs...)
	return nil
}
func (c *testBatchClient) DeleteStale(_ context.Context, msgs []*plugin.MessageDeleteStale) error {
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
	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageMigrateTable{Table: batchTestTables[0]}}); err != nil {
		t.Fatal(err)
	}
	if len(testClient.migrateTables) != 0 {
		t.Fatalf("expected 0 create table messages, got %d", len(testClient.migrateTables))
	}
	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageInsert{Record: record}}); err != nil {
		t.Fatal(err)
	}
	if len(testClient.migrateTables) != 1 {
		t.Fatalf("expected 1 create table messages, got %d", len(testClient.migrateTables))
	}

	if len(testClient.inserts) != 0 {
		t.Fatalf("expected 0 insert messages, got %d", len(testClient.inserts))
	}

	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageMigrateTable{Table: batchTestTables[0]}}); err != nil {
		t.Fatal(err)
	}

	if len(testClient.inserts) != 1 {
		t.Fatalf("expected 1 insert messages, got %d", len(testClient.inserts))
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
	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if len(testClient.inserts) != 0 {
		t.Fatalf("expected 0 create table messages, got %d", len(testClient.inserts))
	}

	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}
	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if len(testClient.inserts) != 2 {
		t.Fatalf("expected 2 create table messages, got %d", len(testClient.inserts))
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
	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}

	if len(testClient.inserts) != 0 {
		t.Fatalf("expected 0 create table messages, got %d", len(testClient.inserts))
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Millisecond * 250)

	if len(testClient.inserts) != 0 {
		t.Fatalf("expected 0 create table messages, got %d", len(testClient.inserts))
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 1)

	if len(testClient.inserts) != 1 {
		t.Fatalf("expected 1 create table messages, got %d", len(testClient.inserts))
	}
}

func TestBatchUpserts(t *testing.T) {
	ctx := context.Background()

	testClient := &testBatchClient{}
	wr, err := NewBatchWriter(testClient)
	if err != nil {
		t.Fatal(err)
	}
	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageInsert{
		Record: record,
		Upsert: true,
	}}); err != nil {
		t.Fatal(err)
	}

	if len(testClient.inserts) != 0 {
		t.Fatalf("expected 0 create table messages, got %d", len(testClient.inserts))
	}

	if err := wr.writeAll(ctx, []plugin.Message{&plugin.MessageInsert{
		Record: record,
	}}); err != nil {
		t.Fatal(err)
	}
	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if len(testClient.inserts) != 1 {
		t.Fatalf("expected 1 create table messages, got %d", len(testClient.inserts))
	}
}
