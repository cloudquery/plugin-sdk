package plugin

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *WriterTestSuite) testUpsertBasic(ctx context.Context) error {
	tableName := s.tableNameForTest("upsert_basic")
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, NotNull: true},
			{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKey: true},
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

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to readAll: %w", err)
	}
	totalItems := TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	records, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems = TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}
	if diff := RecordDiff(records[0], record); diff != "" {
		return fmt.Errorf("record differs: %s", diff)
	}
	return nil
}

func (s *WriterTestSuite) testUpsertAll(ctx context.Context) error {
	tableName := s.tableNameForTest("upsert_all")
	table := schema.TestTable(tableName, s.genDatOptions)
	table.Columns = append(table.Columns, schema.Column{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKey: true})
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	tg := schema.NewTestDataGenerator()
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{MaxRows: 10, TimePrecision: s.genDatOptions.TimePrecision})
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: normalRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	normalRecord = s.handleNulls(normalRecord) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to readAll: %w", err)
	}
	totalItems := TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if diff := RecordDiff(records[0], normalRecord); diff != "" {
		return fmt.Errorf("record differs after insert: %s", diff)
	}

	nullRecord := tg.Generate(table, schema.GenTestDataOptions{MaxRows: 10, TimePrecision: s.genDatOptions.TimePrecision, NullRows: true})
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: nullRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	nullRecord = s.handleNulls(nullRecord) // we process nulls after writing

	records, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	totalItems = TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if diff := RecordDiff(records[0], nullRecord); diff != "" {
		return fmt.Errorf("record differs after upsert (columns should be null): %s", diff)
	}

	return nil
}
