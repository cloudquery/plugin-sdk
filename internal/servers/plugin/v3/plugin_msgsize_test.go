package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
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

func TestSendInsertSplitsOversizedRecord(t *testing.T) {
	t.Parallel()

	const (
		rows       = 40
		valueSize  = 1024
		maxMsgSize = 8 * 1024
	)
	record := newTestRecord(t, "test_table", rows, valueSize)
	s := &Server{Logger: zerolog.Nop(), MaxMsgSize: maxMsgSize}
	stream := &recordingSyncServer{}

	require.NoError(t, s.sendInsert(stream, record))
	require.Greater(t, len(stream.sent), 1, "oversized record should be split into several messages")

	var gotRows int64
	for _, msg := range stream.sent {
		require.LessOrEqual(t, proto.Size(msg), maxMsgSize)

		sentRecord, err := pb.NewRecordFromBytes(msg.GetInsert().GetRecord())
		require.NoError(t, err)
		gotRows += sentRecord.NumRows()
	}
	require.Equal(t, int64(rows), gotRows, "no rows may be dropped when splitting")
}

func TestSendInsertFitsInSingleMessage(t *testing.T) {
	t.Parallel()

	record := newTestRecord(t, "test_table", 10, 16)
	s := &Server{Logger: zerolog.Nop()}
	stream := &recordingSyncServer{}

	require.NoError(t, s.sendInsert(stream, record))
	require.Len(t, stream.sent, 1)
}

func TestSendInsertSingleRowTooLarge(t *testing.T) {
	t.Parallel()

	record := newTestRecord(t, "huge_table", 1, 4*1024)
	s := &Server{Logger: zerolog.Nop(), MaxMsgSize: 1024}
	stream := &recordingSyncServer{}

	err := s.sendInsert(stream, record)
	require.Error(t, err, "a row that cannot be split must fail the sync instead of being dropped")
	require.Contains(t, err.Error(), "huge_table")
	require.Empty(t, stream.sent)
}

func TestMaxRowsPerMessage(t *testing.T) {
	t.Parallel()

	const maxMsgSize = 1000

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
			got := maxRowsPerMessage(tc.size, tc.rows, maxMsgSize)
			require.GreaterOrEqual(t, got, int64(1), "must make progress")
			require.Less(t, got, tc.rows, "must shrink the batch")
		})
	}
}

type migrateTableOnlyClient struct {
	plugin.UnimplementedDestination
	table *schema.Table
}

func (migrateTableOnlyClient) Close(context.Context) error { return nil }

func (c migrateTableOnlyClient) Tables(context.Context, plugin.TableOptions) (schema.Tables, error) {
	return schema.Tables{c.table}, nil
}

func (c migrateTableOnlyClient) Sync(_ context.Context, _ plugin.SyncOptions, res chan<- message.SyncMessage) error {
	res <- &message.SyncMigrateTable{Table: c.table}
	return nil
}

func wideSchemaTable(columns int) *schema.Table {
	table := &schema.Table{Name: "wide_table", Columns: make(schema.ColumnList, columns)}
	for i := range table.Columns {
		table.Columns[i] = schema.Column{
			Name: fmt.Sprintf("column_with_a_deliberately_long_name_%d", i),
			Type: arrow.BinaryTypes.String,
		}
	}
	return table
}

func TestSyncNonInsertExceedingMaxSizeFailsSync(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	table := wideSchemaTable(200)
	s := &Server{
		Logger:     zerolog.Nop(),
		MaxMsgSize: 512,
		Plugin: plugin.NewPlugin("test", "development",
			func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
				return migrateTableOnlyClient{table: table}, nil
			}),
	}
	_, err := s.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	stream := &recordingSyncServer{}
	err = s.Sync(&pb.Sync_Request{}, stream)

	require.Error(t, err, "an oversized non-insert message must fail the sync, not be dropped")
	require.Contains(t, err.Error(), "exceeds max message size")
	require.Contains(t, err.Error(), "SyncMigrateTable")
	require.Empty(t, stream.sent)
}
