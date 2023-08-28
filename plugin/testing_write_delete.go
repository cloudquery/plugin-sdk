package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *WriterTestSuite) testDeleteStale(ctx context.Context) error {
	tableName := s.tableNameForTest("delete")
	syncTime := time.Now().UTC().Round(1 * time.Second)
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			schema.CqSourceNameColumn,
			schema.CqSyncTimeColumn,
		},
	}
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	bldr.Field(1).(*array.TimestampBuilder).AppendTime(syncTime)
	record := bldr.NewRecord()

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	record = s.handleNulls(record) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems := TotalRows(records)

	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	bldr = array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	bldr.Field(1).(*array.TimestampBuilder).AppendTime(syncTime.Add(time.Second))

	if err := s.plugin.writeOne(ctx, &message.WriteDeleteStale{
		TableName:  table.Name,
		SourceName: "test",
		SyncTime:   syncTime,
	}); err != nil {
		return fmt.Errorf("failed to delete stale records: %w", err)
	}

	records, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems = TotalRows(records)

	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}
	if diff := RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record}); diff != "" {
		return fmt.Errorf("record differs: %s", diff)
	}
	return nil
}
