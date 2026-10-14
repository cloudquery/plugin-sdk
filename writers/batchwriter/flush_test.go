package batchwriter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/stretchr/testify/require"
)

type failingBatchClient struct {
	testBatchClient
}

func (*failingBatchClient) WriteTableBatch(context.Context, string, message.WriteInserts) error {
	return errors.New("write failed")
}

func startWrite(t *testing.T, wr *BatchWriter) (chan<- message.WriteMessage, func() error) {
	t.Helper()
	ch := make(chan message.WriteMessage)
	errCh := make(chan error, 1)
	go func() {
		errCh <- wr.Write(context.Background(), ch)
	}()
	return ch, func() error {
		close(ch)
		return <-errCh
	}
}

func flush(ch chan<- message.WriteMessage) error {
	msg := message.NewWriteFlush()
	ch <- msg
	return <-msg.Result()
}

func TestFlushMessageCommitsBatches(t *testing.T) {
	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 3)}
	require.NoError(t, flush(ch))
	require.Equal(t, 1, testClient.MigrateTablesLen())
	require.EqualValues(t, 3, testClient.InsertedRows())

	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 2)}
	require.NoError(t, flush(ch))
	require.EqualValues(t, 5, testClient.InsertedRows())

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}

func TestFlushMessageReportsFailedInserts(t *testing.T) {
	wr, err := New(&failingBatchClient{}, WithBatchTimeout(time.Hour))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 3)}
	require.ErrorContains(t, flush(ch), "write failed")
	// the failure was reported already
	require.NoError(t, flush(ch))

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}

func TestFlushMessageReportsFailuresFromTimeouts(t *testing.T) {
	wr, err := New(&failingBatchClient{}, WithBatchTimeout(time.Millisecond))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 3)}
	require.Eventually(t, func() bool { return wr.insertFailures.Load() > 0 }, 5*time.Second, time.Millisecond)
	require.ErrorContains(t, flush(ch), "write failed")

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}
