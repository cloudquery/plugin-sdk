package streamingbatchwriter

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type messageType int

const (
	messageTypeMigrateTable messageType = iota
	messageTypeInsert
	messageTypeDeleteStale
	messageTypeDeleteRecord
)

type testStreamingBatchClient struct {
	mutex sync.Mutex

	inflight  map[messageType]int
	committed map[messageType]int
	open      map[messageType][]string

	writeErr       error
	writeErrAfter  int64
	writeCounter   map[string]int64 // table name to write counter
	writeCommitErr error
}

func newClient() *testStreamingBatchClient {
	return &testStreamingBatchClient{
		inflight:     make(map[messageType]int),
		committed:    make(map[messageType]int),
		open:         make(map[messageType][]string),
		writeCounter: make(map[string]int64),
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
	if c.writeErr != nil && c.writeErrAfter == -1 {
		return c.writeErr
	}

	key := ""
	for m := range msgs {
		key = c.handleTypeMessage(ctx, messageTypeInsert, m, key)

		c.mutex.Lock()
		c.writeCounter[key]++
		currentCount := c.writeCounter[key]
		c.mutex.Unlock()

		if c.writeErr != nil && currentCount > c.writeErrAfter {
			return c.writeErr // leave msgs open
		}
	}

	if c.writeCommitErr != nil {
		return c.writeCommitErr
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

func (c *testStreamingBatchClient) DeleteRecords(ctx context.Context, msgs <-chan *message.WriteDeleteRecord) error {
	key := ""
	for m := range msgs {
		key = c.handleTypeMessage(ctx, messageTypeDeleteRecord, m, key)
	}
	return c.handleTypeCommit(ctx, messageTypeDeleteRecord, key)
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

var _ Client = (*testStreamingBatchClient)(nil)

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
	wr, err := New(testClient)
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
	record := bldr.NewRecordBatch()

	if l := testClient.MessageLen(messageTypeMigrateTable); l != 0 {
		t.Fatalf("expected 0 migrate table messages, got %d", l)
	}

	ch <- &message.WriteInsert{Record: record}

	waitForLength(t, testClient.MessageLen, messageTypeMigrateTable, 1)

	if l := testClient.MessageLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 insert messages, got %d", l)
	}

	ch <- &message.WriteMigrateTable{Table: streamingBatchTestTable}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 1)

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
	wr, err := New(testClient, WithBatchSizeRows(2))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)
	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)

	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 2)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 0)

	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 2)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)

	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 4)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 0)

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}

	if l := testClient.MessageLen(messageTypeInsert); l != 4 {
		t.Fatalf("expected 3 insert messages, got %d", l)
	}
}

func TestStreamingBatchTimeout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	tickerFn, tickFn := newMockTicker()

	wr, err := New(testClient, withTickerFn(tickerFn))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)
	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)

	time.Sleep(time.Millisecond * 50) // we need to wait for the batch to be flushed

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)

	// flush
	tickFn()
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 1)

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestStreamingBatchNoTimeout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	wr, err := New(testClient, WithBatchTimeout(0), WithBatchSizeRows(2))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)
	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)

	time.Sleep(2 * time.Second)

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)

	ch <- &message.WriteInsert{
		Record: record,
	}
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 2)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 0)

	ch <- &message.WriteInsert{
		Record: record,
	}

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 2)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}

	if l := testClient.MessageLen(messageTypeInsert); l != 3 {
		t.Fatalf("expected 3 insert messages, got %d", l)
	}
}

func TestStreamingBatchUpserts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	tickerFn, tickFn := newMockTicker()
	wr, err := New(testClient, WithBatchSizeRows(2), withTickerFn(tickerFn))
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
	record := bldr.NewRecordBatch()

	ch <- &message.WriteInsert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)

	ch <- &message.WriteInsert{
		Record: record,
	}
	time.Sleep(50 * time.Millisecond)

	// flush the batch
	tickFn()
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 2)

	close(ch)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}

	if l := testClient.OpenLen(messageTypeInsert); l != 0 {
		t.Fatalf("expected 0 open tables, got %d", l)
	}
}

func TestErrorCleanUpBeforeFirstMessage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	testClient.writeErrAfter = -1
	testClient.writeErr = errors.New("test error")

	wr, err := New(testClient, WithBatchTimeout(0), WithBatchSizeRows(100))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			ch <- &message.WriteInsert{
				Record: record,
			}
		}
	}()

	<-done
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)

	close(ch)
	requireErrorCount(t, errCh, 1, 1)
}

