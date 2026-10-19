package mixedbatchwriter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

type failingMixedBatchClient struct {
	testMixedBatchClient
}

func (*failingMixedBatchClient) InsertBatch(context.Context, message.WriteInserts) error {
	return errors.New("insert failed")
}

func TestFlushMessageCommitsBatches(t *testing.T) {
	msgs := getTestMessages()
	client := &testMixedBatchClient{}
	wr, err := New(client, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	ch <- msgs.migrateTable1
	ch <- msgs.insert1
	flush := &message.WriteFlush{}
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}, {msgs.insert1}}, client.receivedBatches)

	ch <- msgs.insert2
	flush = &message.WriteFlush{}
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}, {msgs.insert1}, {msgs.insert2}}, client.receivedBatches)

	close(ch)
	require.NoError(t, receive(t, errCh))
}

func TestFlushMessageReportsFailedInserts(t *testing.T) {
	msgs := getTestMessages()
	wr, err := New(&failingMixedBatchClient{}, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage, 3)
	flush := &message.WriteFlush{}
	ch <- msgs.migrateTable1
	ch <- msgs.insert1
	ch <- flush
	close(ch)

	require.ErrorContains(t, wr.Write(context.Background(), ch), "insert failed")
	err = receive(t, flush.Result())
	require.ErrorContains(t, err, "insert failed")
	require.Equal(t, []string{msgs.migrateTable1.Table.Name}, message.FailedTables(err))
}

func TestFlushMessageCommitsDeletes(t *testing.T) {
	msgs := getTestMessages()
	client := &testMixedBatchClient{}
	wr, err := New(client, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	ch <- msgs.deleteStale1
	flush := &message.WriteFlush{}
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.deleteStale1}}, client.receivedBatches)

	close(ch)
	require.NoError(t, receive(t, errCh))
}

func TestFlushMessageCommitsMigrationsAndDeleteRecords(t *testing.T) {
	msgs := getTestMessages()
	deleteRecord := &message.WriteDeleteRecord{DeleteRecord: message.DeleteRecord{TableName: msgs.migrateTable1.Table.Name}}
	client := &testMixedBatchClient{}
	wr, err := New(client, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	ch <- msgs.migrateTable1
	flush := &message.WriteFlush{}
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}}, client.receivedBatches)

	ch <- deleteRecord
	flush = &message.WriteFlush{}
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}, {deleteRecord}}, client.receivedBatches)

	close(ch)
	require.NoError(t, receive(t, errCh))
}

// receive fails the test instead of hanging when nothing arrives.
func receive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for a value")
		var zero T
		return zero
	}
}

func TestFlushMessageNamesEveryTableOfTheFailedBatch(t *testing.T) {
	table := func(name string) *schema.Table {
		return &schema.Table{Name: name, Columns: schema.ColumnList{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	}
	insert := func(table *schema.Table) *message.WriteInsert {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
		defer bldr.Release()
		bldr.Field(0).(*array.Int64Builder).Append(1)
		return &message.WriteInsert{Record: bldr.NewRecordBatch()}
	}
	wr, err := New(&failingMixedBatchClient{}, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage, 3)
	flush := &message.WriteFlush{}
	ch <- insert(table("table2"))
	ch <- insert(table("table1"))
	ch <- flush
	close(ch)

	require.ErrorContains(t, wr.Write(context.Background(), ch), "insert failed")
	require.Equal(t, []string{"table1", "table2"}, message.FailedTables(receive(t, flush.Result())))
}
