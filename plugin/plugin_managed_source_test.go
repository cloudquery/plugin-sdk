package plugin

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type testExecutionClient struct {
	UnimplementedWriter
}

var _ schema.ClientMeta = &testExecutionClient{}

var deterministicStableUUID = uuid.MustParse("c25355aab52c5b70a4e0c9991f5a3b87")
var randomStableUUID = uuid.MustParse("00000000000040008000000000000000")

var testSyncTime = time.Now()

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
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}
}

func testTableSuccessWithPK() *schema.Table {
	return &schema.Table{
		Name:     "test_table_success",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			{
				Name:       "test_column",
				Type:       arrow.PrimitiveTypes.Int64,
				PrimaryKey: true,
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
				Type: arrow.PrimitiveTypes.Int64,
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
				Type: arrow.PrimitiveTypes.Int64,
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
				Type: arrow.PrimitiveTypes.Int64,
			},
			{
				Name:     "test_column1",
				Type:     arrow.PrimitiveTypes.Int64,
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
				Type: arrow.PrimitiveTypes.Int64,
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

func (*testExecutionClient) Close(context.Context) error {
	return nil
}

func (*testExecutionClient) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- arrow.Record) error {
	return fmt.Errorf("not implemented")
}

func (*testExecutionClient) Sync(ctx context.Context, res chan<- arrow.Record) error {
	return fmt.Errorf("not implemented")
}

func newTestExecutionClient(context.Context, zerolog.Logger, pbPlugin.Spec) (Client, error) {
	return &testExecutionClient{}, nil
}

type syncTestCase struct {
	table             *schema.Table
	stats             Metrics
	data              []scalar.Vector
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
		data: []scalar.Vector{
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
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
		data: []scalar.Vector{
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
			},
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.Int{Value: 3, Valid: true},
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
		data: []scalar.Vector{
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
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
		data: []scalar.Vector{
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
				&scalar.Int{},
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
		data: []scalar.Vector{
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
			},
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.Int{Value: 3, Valid: true},
			},
		},
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
		data: []scalar.Vector{
			{
				&scalar.String{Value: "testSource", Valid: true},
				&scalar.Timestamp{Value: testSyncTime, Valid: true},
				&scalar.UUID{Value: deterministicStableUUID, Valid: true},
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
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

func TestManagedSync(t *testing.T) {
	uuid.SetRand(testRand{})
	for _, scheduler := range plugin.AllSchedulers {
		for _, tc := range syncTestCases {
			tc := tc
			tc.table = tc.table.Copy(nil)
			t.Run(tc.table.Name+"_"+scheduler.String(), func(t *testing.T) {
				testSyncTable(t, tc, scheduler, tc.deterministicCQID)
			})
		}
	}
}

func testSyncTable(t *testing.T, tc syncTestCase, scheduler plugin.Scheduler, deterministicCQID bool) {
	ctx := context.Background()
	tables := []*schema.Table{
		tc.table,
	}

	plugin := NewPlugin(
		"testSourcePlugin",
		"1.0.0",
		newTestExecutionClient,
		WithStaticTables(tables),
	)
	plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(t)))
	sourceName := "testSource"

	if err := plugin.Init(ctx, nil); err != nil {
		t.Fatal(err)
	}

	records, err := plugin.syncAll(ctx, sourceName, testSyncTime, SyncOptions{
		Tables:            []string{"*"},
		Concurrency:       1,
		Scheduler:         scheduler,
		DeterministicCQID: deterministicCQID,
	})
	if err != nil {
		t.Fatal(err)
	}

	var i int
	for _, record := range records {
		if tc.data == nil {
			t.Fatalf("Unexpected resource %v", record)
		}
		if i >= len(tc.data) {
			t.Fatalf("expected %d resources. got %d", len(tc.data), i)
		}
		rec := tc.data[i].ToArrowRecord(record.Schema())
		if !array.RecordEqual(rec, record) {
			t.Fatal(RecordDiff(rec, record))
			// t.Fatalf("expected at i=%d: %v. got %v", i, tc.data[i], record)
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
}

// func TestIgnoredColumns(t *testing.T) {
// 	table := &schema.Table{
// 		Columns: schema.ColumnList{
// 			{
// 				Name:          "a",
// 				Type:          arrow.BinaryTypes.String,
// 				IgnoreInTests: true,
// 			},
// 		},
// 	}
// 	validateResources(t, table, schema.Resources{{
// 		Item: struct{ A *string }{},
// 		Table: &schema.Table{
// 			Columns: schema.ColumnList{
// 				{
// 					Name:          "a",
// 					Type:          arrow.BinaryTypes.String,
// 					IgnoreInTests: true,
// 				},
// 			},
// 		},
// 	}})
// }

var testTable struct {
	PrimaryKey   string
	SecondaryKey string
	TertiaryKey  string
	Quaternary   string
}

// func TestNewPluginPrimaryKeys(t *testing.T) {
// 	testTransforms := []struct {
// 		transformerOptions []transformers.StructTransformerOption
// 		resultKeys         []string
// 	}{
// 		{
// 			transformerOptions: []transformers.StructTransformerOption{transformers.WithPrimaryKeys("PrimaryKey")},
// 			resultKeys:         []string{"primary_key"},
// 		},
// 		{
// 			transformerOptions: []transformers.StructTransformerOption{},
// 			resultKeys:         []string{"_cq_id"},
// 		},
// 	}
// 	for _, tc := range testTransforms {
// 		tables := []*schema.Table{
// 			{
// 				Name: "test_table",
// 				Transform: transformers.TransformWithStruct(
// 					&testTable, tc.transformerOptions...,
// 				),
// 			},
// 		}

// 		plugin := NewPlugin("testSourcePlugin", "1.0.0", tables, newTestExecutionClient)
// 		assert.Equal(t, tc.resultKeys, plugin.tables[0].PrimaryKeys())
// 	}
// }
