package streamingbatchwriter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/stretchr/testify/require"
)

func TestFlushMessageCommitsBatches(t *testing.T) {
	testClient := newClient()
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	ch <- &message.WriteMigrateTable{Table: streamingBatchTestTable}
	ch <- &message.WriteInsert{Record: getRecord(streamingBatchTestTable.ToArrowSchema(), 1)}
	ch <- &message.WriteInsert{Record: getRecord(streamingBatchTestTable.ToArrowSchema(), 1)}
	flush := message.NewWriteFlush()
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, 1, testClient.MessageLen(messageTypeMigrateTable))
	require.Equal(t, 2, testClient.MessageLen(messageTypeInsert))
	require.Equal(t, 0, testClient.OpenLen(messageTypeInsert))

	ch <- &message.WriteInsert{Record: getRecord(streamingBatchTestTable.ToArrowSchema(), 1)}
	flush = message.NewWriteFlush()
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, 3, testClient.MessageLen(messageTypeInsert))

	close(ch)
	require.NoError(t, receive(t, errCh))
}

func TestFlushMessageReportsFailedWrites(t *testing.T) {
	testClient := newClient()
	testClient.writeCommitErr = errors.New("commit failed")
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	ch <- &message.WriteInsert{Record: getRecord(streamingBatchTestTable.ToArrowSchema(), 1)}
	flush := message.NewWriteFlush()
	ch <- flush
	err = receive(t, flush.Result())
	require.ErrorContains(t, err, "failed to write messages for table1")
	require.Equal(t, []string{"table1"}, message.FailedTables(err))

	close(ch)
	require.ErrorContains(t, receive(t, errCh), "commit failed")
}

func TestFlushMessageReportsRejectedMessages(t *testing.T) {
	testClient := newClient()
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	noTableName := arrow.NewSchema([]arrow.Field{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	ch <- &message.WriteInsert{Record: getRecord(noTableName, 1)}
	flush := message.NewWriteFlush()
	ch <- flush
	err = receive(t, flush.Result())
	require.ErrorContains(t, err, "table name not found in metadata")
	require.Nil(t, message.FailedTables(err), "a message without a table may belong to any table")

	close(ch)
	require.ErrorContains(t, receive(t, errCh), "table name not found in metadata")
}

// blockingCommitClient holds each WriteTable call open after its messages, until release is closed.
type blockingCommitClient struct {
	*testStreamingBatchClient
	committing chan struct{}
	release    chan struct{}
}

func (c *blockingCommitClient) WriteTable(_ context.Context, msgs <-chan *message.WriteInsert) error {
	// nolint:revive
	for range msgs {
	}
	c.committing <- struct{}{}
	<-c.release
	return nil
}

func TestFlushMessageCancelledDuringCommit(t *testing.T) {
	testClient := &blockingCommitClient{testStreamingBatchClient: newClient(), committing: make(chan struct{}, 1), release: make(chan struct{})}
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(ctx, ch)
	}()

	ch <- &message.WriteInsert{Record: getRecord(streamingBatchTestTable.ToArrowSchema(), 1)}
	flush := message.NewWriteFlush()
	ch <- flush
	receive(t, testClient.committing)
	cancel()
	require.ErrorIs(t, receive(t, flush.Result()), context.Canceled)
	close(testClient.release)
	close(ch)

	select {
	case <-errCh:
	case <-time.After(5 * time.Second):
		t.Fatal("Write did not return after the flush was cancelled")
	}
}

func TestFlushMessageCommitsDeletes(t *testing.T) {
	testClient := newClient()
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)

	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()

	ch <- &message.WriteDeleteStale{TableName: streamingBatchTestTable.Name, SourceName: "source", SyncTime: time.Now()}
	flush := message.NewWriteFlush()
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, 1, testClient.MessageLen(messageTypeDeleteStale))

	ch <- &message.WriteDeleteRecord{DeleteRecord: message.DeleteRecord{TableName: streamingBatchTestTable.Name}}
	flush = message.NewWriteFlush()
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, 1, testClient.MessageLen(messageTypeDeleteRecord))

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
