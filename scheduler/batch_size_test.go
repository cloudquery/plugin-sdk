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

	// util.TotalRecordSize reports the size of the whole parent buffer for a slice,
	// so measure the per row size on an unsplit batch of the same data instead.
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
