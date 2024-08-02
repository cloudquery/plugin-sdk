package schema

import (
	"testing"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

var testTable = &Table{
	Name:    "test",
	Columns: []Column{},
	Relations: []*Table{
		{
			Name:    "test2",
			Columns: []Column{},
			Parent:  &Table{Name: "test"},
		},
	},
}

func TestTablesFlatten(t *testing.T) {
	srcTables := Tables{testTable}
	tables := srcTables.FlattenTables()
	require.Equal(t, 1, len(srcTables)) // verify that the source Tables were left untouched
	require.Equal(t, 1, len(testTable.Relations))
	require.Equal(t, 2, len(tables))
	for _, table := range tables {
		require.Nil(t, table.Relations)
	}

	srcTables = Tables{testTable, testTable}
	tables = srcTables.FlattenTables()
	require.Equal(t, 2, len(srcTables)) // verify that the source Tables were left untouched
	require.Equal(t, 1, len(testTable.Relations))
	require.Equal(t, 2, len(tables))
	for _, table := range tables {
		require.Nil(t, table.Relations)
	}

	tables = tables.FlattenTables()
	if len(tables) != 2 {
		t.Fatal("expected 2 tables")
	}
}

func TestTablesUnflatten(t *testing.T) {
	srcTables := Tables{testTable}
	tables, err := srcTables.FlattenTables().UnflattenTables()
	require.NoError(t, err)
	require.Equal(t, 1, len(srcTables)) // verify that the source Tables were left untouched
	require.Equal(t, 1, len(tables))    // verify that the tables are equal to what we started with
	require.Equal(t, 1, len(tables[0].Relations))
}

func TestTablesFilterDFS(t *testing.T) {
	tests := []struct {
		name                    string
		tables                  Tables
		configurationTables     []string
		configurationSkipTables []string
		skipDependentTables     bool
		want                    []string
		err                     string
	}{
		{
			name:                "should return all tables when '*' is provided",
			tables:              []*Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}},
			configurationTables: []string{"*"},
			want:                []string{"table1", "table2", "table3"},
		},
		{
			name:                    "should return all tables when '*' is provided, excluding skipped tables",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"table1", "table3"},
			want:                    []string{"table2"},
		},
		{
			name: "should return the parent and all its descendants",
			tables: []*Table{
				{
					Name: "main_table",
					Relations: []*Table{
						{
							Name: "sub_table",
							Relations: []*Table{
								{
									Name: "sub_sub_table",
								},
							},
						},
					},
				},
			},
			configurationTables:     []string{"main_table"},
			configurationSkipTables: []string{},
			want:                    []string{"main_table", "sub_table", "sub_sub_table"},
		},
		{
			name: "should return the parent and all its descendants, but skip children and their descendants if they are skipped",
			tables: []*Table{{Name: "main_table", Relations: []*Table{
				{Name: "sub_table", Parent: &Table{Name: "main_table"}, Relations: []*Table{
					{Name: "sub_sub_table", Parent: &Table{Name: "sub_table"}},
				}}}},
			},
			configurationTables:     []string{"main_table"},
			configurationSkipTables: []string{"sub_table"},
			want:                    []string{"main_table"},
		},
		{
			name:                    "should return only parent table if child table is skipped",
			tables:                  []*Table{{Name: "main_table", Relations: []*Table{{Name: "sub_table", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"sub_table"},
			want:                    []string{"main_table"},
		},
		{
			name:                "should return specific tables when they are provided",
			tables:              []*Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}},
			configurationTables: []string{"table1"},
			want:                []string{"table1"},
		},
		{
			name:                    "should return tables matching glob pattern",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"table*"},
			configurationSkipTables: []string{"table2"},
			want:                    []string{"table1"},
		},
		{
			name:                    "should not return an error when included table is skipped",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"table2", "table1"},
			configurationSkipTables: []string{"table1"},
			want:                    []string{"table2"},
		},
		{
			name:                "should return both tables if child and parent tables are specified",
			tables:              []*Table{{Name: "main_table", Relations: []*Table{{Name: "sub_table", Parent: &Table{Name: "main_table"}}}}},
			configurationTables: []string{"main_table", "sub_table"},
			want:                []string{"main_table", "sub_table"},
		},
		{
			name:                    "should return table only once, even if it is matched by multiple rules",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"*", "table2", "table1", "table*"},
			configurationSkipTables: []string{"table2"},
			want:                    []string{"table1"},
		},
		{
			name:                    "should match prefix globs",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"*2"},
			configurationSkipTables: []string{},
			want:                    []string{"table2"},
		},
		{
			name:                    "should match suffix globs",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"table*"},
			configurationSkipTables: []string{},
			want:                    []string{"table1", "table2"},
		},
		{
			name:                    "should match middle globs",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"t*b*1"},
			configurationSkipTables: []string{},
			want:                    []string{"table1"},
		},
		{
			name:                    "should skip globs",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"t*1"},
			want:                    []string{"table2"},
		},
		{
			name:                    "should skip multiple globs",
			tables:                  []*Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"t*1", "t*2"},
			want:                    []string{"table3"},
		},
		{
			name:                    "should glob match against child tables",
			tables:                  []*Table{{Name: "main_table", Relations: []*Table{{Name: "sub_table", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{},
			want:                    []string{"main_table", "sub_table"},
		},
		{
			name:                    "should skip parent table",
			tables:                  []*Table{{Name: "main_table", Relations: []*Table{{Name: "sub_table", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"main_table"},
			want:                    []string{},
		},
		{
			name:                    "should skip parent table",
			tables:                  []*Table{{Name: "main_table", Relations: []*Table{{Name: "sub_table", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"main_table1"},
			want:                    []string{},
			err:                     "skip_tables include a pattern main_table1 with no matches",
		},
		{
			name:                    "should skip parent table",
			tables:                  []*Table{{Name: "main_table", Relations: []*Table{{Name: "sub_table", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"main_table1"},
			configurationSkipTables: []string{},
			want:                    []string{},
			err:                     "tables include a pattern main_table1 with no matches",
		},
		{
			name: "skip child table but return siblings",
			tables: []*Table{
				{Name: "main_table", Relations: []*Table{
					{Name: "sub_table_1", Parent: &Table{Name: "main_table"}},
					{Name: "sub_table_2", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"main_table"},
			configurationSkipTables: []string{"sub_table_2"},
			want:                    []string{"main_table", "sub_table_1"},
		},
		{
			name: "when specifying a single child table, return only the parent and the specified child",
			tables: []*Table{
				{Name: "main_table", Relations: []*Table{
					{Name: "sub_table_1", Parent: &Table{Name: "main_table"}},
					{Name: "sub_table_2", Parent: &Table{Name: "main_table"}}}}},
			configurationTables: []string{"sub_table_1"},
			want:                []string{"main_table", "sub_table_1"},
		},
		{
			name: "when specifying a leaf table, return only the parents and the leaf",
			tables: []*Table{
				{Name: "0", Relations: []*Table{
					{Name: "0_1", Parent: &Table{Name: "0"}, Relations: []*Table{
						{Name: "0_1_1", Parent: &Table{Name: "0_1"}, Relations: []*Table{
							{Name: "0_1_1_1", Parent: &Table{Name: "0_1_1"}},
							{Name: "0_1_1_2", Parent: &Table{Name: "0_1_1"}, Relations: []*Table{
								{Name: "0_1_1_2_1", Parent: &Table{Name: "0_1_1_2"}},
								{Name: "0_1_1_2_2", Parent: &Table{Name: "0_1_1_2"}},
							}},
							{Name: "0_1_1_3", Parent: &Table{Name: "0_1_1"}},
						}},
						{Name: "0_1_2", Parent: &Table{Name: "0_1"}, Relations: []*Table{
							{Name: "0_1_2_1", Parent: &Table{Name: "0_1_2"}},
							{Name: "0_1_2_2", Parent: &Table{Name: "0_1_2"}},
							{Name: "0_1_2_3", Parent: &Table{Name: "0_1_2"}},
						}},
						{Name: "0_1_3", Parent: &Table{Name: "0_1"}},
					}},
					{Name: "0_2", Parent: &Table{Name: "0"}, Relations: []*Table{
						{Name: "0_2_1", Parent: &Table{Name: "0_2"}},
						{Name: "0_2_2", Parent: &Table{Name: "0_2"}},
					}},
					{Name: "0_3", Parent: &Table{Name: "0"}},
				}},
				{Name: "1", Relations: []*Table{
					{Name: "1_1", Parent: &Table{Name: "1"}, Relations: []*Table{
						{Name: "1_1_1", Parent: &Table{Name: "1_1"}, Relations: []*Table{
							{Name: "1_1_1_1", Parent: &Table{Name: "1_1_1"}},
							{Name: "1_1_1_2", Parent: &Table{Name: "1_1_1"}, Relations: []*Table{
								{Name: "1_1_1_2_1", Parent: &Table{Name: "1_1_1_2"}},
								{Name: "1_1_1_2_2", Parent: &Table{Name: "1_1_1_2"}},
							}},
							{Name: "1_1_1_3", Parent: &Table{Name: "1_1_1"}},
						}},
						{Name: "1_1_2", Parent: &Table{Name: "1_1"}, Relations: []*Table{
							{Name: "1_1_2_1", Parent: &Table{Name: "1_1_2"}},
							{Name: "1_1_2_2", Parent: &Table{Name: "1_1_2"}},
							{Name: "1_1_2_3", Parent: &Table{Name: "1_1_2"}},
						}},
						{Name: "1_1_3", Parent: &Table{Name: "1_1"}},
					}},
					{Name: "1_2", Parent: &Table{Name: "1"}, Relations: []*Table{
						{Name: "1_2_1", Parent: &Table{Name: "1_2"}},
						{Name: "1_2_2", Parent: &Table{Name: "1_2"}},
					}},
					{Name: "1_3", Parent: &Table{Name: "1"}},
				}},
			},
			configurationTables: []string{"0_1_1_2_2", "0_1_2_3", "1_1_2_3"},
			want:                []string{"0", "0_1", "0_1_1", "0_1_1_2", "0_1_1_2_2", "0_1_2", "0_1_2_3", "1", "1_1", "1_1_2", "1_1_2_3"},
		},
		{
			name: "when specifying a descendant table, return the parents, the specified descendant and all its descendant if skip_dependent_tables is false",
			tables: []*Table{
				{Name: "0", Relations: []*Table{
					{Name: "0_1", Parent: &Table{Name: "0"}, Relations: []*Table{
						{Name: "0_1_1", Parent: &Table{Name: "0_1"}, Relations: []*Table{
							{Name: "0_1_1_1", Parent: &Table{Name: "0_1_1"}},
							{Name: "0_1_1_2", Parent: &Table{Name: "0_1_1"}, Relations: []*Table{
								{Name: "0_1_1_2_1", Parent: &Table{Name: "0_1_1_2"}},
								{Name: "0_1_1_2_2", Parent: &Table{Name: "0_1_1_2"}},
							}},
							{Name: "0_1_1_3", Parent: &Table{Name: "0_1_1"}},
						}},
						{Name: "0_1_2", Parent: &Table{Name: "0_1"}, Relations: []*Table{
							{Name: "0_1_2_1", Parent: &Table{Name: "0_1_2"}},
							{Name: "0_1_2_2", Parent: &Table{Name: "0_1_2"}},
							{Name: "0_1_2_3", Parent: &Table{Name: "0_1_2"}},
						}},
						{Name: "0_1_3", Parent: &Table{Name: "0_1"}},
					}},
					{Name: "0_2", Parent: &Table{Name: "0"}, Relations: []*Table{
						{Name: "0_2_1", Parent: &Table{Name: "0_2"}},
						{Name: "0_2_2", Parent: &Table{Name: "0_2"}},
					}},
					{Name: "0_3", Parent: &Table{Name: "0"}},
				}},
				{Name: "1", Relations: []*Table{
					{Name: "1_1", Parent: &Table{Name: "1"}, Relations: []*Table{
						{Name: "1_1_1", Parent: &Table{Name: "1_1"}, Relations: []*Table{
							{Name: "1_1_1_1", Parent: &Table{Name: "1_1_1"}},
							{Name: "1_1_1_2", Parent: &Table{Name: "1_1_1"}, Relations: []*Table{
								{Name: "1_1_1_2_1", Parent: &Table{Name: "1_1_1_2"}},
								{Name: "1_1_1_2_2", Parent: &Table{Name: "1_1_1_2"}},
							}},
							{Name: "1_1_1_3", Parent: &Table{Name: "1_1_1"}},
						}},
						{Name: "1_1_2", Parent: &Table{Name: "1_1"}, Relations: []*Table{
							{Name: "1_1_2_1", Parent: &Table{Name: "1_1_2"}},
							{Name: "1_1_2_2", Parent: &Table{Name: "1_1_2"}},
							{Name: "1_1_2_3", Parent: &Table{Name: "1_1_2"}},
						}},
						{Name: "1_1_3", Parent: &Table{Name: "1_1"}},
					}},
					{Name: "1_2", Parent: &Table{Name: "1"}, Relations: []*Table{
						{Name: "1_2_1", Parent: &Table{Name: "1_2"}},
						{Name: "1_2_2", Parent: &Table{Name: "1_2"}},
					}},
					{Name: "1_3", Parent: &Table{Name: "1"}},
				}},
			},
			configurationTables: []string{"1_1_1_2"},
			want:                []string{"1", "1_1", "1_1_1", "1_1_1_2", "1_1_1_2_1", "1_1_1_2_2"},
		},
		{
			name: "when specifying a descendant table, return the parents and only the specified descendant if skip_dependent_tables is true",
			tables: []*Table{
				{Name: "0", Relations: []*Table{
					{Name: "0_1", Parent: &Table{Name: "0"}, Relations: []*Table{
						{Name: "0_1_1", Parent: &Table{Name: "0_1"}, Relations: []*Table{
							{Name: "0_1_1_1", Parent: &Table{Name: "0_1_1"}},
							{Name: "0_1_1_2", Parent: &Table{Name: "0_1_1"}, Relations: []*Table{
								{Name: "0_1_1_2_1", Parent: &Table{Name: "0_1_1_2"}},
								{Name: "0_1_1_2_2", Parent: &Table{Name: "0_1_1_2"}},
							}},
							{Name: "0_1_1_3", Parent: &Table{Name: "0_1_1"}},
						}},
						{Name: "0_1_2", Parent: &Table{Name: "0_1"}, Relations: []*Table{
							{Name: "0_1_2_1", Parent: &Table{Name: "0_1_2"}},
							{Name: "0_1_2_2", Parent: &Table{Name: "0_1_2"}},
							{Name: "0_1_2_3", Parent: &Table{Name: "0_1_2"}},
						}},
						{Name: "0_1_3", Parent: &Table{Name: "0_1"}},
					}},
					{Name: "0_2", Parent: &Table{Name: "0"}, Relations: []*Table{
						{Name: "0_2_1", Parent: &Table{Name: "0_2"}},
						{Name: "0_2_2", Parent: &Table{Name: "0_2"}},
					}},
					{Name: "0_3", Parent: &Table{Name: "0"}},
				}},
				{Name: "1", Relations: []*Table{
					{Name: "1_1", Parent: &Table{Name: "1"}, Relations: []*Table{
						{Name: "1_1_1", Parent: &Table{Name: "1_1"}, Relations: []*Table{
							{Name: "1_1_1_1", Parent: &Table{Name: "1_1_1"}},
							{Name: "1_1_1_2", Parent: &Table{Name: "1_1_1"}, Relations: []*Table{
								{Name: "1_1_1_2_1", Parent: &Table{Name: "1_1_1_2"}},
								{Name: "1_1_1_2_2", Parent: &Table{Name: "1_1_1_2"}},
							}},
							{Name: "1_1_1_3", Parent: &Table{Name: "1_1_1"}},
						}},
						{Name: "1_1_2", Parent: &Table{Name: "1_1"}, Relations: []*Table{
							{Name: "1_1_2_1", Parent: &Table{Name: "1_1_2"}},
							{Name: "1_1_2_2", Parent: &Table{Name: "1_1_2"}},
							{Name: "1_1_2_3", Parent: &Table{Name: "1_1_2"}},
						}},
						{Name: "1_1_3", Parent: &Table{Name: "1_1"}},
					}},
					{Name: "1_2", Parent: &Table{Name: "1"}, Relations: []*Table{
						{Name: "1_2_1", Parent: &Table{Name: "1_2"}},
						{Name: "1_2_2", Parent: &Table{Name: "1_2"}},
					}},
					{Name: "1_3", Parent: &Table{Name: "1"}},
				}},
			},
			configurationTables: []string{"1_1_1_2"},
			skipDependentTables: true,
			want:                []string{"1", "1_1", "1_1_1", "1_1_1_2"},
		},
		{
			name: "skip child tables if skip_dependent_tables is true",
			tables: []*Table{
				{Name: "main_table", Relations: []*Table{
					{Name: "sub_table_1", Parent: &Table{Name: "main_table"}},
					{Name: "sub_table_2", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"main_table"},
			configurationSkipTables: []string{},
			skipDependentTables:     true,
			want:                    []string{"main_table"},
		},
		{
			name: "skip child tables if skip_dependent_tables is true, but not if explicitly included",
			tables: []*Table{
				{Name: "main_table_1", Relations: []*Table{
					{Name: "sub_table_1"},
				}},
				{Name: "main_table_2", Relations: []*Table{
					{Name: "sub_table_2", Parent: &Table{Name: "main_table"}},
					{Name: "sub_table_3", Parent: &Table{Name: "main_table"}}}}},
			configurationTables:     []string{"main_table_1", "sub_table_2"},
			configurationSkipTables: []string{},
			skipDependentTables:     true,
			want:                    []string{"main_table_1", "main_table_2", "sub_table_2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTables, err := tt.tables.FilterDfs(tt.configurationTables, tt.configurationSkipTables, tt.skipDependentTables)
			// nolint:gocritic
			if err != nil && tt.err == "" {
				t.Errorf("got error %v, want nil", err)
			} else if err != nil && tt.err != "" && err.Error() != tt.err {
				t.Errorf("got error %v, want %v", err, tt.err)
			} else if err == nil && tt.err != "" {
				t.Errorf("got nil, want error %v", tt.err)
			}
			gotTables = gotTables.FlattenTables()
			gotNames := make([]string, len(gotTables))
			for i := range gotTables {
				gotNames[i] = gotTables[i].Name
			}
			if diff := cmp.Diff(gotNames, tt.want); diff != "" {
				t.Errorf("diff (+got, -want): %v", diff)
			}
		})
	}
}

type testTableGetChangeTestCase struct {
	name            string
	target          *Table
	source          *Table
	expectedChanges []TableColumnChange
}

var testTableGetChangeTestCases = []testTableGetChangeTestCase{
	{
		name: "no changes",
		target: &Table{
			Name: "test",
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
		source: &Table{
			Name: "test",
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
		expectedChanges: nil,
	},
	{
		name: "add column",
		target: &Table{
			Name: "test",
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
				{Name: "bool1", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
		source: &Table{
			Name: "test",
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
		expectedChanges: []TableColumnChange{
			{
				Type:       TableColumnChangeTypeAdd,
				ColumnName: "bool1",
				Current:    Column{Name: "bool1", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
	},
	{
		name: "remove column",
		target: &Table{
			Name: "test",
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
		source: &Table{
			Name: "test",
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
				{Name: "bool1", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
		expectedChanges: []TableColumnChange{
			{
				Type:       TableColumnChangeTypeRemove,
				ColumnName: "bool1",
				Previous:   Column{Name: "bool1", Type: arrow.FixedWidthTypes.Boolean},
			},
		},
	},

	{
		name: "move to cq_id as primary key",
		target: &Table{
			Name: "test",
			Columns: []Column{
				{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
					PrimaryKey:  true,
				},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: false},
			},
		},
		source: &Table{
			Name: "test",
			Columns: []Column{
				{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: true},
			},
		},
		expectedChanges: []TableColumnChange{
			{
				Type: TableColumnChangeTypeMoveToCQOnly,
			},
			{
				Type:       TableColumnChangeTypeUpdate,
				ColumnName: "_cq_id",
				Current: Column{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
					PrimaryKey:  true,
				},
				Previous: Column{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
			},
			{
				Type:       TableColumnChangeTypeUpdate,
				ColumnName: "bool",
				Current:    Column{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: false},
				Previous:   Column{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: true},
			},
		},
	},

	{
		name: "move to cq_id as primary key and drop unique constraint",
		target: &Table{
			Name: "test",
			Columns: []Column{
				{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					// Unique:      true,
					PrimaryKey: true,
				},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: false},
			},
		},
		source: &Table{
			Name: "test",
			Columns: []Column{
				{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: true},
			},
		},
		expectedChanges: []TableColumnChange{
			{
				Type: TableColumnChangeTypeMoveToCQOnly,
			},
			{
				Type:       TableColumnChangeTypeUpdate,
				ColumnName: "_cq_id",
				Current: Column{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      false,
					PrimaryKey:  true,
				},
				Previous: Column{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
			},
			{
				Type:       TableColumnChangeTypeRemoveUniqueConstraint,
				ColumnName: "_cq_id",
				Previous: Column{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
			},
			{
				Type:       TableColumnChangeTypeUpdate,
				ColumnName: "bool",
				Current:    Column{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: false},
				Previous:   Column{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: true},
			},
		},
	},

	{
		name: "drop unique constraint",
		target: &Table{
			Name: "test",
			Columns: []Column{
				{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
				},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: true},
			},
		},
		source: &Table{
			Name: "test",
			Columns: []Column{
				{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, PrimaryKey: true},
			},
		},
		expectedChanges: []TableColumnChange{
			{
				Type:       TableColumnChangeTypeRemoveUniqueConstraint,
				ColumnName: "_cq_id",
				Previous: Column{
					Name:        "_cq_id",
					Type:        types.ExtensionTypes.UUID,
					Description: "Internal CQ ID of the row",
					NotNull:     true,
					Unique:      true,
				},
			},
		},
	},
}

func TestTableGetChanges(t *testing.T) {
	for _, tc := range testTableGetChangeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			changes := tc.target.GetChanges(tc.source)
			if diff := cmp.Diff(changes, tc.expectedChanges); diff != "" {
				t.Errorf("diff (+got, -want): %v", diff)
			}
		})
	}
}

func TestTablesToAndFromArrow(t *testing.T) {
	// The attributes in this table should all be preserved when converting to and from Arrow.
	tablesToTest := Tables{
		// Test empty table
		&Table{
			Columns: []Column{},
		},
		// Test table with attributes
		&Table{
			Name:        "test_table",
			Description: "Test table description",
			Title:       "Test Table",
			Parent: &Table{
				Name: "parent_table",
			},
			IsIncremental: true,
			Columns: []Column{
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
				{Name: "int", Type: arrow.PrimitiveTypes.Int64},
				{Name: "float", Type: arrow.PrimitiveTypes.Float64},
				{Name: "string", Type: arrow.BinaryTypes.String},
				{Name: "json", Type: types.ExtensionTypes.JSON},
				{Name: "unique", Type: arrow.BinaryTypes.String, Unique: true},
				{Name: "primary_key", Type: arrow.BinaryTypes.String, PrimaryKey: true},
				{Name: "not_null", Type: arrow.BinaryTypes.String, NotNull: true},
				{Name: "incremental_key", Type: arrow.BinaryTypes.String, IncrementalKey: true},
				{Name: "multiple_attributes", Type: arrow.BinaryTypes.String, PrimaryKey: true, IncrementalKey: true, NotNull: true, Unique: true},
			},
			PermissionsNeeded: []string{"storage.buckets.list", "compute.acceleratorTypes.list", "test,test"},
		},
	}

	for _, table := range tablesToTest {
		arrowSchema := table.ToArrowSchema()
		tableFromArrow, err := NewTableFromArrowSchema(arrowSchema)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(table, tableFromArrow); diff != "" {
			t.Errorf("diff (+got, -want): %v", diff)
		}
	}
}

func TestValidateDuplicateTables(t *testing.T) {
	tests := []struct {
		name   string
		tables Tables
		err    string
	}{
		{
			name:   "should return error when duplicate tables are found",
			tables: Tables{{Name: "table1"}, {Name: "table1"}},
			err:    "duplicate table table1",
		},
		{
			name:   "should return error when duplicate relational tables are found",
			tables: Tables{{Name: "table1", Relations: []*Table{{Name: "table2"}, {Name: "table2"}}}},
			err:    "duplicate table table2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.tables.ValidateDuplicateTables()
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
