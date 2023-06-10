package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type testExecutionClient struct {
}

func (t *testExecutionClient) ID() string {
	return "test"
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

type syncTestCase struct {
	table             *schema.Table
	data              []scalar.Vector
	deterministicCQID bool
}

var syncTestCases = []syncTestCase{
	{
		table: testTableSuccess(),
		data: []scalar.Vector{
			{
				&scalar.Int64{Value: 3, Valid: true},
			},
		},
	},
	{
		table: testTableResolverPanic(),
		data:  nil,
	},
	{
		table: testTablePreResourceResolverPanic(),
		data:  nil,
	},

	{
		table: testTableRelationSuccess(),
		data: []scalar.Vector{
			{
				&scalar.Int64{Value: 3, Valid: true},
			},
			{
				&scalar.Int64{Value: 3, Valid: true},
			},
		},
	},
	{
		table: testTableSuccess(),
		data: []scalar.Vector{
			{
				// &scalar.String{Value: "testSource", Valid: true},
				// &scalar.Timestamp{Value: testSyncTime, Valid: true},
				// &scalar.UUID{Value: randomStableUUID, Valid: true},
				// &scalar.UUID{},
				&scalar.Int64{Value: 3, Valid: true},
			},
		},
		deterministicCQID: true,
	},
	{
		table: testTableColumnResolverPanic(),
		data: []scalar.Vector{
			{
				&scalar.Int64{Value: 3, Valid: true},
				&scalar.Int64{},
			},
		},
		// deterministicCQID: true,
	},
	{
		table: testTableRelationSuccess(),
		data: []scalar.Vector{
			{
				// &scalar.String{Value: "testSource", Valid: true},
				// &scalar.Timestamp{Value: testSyncTime, Valid: true},
				// &scalar.UUID{Value: randomStableUUID, Valid: true},
				// &scalar.UUID{},
				&scalar.Int64{Value: 3, Valid: true},
			},
			{
				// &scalar.String{Value: "testSource", Valid: true},
				// &scalar.Timestamp{Value: testSyncTime, Valid: true},
				// &scalar.UUID{Value: randomStableUUID, Valid: true},
				// &scalar.UUID{Value: randomStableUUID, Valid: true},
				&scalar.Int64{Value: 3, Valid: true},
			},
		},
		// deterministicCQID: true,
	},
	{
		table: testTableSuccessWithPK(),
		data: []scalar.Vector{
			{
				// &scalar.String{Value: "testSource", Valid: true},
				// &scalar.Timestamp{Value: testSyncTime, Valid: true},
				// &scalar.UUID{Value: deterministicStableUUID, Valid: true},
				// &scalar.UUID{},
				&scalar.Int64{Value: 3, Valid: true},
			},
		},
		// deterministicCQID: true,
	},
}

func TestScheduler(t *testing.T) {
	// uuid.SetRand(testRand{})
	for _, scheduler := range AllSchedulers {
		for _, tc := range syncTestCases {
			tc := tc
			tc.table = tc.table.Copy(nil)
			t.Run(tc.table.Name+"_"+scheduler.String(), func(t *testing.T) {
				testSyncTable(t, tc, scheduler, tc.deterministicCQID)
			})
		}
	}
}

func testSyncTable(t *testing.T, tc syncTestCase, strategy SchedulerStrategy, deterministicCQID bool) {
	ctx := context.Background()
	tables := []*schema.Table{
		tc.table,
	}
	c := testExecutionClient{}
	opts := []Option{
		WithLogger(zerolog.New(zerolog.NewTestWriter(t))),
		WithSchedulerStrategy(strategy),
		// WithDeterministicCQId(deterministicCQID),
	}
	sc := NewScheduler(tables, &c, opts...)
	records := make(chan arrow.Record, 10)
	if err := sc.Sync(ctx, records); err != nil {
		t.Fatal(err)
	}
	close(records)

	var i int
	for record := range records {
		if tc.data == nil {
			t.Fatalf("Unexpected resource %v", record)
		}
		if i >= len(tc.data) {
			t.Fatalf("expected %d resources. got %d", len(tc.data), i)
		}
		rec := tc.data[i].ToArrowRecord(record.Schema())
		if !array.RecordEqual(rec, record) {
			t.Fatalf("expected at i=%d: %v. got %v", i, tc.data[i], record)
		}
		i++
	}
	if len(tc.data) != i {
		t.Fatalf("expected %d resources. got %d", len(tc.data), i)
	}
}
