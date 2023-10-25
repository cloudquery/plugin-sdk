package premium

import "github.com/cloudquery/plugin-sdk/v4/schema"

// ContainsPaidTables returns true if any of the tables are paid
func ContainsPaidTables(tables schema.Tables) bool {
	for _, t := range tables {
		if t.IsPaid {
			return true
		}
	}
	return false
}
