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
	require.NoError(t, <-flush.Result())
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}, {msgs.insert1}}, client.receivedBatches)

	ch <- msgs.insert2
	flush = message.NewWriteFlush()
	ch <- flush
	require.NoError(t, <-flush.Result())
	require.Equal(t, [][]message.WriteMessage{{msgs.migrateTable1}, {msgs.insert1}, {msgs.insert2}}, client.receivedBatches)

	close(ch)
	require.NoError(t, <-errCh)
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
	require.ErrorContains(t, <-flush.Result(), "insert failed")
}
