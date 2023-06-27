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

func (s *WriterTestSuite) testUpsert(ctx context.Context) error {
	tableName := fmt.Sprintf("cq_test_upsert_%d", time.Now().Unix())
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKey: true},
		},
	}
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("foo")
	record := bldr.NewRecord()

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

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

	return nil
}
