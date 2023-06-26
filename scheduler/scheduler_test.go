package scheduler

import (
	"context"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type testExecutionClient struct {
}

func (*testExecutionClient) ID() string {
	return "test"
}

var _ schema.ClientMeta = &testExecutionClient{}

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
				&scalar.Int{Value: 3, Valid: true},
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
				&scalar.Int{Value: 3, Valid: true},
			},
			{
				&scalar.Int{Value: 3, Valid: true},
			},
		},
	},
	{
		table: testTableSuccess(),
		data: []scalar.Vector{
			{
				&scalar.Int{Value: 3, Valid: true},
			},
		},
		deterministicCQID: true,
	},
	{
		table: testTableColumnResolverPanic(),
		data: []scalar.Vector{
			{
				&scalar.Int{Value: 3, Valid: true},
				&scalar.Int{},
			},
		},
		// deterministicCQID: true,
	},
	{
		table: testTableRelationSuccess(),
		data: []scalar.Vector{
			{
				&scalar.Int{Value: 3, Valid: true},
			},
			{
				&scalar.Int{Value: 3, Valid: true},
			},
		},
		// deterministicCQID: true,
	},
	{
		table: testTableSuccessWithPK(),
		data: []scalar.Vector{
			{
				&scalar.Int{Value: 3, Valid: true},
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

func testSyncTable(t *testing.T, tc syncTestCase, strategy Strategy, deterministicCQID bool) {
	ctx := context.Background()
	tables := []*schema.Table{
		tc.table,
	}
	c := testExecutionClient{}
	opts := []Option{
		WithLogger(zerolog.New(zerolog.NewTestWriter(t))),
		WithSchedulerStrategy(strategy),
	}
	sc := NewScheduler(&c, opts...)
	msgs := make(chan message.Message, 10)
	if err := sc.Sync(ctx, tables, msgs, WithSyncDeterministicCQID(deterministicCQID)); err != nil {
		t.Fatal(err)
	}
	close(msgs)

	var i int
	for msg := range msgs {
		switch v := msg.(type) {
		case *message.Insert:
			record := v.Record
			rec := tc.data[i].ToArrowRecord(record.Schema())
			if !array.RecordEqual(rec, record) {
				t.Fatalf("expected at i=%d: %v. got %v", i, tc.data[i], record)
			}
			i++
		case *message.MigrateTable:
			// ignore
		default:
			t.Fatalf("expected insert message. got %T", msg)
		}
	}
	if len(tc.data) != i {
		t.Fatalf("expected %d resources. got %d", len(tc.data), i)
	}
}
