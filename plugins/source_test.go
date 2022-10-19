package plugins

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type testExecutionClient struct {
	logger zerolog.Logger
}

var _ schema.ClientMeta = &testExecutionClient{}

func testTable() *schema.Table {
	return &schema.Table{
		Name: "test_table",
		Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
	}
}

func (c *testExecutionClient) Logger() *zerolog.Logger {
	return &c.logger
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

func TestSync(t *testing.T) {
	ctx := context.Background()
	plugin := NewSourcePlugin(
		"testSourcePlugin",
		"1.0.0",
		[]*schema.Table{testTable()},
		newTestExecutionClient,
	)

	spec := specs.Source{
		Name:         "testSource",
		Tables:       []string{"*"},
		Version:      "v1.0.0",
		Destinations: []string{"test"},
	}

	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		_, err := plugin.Sync(ctx,
			zerolog.New(zerolog.NewTestWriter(t)),
			spec,
			resources)
		return err
	})

	for resource := range resources {
		if resource.Table.Name != "test_table" {
			t.Fatalf("unexpected resource table name: %s", resource.Table.Name)
		}
		obj := resource.Get("test_column")
		val, ok := obj.(int)
		if !ok {
			t.Fatalf("unexpected resource column value (expected int): %v", obj)
		}

		if val != 3 {
			t.Fatalf("unexpected resource column value: %v", val)
		}
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

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
			name:                "should return an error if child table is without its parent",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table 1", Parent: &schema.Table{Name: "table 2"}}, {Name: "table 2"}}},
			configurationTables: []string{"table 1"},
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.listAndValidateTables(tt.configurationTables, tt.configurationSkipTables)
			if (err != nil) != tt.wantErr {
				t.Errorf("SourcePlugin.interpolateAllResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SourcePlugin.interpolateAllResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
