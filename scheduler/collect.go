package scheduler

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *Scheduler) collectSend(resources <-chan *schema.Resource, res chan<- message.SyncMessage) {
	data := make(map[string]*tableBuilder)
	for resource := range resources {
		tb, ok := data[resource.Table.Name]
		if !ok {
			tb = &tableBuilder{
				RecordBuilder: array.NewRecordBuilder(memory.DefaultAllocator, resource.Table.ToArrowSchema()),
			}
			data[resource.Table.Name] = tb
		}
		tb.Append(resource.GetValues())
		if tb.rows == s.rowsPerRecord {
			res <- &message.SyncInsert{Record: tb.NewRecord()}
		}
	}
	for _, tb := range data {
		if tb.rows > 0 {
			res <- &message.SyncInsert{Record: tb.NewRecord()}
		}
	}
}

type tableBuilder struct {
	*array.RecordBuilder
	rows int
}

func (tb *tableBuilder) NewRecord() arrow.Record {
	tb.rows = 0
	return tb.RecordBuilder.NewRecord()
}

func (tb *tableBuilder) Append(vector scalar.Vector) {
	tb.rows++
	scalar.AppendToRecordBuilder(tb.RecordBuilder, vector)
}
