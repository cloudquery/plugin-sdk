package writers

import (
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// FailedTables names the tables of msgs in err, so that a flush that fails with it reports them. It returns err as is
// if it is nil, or if a message has no table name.
func FailedTables[T message.WriteMessage](err error, msgs []T) error {
	if err == nil {
		return nil
	}
	tables := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		name := tableName(msg)
		if name == "" {
			return err
		}
		tables = append(tables, name)
	}
	return &message.FailedTablesError{Tables: tables, Err: err}
}

func tableName(msg message.WriteMessage) string {
	switch m := msg.(type) {
	case *message.WriteInsert:
		name, _ := m.Record.Schema().Metadata().GetValue(schema.MetadataTableName)
		return name
	case *message.WriteMigrateTable:
		return m.Table.Name
	case *message.WriteDeleteStale:
		return m.TableName
	case *message.WriteDeleteRecord:
		return m.TableName
	default:
		return ""
	}
}
