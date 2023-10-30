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

// MakeAllTablesPaid sets all tables to paid
func MakeAllTablesPaid(tables schema.Tables) schema.Tables {
	for _, table := range tables {
		MakeTablePaid(table)
	}
	return tables
}

// MakeTablePaid sets the table to paid
func MakeTablePaid(table *schema.Table) *schema.Table {
	table.IsPaid = true
	return table
}
