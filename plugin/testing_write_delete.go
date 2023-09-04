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
	syncTime := time.Now().UTC().Truncate(s.genDatOptions.TimePrecision).
		Truncate(time.Microsecond) // https://github.com/golang/go/issues/41087
	table := &schema.Table{
		Name: tableName,
		Columns: schema.ColumnList{
			schema.Column{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true, NotNull: true},
			schema.CqSourceNameColumn,
			schema.CqSyncTimeColumn,
		},
	}
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")
	const sourceName = "source-test"

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(0)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime)
	record1 := bldr.NewRecord()

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteInsert{Record: record1}), "failed to insert record")
	record1 = s.handleNulls(record1) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read")
	require.EqualValuesf(s.t, 1, TotalRows(records), "unexpected amount of items")

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: sourceName,
		SyncTime:   syncTime,
	}), "failed to delete stale records")

	records, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read after delete stale")
	require.EqualValuesf(s.t, 1, TotalRows(records), "unexpected amount of items after delete stale")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record1}), "record differs after delete stale")

	bldr.Field(0).(*array.Int64Builder).Append(1)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime.Add(time.Second))
	record2 := bldr.NewRecord()

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteInsert{Record: record2}), "failed to insert second record")
	record2 = s.handleNulls(record2) // we process nulls after writing

	records, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read second time")
	sortRecords(table, records, "id")
	require.EqualValuesf(s.t, 2, TotalRows(records), "unexpected amount of items second time")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record1, record2}), "record differs after delete stale")

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: sourceName,
		SyncTime:   syncTime.Add(time.Second),
	}), "failed to delete stale records second time")

	records, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read after second delete stale")
	require.EqualValuesf(s.t, 1, TotalRows(records), "unexpected amount of items after second delete stale")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record2}), "record differs after second delete stale")
}

func (s *WriterTestSuite) testDeleteStaleAll(ctx context.Context) {
	const rowsPerRecord = 10

	tableName := s.tableNameForTest("delete_all")
	// https://github.com/golang/go/issues/41087
	syncTime := time.Now().UTC().Truncate(time.Microsecond)
	table := schema.TestTable(tableName, s.genDatOptions)
	table.Columns = append(schema.ColumnList{schema.CqSourceNameColumn, schema.CqSyncTimeColumn}, table.Columns...)
	table.Columns[table.Columns.Index("id")].PrimaryKey = true
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")

	tg := schema.NewTestDataGenerator()
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		SourceName:    "test",
		SyncTime:      syncTime, // Generate call may truncate the value further based on the options
	})
	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteInsert{Record: normalRecord}), "failed to insert record")
	normalRecord = s.handleNulls(normalRecord) // we process nulls after writing

	readRecords, err := s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read")
	require.EqualValuesf(s.t, rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after read")

	require.NoErrorf(s.t, s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime, // Generate call may truncate the value further based on the options
	}), "failed to delete stale records")

	readRecords, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read after delete stale")
	require.EqualValuesf(s.t, rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after delete stale")

	// https://github.com/golang/go/issues/41087
	syncTime = time.Now().UTC().Truncate(time.Microsecond)
	nullRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		NullRows:      true,
		SourceName:    "test",
		SyncTime:      syncTime, // Generate call may truncate the value further based on the options
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
		SyncTime:   syncTime, // Generate call may truncate the value further based on the options
	}), "failed to delete stale records second time")

	readRecords, err = s.plugin.readAll(ctx, table)
	require.NoErrorf(s.t, err, "failed to read after second delete stale")
	sortRecords(table, readRecords, "id")
	require.EqualValuesf(s.t, rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after second delete stale")
	require.Emptyf(s.t, RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{nullRecord}), "record differs")
}