func TestErrorCleanUpFirstMessage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	testClient.writeErrAfter = 0
	testClient.writeErr = errors.New("test error")

	wr, err := New(testClient, WithBatchTimeout(0), WithBatchSizeRows(100))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			ch <- &message.WriteInsert{
				Record: record,
			}
		}
	}()

	<-done
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)
	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1)

	close(ch)
	requireErrorCount(t, errCh, 1, 1)
}

func TestErrorCleanUpSecondMessage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	testClient.writeErrAfter = 1
	testClient.writeErr = errors.New("test error")

	wr, err := New(testClient, WithBatchTimeout(0), WithBatchSizeRows(2))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			ch <- &message.WriteInsert{
				Record: record,
			}
		}
	}()

	<-done

	close(ch)
	numErrs := requireErrorCount(t, errCh, 1, 2) // can have 2 errors depending on processing order

	waitForLength(t, testClient.InflightLen, messageTypeInsert, 1+numErrs) // testStreamingBatchClient doesn't commit the batch before erroring
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0)
}

func TestErrorCleanUpAfterClose(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ch := make(chan message.WriteMessage)

	testClient := newClient()
	testClient.writeCommitErr = errors.New("test error")

	wr, err := New(testClient, WithBatchTimeout(0), WithBatchSizeRows(100))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	table := schema.Table{Name: "table1", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	record := getRecord(table.ToArrowSchema(), 1)

	for i := 0; i < 10; i++ {
		ch <- &message.WriteInsert{
			Record: record,
		}
	}

	waitForLength(t, testClient.InflightLen, messageTypeInsert, 10)
	close(ch)

	requireErrorCount(t, errCh, 1, 1)

	waitForLength(t, testClient.MessageLen, messageTypeInsert, 0) // batch size 1
}

func waitForLength(t *testing.T, checkLen func(messageType) int, msgType messageType, want int) {
	t.Helper()
	lastValue := -1
	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatalf("timed out waiting for %v message length %d (last value: %d)", msgType, want, lastValue)
		default:
			if lastValue = checkLen(msgType); lastValue == want {
				return
			}
		}
	}
}

// nolint:unparam
func getRecord(sc *arrow.Schema, rows int) arrow.RecordBatch {
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	defer builder.Release()

	for _, f := range builder.Fields() {
		f.AppendEmptyValues(rows)
	}

	return builder.NewRecordBatch()
}

// nolint:unparam
func requireErrorCount(t *testing.T, errCh chan error, expectedMin, expectedMax int) int {
	t.Helper()
	select {
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for errCh")
	case err := <-errCh:
		jointErrs, ok := err.(interface{ Unwrap() []error })
		if !ok {
			t.Fatalf("errCh did not contain joint errors: %T", err)
		}

		errs := jointErrs.Unwrap()
		l := len(errs)
		if expectedMin == expectedMax && l != expectedMin {
			t.Fatalf("expected %d errors, got %d: %v", expectedMin, l, errs)
		} else if l < expectedMin || l > expectedMax {
			t.Fatalf("expected between %d and %d errors, got %d: %v", expectedMin, expectedMax, l, errs)
		}
		return l
	}
	return -1
}

func TestDeleteRecordFlushesPendingInserts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	errCh := make(chan error, 10)

	testClient := newClient()
	wr, err := New(testClient, WithBatchSizeRows(1000000)) // large batch to avoid auto-flush
	if err != nil {
		t.Fatal(err)
	}

	// Create a table for insert
	insertTable := &schema.Table{
		Name: "child_table",
		Columns: []schema.Column{
			{
				Name: "id",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}

	// Build insert record
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, insertTable.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	record := bldr.NewRecord()

	md := arrow.NewMetadata(
		[]string{schema.MetadataTableName},
		[]string{insertTable.Name},
	)
	newSchema := arrow.NewSchema(
		record.Schema().Fields(),
		&md,
	)

	record = array.NewRecord(newSchema, record.Columns(), record.NumRows())

	// Send insert 
	if err := wr.startWorker(ctx, errCh, &message.WriteInsert{Record: record}); err != nil {
		t.Fatal(err)
	}

	// send delete record to trigger flush
	del := &message.WriteDeleteRecord{
		DeleteRecord: message.DeleteRecord{
			TableName: insertTable.Name,
		},
	}

	if err := wr.startWorker(ctx, errCh, del); err != nil {
		t.Fatal(err)
	}
	waitForLength(t, testClient.MessageLen, messageTypeInsert, 1)
	_ = wr.Close(ctx)
}
