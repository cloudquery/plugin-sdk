package plugin

import (
	"context"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func (s *WriterTestSuite) testDeleteStaleBasic(ctx context.Context) {
	tableName := s.tableNameForTest("delete_basic")
	syncTime := time.Now().UTC().Round(1 * time.Second)
	table := &schema.Table{
		Name:    tableName,
		Columns: schema.ColumnList{schema.CqSourceNameColumn, schema.CqSyncTimeColumn},
	}
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	bldr.Field(1).(*array.TimestampBuilder).AppendTime(syncTime)
	record := bldr.NewRecord()

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteInsert{Record: record}), "failed to insert record")
	record = s.handleNulls(record) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read")
	require.EqualValuesf(s.t, 1, TotalRows(records), "unexpected amount of items")

	bldr = array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	bldr.Field(1).(*array.TimestampBuilder).AppendTime(syncTime.Add(time.Second))

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime,
	}), "failed to delete stale records")

	records, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read second time")
	require.EqualValuesf(s.t, 1, TotalRows(records), "unexpected amount of items")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record}), "record differs")
}

func (s *WriterTestSuite) testDeleteStaleAll(ctx context.Context) {
	const rowsPerRecord = 10

	tableName := s.tableNameForTest("delete_all")
	syncTime := time.Now().UTC().Truncate(s.genDatOptions.TimePrecision)
	table := schema.TestTable(tableName, s.genDatOptions)
	table.Columns = append(schema.ColumnList{schema.CqSourceNameColumn, schema.CqSyncTimeColumn}, table.Columns...)
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")

	tg := schema.NewTestDataGenerator()
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		SourceName:    "test",
		SyncTime:      syncTime,
	})
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteInsert{Record: normalRecord}), "failed to insert record")
	normalRecord = s.handleNulls(normalRecord) // we process nulls after writing

	readRecords, err := s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read")
	require.EqualValuesf(s.t, rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after read")

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime,
	}), "failed to delete stale records")

	readRecords, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read after delete stale")
	require.EqualValuesf(s.t, rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after delete stale")

	syncTime = time.Now().UTC().Truncate(s.genDatOptions.TimePrecision) // bump sync time
	nullRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		NullRows:      true,
		SourceName:    "test",
		SyncTime:      syncTime,
	})
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteInsert{Record: nullRecord}), "failed to insert record second time")
	nullRecord = s.handleNulls(nullRecord) // we process nulls after writing

	readRecords, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read second time")
	sortRecords(table, readRecords, "id")
	require.EqualValuesf(s.t, 2*rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after second read")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{normalRecord, nullRecord}), "record differs")

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime,
	}), "failed to delete stale records second time")

	readRecords, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read after second delete stale")
	sortRecords(table, readRecords, "id")
	require.EqualValuesf(s.t, rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after second delete stale")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{nullRecord}), "record differs")
}
