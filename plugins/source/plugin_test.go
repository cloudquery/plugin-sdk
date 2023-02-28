package source

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

var deterministicStableUUID = uuid.MustParse("5cdb0df90f0a5215adc56f15aa74a902")
var randomStableUUID = uuid.MustParse("00000000000040008000000000000000")

func testResolverSuccess(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
	res <- map[string]any{
		"TestColumn": 3,
	}
	return nil
}

func testResolverPanic(context.Context, schema.ClientMeta, *schema.Resource, chan<- any) error {
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

func (*testExecutionClient) ID() string {
	return "testExecutionClient"
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source, Options) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

type syncTestCase struct {
	table             *schema.Table
	stats             Metrics
	data              []schema.CQTypes
	deterministicCQID bool
}

var syncTestCases = []syncTestCase{
	{
		table: testTableSuccess(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
				"test_table_success": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data: []schema.CQTypes{
			{
				&schema.UUID{Bytes: randomStableUUID, Status: schema.Present},
				&schema.UUID{Status: schema.Null},
				&schema.Int8{Int: 3, Status: schema.Present},
			},
		},
	},
	{
		table: testTableResolverPanic(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
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
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
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
		table: testTableRelationSuccess(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
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
				&schema.UUID{Bytes: randomStableUUID, Status: schema.Present},
				&schema.UUID{Status: schema.Null},
				&schema.Int8{Int: 3, Status: schema.Present},
			},
			{
				&schema.UUID{Bytes: randomStableUUID, Status: schema.Present},
				&schema.UUID{Bytes: randomStableUUID, Status: schema.Present},
				&schema.Int8{Int: 3, Status: schema.Present},
			},
		},
	},
	{
		table: testTableSuccess(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
				"test_table_success": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data: []schema.CQTypes{
			{
				&schema.UUID{Bytes: deterministicStableUUID, Status: schema.Present},
				&schema.UUID{Status: schema.Null},
				&schema.Int8{Int: 3, Status: schema.Present},
			},
		},
		deterministicCQID: true,
	},
	{
		table: testTableColumnResolverPanic(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
				"test_table_column_resolver_panic": {
					"testExecutionClient": {
						Panics:    1,
						Resources: 1,
					},
				},
			},
		},
		data: []schema.CQTypes{
			{
				&schema.UUID{Bytes: deterministicStableUUID, Status: schema.Present},
				&schema.UUID{Status: schema.Null},
				&schema.Int8{Int: 3, Status: schema.Present},
				&schema.Int8{Status: schema.Undefined},
			},
		},
		deterministicCQID: true,
	},
	{
		table: testTableRelationSuccess(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
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
				&schema.UUID{Bytes: deterministicStableUUID, Status: schema.Present},
				&schema.UUID{Status: schema.Null},
				&schema.Int8{Int: 3, Status: schema.Present},
			},
			{
				&schema.UUID{Bytes: deterministicStableUUID, Status: schema.Present},
				&schema.UUID{Bytes: deterministicStableUUID, Status: schema.Present},
				&schema.Int8{Int: 3, Status: schema.Present},
			},
		},
		deterministicCQID: true,
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
	for _, scheduler := range specs.AllSchedulers {
		for _, tc := range syncTestCases {
			tc := tc
			tc.table = tc.table.Copy(nil)
			t.Run(tc.table.Name+"_"+scheduler.String(), func(t *testing.T) {
				testSyncTable(t, tc, scheduler, tc.deterministicCQID)
			})
		}
	}
}

func testSyncTable(t *testing.T, tc syncTestCase, scheduler specs.Scheduler, deterministicCQID bool) {
	ctx := context.Background()
	tables := []*schema.Table{
		tc.table,
	}

	plugin := NewPlugin(
		"testSourcePlugin",
		"1.0.0",
		tables,
		newTestExecutionClient,
	)
	plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(t)))
	spec := specs.Source{
		Name:              "testSource",
		Path:              "cloudquery/testSource",
		Tables:            []string{"*"},
		Version:           "v1.0.0",
		Destinations:      []string{"test"},
		Concurrency:       1, // choose a very low value to check that we don't run into deadlocks
		Scheduler:         scheduler,
		DeterministicCQID: deterministicCQID,
	}
	if err := plugin.Init(ctx, spec); err != nil {
		t.Fatal(err)
	}

	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		return plugin.Sync(ctx,
			resources)
	})

	var i int
	for resource := range resources {
		if tc.data == nil {
			t.Fatalf("Unexpected resource %v", resource)
		}
		if i >= len(tc.data) {
			t.Fatalf("expected %d resources. got %d", len(tc.data), i)
		}
		if !resource.GetValues().Equal(tc.data[i]) {
			t.Fatalf("expected at i=%d: %v. got %v", i, tc.data[i], resource.GetValues())
		}
		i++
	}
	if len(tc.data) != i {
		t.Fatalf("expected %d resources. got %d", len(tc.data), i)
	}

	stats := plugin.Metrics()
	if !tc.stats.Equal(stats) {
		t.Fatalf("unexpected stats: %v", cmp.Diff(tc.stats, stats))
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestIgnoredColumns(t *testing.T) {
	validateResources(t, schema.Resources{{
		Item: struct{ A *string }{},
		Table: &schema.Table{
			Columns: schema.ColumnList{
				{
					Name:          "a",
					Type:          schema.TypeString,
					IgnoreInTests: true,
				},
			},
		},
	}})
}

var testTable struct {
	PrimaryKey   string
	SecondaryKey string
	TertiaryKey  string
	Quaternary   string
}

func TestNewPluginPrimaryKeys(t *testing.T) {
	testTransforms := []struct {
		transformerOptions []transformers.StructTransformerOption
		resultKeys         []string
	}{
		{
			transformerOptions: []transformers.StructTransformerOption{transformers.WithPrimaryKeys("PrimaryKey")},
			resultKeys:         []string{"primary_key"},
		},
		{
			transformerOptions: []transformers.StructTransformerOption{},
			resultKeys:         []string{"_cq_id"},
		},
	}
	for _, tc := range testTransforms {
		tables := []*schema.Table{
			{
				Name: "test_table",
				Transform: transformers.TransformWithStruct(
					&testTable, tc.transformerOptions...,
				),
			},
		}

		plugin := NewPlugin("testSourcePlugin", "1.0.0", tables, newTestExecutionClient)
		assert.Equal(t, tc.resultKeys, plugin.tables[0].PrimaryKeys())
	}
}
