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

func (s *WriterTestSuite) testInsert(ctx context.Context) error {
	tableName := fmt.Sprintf("cq_test_insert_%d", time.Now().Unix())
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			{Name: "name", Type: arrow.BinaryTypes.String},
		},
	}
	if err := s.plugin.writeOne(ctx, WriteOptions{}, &message.MigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("foo")
	record := bldr.NewRecord()

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &message.Insert{
		Record: record,
		Upsert: false,
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

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &message.Insert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	readRecords, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	totalItems = TotalRows(readRecords)
	if totalItems != 2 {
		return fmt.Errorf("expected 2 item, got %d", totalItems)
	}

	return nil
}
