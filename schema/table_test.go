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

var testTable1 = &Table{
	Name: "test",
	Columns: []Column{
		{Name: "bool", Type: TypeBool},
	},
}

var testTable2 = &Table{
	Name: "test",
	Columns: []Column{
		{Name: "bool", Type: TypeBool},
		{Name: "bool1", Type: TypeBool},
	},
}

var testTable3 = &Table{
	Name: "test",
	Columns: []Column{
		{Name: "bool", Type: TypeString},
	},
}

var testTable4 = &Table{
	Name: "test",
	Columns: []Column{
		{Name: "bool", Type: TypeBool, CreationOptions: ColumnCreationOptions{PrimaryKey: true, NotNull: true}},
	},
}

func TestGetAddedColumns(t *testing.T) {
	columns := testTable1.GetAddedColumns(testTable1)
	if columns != nil {
		t.Fatalf("got %v want nil", columns)
	}

	columns = testTable2.GetAddedColumns(testTable1)
	if len(columns) != 1 {
		t.Fatalf("got %v want 1", columns)
	}
	if columns[0].Name != "bool1" {
		t.Fatalf("got %v want bool1", columns[0].Name)
	}
}

func TestGetChangedColumns(t *testing.T) {
	columns, _ := testTable1.GetChangedColumns(testTable1)
	if columns != nil {
		t.Fatalf("got %v want nil", columns)
	}

	columns, got := testTable3.GetChangedColumns(testTable2)
	if len(columns) != 1 {
		t.Fatalf("got %v want 1", columns)
	}
	if columns[0].Name != "bool" {
		t.Fatalf("got %v want bool", columns[0].Name)
	}
	if columns[0].Type != TypeString {
		t.Fatalf("got %v want TypeString", columns[0].Type)
	}
	if got[0].Type != TypeBool {
		t.Fatalf("got %v want TypeBool", got[0].Type)
	}

	columns, _ = testTable4.GetChangedColumns(testTable2)
	if len(columns) != 1 {
		t.Fatalf("got %v want 1", columns)
	}
	if columns[0].Name != "bool" {
		t.Fatalf("got %v want bool", columns[0].Name)
	}
}

func TestIsPkEqual(t *testing.T) {
	if !testTable1.IsPrimaryKeyEqual(testTable1) {
		t.Fatalf("got false want true")
	}
	if testTable4.IsPrimaryKeyEqual(testTable2) {
		t.Fatalf("got true want false")
	}
}
