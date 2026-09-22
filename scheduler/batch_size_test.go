package scheduler

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func collectBatches(t *testing.T, settings BatchSettings, rows int, valueSize int) message.SyncInserts {
	t.Helper()

	table := &schema.Table{
		Name:    "test_table",
		Columns: []schema.Column{{Name: "data", Type: arrow.BinaryTypes.String}},
	}

	res := make(chan message.SyncMessage, rows)
	batcher := settings.getBatcher(context.Background(), res, zerolog.Nop())
	for range rows {
		resource := schema.NewResourceData(table, nil, nil)
		require.NoError(t, resource.Set("data", strings.Repeat("a", valueSize)))
		batcher.process(resource)
	}
	batcher.close()
	close(res)

	var inserts message.SyncInserts
	for msg := range res {
		insert, ok := msg.(*message.SyncInsert)
		require.True(t, ok)
		inserts = append(inserts, insert)
	}
	return inserts
}

func TestBatcherMaxSizeBytes(t *testing.T) {
	const (
		rows         = 40
		valueSize    = 1024
		maxSizeBytes = 8 * 1024
	)
	settings := BatchSettings{
		MaxRows:      rows,
		MaxSizeBytes: maxSizeBytes,
		Timeout:      time.Hour,
	}

	unsplit := collectBatches(t, BatchSettings{MaxRows: rows, Timeout: time.Hour}, rows, valueSize)
	require.Len(t, unsplit, 1)
	bytesPerRow := util.TotalRecordSize(unsplit[0].Record) / unsplit[0].Record.NumRows()

	inserts := collectBatches(t, settings, rows, valueSize)
	require.Greater(t, len(inserts), 1, "oversized batch should be split into several messages")

	var gotRows int64
	for _, insert := range inserts {
		require.LessOrEqual(t, insert.Record.NumRows()*bytesPerRow, int64(maxSizeBytes))
		gotRows += insert.Record.NumRows()
	}
	require.Equal(t, int64(rows), gotRows, "no rows may be dropped when splitting")
}

func TestBatcherWithoutMaxSizeBytes(t *testing.T) {
	const rows = 40
	settings := BatchSettings{MaxRows: rows, Timeout: time.Hour}

	inserts := collectBatches(t, settings, rows, 1024)
	require.Len(t, inserts, 1)
	require.Equal(t, int64(rows), inserts[0].Record.NumRows())
}

func TestBatcherTableWithoutColumns(t *testing.T) {
	table := &schema.Table{Name: "no_columns"}
	res := make(chan message.SyncMessage, 16)
	batcher := (&BatchSettings{MaxRows: 2, MaxSizeBytes: 1024, Timeout: time.Hour}).
		getBatcher(context.Background(), res, zerolog.Nop())

	batcher.process(schema.NewResourceData(table, nil, nil))
	batcher.close()
	close(res)

	for msg := range res {
		insert, ok := msg.(*message.SyncInsert)
		require.True(t, ok)
		require.Equal(t, int64(0), insert.Record.NumRows())
	}
}

func TestBatcherMoreRowsThanBytes(t *testing.T) {
	const rows = 20000
	table := &schema.Table{
		Name:    "bools",
		Columns: []schema.Column{{Name: "b", Type: arrow.FixedWidthTypes.Boolean}},
	}
	res := make(chan message.SyncMessage, 4*rows)
	batcher := (&BatchSettings{MaxRows: rows, MaxSizeBytes: 64, Timeout: time.Hour}).
		getBatcher(context.Background(), res, zerolog.Nop())

	for range rows {
		resource := schema.NewResourceData(table, nil, nil)
		require.NoError(t, resource.Set("b", true))
		batcher.process(resource)
	}
	batcher.close()
	close(res)

	var gotRows int64
	for msg := range res {
		insert, ok := msg.(*message.SyncInsert)
		require.True(t, ok)
		gotRows += insert.Record.NumRows()
	}
	require.Equal(t, int64(rows), gotRows)
}

func TestWorkerReachedSizeLimit(t *testing.T) {
	for _, tc := range []struct {
		name         string
		maxSizeBytes int64
		bytesPerRow  int64
		curRows      int
		want         bool
	}{
		{name: "no cap", maxSizeBytes: 0, bytesPerRow: 1000, curRows: 1_000_000},
		{name: "unmeasured below cap", maxSizeBytes: 4096, curRows: unmeasuredBatchMaxRows - 1},
		{name: "unmeasured at cap", maxSizeBytes: 4096, curRows: unmeasuredBatchMaxRows, want: true},
		{name: "measured below cap", maxSizeBytes: 4096, bytesPerRow: 100, curRows: 40},
		{name: "measured at cap", maxSizeBytes: 4096, bytesPerRow: 100, curRows: 41, want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := &worker{maxSizeBytes: tc.maxSizeBytes, bytesPerRow: tc.bytesPerRow, curRows: tc.curRows}
			require.Equal(t, tc.want, w.reachedSizeLimit())
		})
	}
}
