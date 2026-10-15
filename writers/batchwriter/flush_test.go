package batchwriter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
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
		return receive(t, errCh)
	}
}

func flush(t *testing.T, ch chan<- message.WriteMessage) error {
	t.Helper()
	msg := message.NewWriteFlush()
	ch <- msg
	return receive(t, msg.Result())
}

func TestFlushMessageCommitsBatches(t *testing.T) {
	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 3)}
	require.NoError(t, flush(t, ch))
	require.Equal(t, 1, testClient.MigrateTablesLen())
	require.EqualValues(t, 3, testClient.InsertedRows())

	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 2)}
	require.NoError(t, flush(t, ch))
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
	require.ErrorContains(t, flush(t, ch), "write failed")
	// the failure was reported already
	require.NoError(t, flush(t, ch))

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}

func TestFlushMessageReportsFailuresFromTimeouts(t *testing.T) {
	wr, err := New(&failingBatchClient{}, WithBatchTimeout(time.Millisecond))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 3)}
	require.Eventually(t, func() bool {
		wr.failuresLock.Lock()
		defer wr.failuresLock.Unlock()
		return wr.failures[batchTestTables[0].Name].count > 0
	}, 5*time.Second, time.Millisecond)
	require.ErrorContains(t, flush(t, ch), "write failed")

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}

func TestFlushMessageCommitsDeletes(t *testing.T) {
	testClient := &testBatchClient{}
	wr, err := New(testClient, WithBatchTimeout(time.Hour))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteDeleteStale{TableName: batchTestTables[0].Name, SourceName: "source", SyncTime: time.Now()}
	require.NoError(t, flush(t, ch))
	require.Equal(t, 1, testClient.DeleteStalesLen())

	ch <- &message.WriteDeleteRecord{DeleteRecord: message.DeleteRecord{TableName: batchTestTables[0].Name}}
	require.NoError(t, flush(t, ch))
	testClient.mutex.Lock()
	require.Len(t, testClient.deleteRecords, 1)
	testClient.mutex.Unlock()

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}

func TestFlushMessageErrorNamesTable(t *testing.T) {
	wr, err := New(&failingBatchClient{}, WithBatchTimeout(time.Hour))
	require.NoError(t, err)
	ch, stop := startWrite(t, wr)

	ch <- &message.WriteMigrateTable{Table: batchTestTables[0]}
	ch <- &message.WriteInsert{Record: getRecord(batchTestTables[0].ToArrowSchema(), 1)}
	err = flush(t, ch)
	require.ErrorContains(t, err, "table "+batchTestTables[0].Name+": 1 batch writes failed since the previous flush, last error: write failed")
	require.Equal(t, []string{batchTestTables[0].Name}, message.FailedTables(err))

	require.NoError(t, stop())
	require.NoError(t, wr.Close(context.Background()))
}

type failingOpsClient struct {
	testBatchClient
	migrateErr, deleteStaleErr, deleteRecordErr error
}

func (c *failingOpsClient) MigrateTables(ctx context.Context, msgs message.WriteMigrateTables) error {
	if c.migrateErr != nil {
		return c.migrateErr
	}
	return c.testBatchClient.MigrateTables(ctx, msgs)
}

func (c *failingOpsClient) DeleteStale(ctx context.Context, msgs message.WriteDeleteStales) error {
	if c.deleteStaleErr != nil {
		return c.deleteStaleErr
	}
	return c.testBatchClient.DeleteStale(ctx, msgs)
}

func (c *failingOpsClient) DeleteRecord(ctx context.Context, msgs message.WriteDeleteRecords) error {
	if c.deleteRecordErr != nil {
		return c.deleteRecordErr
	}
	return c.testBatchClient.DeleteRecord(ctx, msgs)
}

func TestFlushMessageReportsFailedMigrationsAndDeletes(t *testing.T) {
	table := batchTestTables[0]
	for _, tc := range []struct {
		name   string
		client *failingOpsClient
		msgs   []message.WriteMessage
	}{
		{
			name:   "migrate",
			client: &failingOpsClient{migrateErr: errors.New("migrate failed")},
			msgs:   []message.WriteMessage{&message.WriteMigrateTable{Table: table}},
		},
		{
			name:   "delete stale",
			client: &failingOpsClient{deleteStaleErr: errors.New("delete stale failed")},
			msgs: []message.WriteMessage{
				&message.WriteMigrateTable{Table: table},
				&message.WriteDeleteStale{TableName: table.Name, SourceName: "source", SyncTime: time.Now()},
			},
		},
		{
			name:   "delete record",
			client: &failingOpsClient{deleteRecordErr: errors.New("delete record failed")},
			msgs: []message.WriteMessage{
				&message.WriteMigrateTable{Table: table},
				&message.WriteDeleteRecord{DeleteRecord: message.DeleteRecord{TableName: table.Name}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wr, err := New(tc.client, WithBatchTimeout(time.Hour))
			require.NoError(t, err)
			ch, stop := startWrite(t, wr)
			for _, msg := range tc.msgs {
				ch <- msg
			}
			err = flush(t, ch)
			require.ErrorContains(t, err, tc.name+" failed")
			require.Equal(t, []string{table.Name}, message.FailedTables(err))
			require.NoError(t, stop())
		})
	}
}

// failingTableClient fails to write the batches of one table.
type failingTableClient struct {
	testBatchClient
	table string
}

func (c *failingTableClient) WriteTableBatch(ctx context.Context, name string, msgs message.WriteInserts) error {
	if name == c.table {
		return errors.New("write failed")
	}
	return c.testBatchClient.WriteTableBatch(ctx, name, msgs)
}

func TestFlushMessageReportsOnlyTheFailuresOfItsOwnTables(t *testing.T) {
	table := func(name string) *schema.Table {
		return &schema.Table{Name: name, Columns: schema.ColumnList{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	}
	failing, healthy := table("failing"), table("healthy")
	wr, err := New(&failingTableClient{table: failing.Name}, WithBatchTimeout(time.Hour))
	require.NoError(t, err)
	// two streams share the writer, as the state backend's does with the sync's when both use one destination
	failingCh, stopFailing := startWrite(t, wr)
	healthyCh, stopHealthy := startWrite(t, wr)

	failingCh <- &message.WriteMigrateTable{Table: failing}
	healthyCh <- &message.WriteMigrateTable{Table: healthy}
	failingCh <- &message.WriteInsert{Record: getRecord(failing.ToArrowSchema(), 1)}
	healthyCh <- &message.WriteInsert{Record: getRecord(healthy.ToArrowSchema(), 1)}

	err = flush(t, failingCh)
	require.ErrorContains(t, err, "write failed")
	require.Equal(t, []string{failing.Name}, message.FailedTables(err))
	require.NoError(t, flush(t, healthyCh), "the failure of another stream's table is not reported")

	// a stream that starts writing a table later does not report its earlier failures
	healthyCh <- &message.WriteInsert{Record: getRecord(failing.ToArrowSchema(), 1)}
	err = flush(t, healthyCh)
	require.Equal(t, []string{failing.Name}, message.FailedTables(err), "only the failure of its own insert")
	require.ErrorContains(t, err, "1 batch writes failed")

	require.NoError(t, stopFailing())
	require.NoError(t, stopHealthy())
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
