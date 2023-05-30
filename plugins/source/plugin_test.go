package source

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v2/types"
	"github.com/cloudquery/plugin-sdk/v3/scalar"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/transformers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

var deterministicStableUUID = uuid.MustParse("0f922ba18f665741b093b9f560ee486f")
var randomStableUUID = uuid.MustParse("00000000000040008000000000000000")

var testSyncTime = time.Now()

var testColumns = []schema.Column{
	{
		Name: "int_column",
		Type: arrow.PrimitiveTypes.Int64,
	},
	{
		Name: "json_column",
		Type: types.ExtensionTypes.JSON,
	},
	{
		Name: "string_column",
		Type: arrow.BinaryTypes.String,
	},
	{
		Name: "float_column",
		Type: arrow.PrimitiveTypes.Float64,
	},
	{
		Name: "bool_column",
		Type: arrow.FixedWidthTypes.Boolean,
	},
	{
		Name: "time_column",
		Type: arrow.FixedWidthTypes.Timestamp_ms,
	},
	{
		Name: "uuid_column",
		Type: types.ExtensionTypes.UUID,
	},
	{
		Name: "array_column",
		Type: arrow.ListOf(arrow.PrimitiveTypes.Int64),
	},
	// TODO: Add support for map and struct types to scalar package, then test here
	//{
	//	Name: "map_column",
	//	Type: arrow.MapOf(arrow.BinaryTypes.String, arrow.PrimitiveTypes.Int64),
	//},
	//{
	//	Name: "struct_column",
	//	Type: arrow.StructOf([]arrow.Field{
	//		{Name: "int_col", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
	//		{Name: "string_col", Type: arrow.BinaryTypes.String, Nullable: true},
	//	}...),
	//},
}

func testResolverSuccess(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
	res <- map[string]any{
		"IntColumn":    3,
		"JsonColumn":   []byte(`{"test": "json"}`),
		"StringColumn": "test",
		"FloatColumn":  3.14,
		"BoolColumn":   true,
		"TimeColumn":   testSyncTime,
		"UuidColumn":   "00000000000040008000000000000000",
		"ArrayColumn":  []int64{},
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
		Columns:  testColumns,
	}
}

func testTableSuccessWithPK() *schema.Table {
	newCols := make([]schema.Column, len(testColumns))
	copy(newCols, testColumns)
	newCols[0].PrimaryKey = true
	return &schema.Table{
		Name:     "test_table_success",
		Resolver: testResolverSuccess,
		Columns:  newCols,
	}
}

func testTableResolverPanic() *schema.Table {
	return &schema.Table{
		Name:     "test_table_resolver_panic",
		Resolver: testResolverPanic,
		Columns:  testColumns,
	}
}

func testTablePreResourceResolverPanic() *schema.Table {
	return &schema.Table{
		Name:                "test_table_pre_resource_resolver_panic",
		PreResourceResolver: testPreResourceResolverPanic,
		Resolver:            testResolverSuccess,
		Columns:             testColumns,
	}
}

func testTableColumnResolverPanic() *schema.Table {
	return &schema.Table{
		Name:     "test_table_column_resolver_panic",
		Resolver: testResolverSuccess,
		Columns: append(testColumns,
			[]schema.Column{
				{
					Name:     "int_column1",
					Type:     arrow.PrimitiveTypes.Int64,
					Resolver: testColumnResolverPanic,
				},
			}...),
	}
}

func testTableRelationSuccess() *schema.Table {
	return &schema.Table{
		Name:     "test_table_relation_success",
		Resolver: testResolverSuccess,
		Columns:  testColumns,
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
	data              []scalar.Vector
	deterministicCQID bool
	typeSupport       specs.TypeSupport
}

var exampleData = scalar.Vector{
	&scalar.String{Value: "testSource", Valid: true},
	&scalar.Timestamp{Value: testSyncTime.UTC(), Valid: true},
	&scalar.UUID{Value: randomStableUUID, Valid: true},
	&scalar.UUID{},
	&scalar.Int64{Value: 3, Valid: true},
	&scalar.JSON{Value: []byte(`{"test": "json"}`), Valid: true},
	&scalar.String{Value: "test", Valid: true},
	&scalar.Float64{Value: 3.14, Valid: true},
	&scalar.Bool{Value: true, Valid: true},
	&scalar.Timestamp{Value: testSyncTime.UTC(), Valid: true},
	&scalar.UUID{Value: randomStableUUID, Valid: true},
	&scalar.List{Value: scalar.Vector{}, Valid: true},
	// TODO: map, struct
}

var exampleDataForRelation = scalar.Vector{
	&scalar.String{Value: "testSource", Valid: true},
	&scalar.Timestamp{Value: testSyncTime.UTC(), Valid: true},
	&scalar.UUID{Value: randomStableUUID, Valid: true},
	&scalar.UUID{Value: randomStableUUID, Valid: true},
	&scalar.Int64{Value: 3, Valid: true},
	&scalar.JSON{Value: []byte(`{"test": "json"}`), Valid: true},
	&scalar.String{Value: "test", Valid: true},
	&scalar.Float64{Value: 3.14, Valid: true},
	&scalar.Bool{Value: true, Valid: true},
	&scalar.Timestamp{Value: testSyncTime.UTC(), Valid: true},
	&scalar.UUID{Value: randomStableUUID, Valid: true},
	&scalar.List{Value: scalar.Vector{}, Valid: true},
	// TODO: map, struct
}

var exampleDataWithDeterministicStableUUID = scalar.Vector{
	&scalar.String{Value: "testSource", Valid: true},
	&scalar.Timestamp{Value: testSyncTime.UTC(), Valid: true},
	&scalar.UUID{Value: deterministicStableUUID, Valid: true},
	&scalar.UUID{},
	&scalar.Int64{Value: 3, Valid: true},
	&scalar.JSON{Value: []byte(`{"test": "json"}`), Valid: true},
	&scalar.String{Value: "test", Valid: true},
	&scalar.Float64{Value: 3.14, Valid: true},
	&scalar.Bool{Value: true, Valid: true},
	&scalar.Timestamp{Value: testSyncTime.UTC(), Valid: true},
	&scalar.UUID{Value: randomStableUUID, Valid: true},
	&scalar.List{Value: scalar.Vector{}, Valid: true},
	// TODO: map, struct
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
		data: []scalar.Vector{exampleData},
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
		data: []scalar.Vector{exampleData, exampleDataForRelation},
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
		data:              []scalar.Vector{exampleData},
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
		data:              []scalar.Vector{append(exampleData, &scalar.Int64{})},
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
		data:              []scalar.Vector{exampleData, exampleDataForRelation},
		deterministicCQID: true,
	},
	{
		table: testTableSuccessWithPK(),
		stats: Metrics{
			TableClient: map[string]map[string]*TableClientMetrics{
				"test_table_success": {
					"testExecutionClient": {
						Resources: 1,
					},
				},
			},
		},
		data:              []scalar.Vector{exampleDataWithDeterministicStableUUID},
		deterministicCQID: true,
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
		data:              []scalar.Vector{exampleData},
		deterministicCQID: false,
		typeSupport:       specs.TypeSupportFull,
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
			name := fmt.Sprintf("%s_%s_%s_types", tc.table.Name, scheduler.String(), tc.typeSupport.String())
			if tc.deterministicCQID {
				name += "_deterministic_cqid"
			}
			t.Run(name, func(t *testing.T) {
				testSyncTable(t, tc, scheduler, tc.deterministicCQID, tc.typeSupport)
			})
		}
	}
}

