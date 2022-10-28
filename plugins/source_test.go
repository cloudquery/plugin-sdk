package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/cqtypes"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

func testResolverSuccess(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- interface{}) error {
	res <- map[string]interface{}{
		"TestColumn": 3,
	}
	return nil
}

func testResolverPanic(context.Context, schema.ClientMeta, *schema.Resource, chan<- interface{}) error {
	panic("Resolver")
}

func testPreResourceResolverPanic(context.Context, schema.ClientMeta, *schema.Resource) error {
	panic("PreResourceResolver")
}

func testColumnResolverPanic(context.Context, schema.ClientMeta, *schema.Resource, schema.Column) error {
	panic("ColumnResolver")
}

func testTableSuccess() *schema.Table {
	return &schema.Table{
		Name:     "test_table_success",
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
		Name:     "test_table_resolver_panic",
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
		Name:                "test_table_pre_resource_resolver_panic",
		PreResourceResolver: testPreResourceResolverPanic,
		Resolver:            testResolverSuccess,
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
		Name:     "test_table_column_resolver_panic",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
			{
				Name:     "test_column1",
				Type:     schema.TypeInt,
				Resolver: testColumnResolverPanic,
			},
		},
	}
}

func testTableRelationSuccess() *schema.Table {
	return &schema.Table{
		Name:     "test_table_relation_success",
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
	data  []schema.CQTypes
}

var syncTestCases = []syncTestCase{
	{
		table: testTableSuccess(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"test_table_success": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data: []schema.CQTypes{
			{
				&cqtypes.UUID{Bytes: [16]byte{1}, Status: cqtypes.Present},
				nil,
				&cqtypes.Int8{Int: 3, Status: cqtypes.Present},
			},
		},
	},
	{
		table: testTableResolverPanic(),
		stats: SourceStats{
			TableClient: map[string]map[string]*TableClientStats{
				"test_table_resolver_panic": {
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
				"test_table_pre_resource_resolver_panic": {
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
				"test_table_column_resolver_panic": {
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
				"test_table_relation_success": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
				"test_table_success": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data: []schema.CQTypes{
			{
				&cqtypes.UUID{Bytes: [16]byte{1}, Status: cqtypes.Present},
				nil,
				&cqtypes.Int8{Int: 3, Status: cqtypes.Present},
			},
			{
				&cqtypes.UUID{Bytes: [16]byte{1}, Status: cqtypes.Present},
				&cqtypes.UUID{Bytes: [16]byte{1}, Status: cqtypes.Present},
				&cqtypes.Int8{Int: 3, Status: cqtypes.Present},
			},
		},
	},
}

type testRand struct{}

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
	plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(t)))
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
		// if !tc.data[i].Equal(resource.GetValues()) {
		// 	t.Fatalf("expected in item %d %v. got %v", i, tc.data[i], resource.GetValues())
		// }
		i++
	}
	if len(tc.data) != i {
		t.Fatalf("expected %d resources. got %d", len(tc.data), i)
	}
	stats := plugin.Stats()
	if !tc.stats.Equal(&stats) {
		t.Fatalf("unexpected stats: %v", cmp.Diff(tc.stats, stats))
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}
