package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
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

func (s *WriterTestSuite) testInsertBasic(ctx context.Context) error {
	tableName := fmt.Sprintf("cq_insert_basic_%d", time.Now().Unix())
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

	if diff := RecordDiff(readRecords[0], record); diff != "" {
		return fmt.Errorf("record[0] differs: %s", diff)
	}
	if diff := RecordDiff(readRecords[1], record); diff != "" {
		return fmt.Errorf("record[1] differs: %s", diff)
	}

	return nil
}

func (s *WriterTestSuite) testInsertAll(ctx context.Context) error {
	tableName := fmt.Sprintf("cq_insert_all_%d", time.Now().Unix())
	table := schema.TestTable(tableName, s.genDatOptions)
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	tg := schema.NewTestDataGenerator()
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows: 1,
	})[0]
	nullRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:  1,
		NullRows: true,
	})[0]

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: normalRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	readRecords, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	totalItems := TotalRows(readRecords)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: nullRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	readRecords, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	totalItems = TotalRows(readRecords)
	if totalItems != 2 {
		return fmt.Errorf("expected 2 items, got %d", totalItems)
	}

	wantNormalRecord := s.allowNull.replaceNullsWithEmpty(normalRecord)
	if s.ignoreNullsInLists {
		wantNormalRecord = stripNullsFromLists(wantNormalRecord)
	}

	if diff := RecordDiff(readRecords[0], wantNormalRecord); diff != "" {
		return fmt.Errorf("record[0] differs: %s", diff)
	}

	wantNullRecord := s.allowNull.replaceNullsWithEmpty(nullRecord)
	if s.ignoreNullsInLists {
		wantNullRecord = stripNullsFromLists(wantNullRecord)
	}

	if diff := RecordDiff(readRecords[1], wantNullRecord); diff != "" {
		return fmt.Errorf("record[1] differs: %s", diff)
	}
	return nil
}
