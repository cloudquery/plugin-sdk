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
	require.NoError(t, <-flush.Result())
	require.Equal(t, 1, testClient.MessageLen(messageTypeMigrateTable))
	require.Equal(t, 2, testClient.MessageLen(messageTypeInsert))
	require.Equal(t, 0, testClient.OpenLen(messageTypeInsert))

	ch <- &message.WriteInsert{Record: getRecord(streamingBatchTestTable.ToArrowSchema(), 1)}
	flush = message.NewWriteFlush()
	ch <- flush
	require.NoError(t, <-flush.Result())
	require.Equal(t, 3, testClient.MessageLen(messageTypeInsert))

	close(ch)
	require.NoError(t, <-errCh)
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
	require.ErrorContains(t, <-flush.Result(), "failed to write messages for table1")

	close(ch)
	require.ErrorContains(t, <-errCh, "commit failed")
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
	require.ErrorContains(t, <-flush.Result(), "table name not found in metadata")

	close(ch)
	require.ErrorContains(t, <-errCh, "table name not found in metadata")
}

// blockingCommitClient holds each WriteTable call open after its messages, until release is closed.
type blockingCommitClient struct {
	*testStreamingBatchClient
	committing chan struct{}
	release    chan struct{}
}

func (c *blockingCommitClient) WriteTable(_ context.Context, msgs <-chan *message.WriteInsert) error {
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
	<-testClient.committing
	cancel()
	require.ErrorIs(t, <-flush.Result(), context.Canceled)
	close(testClient.release)
	close(ch)

	select {
	case <-errCh:
	case <-time.After(5 * time.Second):
		t.Fatal("Write did not return after the flush was cancelled")
	}
}
