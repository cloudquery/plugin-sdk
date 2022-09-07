package plugins

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"strings"
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

const wantSourceConfig = `
kind: source
spec:
  # Name of the plugin.
  name: "testSourcePlugin"

  # Version of the plugin to use.
  version: "1.0.0"

  # Registry to use (one of "github", "local" or "grpc").
  registry: "grpc"

  # Path to plugin. Required format depends on the registry.
  path: "testSourcePlugin"

  # List of tables to sync.
  tables: ["*"]

  ## Tables to skip during sync. Optional.
  # skip_tables: []

  # Names of destination plugins to sync to.
  destinations: ["*"]

  ## Approximate cap on number of requests to perform concurrently. Optional.
  # max_goroutines: 5

  # Plugin-specific configuration.
  spec:
    
    # This is an example config file for the test plugin.
    key: value
`

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
		WithSourceExampleConfig(`
# This is an example config file for the test plugin.
key: value`))

	// test round trip: get example config -> sync with example config -> success
	opts := SourceExampleConfigOptions{
		Registry: specs.RegistryGrpc,
	}
	exampleConfig, err := plugin.ExampleConfig(opts)
	if err != nil {
		t.Fatalf("unexpected error calling ExampleConfig: %v", err)
	}
	want := strings.TrimSpace(wantSourceConfig)
	if diff := cmp.Diff(exampleConfig, want); diff != "" {
		t.Fatalf("generated source config not as expected (-got, +want): %v", diff)
	}

	var spec specs.Spec
	if err := specs.SpecUnmarshalYamlStrict([]byte(exampleConfig), &spec); err != nil {
		t.Fatal(err)
	}

	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		err := plugin.Sync(ctx,
			*spec.Spec.(*specs.Source),
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
