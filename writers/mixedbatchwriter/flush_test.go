package mixedbatchwriter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/message"
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
	flush := message.NewWriteFlush()
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}, {msgs.insert1}}, client.receivedBatches)

	ch <- msgs.insert2
	flush = message.NewWriteFlush()
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
	flush := message.NewWriteFlush()
	ch <- msgs.migrateTable1
	ch <- msgs.insert1
	ch <- flush
	close(ch)

	require.ErrorContains(t, wr.Write(context.Background(), ch), "insert failed")
	require.ErrorContains(t, receive(t, flush.Result()), "insert failed")
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
	flush := message.NewWriteFlush()
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
	flush := message.NewWriteFlush()
	ch <- flush
	require.NoError(t, receive(t, flush.Result()))
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}}, client.receivedBatches)

	ch <- deleteRecord
	flush = message.NewWriteFlush()
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
