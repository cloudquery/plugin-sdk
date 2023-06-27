package writers

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type testStreamingBatchClient struct {
	mutex         sync.Mutex
	migrateTables []*message.MigrateTable
	inserts       []*message.Insert
	deleteStales  []*message.DeleteStale

	openTables []string
}

func (c *testStreamingBatchClient) MigrateTablesLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.migrateTables)
}

func (c *testStreamingBatchClient) InsertsLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.inserts)
}

func (c *testStreamingBatchClient) DeleteStalesLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.deleteStales)
}

func (c *testStreamingBatchClient) OpenTablesLen() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.openTables)
}

func (c *testStreamingBatchClient) MigrateTables(_ context.Context, msgs []*message.MigrateTable) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.migrateTables = append(c.migrateTables, msgs...)
	return nil
}

func (c *testStreamingBatchClient) OpenTable(_ context.Context, sourceName string, table *schema.Table, syncTime time.Time) (any, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := sourceName + ":" + table.Name
	c.openTables = append(c.openTables, key)
	return key, nil
}

func (c *testStreamingBatchClient) CloseTable(_ context.Context, handle any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := handle.(string)
	for i, openTable := range c.openTables {
		if openTable == key {
			c.openTables = append(c.openTables[:i], c.openTables[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("CloseTable: table not open")
}

func (c *testStreamingBatchClient) WriteTableStream(_ context.Context, handle any, msgs []*message.Insert) error {
	if len(msgs) == 0 {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := handle.(string)
	sn := (&StreamingBatchWriter{}).getSourceName(msgs[0].Record)
	tn := msgs[0].GetTable().Name
	if key != sn+":"+tn {
		return fmt.Errorf("WriteTableStream: table %s not open", tn)
	}
	c.inserts = append(c.inserts, msgs...)
	return nil
}

func (c *testStreamingBatchClient) DeleteStale(_ context.Context, msgs []*message.DeleteStale) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.deleteStales = append(c.deleteStales, msgs...)
	return nil
}

var streamingBatchTestTables = schema.Tables{
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

// TestBatchStreamFlushDifferentMessages tests that if writer receives a message of a new type all other pending batches are flushed.
func TestBatchStreamFlushDifferentMessages(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.Message)

	testClient := &testStreamingBatchClient{}
	wr, err := NewStreamingBatchWriter(testClient)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	ch <- &message.MigrateTable{Table: streamingBatchTestTables[0]}
	time.Sleep(50 * time.Millisecond)

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, streamingBatchTestTables[0].ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	if l := testClient.MigrateTablesLen(); l != 0 {
		t.Fatalf("expected 0 migrate table messages, got %d", l)
	}

	ch <- &message.Insert{Record: record}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.MigrateTablesLen(); l != 1 {
		t.Fatalf("expected 1 migrate table message, got %d", l)
	}

	if l := testClient.InsertsLen(); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.MigrateTable{Table: streamingBatchTestTables[0]}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.InsertsLen(); l != 1 {
		t.Fatalf("expected 1 insert message, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenTablesLen(); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchSize(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.Message)

	testClient := &testStreamingBatchClient{}
	wr, err := NewStreamingBatchWriter(testClient, WithStreamingBatchWriterBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	ch <- &message.Insert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.InsertsLen(); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.Insert{
		Record: record,
	}
	ch <- &message.Insert{ // third message, because we flush before exceeding the limit and then save the third one
		Record: record,
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if l := testClient.InsertsLen(); l != 2 {
		t.Fatalf("expected 2 insert messages, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenTablesLen(); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchTimeout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.Message)

	testClient := &testStreamingBatchClient{}
	wr, err := NewStreamingBatchWriter(testClient, WithStreamingBatchWriterBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	ch <- &message.Insert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.InsertsLen(); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Millisecond * 250)

	if l := testClient.InsertsLen(); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 1)

	if l := testClient.InsertsLen(); l != 1 {
		t.Fatalf("expected 1 insert message, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenTablesLen(); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchUpserts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.Message)

	testClient := &testStreamingBatchClient{}
	wr, err := NewStreamingBatchWriter(testClient, WithStreamingBatchWriterBatchSize(2), WithStreamingBatchWriterBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true}}}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	ch <- &message.Insert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.InsertsLen(); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.Insert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if l := testClient.InsertsLen(); l != 2 {
		t.Fatalf("expected 2 insert messages, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenTablesLen(); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}
