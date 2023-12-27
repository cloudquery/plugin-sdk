package premium

import "github.com/cloudquery/plugin-sdk/v4/schema"

// ContainsPaidTables returns true if any of the tables are paid
func ContainsPaidTables(tables schema.Tables) bool {
	if tables == nil {
		return false
	}
	for _, t := range tables {
		if t.IsPaid || ContainsPaidTables(t.Relations) {
			return true
		}
	}
	return false
}

// MakeAllTablesPaid sets all tables to paid (including relations)
func MakeAllTablesPaid(tables schema.Tables) schema.Tables {
	for _, table := range tables {
		table.IsPaid = true
		MakeAllTablesPaid(table.Relations)
	}
	return tables
}
