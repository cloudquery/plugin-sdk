package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func testTableSuccessWithData(data []any) *schema.Table {
	return &schema.Table{
		Name: "test_table_success_with_data",
		Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- data
			return nil
		},
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}
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
		Name:     "test_table_success_pk",
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

func testTableSuccessWithCQIDPK() *schema.Table {
	return &schema.Table{
		Name:     "test_table_success_cq_id",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			schema.CqIDColumn,
			{
				Name:       "test_column",
				Type:       arrow.PrimitiveTypes.Int64,
				PrimaryKey: true,
			},
		},
	}
}

func testTableSuccessWithPKComponents() *schema.Table {
	cqID := schema.CqIDColumn
	cqID.PrimaryKey = true
	return &schema.Table{
		Name:     "test_table_succes_vpk__cq_id",
		Resolver: testResolverSuccess,
		Columns: []schema.Column{
			cqID,
			{
				Name:                "test_column",
				Type:                arrow.PrimitiveTypes.Int64,
				PrimaryKeyComponent: true,
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

func testNoTables() *schema.Table {
	return nil
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
	err               error
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
		table: testNoTables(),
		data:  nil,
		err:   ErrNoTables,
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
	{
		table: testTableSuccessWithCQIDPK(),
		data: []scalar.Vector{
			{
				// This value will be validated because deterministicCQID is true
				&scalar.UUID{Value: [16]byte{194, 83, 85, 170, 181, 44, 91, 112, 164, 224, 201, 153, 31, 90, 59, 135}, Valid: true},
				&scalar.Int{Value: 3, Valid: true},
			},
		},
		deterministicCQID: true,
	},
	{
		table: testTableSuccessWithCQIDPK(),
		data: []scalar.Vector{
			{
				// This value will not be validated as it will be randomly set by the scheduler
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
			},
		},
		deterministicCQID: false,
	},
	{
		table: testTableSuccessWithPKComponents(),
		data: []scalar.Vector{
			{
				// This value will not be validated as it will be randomly set by the scheduler
				&scalar.UUID{},
				&scalar.Int{Value: 3, Valid: true},
			},
		},
	},
}

type optionsTestCase struct {
	name    string
	options []Option
}

var allOptionsTestCases = []optionsTestCase{
	{name: "default_batching"},
	{name: "without_batching", options: []Option{WithoutBatching()}},
	{
		name:    "10 rows, 2s",
		options: []Option{WithBatchOptions(WithBatchTimeout(2*time.Second), WithBatchMaxRows(10))},
	},
}

func TestScheduler(t *testing.T) {
	for _, strategy := range AllStrategies {
		t.Run(strategy.String(), func(t *testing.T) {
			for _, opts := range allOptionsTestCases {
				t.Run(opts.name, func(t *testing.T) {
					for _, tc := range syncTestCases {
						testName := "No table_" + strategy.String()
						if tc.table != nil {
							tc.table = tc.table.Copy(nil)
							testName = tc.table.Name + "_" + strategy.String()
						}
						t.Run(testName, func(t *testing.T) {
							testSyncTable(t, tc, strategy, tc.deterministicCQID, opts.options...)
						})
					}
				})
			}
		})
	}
}

// nolint:revive
func testSyncTable(t *testing.T, tc syncTestCase, strategy Strategy, deterministicCQID bool, extra ...Option) {
	ctx := context.Background()
	var tables schema.Tables
	if tc.table != nil {
		tables = append(tables, tc.table)
	}
	c := testExecutionClient{}
	opts := append([]Option{
		WithLogger(zerolog.New(zerolog.NewTestWriter(t)).Level(zerolog.DebugLevel)),
		WithStrategy(strategy),
	}, extra...)
	sc := NewScheduler(opts...)
	msgs := make(chan message.SyncMessage, 10)
	err := sc.Sync(ctx, &c, tables, msgs, WithSyncDeterministicCQID(deterministicCQID))
	require.ErrorIs(t, err, tc.err)
	close(msgs)

	var i int
	for msg := range msgs {
		switch v := msg.(type) {
		case *message.SyncInsert:
			record := v.Record
			rec := tc.data[i].ToArrowRecord(record.Schema())
			if !array.RecordEqual(rec, record) {
				// For records that include CqIDColumn, we can't verify equality because it is generated by the scheduler, unless deterministicCQID is true
				onlyCqIDInequality := false
				for col := range rec.Columns() {
					if !deterministicCQID && rec.ColumnName(col) == schema.CqIDColumn.Name {
						onlyCqIDInequality = true
						continue
					}
					lc := rec.Column(col)
					rc := record.Column(col)
					if !array.Equal(lc, rc) {
						onlyCqIDInequality = false
					}
				}
				if !onlyCqIDInequality {
					t.Fatalf("expected at i=%d: %v. got %v", i, tc.data[i], record)
				}
			}
			i++
		case *message.SyncMigrateTable:
			migratedTable := v.Table

			initialTable := tables.Get(v.Table.Name)

			pks := migratedTable.PrimaryKeys()
			if (deterministicCQID || len(migratedTable.PrimaryKeyComponents()) > 0) && initialTable.Columns.Get(schema.CqIDColumn.Name) != nil {
				if len(pks) != 1 {
					t.Fatalf("expected 1 pk. got %d", len(pks))
				}
				if pks[0] != schema.CqIDColumn.Name {
					t.Fatalf("expected pk name %s. got %s", schema.CqIDColumn.Name, pks[0])
				}
			} else if len(pks) != len(initialTable.PrimaryKeys()) {
				t.Fatalf("expected 0 pk. got %d", len(pks))
			}

			if len(pks) == 0 {
				continue
			}
		default:
			t.Fatalf("expected insert message. got %T", msg)
		}
	}
	if len(tc.data) != i {
		t.Fatalf("expected %d resources. got %d", len(tc.data), i)
	}
}

func TestScheduler_Cancellation(t *testing.T) {
	data := make([]any, 100)

	tests := []struct {
		name           string
		data           []any
		cancel         bool
		messagesOrRows int
	}{
		{
			name:           "should consume all message",
			data:           data,
			cancel:         false,
			messagesOrRows: len(data) + 1, // 9 data + 1 migration message
		},
		{
			name:           "should not consume all message on cancel",
			data:           data,
			cancel:         true,
			messagesOrRows: len(data) + 1, // 9 data + 1 migration message
		},
	}

	for _, strategy := range AllStrategies {
		strategy := strategy
		for _, tc := range tests {
			tc := tc
			t.Run(fmt.Sprintf("%s_%s", tc.name, strategy.String()), func(t *testing.T) {
				logger := zerolog.New(zerolog.NewTestWriter(t))
				if tc.cancel {
					logger = zerolog.Nop() // FIXME without this, zerolog usage causes a race condition when tests are run with `-race -count=100`
				}
				sc := NewScheduler(WithLogger(logger), WithStrategy(strategy))

				messages := make(chan message.SyncMessage)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				go func() {
					err := sc.Sync(
						ctx,
						&testExecutionClient{},
						[]*schema.Table{testTableSuccessWithData(tc.data)},
						messages,
					)
					if tc.cancel {
						assert.Equal(t, err, context.Canceled)
					} else {
						require.NoError(t, err)
					}
					close(messages)
				}()

				messagesOrRows := 0
				for msg := range messages {
					if tc.cancel {
						cancel()
					}
					if r, ok := msg.(*message.SyncInsert); ok {
						messagesOrRows += int(r.Record.NumRows())
					} else {
						messagesOrRows++
					}
				}

				if tc.cancel {
					assert.NotEqual(t, tc.messagesOrRows, messagesOrRows)
				} else {
					assert.Equal(t, tc.messagesOrRows, messagesOrRows)
				}
			})
		}
	}
}
