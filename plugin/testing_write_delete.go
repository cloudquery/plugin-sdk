package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	// "github.com/cloudquery/plugin-sdk/v4/types"
)

func (s *WriterTestSuite) testDeleteStale(ctx context.Context) error {
	tableName := fmt.Sprintf("cq_delete_%d", time.Now().Unix())
	syncTime := time.Now().UTC().Round(1 * time.Second)
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			schema.CqSourceNameColumn,
			schema.CqSyncTimeColumn,
		},
	}
	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageCreateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	bldr.Field(1).(*array.TimestampBuilder).AppendTime(syncTime)
	record := bldr.NewRecord()

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	messages, err := s.plugin.SyncAll(ctx, SyncOptions{
		Tables: []string{tableName},
	})
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems := messages.InsertItems()

	if totalItems != 1 {
		return fmt.Errorf("expected 1 items, got %d", totalItems)
	}

	bldr = array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	bldr.Field(1).(*array.TimestampBuilder).AppendTime(syncTime.Add(time.Second))

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageDeleteStale{
		Table:      table,
		SourceName: "test",
		SyncTime:   syncTime,
	}); err != nil {
		return fmt.Errorf("failed to delete stale records: %w", err)
	}

	messages, err = s.plugin.SyncAll(ctx, SyncOptions{
		Tables: []string{tableName},
	})
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems = messages.InsertItems()

	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	return nil
}
