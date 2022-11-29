package helpers

import "github.com/cloudquery/plugin-sdk/schema"

// GetFlatTablesList returns a flat list of *Table objects from a non-flat *Table slice.
// This doesn't modify the Table objects themselves! Neither should the caller!
func GetFlatTableList(tables []*schema.Table) []*schema.Table {
	flatTables := make([]*schema.Table, 0)

	for _, table := range tables {
		flatTables = append(flatTables, table)
		flatTables = append(flatTables, GetFlatTableList(table.Relations)...)
	}

	return flatTables
}
