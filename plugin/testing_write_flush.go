package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// flushTestTimeout bounds each wait in testFlush, so that a plugin that never completes a flush fails the test
// instead of hanging it.
const flushTestTimeout = time.Minute

// testFlush checks that the changes before each flush are readable once the flush completes, while Write is still
// running.
func (s *WriterTestSuite) testFlush(ctx context.Context) error {
	table := &schema.Table{
		Name: s.tableNameForTest("flush"),
		Columns: schema.ColumnList{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true, NotNull: true},
			schema.CqSourceNameColumn,
			schema.CqSyncTimeColumn,
		},
	}
	const sourceName = "source-test"
	syncTime := time.Now().UTC().Truncate(s.genDatOptions.TimePrecision).Truncate(time.Microsecond)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	msgs := make(chan message.WriteMessage)
	writeErr := make(chan error, 1)
	writeDone := make(chan struct{})
	go func() {
		defer close(writeDone)
		writeErr <- s.plugin.Write(ctx, msgs)
	}()
	msgsOpen := true
	defer func() {
		if msgsOpen {
			close(msgs)
		}
		// a failed step must not leave Write running into the next test
		cancel()
		select {
		case <-writeDone:
		case <-time.After(flushTestTimeout):
		}
	}()

	send := func(msg message.WriteMessage) error {
		select {
		case msgs <- msg:
			return nil
		case err := <-writeErr:
			return fmt.Errorf("write returned early: %w", err)
		case <-time.After(flushTestTimeout):
			return fmt.Errorf("write did not take a %T message within %s", msg, flushTestTimeout)
		}
	}
	flushAndCount := func(step string, wantRows int64) error {
		msg := &message.WriteFlush{}
		if err := send(msg); err != nil {
			return fmt.Errorf("%s: %w", step, err)
		}
		select {
		case err := <-msg.Result():
			if err != nil {
				return fmt.Errorf("%s: flush failed: %w", step, err)
			}
		case err := <-writeErr:
			return fmt.Errorf("%s: write returned before the flush completed: %w", step, err)
		case <-time.After(flushTestTimeout):
			return fmt.Errorf("%s: flush did not complete within %s", step, flushTestTimeout)
		}
		records, err := s.plugin.readAll(ctx, table)
		if err != nil {
			return fmt.Errorf("%s: failed to read: %w", step, err)
		}
		if rows := TotalRows(records); rows != wantRows {
			return fmt.Errorf("%s: expected %d rows after the flush, got %d", step, wantRows, rows)
		}
		return nil
	}
	insert := func(syncTime time.Time, ids ...int64) error {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
		defer bldr.Release()
		for _, id := range ids {
			bldr.Field(0).(*array.Int64Builder).Append(id)
			bldr.Field(1).(*array.StringBuilder).Append(sourceName)
			bldr.Field(2).(*array.TimestampBuilder).AppendTime(syncTime)
		}
		return send(&message.WriteInsert{Record: bldr.NewRecordBatch()})
	}

	if err := send(&message.WriteMigrateTable{Table: table}); err != nil {
		return err
	}
	if err := insert(syncTime, 1, 2, 3); err != nil {
		return err
	}
	if err := flushAndCount("first insert", 3); err != nil {
		return err
	}
	syncTime = syncTime.Add(time.Second)
	if err := insert(syncTime, 4, 5); err != nil {
		return err
	}
	if err := flushAndCount("second insert", 5); err != nil {
		return err
	}

	rows := int64(5)
	if !s.tests.SkipDeleteStale {
		if err := send(&message.WriteDeleteStale{TableName: table.Name, SourceName: sourceName, SyncTime: syncTime}); err != nil {
			return err
		}
		rows = 2
		if err := flushAndCount("delete stale", rows); err != nil {
			return err
		}
	}
	if !s.tests.SkipDeleteRecord {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, (&schema.Table{
			Name:    table.Name,
			Columns: schema.ColumnList{{Name: "id", Type: arrow.PrimitiveTypes.Int64}},
		}).ToArrowSchema())
		defer bldr.Release()
		bldr.Field(0).(*array.Int64Builder).Append(4)
		if err := send(&message.WriteDeleteRecord{DeleteRecord: message.DeleteRecord{
			TableName: table.Name,
			WhereClause: message.PredicateGroups{{
				GroupingType: "AND",
				Predicates:   message.Predicates{{Operator: "eq", Column: "id", Record: bldr.NewRecordBatch()}},
			}},
		}}); err != nil {
			return err
		}
		rows--
		if err := flushAndCount("delete record", rows); err != nil {
			return err
		}
	}

	close(msgs)
	msgsOpen = false
	select {
	case err := <-writeErr:
		return err
	case <-time.After(flushTestTimeout):
		return fmt.Errorf("write did not return within %s after its messages ended", flushTestTimeout)
	}
}
