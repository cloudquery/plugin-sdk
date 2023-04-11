package destination

import (
	"github.com/cloudquery/plugin-sdk/schema"
)

func setupTables() schema.Table {
	topLevelTable := schema.Table{
		Name: "test_table",
		Columns: []schema.Column{
			schema.CqIDColumn,
			schema.CqParentIDColumn,
			{
				Name: "id",
				Type: schema.TypeUUID,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
	}
	nestedTable := schema.Table{
		Name: "test_relation_table",
		Columns: []schema.Column{
			schema.CqIDColumn,
			schema.CqParentIDColumn,
			{
				Name: "id",
				Type: schema.TypeUUID,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
		Parent: &topLevelTable,
	}
	topLevelTable.Relations = []*schema.Table{&nestedTable}
	return topLevelTable
}
