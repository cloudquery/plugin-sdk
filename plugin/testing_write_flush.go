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

// testFlush checks that rows are readable once each flush completes, while Write is still running.
func (s *WriterTestSuite) testFlush(ctx context.Context) error {
	table := &schema.Table{
		Name: s.tableNameForTest("flush"),
		Columns: []schema.Column{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, NotNull: true},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	msgs := make(chan message.WriteMessage)
	writeErr := make(chan error, 1)
	go func() {
		writeErr <- s.plugin.Write(ctx, msgs)
	}()
	msgsOpen := true
	defer func() {
		if msgsOpen {
			close(msgs)
		}
	}()

	send := func(msg message.WriteMessage) error {
		select {
		case msgs <- msg:
			return nil
		case err := <-writeErr:
			return fmt.Errorf("write returned early: %w", err)
		}
	}
	flush := func() error {
		msg := message.NewWriteFlush()
		if err := send(msg); err != nil {
			return err
		}
		select {
		case err := <-msg.Result():
			if err != nil {
				return fmt.Errorf("flush failed: %w", err)
			}
			return nil
		case err := <-writeErr:
			return fmt.Errorf("write returned before the flush completed: %w", err)
		}
	}
	insertAndFlush := func(ids []int64, wantRows int64) error {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
		defer bldr.Release()
		bldr.Field(0).(*array.Int64Builder).AppendValues(ids, nil)
		if err := send(&message.WriteInsert{Record: bldr.NewRecordBatch()}); err != nil {
			return err
		}
		if err := flush(); err != nil {
			return err
		}
		records, err := s.plugin.readAll(ctx, table)
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}
		if rows := TotalRows(records); rows != wantRows {
			return fmt.Errorf("expected %d rows after the flush, got %d", wantRows, rows)
		}
		return nil
	}

	if err := send(&message.WriteMigrateTable{Table: table}); err != nil {
		return err
	}
	if err := insertAndFlush([]int64{1, 2, 3}, 3); err != nil {
		return err
	}
	if err := insertAndFlush([]int64{4, 5}, 5); err != nil {
		return err
	}

	close(msgs)
	msgsOpen = false
	return <-writeErr
}
