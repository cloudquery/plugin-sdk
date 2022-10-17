package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

func testResolverSuccess(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	res <- map[string]interface{}{
		"TestColumn": 3,
	}
	return nil
}

func testResolverPanic(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	panic("Resolver")
}

func testPreResourceResolverPanic(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	panic("PreResourceResolver")
}

func testColumnResolverPanic(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	panic("ColumnResolver")
}

func testTableSuccess() *schema.Table {
	return &schema.Table{
		Name: "testTableSuccess",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
	}
}

func testTableResolverPanic() *schema.Table {
	return &schema.Table{
		Name: "testTableResolverPanic",
		Resolver: testResolverPanic,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
	}
}

func testTablePreResourceResolverPanic() *schema.Table {
	return &schema.Table{
		Name: "testTablePreResourceResolverPanic",
		PreResourceResolver: testPreResourceResolverPanic,
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
	}
}

func testTableColumnResolverPanic() *schema.Table {
	return &schema.Table{
		Name: "testTableColumnResolverPanic",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
			{
				Name: "test_column1",
				Type: schema.TypeInt,
				Resolver: testColumnResolverPanic,
			},
		},
	}
}


func testTableRelationSuccess() *schema.Table {
	return &schema.Table{
		Name: "testTableRelationSuccess",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
		Relations: []*schema.Table{
			testTableSuccess(),
		},
	}
}

func (*testExecutionClient) Name() string {
	return "testExecutionClient"
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

type syncTestCase struct {
	table *schema.Table
	stats SourceStats
	data [][]interface{}
}

var testUUID = uuid.MustParse("00000000-0000-4000-8000-000000000000")

var syncTestCases = []syncTestCase{
	{
		table: testTableSuccess(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"testTableSuccess": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data: [][]interface{}{
			{
				testUUID,
				nil,
				3,
			},
		},
	},
	{
		table: testTableResolverPanic(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"testTableResolverPanic": {
					"testExecutionClient": {
						Panics: 1,
					},
				},
			},
		},
		data: nil,
	},
	{
		table: testTablePreResourceResolverPanic(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"testTablePreResourceResolverPanic": {
					"testExecutionClient": {
						Panics: 1,
					},
				},
			},
		},
		data: nil,
	},
	{
		table: testTableColumnResolverPanic(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"testTableColumnResolverPanic": {
					"testExecutionClient": {
						Panics: 1,
					},
				},
			},
		},
		data: nil,
	},
	{
		table: testTableRelationSuccess(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"testTableRelationSuccess": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
				"testTableSuccess": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data: [][]interface{}{
			{
				testUUID,
				nil,
				3,
			},
			{
				testUUID,
				testUUID,
				3,
			},
		},
	},
}


type testRand struct {}
func (testRand) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(0)
	}
	return len(p), nil
}

func TestSync(t *testing.T) {
	uuid.SetRand(testRand{})
	for _, tc := range syncTestCases {
		t.Run(tc.table.Name, func(t *testing.T) {
			testSyncTable(t, tc)
		})
	}
	
}


func testSyncTable(t *testing.T, tc syncTestCase) {
	ctx := context.Background()
	tables := []*schema.Table{
		tc.table,
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
	
	var i int
	for resource := range resources {
		if tc.data == nil {
			t.Fatalf("Unexpected resource %v", resource)
		}
		if i > len(tc.data) {
			t.Fatalf("expected %d resources. got %d", len(tc.data), i)
		}
		if diff := cmp.Diff(tc.data[i], resource.GetValues()); diff != "" {
			t.Fatalf("unexpected resource %s", diff)
		}
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
		i++
	}
	stats := plugin.Stats()
	if !tc.stats.Equal(&stats) {
		t.Fatalf("unexpected stats: %v", cmp.Diff(tc.stats, stats))
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

}