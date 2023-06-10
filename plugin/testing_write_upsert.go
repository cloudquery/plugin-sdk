package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *PluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *Plugin) error {
	tableName := fmt.Sprintf("cq_test_upsert_%d", time.Now().Unix())
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKey: true},
		},
	}
	if err := p.writeOne(ctx, WriteOptions{}, &MessageCreateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("foo")

	if err := p.writeOne(ctx, WriteOptions{}, &MessageInsert{
		Record: bldr.NewRecord(),
		Upsert: true,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	messages, err := p.syncAll(ctx, SyncOptions{
		Tables: []string{tableName},
	})
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems := messages.InsertItems()
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if err := p.writeOne(ctx, WriteOptions{}, &MessageInsert{
		Record: bldr.NewRecord(),
		Upsert: true,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	messages, err = p.syncAll(ctx, SyncOptions{
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
