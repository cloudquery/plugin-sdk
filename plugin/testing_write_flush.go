package plugin

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// testFlush checks that rows are readable once a flush completes, while Write is still running.
func (s *WriterTestSuite) testFlush(ctx context.Context) error {
	table := &schema.Table{
		Name: s.tableNameForTest("flush"),
		Columns: []schema.Column{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, NotNull: true},
		},
	}

	msgs := make(chan message.WriteMessage)
	writeErr := make(chan error, 1)
	go func() {
		writeErr <- s.plugin.Write(ctx, msgs)
	}()
	send := func(msg message.WriteMessage) error {
		select {
		case msgs <- msg:
			return nil
		case err := <-writeErr:
			return fmt.Errorf("write returned early: %w", err)
		}
	}

	if err := send(&message.WriteMigrateTable{Table: table}); err != nil {
		return err
	}
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	defer bldr.Release()
	bldr.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3}, nil)
	if err := send(&message.WriteInsert{Record: bldr.NewRecordBatch()}); err != nil {
		return err
	}
	flush := message.NewWriteFlush()
	if err := send(flush); err != nil {
		return err
	}
	select {
	case err := <-flush.Result():
		if err != nil {
			return fmt.Errorf("flush failed: %w", err)
		}
	case err := <-writeErr:
		return fmt.Errorf("write returned before the flush completed: %w", err)
	}

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}
	if rows := TotalRows(records); rows != 3 {
		return fmt.Errorf("expected 3 rows after the flush, got %d", rows)
	}

	close(msgs)
	return <-writeErr
}