func testSyncTable(t *testing.T, tc syncTestCase, scheduler specs.Scheduler, deterministicCQID bool, typeSupport specs.TypeSupport) {
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
		TypeSupport:       typeSupport,
	}
	if err := plugin.Init(ctx, spec); err != nil {
		t.Fatal(err)
	}

	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		return plugin.Sync(ctx,
			testSyncTime,
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
					Type:          arrow.BinaryTypes.String,
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

func TestGetDynamicTables(t *testing.T) {
	addedColumns := []schema.Column{
		{
			Name: "_cq_source_name",
			Type: arrow.BinaryTypes.String,
		},
		{
			Name: "_cq_sync_time",
			Type: arrow.FixedWidthTypes.Timestamp_us,
		},
		schema.CqSourceNameColumn,
		schema.CqSyncTimeColumn,
		{
			Name:        "_cq_id",
			Type:        types.ExtensionTypes.UUID,
			Description: "Internal CQ ID of the row",
			NotNull:     true,
			Unique:      true,
			PrimaryKey:  true,
		},
		schema.CqParentIDColumn,
	}
	cases := []struct {
		TypeSupport specs.TypeSupport
		Want        schema.Tables
	}{
		{specs.TypeSupportLimited, schema.Tables{
			{Name: "test_table_success",
				Resolver: testResolverSuccess,
				Columns: append([]schema.Column{
					{
						Name: "int_column",
						Type: arrow.PrimitiveTypes.Int64,
					},
					{
						Name: "json_column",
						Type: types.ExtensionTypes.JSON,
					},
					{
						Name: "string_column",
						Type: arrow.BinaryTypes.String,
					},
					{
						Name: "float_column",
						Type: arrow.PrimitiveTypes.Float64,
					},
					{
						Name: "bool_column",
						Type: arrow.FixedWidthTypes.Boolean,
					},
					{
						Name: "time_column",
						Type: arrow.FixedWidthTypes.Timestamp_us,
					},
					{
						Name: "uuid_column",
						Type: types.ExtensionTypes.UUID,
					},
					{
						Name: "array_column",
						Type: arrow.ListOf(arrow.PrimitiveTypes.Int64),
					},
					//{
					//	Name: "map_column",
					//	Type: types.ExtensionTypes.JSON,
					//},
					//{
					//	Name: "struct_column",
					//	Type: types.ExtensionTypes.JSON,
					//},
				}, addedColumns...),
			},
		}},
		{specs.TypeSupportFull, schema.Tables{
			{Name: "test_table_success",
				Resolver: testResolverSuccess,
				Columns:  append(testTableSuccess().Columns, addedColumns...),
			},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.TypeSupport.String(), func(t *testing.T) {
			ctx := context.Background()
			tables := []*schema.Table{
				testTableSuccess(),
			}
			plugin := NewPlugin(
				"testSourcePlugin",
				"1.0.0",
				tables,
				newTestExecutionClient,
			)
			plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(t)))
			spec := specs.Source{
				Name:         "testSource",
				Path:         "cloudquery/testSource",
				Tables:       []string{"*"},
				Version:      "v1.0.0",
				Destinations: []string{"test"},
				Concurrency:  1,
				TypeSupport:  tc.TypeSupport,
			}
			if err := plugin.Init(ctx, spec); err != nil {
				t.Fatal(err)
			}
			got := plugin.GetDynamicTables()
			if len(got) != 1 {
				t.Fatalf("expected 1 table got %d", len(got))
			}
			changes := got[0].GetChanges(tc.Want[0])
			if len(changes) != 0 {
				s := ""
				for _, c := range changes {
					s += fmt.Sprintf("%s [%s -> %s] (%s)\n", c.ColumnName, c.Previous, c.Current, c.Type.String())
				}
				t.Fatalf("unexpected changes to table:\n%s", s)
			}
		})
	}
}
