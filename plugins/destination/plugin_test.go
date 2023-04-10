package destination

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/stretchr/testify/require"
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

func TestSetCQIDAsPrimaryKeysForTables(t *testing.T) {
	topLevelTable := setupTables()
	// Prior to executing setCQIDAsPrimaryKeysForTables only the id column should be a primary key
	require.False(t, topLevelTable.Columns[0].CreationOptions.PrimaryKey)
	require.False(t, topLevelTable.Columns[1].CreationOptions.PrimaryKey)
	require.True(t, topLevelTable.Columns[2].CreationOptions.PrimaryKey)

	sch := topLevelTable.ToArrowSchema()
	newSchemas := setCQIDAsPrimaryKeysForTables(schema.Schemas{sch})
	got := newSchemas[0]

	// After executing setCQIDAsPrimaryKeysForTables all cq_id columns should be primary keys
	require.True(t, schema.IsPk(got.Field(0)))
	require.False(t, schema.IsPk(got.Field(1)))
	require.False(t, schema.IsPk(got.Field(2)))
}

func TestSetDestinationManagedCqColumns(t *testing.T) {
	topLevelTable := setupTables()
	sc := topLevelTable.ToArrowSchema()

	require.False(t, sc.HasField(schema.CqSyncTimeColumn.Name))
	require.False(t, sc.HasField(schema.CqSourceNameColumn.Name))

	newSchemas := SetDestinationManagedCqColumns(schema.Schemas{sc})
	newSchema := newSchemas[0]

	require.True(t, newSchema.HasField(schema.CqIDColumn.Name))
	require.True(t, newSchema.HasField(schema.CqSyncTimeColumn.Name))
	require.True(t, newSchema.HasField(schema.CqSourceNameColumn.Name))
}
