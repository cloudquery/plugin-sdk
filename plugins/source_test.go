package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

func testSimpleResolver(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	res <- map[string]interface{}{
		"TestColumn": 3,
	}
	return nil
}

func testSimpleTable() *schema.Table {
	return &schema.Table{
		Name: "testSimpleTable",
		Resolver: testSimpleResolver,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
	}
}

func testRelationalTable() *schema.Table {
	return &schema.Table{
		Name: "testRelationalTableParent",
		Resolver: testSimpleResolver,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "testRelationalTableChild",
				Resolver: testSimpleResolver,
				Columns: []schema.Column{
					{
						Name: "test_column",
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}

func (*testExecutionClient) Name() string {
	return "testExecutionClient"
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

func TestSync(t *testing.T) {
	ctx := context.Background()
	tables := []*schema.Table{
		testSimpleTable(),
		testRelationalTable(),
	}
	plugin := NewSourcePlugin(
		"testSourcePlugin",
		"1.0.0",
		tables,
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
		return plugin.Sync(ctx,
			zerolog.New(zerolog.NewTestWriter(t)),
			spec,
			resources)
	})

	for _ = range resources {
		// if resource.Table.Name != "testTable" {
		// 	t.Fatalf("unexpected resource table name: %s", resource.Table.Name)
		// }
		// obj := resource.Get("test_column")
		// val, ok := obj.(int)
		// if !ok {
		// 	t.Fatalf("unexpected resource column value (expected int): %v", obj)
		// }

		// if val != 3 {
		// 	t.Fatalf("unexpected resource column value: %v", val)
		// }
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

