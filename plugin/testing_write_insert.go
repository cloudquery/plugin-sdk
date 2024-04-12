package plugin

import (
	"context"
	"fmt"
	"slices"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func TotalRows(records []arrow.Record) int64 {
	totalRows := int64(0)
	for _, record := range records {
		totalRows += record.NumRows()
	}
	return totalRows
}

func SkipRows(records []arrow.Record, toSkip int64) []arrow.Record {
	if toSkip < 0 {
		panic(fmt.Sprintf("toSkip(%d) < 0", toSkip))

	}
	if total := TotalRows(records); total < toSkip {
		panic(fmt.Sprintf("total(%d) < toSkip(%d)", total, toSkip))
	}
	var first arrow.Record
	result := slices.Clone(records)
	for toSkip > 0 {
		first, result = result[0], result[1:]
		rows := first.NumRows()
		if rows > toSkip {
			// we need to split rows
			return append([]arrow.Record{first.NewSlice(toSkip, rows)}, result...)
		}
		toSkip -= rows
	}
	// shouldn't even get here
	return nil
}

func (s *WriterTestSuite) testInsertBasic(ctx context.Context) error {
	tableName := s.tableNameForTest("insert_basic")
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, NotNull: true},
			{Name: "name", Type: arrow.BinaryTypes.String},
		},
	}
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	bldr.Field(1).(*array.StringBuilder).Append("foo")
	record := bldr.NewRecord()

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	record = s.handleNulls(record) // we process nulls after writing

	readRecords, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	totalItems := TotalRows(readRecords)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	readRecords, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	sortRecords(table, readRecords, "id")

	totalItems = TotalRows(readRecords)
	if totalItems != 2 {
		return fmt.Errorf("expected 2 items, got %d", totalItems)
	}

	if diff := RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{record, record}); diff != "" {
		return fmt.Errorf("record[0] differs: %s", diff)
	}

	return nil
}

func (s *WriterTestSuite) testInsertAll(ctx context.Context) error {
	const rowsPerRecord = 10
	tableName := s.tableNameForTest("insert_all")
	table := schema.TestTable(tableName, s.genDatOptions)
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	tg := schema.NewTestDataGenerator(0)
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:            rowsPerRecord,
		TimePrecision:      s.genDatOptions.TimePrecision,
		UseHomogeneousType: s.useHomogeneousTypes,
	})
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: normalRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	normalRecord = s.handleNulls(normalRecord) // we process nulls after writing

	readRecords, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	totalItems := TotalRows(readRecords)
	if totalItems != rowsPerRecord {
		return fmt.Errorf("items expected after first insert: %d, got: %d", rowsPerRecord, totalItems)
	}

	nullRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:       rowsPerRecord,
		TimePrecision: s.genDatOptions.TimePrecision,
		NullRows:      true,
	})
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: nullRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	nullRecord = s.handleNulls(nullRecord) // we process nulls after writing

	readRecords, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	sortRecords(table, readRecords, "id")

	totalItems = TotalRows(readRecords)
	if totalItems != 2*rowsPerRecord {
		return fmt.Errorf("items expected after second insert: %d, got: %d", 2*rowsPerRecord, totalItems)
	}
	if diff := RecordsDiff(table.ToArrowSchema(), readRecords, []arrow.Record{normalRecord, nullRecord}); diff != "" {
		return fmt.Errorf("record[0] differs: %s", diff)
	}
	return nil
}
