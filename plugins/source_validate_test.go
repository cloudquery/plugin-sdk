package plugins

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/google/go-cmp/cmp"
)

func TestSourcePlugin_listAndValidateAllResources(t *testing.T) {
	tests := []struct {
		name                    string
		plugin                  SourcePlugin
		configurationTables     []string
		configurationSkipTables []string
		want                    []string
		wantErr                 bool
	}{
		{
			name:                "should return all tables when '*' is provided",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}}},
			configurationTables: []string{"*"},
			want:                []string{"table1", "table2", "table3"},
			wantErr:             false,
		},
		{
			name:                    "should return all tables when '*' is provided, excluding skipped tables",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"table1", "table3"},
			want:                    []string{"table2"},
			wantErr:                 false,
		},
		{
			name:                    "should return error when invalid table is skipped",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"table4"},
			wantErr:                 true,
		},
		{
			name:                    "should return only the exact set of tables specified",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables:     []string{"main_table"},
			configurationSkipTables: []string{},
			want:                    []string{"main_table"},
			wantErr:                 false,
		},
		{
			name:                    "should return only parent table if child table is skipped",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"sub_table"},
			want:                    []string{"main_table"},
			wantErr:                 false,
		},
		{
			name:                "should return specific tables when they are provided",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}}},
			configurationTables: []string{"table1"},
			want:                []string{"table1"},
			wantErr:             false,
		},
		{
			name:    "should return an error when nil is provided",
			plugin:  SourcePlugin{tables: []*schema.Table{{Name: "table1"}}},
			wantErr: true,
		},
		{
			name:                    "should return an error if glob-matching is attempted in tables",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"table*"},
			configurationSkipTables: []string{""},
			wantErr:                 true,
		},
		{
			name:                    "should return tables matching glob pattern",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"table*"},
			configurationSkipTables: []string{"table2"},
			want:                    []string{"table1"},
			wantErr:                 false,
		},
		{
			name:                    "should not return an error when included table is skipped",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"table2", "table1"},
			configurationSkipTables: []string{"table1"},
			want:                    []string{"table2"},
			wantErr:                 false,
		},
		{
			name:                    "should return an error if table is unmatched",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}}},
			configurationTables:     []string{"table2"},
			configurationSkipTables: []string{"table1"},
			wantErr:                 true,
		},
		{
			name:                "should return an error if child table is specified without its ancestors",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables: []string{"sub_table"},
			wantErr:             true,
		},
		{
			name:                "should return both tables if child and parent tables are specified",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables: []string{"main_table", "sub_table"},
			want:                []string{"main_table", "sub_table"},
			wantErr:             false,
		},
		{
			name:                    "should return table only once, even if it is matched by multiple rules",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"*", "table2", "table1", "table*"},
			configurationSkipTables: []string{"table2"},
			want:                    []string{"table1"},
			wantErr:                 false,
		},
		{
			name:                    "should match prefix globs",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"*2"},
			configurationSkipTables: []string{},
			want:                    []string{"table2"},
			wantErr:                 false,
		},
		{
			name:                    "should match suffix globs",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"table*"},
			configurationSkipTables: []string{},
			want:                    []string{"table1", "table2"},
			wantErr:                 false,
		},
		{
			name:                    "should match middle globs",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"t*b*1"},
			configurationSkipTables: []string{},
			want:                    []string{"table1"},
			wantErr:                 false,
		},
		{
			name:                    "should skip globs",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"t*1"},
			want:                    []string{"table2"},
			wantErr:                 false,
		},
		{
			name:                    "should skip multiple globs",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table1"}, {Name: "table2"}, {Name: "table3"}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"t*1", "t*2"},
			want:                    []string{"table3"},
			wantErr:                 false,
		},
		{
			name:                    "should glob match against child tables",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{},
			want:                    []string{"main_table", "sub_table"},
			wantErr:                 false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTables, err := tt.plugin.listAndValidateTables(tt.configurationTables, tt.configurationSkipTables)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePlugin.listAndValidateTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			gotNames := make([]string, len(gotTables))
			for i := range gotTables {
				gotNames[i] = gotTables[i].Name
			}
			if diff := cmp.Diff(gotNames, tt.want); diff != "" {
				t.Errorf("SourcePlugin.listAndValidateTables() diff (+got, -want): %v", diff)
			}
		})
	}
}
