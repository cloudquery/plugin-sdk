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

type Account struct {
	Name    string   `json:"name,omitempty"`
	Regions []string `json:"regions,omitempty"`
}

// type testSourceSpec struct {
// 	Accounts []Account `json:"accounts,omitempty"`
// 	Regions  []string  `json:"regions,omitempty"`
// }

// func newTestSourceSpec() interface{} {
// 	return &testSourceSpec{}
// }

var _ schema.ClientMeta = &testExecutionClient{}

func testTable() *schema.Table {
	return &schema.Table{
		Name: "testTable",
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
		WithSourceLogger(zerolog.New(zerolog.NewTestWriter(t))),
	)

	spec := specs.Source{
		Name:   "testSource",
		Tables: []string{"*"},
	}

	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		err := plugin.Sync(ctx,
			spec,
			resources)
		return err
	})

	for resource := range resources {
		if resource.Table.Name != "testTable" {
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

func TestSourcePlugin_interpolateAllResources(t *testing.T) {
	tests := []struct {
		name                string
		plugin              SourcePlugin
		configurationTables []string
		want                []string
		wantErr             bool
	}{
		{
			name:                "should return all tables when '*' is provided",
			plugin:              SourcePlugin{tables: []*schema.Table{{Name: "table 1"}, {Name: "table 2"}, {Name: "table 3"}}},
			configurationTables: []string{"*"},
			want:                []string{"table 1", "table 2", "table 3"},
			wantErr:             false,
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
			name:    "should return empty array when nil is provided",
			plugin:  SourcePlugin{tables: []*schema.Table{{Name: "table 1"}}},
			want:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.interpolateAllResources(tt.configurationTables)
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
