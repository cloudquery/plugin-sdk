package plugins

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

var _ schema.ClientMeta = &testExecutionClient{}

const testSourcePluginExampleConfig = `# specify all accounts you want to sync
accounts: []
`

type testExecutionClient struct {
	logger zerolog.Logger
}

type Account struct {
	Name    string   `json:"name,omitempty"`
	Regions []string `json:"regions"`
}

type TestConfig struct {
	Accounts []Account `json:"accounts"`
	Regions  []string  `json:"regions"`
}

func (TestConfig) Example() string {
	return ""
}

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

func newTestExecutionClient(context.Context, *SourcePlugin, specs.Source) (schema.ClientMeta, error) {
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
		WithSourceExampleConfig(testSourcePluginExampleConfig),
	)

	// test round trip: get example config -> sync with example config -> success
	exampleConfig, err := plugin.ExampleConfig()
	if err != nil {
		t.Fatal(err)
	}
	var spec specs.Source
	if err := yaml.Unmarshal([]byte(exampleConfig), &spec); err != nil {
		t.Fatal(err)
	}

	a := json.NewDecoder(strings.NewReader(exampleConfig))
	json.Strin

	d := yaml.NewDecoder(strings.NewReader(exampleConfig))
	d.KnownFields(true)
	if err := d.Decode(&spec); err != nil {
		t.Fatal(err)
	}

	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		_, err = plugin.Sync(ctx,
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
