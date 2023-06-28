package streamingbatchwriter_test

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
	"github.com/cloudquery/plugin-sdk/v4/writers/streamingbatchwriter"
)

type messageType int

const (
	messageTypeMigrateTable messageType = iota
	messageTypeInsert
	messageTypeDeleteStale
)

type testStreamingBatchClient struct {
	mutex sync.Mutex

	inflight  map[messageType]int
	committed map[messageType]int
	open      map[messageType][]string
}

func newClient() *testStreamingBatchClient {
	return &testStreamingBatchClient{
		inflight:  make(map[messageType]int),
		committed: make(map[messageType]int),
		open:      make(map[messageType][]string),
	}
}

func (c *testStreamingBatchClient) MessageLen(t messageType) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.committed[t]
}

func (c *testStreamingBatchClient) InflightLen(t messageType) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.inflight[t]
}

func (c *testStreamingBatchClient) OpenLen(t messageType) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.open[t])
}

func (c *testStreamingBatchClient) MigrateTable(ctx context.Context, msgs <-chan *message.WriteMigrateTable) error {
	key := ""
	for m := range msgs {
		key = c.handleTypeMessage(ctx, messageTypeMigrateTable, m, key)
	}
	return c.handleTypeCommit(ctx, messageTypeMigrateTable, key)
}

func (c *testStreamingBatchClient) WriteTable(ctx context.Context, msgs <-chan *message.WriteInsert) error {
	key := ""
	for m := range msgs {
		key = c.handleTypeMessage(ctx, messageTypeInsert, m, key)
	}
	return c.handleTypeCommit(ctx, messageTypeInsert, key)
}

func (c *testStreamingBatchClient) DeleteStale(ctx context.Context, msgs <-chan *message.WriteDeleteStale) error {
	key := ""
	for m := range msgs {
		key = c.handleTypeMessage(ctx, messageTypeDeleteStale, m, key)
	}
	return c.handleTypeCommit(ctx, messageTypeDeleteStale, key)
}

func (c *testStreamingBatchClient) handleTypeMessage(_ context.Context, t messageType, msg message.WriteMessage, key string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if key == "" {
		key = msg.GetTable().Name
		c.open[t] = append(c.open[t], key)
	}
	c.inflight[t]++

	return key
}

func (c *testStreamingBatchClient) handleTypeCommit(_ context.Context, t messageType, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.committed[t] += c.inflight[t]
	c.inflight[t] = 0

	for i, openTable := range c.open[t] {
		if openTable == key {
			c.open[t] = append(c.open[t][:i], c.open[t][i+1:]...)
			break
		}
	}

	return nil
}

var _ streamingbatchwriter.Client = (*testStreamingBatchClient)(nil)

var streamingBatchTestTable = &schema.Table{
	Name: "table1",
	Columns: []schema.Column{
		{
			Name: "id",
			Type: arrow.PrimitiveTypes.Int64,
		},
	},
}

// TestBatchStreamFlushDifferentMessages tests that if writer receives a message of a new type all other pending batches are flushed.
func TestBatchStreamFlushDifferentMessages(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	wr, err := streamingbatchwriter.New(testClient)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	ch <- &message.WriteMigrateTable{Table: streamingBatchTestTable}
	time.Sleep(50 * time.Millisecond)

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, streamingBatchTestTable.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	if l := testClient.MessageLen(messageTypeMigrateTable); l != 0 {
		t.Fatalf("expected 0 migrate table messages, got %d", l)
	}

	ch <- &message.WriteInsert{Record: record}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.MessageLen(messageTypeMigrateTable); l != 1 {
		t.Fatalf("expected 1 migrate table message, got %d", l)
	}

	if l := testClient.MessageLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.WriteMigrateTable{Table: streamingBatchTestTable}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.MessageLen(messageTypeInsert); l != 1 {
		t.Fatalf("expected 1 insert message, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchSizeRows(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	wr, err := streamingbatchwriter.New(testClient, streamingbatchwriter.WithBatchSizeRows(2))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	ch <- &message.WriteInsert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.MessageLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.WriteInsert{
		Record: record,
	}
	ch <- &message.WriteInsert{ // third message, because we flush before exceeding the limit and then save the third one
		Record: record,
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if l := testClient.MessageLen(messageTypeInsert); l != 2 {
		t.Fatalf("expected 2 insert messages, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchTimeout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	wr, err := streamingbatchwriter.New(testClient, streamingbatchwriter.WithBatchTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := array.NewRecord(table.ToArrowSchema(), nil, 0)
	ch <- &message.WriteInsert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.MessageLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Millisecond * 250)

	if l := testClient.MessageLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 1)

	if l := testClient.MessageLen(messageTypeInsert); l != 1 {
		t.Fatalf("expected 1 insert message, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchUpserts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	wr, err := streamingbatchwriter.New(testClient, streamingbatchwriter.WithBatchSizeRows(2), streamingbatchwriter.WithBatchTimeout(time.Second))
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

	ch <- &message.WriteInsert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	if l := testClient.InflightLen(messageTypeInsert); l != 1 {
		t.Fatalf("expected 1 inflight insert message, got %d", l)
	}

	if l := testClient.MessageLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.WriteInsert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	// we need to wait for the batch to be flushed
	time.Sleep(time.Second * 2)

	if l := testClient.MessageLen(messageTypeInsert); l != 2 {
		t.Fatalf("expected 2 insert messages, got %d", l)
	}

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}
