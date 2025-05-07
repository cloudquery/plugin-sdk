package plugin

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func ValidateNoEmptyColumns(t *testing.T, tables schema.Tables, messages message.SyncMessages) {
	for _, table := range tables.FlattenTables() {
		records := messages.GetInserts().GetRecordsForTable(table)
		emptyColumns := schema.FindEmptyColumns(table, records)
		if len(emptyColumns) > 0 {
			t.Fatalf("found empty column(s): %v in %s", emptyColumns, table.Name)
		}
		nonMatchingColumns, nonMatchingJSONColumns := schema.FindNotMatchingSensitiveColumns(table)
		if len(nonMatchingColumns) > 0 {
			t.Fatalf("found non-matching sensitive column(s): %v in %s", nonMatchingColumns, table.Name)
		}
		if len(nonMatchingJSONColumns) > 0 {
			t.Fatalf("found non-matching sensitive JSON column(s): %v in %s", nonMatchingJSONColumns, table.Name)
		}
	}
}
