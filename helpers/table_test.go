package helpers

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/stretchr/testify/assert"
)

func TestGetFlatTableList(t *testing.T) {
	var tables []*schema.Table = []*schema.Table{
		{
			Name: "table1",
			Relations: []*schema.Table{
				{
					Name: "table1.1",
				},
				{
					Name: "table1.2",
				},
			},
		},
		{
			Name: "table2",
			Relations: []*schema.Table{
				{
					Name: "table2.1",
				},
			},
		},
	}

	flatTables := GetFlatTableList(tables)
	flatTableNames := tableNames(flatTables)
	expectedTableNames := []string{"table1", "table1.1", "table1.2", "table2", "table2.1"}
	assert.ElementsMatch(t, flatTableNames, expectedTableNames)
}

func tableNames(tables []*schema.Table) []string {
	names := make([]string, 0)
	for _, table := range tables {
		names = append(names, table.Name)
	}
	return names
}
