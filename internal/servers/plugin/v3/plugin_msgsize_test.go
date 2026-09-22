package plugin

import (
	"strings"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type recordingSyncServer struct {
	mockSyncServer
	sent []*pb.Sync_Response
}

func (s *recordingSyncServer) Send(msg *pb.Sync_Response) error {
	s.sent = append(s.sent, msg)
	return nil
}

func newTestRecord(t *testing.T, tableName string, rows int, valueSize int) arrow.RecordBatch {
	t.Helper()

	table := schema.Table{
		Name:    tableName,
		Columns: []schema.Column{{Name: "data", Type: arrow.BinaryTypes.String}},
	}
	builder := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	defer builder.Release()

	for i := range rows {
		builder.Field(0).(*array.StringBuilder).Append(strings.Repeat(string(rune('a'+i%26)), valueSize))
	}

	return builder.NewRecordBatch()
}

func setMaxMsgSize(t *testing.T, size int) {
	t.Helper()

	original := MaxMsgSize
	MaxMsgSize = size
	t.Cleanup(func() { MaxMsgSize = original })
}

func TestSendInsertSplitsOversizedRecord(t *testing.T) {
	setMaxMsgSize(t, 8*1024)

	const (
		rows      = 40
		valueSize = 1024
	)
	record := newTestRecord(t, "test_table", rows, valueSize)
	s := &Server{Logger: zerolog.Nop()}
	stream := &recordingSyncServer{}

	require.NoError(t, s.sendInsert(stream, record))
	require.Greater(t, len(stream.sent), 1, "oversized record should be split into several messages")

	var gotRows int64
	for _, msg := range stream.sent {
		require.LessOrEqual(t, proto.Size(msg), MaxMsgSize)

		sentRecord, err := pb.NewRecordFromBytes(msg.GetInsert().GetRecord())
		require.NoError(t, err)
		gotRows += sentRecord.NumRows()
	}
	require.Equal(t, int64(rows), gotRows, "no rows may be dropped when splitting")
}

func TestSendInsertFitsInSingleMessage(t *testing.T) {
	record := newTestRecord(t, "test_table", 10, 16)
	s := &Server{Logger: zerolog.Nop()}
	stream := &recordingSyncServer{}

	require.NoError(t, s.sendInsert(stream, record))
	require.Len(t, stream.sent, 1)
}

func TestSendInsertSingleRowTooLarge(t *testing.T) {
	setMaxMsgSize(t, 1024)

	record := newTestRecord(t, "huge_table", 1, 4*1024)
	s := &Server{Logger: zerolog.Nop()}
	stream := &recordingSyncServer{}

	err := s.sendInsert(stream, record)
	require.Error(t, err, "a row that cannot be split must fail the sync instead of being dropped")
	require.Contains(t, err.Error(), "huge_table")
	require.Empty(t, stream.sent)
}

func TestMaxRowsPerMessage(t *testing.T) {
	setMaxMsgSize(t, 1000)

	for _, tc := range []struct {
		name string
		size int
		rows int64
	}{
		{name: "slightly over", size: 1001, rows: 100},
		{name: "twice over", size: 2000, rows: 100},
		{name: "far over", size: 1_000_000, rows: 100},
		{name: "two rows", size: 5000, rows: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := maxRowsPerMessage(tc.size, tc.rows)
			require.GreaterOrEqual(t, got, int64(1), "must make progress")
			require.Less(t, got, tc.rows, "must shrink the batch")
		})
	}
}
