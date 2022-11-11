package plugins

import (
	"reflect"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
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
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}, {Name: "table 3"}}},
			configurationTables: []string{"*"},
			want:                []string{"table 1", "table 2", "table 3"},
			wantErr:             false,
		},
		{
			name:                    "should return all tables when '*' is provided, excluding skipped tables",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}, {Name: "table 3"}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"table 1", "table 3"},
			want:                    []string{"table 2"},
			wantErr:                 false,
		},
		{
			name:                    "should return error when invalid table is skipped",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}, {Name: "table 3"}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"table 4"},
			wantErr:                 true,
		},
		{
			name:                    "should return an error if child table is skipped",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables:     []string{"*"},
			configurationSkipTables: []string{"sub_table"},
			wantErr:                 true,
		},

		{
			name:                "should return specific tables when they are provided",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}, {Name: "table 3"}}},
			configurationTables: []string{"table 1"},
			want:                []string{"table 1"},
			wantErr:             false,
		},
		{
			name:                "should return error when '*' is provided with other tables",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}, {Name: "table 3"}}},
			configurationTables: []string{"table 1", "*"},
			wantErr:             true,
		},
		{
			name:    "should return an error when nil is provided",
			plugin:  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}}},
			wantErr: true,
		},
		{
			name:                    "should return an error if glob-matching is attempted in tables",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}}},
			configurationTables:     []string{"table*"},
			configurationSkipTables: []string{""},
			wantErr:                 true,
		},
		{
			name:                    "should return an error if glob-matching is attempted in skipped tables",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}}},
			configurationTables:     []string{"table 1"},
			configurationSkipTables: []string{"table *"},
			wantErr:                 true,
		},
		{
			name:                    "should return an error when included table is skipped",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}}},
			configurationTables:     []string{"table 2", "table 1"},
			configurationSkipTables: []string{"table 1"},
			wantErr:                 true,
		},
		{
			name:                    "should return an error if table is unmatched",
			plugin:                  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}}},
			configurationTables:     []string{"table 2"},
			configurationSkipTables: []string{"table 1"},
			wantErr:                 true,
		},
		{
			name:                "should return an error if child table is specified",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables: []string{"sub_table"},
			wantErr:             true,
		},
		{
			name:                "should return an error if both child and parent tables are specified",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "main_table", Relations: []*schema.Table{{Name: "sub_table", Parent: &schema.Table{Name: "main_table"}}}}}},
			configurationTables: []string{"main_table", "sub_table"},
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.listAndValidateTables(tt.configurationTables, tt.configurationSkipTables)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePlugin.listAndValidateTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcePlugin.listAndValidateTables() = %v, want %v", got, tt.want)
			}
		})
	}
}
