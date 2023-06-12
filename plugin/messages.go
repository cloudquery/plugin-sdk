package plugin

import (
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type MessageCreateTable struct {
	Table        *schema.Table
	MigrateForce bool
}

type MessageInsert struct {
	Record arrow.Record
	Upsert bool
}

// MessageDeleteStale is a pretty specific message which requires the destination to be aware of a CLI use-case
// thus it might be deprecated in the future
// in favour of MessageDelete or MessageRawQuery
// The message indeciates that the destination needs to run something like "DELETE FROM table WHERE _cq_source_name=$1 and sync_time < $2"
type MessageDeleteStale struct {
	Table      *schema.Table
	SourceName string
	SyncTime   time.Time
}

type Message any

type Messages []Message

func (messages Messages) InsertItems() int64 {
	items := int64(0)
	for _, msg := range messages {
		switch m := msg.(type) {
		case *MessageInsert:
			items += m.Record.NumRows()
		}
	}
	return items
}
