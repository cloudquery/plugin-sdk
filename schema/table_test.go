package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testTable = &Table{
	Name:    "test",
	Columns: []Column{},
	Relations: []*Table{
		{
			Name:    "test2",
			Columns: []Column{},
		},
	},
}

func TestTablesFlatten(t *testing.T) {
	tables := Tables{testTable}.FlattenTables()
	if len(tables) != 2 {
		t.Fatal("expected 2 tables")
	}
}

func TestTablesFilterDFS(t *testing.T) {
	tests := []struct {
		name                    string
		tables                  Tables
		configurationTables     []string
		configurationSkipTables []string
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
			configurationTables:     []string{"sub_sub_table"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTables, err := tt.tables.FilterDfs(tt.configurationTables, tt.configurationSkipTables)
			//nolint:gocritic
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
