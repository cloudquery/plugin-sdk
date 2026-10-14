package streamingbatchwriter

import (
	"context"
	"errors"
	"testing"
	"time"

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
