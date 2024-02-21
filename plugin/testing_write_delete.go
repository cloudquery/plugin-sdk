package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func (s *WriterTestSuite) testDeleteStaleBasic(ctx context.Context, t *testing.T) {
	r := require.New(t)
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
	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")
	const sourceName = "source-test"

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(0)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime)
	record1 := bldr.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: record1}), "failed to insert record")
	record1 = s.handleNulls(record1) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read")
	r.EqualValuesf(1, TotalRows(records), "unexpected amount of items")

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: sourceName,
		SyncTime:   syncTime,
	}), "failed to delete stale records")

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after delete stale")
	r.EqualValuesf(1, TotalRows(records), "unexpected amount of items after delete stale")
	r.Emptyf(RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record1}), "record differs after delete stale")

	bldr.Field(0).(*array.Int64Builder).Append(1)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime.Add(time.Second))
	record2 := bldr.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: record2}), "failed to insert second record")
	record2 = s.handleNulls(record2) // we process nulls after writing

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read second time")
	sortRecords(table, records, "id")
	r.EqualValuesf(2, TotalRows(records), "unexpected amount of items second time")
	r.Emptyf(RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record1, record2}), "record differs after delete stale")

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: sourceName,
		SyncTime:   syncTime.Add(time.Second),
	}), "failed to delete stale records second time")

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after second delete stale")
	r.EqualValuesf(1, TotalRows(records), "unexpected amount of items after second delete stale")
	r.Emptyf(RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record2}), "record differs after second delete stale")
}

func (s *WriterTestSuite) testDeleteStaleAll(ctx context.Context, t *testing.T) {
	const rowsPerRecord = 10

	r := require.New(t)
	tableName := s.tableNameForTest("delete_all")
	// https://github.com/golang/go/issues/41087
	syncTime := time.Now().UTC().Truncate(time.Microsecond)
	table := schema.TestTable(tableName, s.genDatOptions)
	table.Columns = append(schema.ColumnList{schema.CqSourceNameColumn, schema.CqSyncTimeColumn}, table.Columns...)
	table.Columns[table.Columns.Index("id")].PrimaryKey = true
	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")

	tg := schema.NewTestDataGenerator(0)
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		SourceName:    "test",
		SyncTime:      syncTime, // Generate call may truncate the value further based on the options
	})
	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: normalRecord}), "failed to insert record")
	normalRecord = s.handleNulls(normalRecord) // we process nulls after writing

	readRecords, err := s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read")
	r.EqualValuesf(rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after read")

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime, // Generate call may truncate the value further based on the options
	}), "failed to delete stale records")

	readRecords, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after delete stale")
	r.EqualValuesf(rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after delete stale")

	// https://github.com/golang/go/issues/41087
	syncTime = time.Now().UTC().Truncate(time.Microsecond)
	nullRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		NullRows:      true,
		SourceName:    "test",
		SyncTime:      syncTime, // Generate call may truncate the value further based on the options
	})
	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: nullRecord}), "failed to insert record second time")
	nullRecord = s.handleNulls(nullRecord) // we process nulls after writing

	readRecords, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read second time")
	sortRecords(table, readRecords, "id")
	r.EqualValuesf(2*rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after second read")
	r.Emptyf(RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{normalRecord, nullRecord}), "record differs")

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime, // Generate call may truncate the value further based on the options
	}), "failed to delete stale records second time")

	readRecords, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after second delete stale")
	sortRecords(table, readRecords, "id")
	r.EqualValuesf(rowsPerRecord, TotalRows(readRecords), "unexpected amount of items after second delete stale")
	r.Emptyf(RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{nullRecord}), "record differs")
}

func (s *WriterTestSuite) testDeleteRecordBasic(ctx context.Context, t *testing.T) {
	r := require.New(t)
	tableName := s.tableNameForTest("delete_all_rows")
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
	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")
	const sourceName = "source-test"

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(0)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime)
	record1 := bldr.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: record1}), "failed to insert record")
	record1 = s.handleNulls(record1) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read")
	r.EqualValuesf(1, TotalRows(records), "unexpected amount of items")

	// create value for delete statement but nothing will be deleted because ID value isn't present
	bldrDeleteNoMatch := array.NewRecordBuilder(memory.DefaultAllocator, (&schema.Table{
		Name: tableName,
		Columns: schema.ColumnList{
			schema.Column{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		},
	}).ToArrowSchema())
	bldrDeleteNoMatch.Field(0).(*array.Int64Builder).Append(1)
	deleteValue := bldrDeleteNoMatch.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteRecord{
		DeleteRecord: message.DeleteRecord{
			TableName: table.Name,
			WhereClause: message.PredicateGroups{
				{
					GroupingType: "AND",
					Predicates: []message.Predicate{
						{
							Operator: "eq",
							Column:   "id",
							Record:   deleteValue,
						},
					},
				},
			},
		},
	}), "failed to delete record no match")

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after delete with no match")
	r.EqualValuesf(1, TotalRows(records), "unexpected amount of items after delete with no match")
	r.Emptyf(RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record1}), "record differs after delete with no match")

	// create value for delete statement will be delete One record
	bldrDeleteMatch := array.NewRecordBuilder(memory.DefaultAllocator, (&schema.Table{
		Name: tableName,
		Columns: schema.ColumnList{
			schema.Column{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		},
	}).ToArrowSchema())
	bldrDeleteMatch.Field(0).(*array.Int64Builder).Append(0)
	deleteValue = bldrDeleteMatch.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteRecord{
		DeleteRecord: message.DeleteRecord{
			TableName: table.Name,
			WhereClause: message.PredicateGroups{
				{
					GroupingType: "AND",
					Predicates: []message.Predicate{
						{
							Operator: "eq",
							Column:   "id",
							Record:   deleteValue,
						},
					},
				},
			},
		},
	}), "failed to delete record no match")

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after delete with match")
	r.EqualValuesf(0, TotalRows(records), "unexpected amount of items after delete with match")
}

func (s *WriterTestSuite) testDeleteAllRecords(ctx context.Context, t *testing.T) {
	r := require.New(t)
	tableName := s.tableNameForTest("delete_all_records")
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
	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteMigrateTable{Table: table}), "failed to create table")
	const sourceName = "source-test"

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(0)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime)
	record1 := bldr.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: record1}), "failed to insert record")

	records, err := s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read")
	r.EqualValuesf(1, TotalRows(records), "unexpected amount of items")

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteRecord{
		DeleteRecord: message.DeleteRecord{
			TableName: table.Name,
		},
	}), "failed to delete records")

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read after delete all records")
	r.EqualValuesf(0, TotalRows(records), "unexpected amount of items after delete stale")

	bldr.Field(0).(*array.Int64Builder).Append(1)
	bldr.Field(1).(*array.StringBuilder).Append(sourceName)
	bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime.Add(time.Second))
	record2 := bldr.NewRecord()

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteInsert{Record: record2}), "failed to insert second record")

	r.NoErrorf(s.plugin.writeOne(ctx, &message.WriteDeleteRecord{
		DeleteRecord: message.DeleteRecord{
			TableName: table.Name,
		},
	}), "failed to delete records second time")

	records, err = s.plugin.readAll(ctx, table)
	r.NoErrorf(err, "failed to read second time")
	sortRecords(table, records, "id")
	r.EqualValuesf(0, TotalRows(records), "unexpected amount of items second time")
}
