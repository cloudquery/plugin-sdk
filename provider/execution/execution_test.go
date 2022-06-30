package execution

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/cloudquery/cq-provider-sdk/helpers/limit"
	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/cq-provider-sdk/testlog"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/semaphore"
)

type ExecutionTestCase struct {
	Name  string
	Table *schema.Table

	SetupStorage          func(t *testing.T) Storage
	ExpectedResourceCount uint64
	ErrorExpected         bool
	ExpectedDiags         []diag.FlatDiag
}

type executionClient struct {
	l hclog.Logger
}

type zeroValuedStruct struct {
	ZeroBool      bool   `default:"false"`
	ZeroInt       int    `default:"0"`
	NotZeroInt    int    `default:"5"`
	NotZeroBool   bool   `default:"true"`
	ZeroIntPtr    *int   `default:"0"`
	NotZeroIntPtr *int   `default:"5"`
	ZeroString    string `default:""`
}

type resolveColumnsTestCase struct {
	Name         string
	Table        *schema.Table
	ResourceData interface{}
	MetaData     map[string]interface{}

	SetupStorage   func(t *testing.T) Storage
	CompareValues  func(t *testing.T, r *schema.Resource, want []interface{})
	ExpectedValues []interface{}
	ExpectedDiags  []diag.FlatDiag
}

var (
	commonColumns     = []schema.Column{{Name: "name", Type: schema.TypeString}}
	doNothingResolver = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		return nil
	}
	returnErrorResolver = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		return fmt.Errorf("some error")
	}

	returnValueResolver = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		res <- map[string]string{"name": "test"}
		return nil
	}

	returnWrapErrorResolver = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		return diag.WrapError(fmt.Errorf("some error"))
	}

	panicResolver = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		panic("resolver panic")
	}

	simpleMultiplexer = func(meta schema.ClientMeta) []schema.ClientMeta {
		return []schema.ClientMeta{meta, meta}
	}

	emptyMultiplexer = func(meta schema.ClientMeta) []schema.ClientMeta {
		return nil
	}

	postResourceResolver = func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
		return resource.Set("name", "data")
	}
	postResourceResolverError = func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
		return diag.Diagnostics{
			diag.NewBaseError(nil, diag.RESOLVING, diag.WithResourceName(resource.TableName()), diag.WithSummary("some error")),
			diag.NewBaseError(nil, diag.RESOLVING, diag.WithResourceName(resource.TableName()), diag.WithSummary("some error 2")),
		}
	}
	postResourceResolverWarning = func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
		return diag.Diagnostics{
			diag.NewBaseError(nil, diag.RESOLVING, diag.WithResourceName(resource.TableName()), diag.WithSummary("some warning"), diag.WithSeverity(diag.WARNING)),
		}
	}

	timeoutResolver = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Minute):
			panic("timeoutResolver timed out unexpectedly")
		}
	}
	testZeroTable = &schema.Table{
		Name: "test_zero_table",
		Columns: []schema.Column{
			{
				Name: "zero_bool",
				Type: schema.TypeBool,
			},
			{
				Name: "zero_int",
				Type: schema.TypeBigInt,
			},
			{
				Name: "not_zero_bool",
				Type: schema.TypeBool,
			},
			{
				Name: "not_zero_int",
				Type: schema.TypeBigInt,
			},
			{
				Name: "zero_int_ptr",
				Type: schema.TypeBigInt,
			},
			{
				Name: "not_zero_int_ptr",
				Type: schema.TypeBigInt,
			},
			{
				Name: "zero_string",
				Type: schema.TypeString,
			},
		},
	}
)

func (e executionClient) Logger() hclog.Logger {
	return e.l
}

func TestTableExecutor_Resolve(t *testing.T) {
	testCases := []ExecutionTestCase{
		{
			Name: "simple",
			Table: &schema.Table{
				Name:     "simple",
				Resolver: doNothingResolver,
				Columns:  commonColumns,
			},
		},
		{
			Name: "multiplex",
			Table: &schema.Table{
				Name:      "simple",
				Multiplex: simpleMultiplexer,
				Resolver:  doNothingResolver,
				Columns:   commonColumns,
			},
		},
		{
			Name: "multiplex_relation",
			Table: &schema.Table{
				Name:      "multiplex_relation",
				Multiplex: simpleMultiplexer,
				Resolver:  returnValueResolver,
				Columns:   commonColumns,
				Relations: []*schema.Table{
					{
						Resolver:  doNothingResolver,
						Multiplex: simpleMultiplexer,
						Name:      "relation_with_multiplex",
						Columns:   commonColumns,
					},
				},
			},
			ExpectedResourceCount: 2,
		},
		{
			Name: "multiplex_empty",
			Table: &schema.Table{
				Name:      "multiplex_empty",
				Multiplex: emptyMultiplexer,
				Resolver:  returnValueResolver,
				Columns:   commonColumns,
			},
			ExpectedResourceCount: 0,
		},
		{
			// if tables don't define a resolver, an execution error by execution
			Name: "missing_table_resolver",
			Table: &schema.Table{
				Name:    "no-resolver",
				Columns: commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      `table "no-resolver" missing resolver, make sure table implements the resolver`,
					Resource: "missing_table_resolver",
					Severity: diag.ERROR,
					Summary:  `table "no-resolver" missing resolver, make sure table implements the resolver`,
					Type:     diag.SCHEMA,
				},
			},
		},
		{
			// if tables don't define a resolver, an execution error by execution, we check here that a relation will cause an error
			Name: "missing_table_relation_resolver",
			Table: &schema.Table{
				Name:     "no-resolver",
				Resolver: returnValueResolver,
				Columns:  commonColumns,
				Relations: []*schema.Table{
					{
						Name:    "relation-no-resolver",
						Columns: commonColumns,
					},
				},
			},
			ExpectedResourceCount: 1,
			ErrorExpected:         true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      `table "relation-no-resolver" missing resolver, make sure table implements the resolver`,
					Resource: "missing_table_relation_resolver",
					Severity: diag.ERROR,
					Summary:  `table "relation-no-resolver" missing resolver, make sure table implements the resolver`,
					Type:     diag.SCHEMA,
				},
			},
		},
		{
			Name: "panic_resolver",
			Table: &schema.Table{
				Name:     "panic_resolver",
				Resolver: panicResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "table resolver panic: resolver panic",
					Resource: "panic_resolver",
					Severity: diag.PANIC,
					Summary:  `panic on resource table "panic_resolver" fetch: table resolver panic: resolver panic`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "panic_relation_resolver",
			Table: &schema.Table{
				Name:     "panic_resolver",
				Resolver: returnValueResolver,
				Columns:  commonColumns,
				Relations: []*schema.Table{
					{
						Name:     "relation_panic_resolver",
						Resolver: panicResolver,
						Columns:  commonColumns,
					},
				},
			},
			ExpectedResourceCount: 1,
			ErrorExpected:         true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "table resolver panic: resolver panic",
					Resource: "panic_relation_resolver",
					Severity: diag.PANIC,
					Summary:  `panic on resource table "relation_panic_resolver" fetch: table resolver panic: resolver panic`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "error_returning",
			Table: &schema.Table{
				Name:     "simple",
				Resolver: returnErrorResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "some error",
					Resource: "error_returning",
					Severity: diag.ERROR,
					Summary:  `failed to resolve table "simple": some error`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "error_returning_ignore_fail",
			Table: &schema.Table{
				IgnoreError: func(err error) bool {
					return false
				},
				Name:     "simple",
				Resolver: returnErrorResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "some error",
					Resource: "error_returning_ignore_fail",
					Severity: diag.ERROR,
					Summary:  `failed to resolve table "simple": some error`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "error_returning_ignore",
			Table: &schema.Table{
				IgnoreError: func(err error) bool {
					return true
				},
				Name:     "simple",
				Resolver: returnErrorResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "some error",
					Resource: "error_returning_ignore",
					Severity: diag.IGNORE,
					Summary:  `table "simple" resolver ignored error: some error`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "cleanup_stale_data_fail",
			SetupStorage: func(t *testing.T) Storage {
				db := new(DatabaseMock)
				db.On("RemoveStaleData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(fromError(errors.New("failed delete"), diag.WithResourceName("cleanup_stale_data_fail"), diag.WithType(diag.DATABASE)))
				db.On("Dialect").Return(noopDialect{})
				return db
			},
			Table: &schema.Table{
				Name: "cleanup_delete",
				DeleteFilter: func(meta schema.ClientMeta, parent *schema.Resource) []interface{} {
					return []interface{}{}
				},
				Resolver: doNothingResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "failed delete",
					Resource: "cleanup_stale_data_fail",
					Severity: diag.ERROR,
					Type:     diag.DATABASE,
					Summary:  `failed to cleanup stale data on table "cleanup_delete": failed delete`,
				},
			},
		},
		{
			Name: "post_resource_resolver",
			SetupStorage: func(t *testing.T) Storage {
				db := new(DatabaseMock)
				db.On("RemoveStaleData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				db.On("Dialect").Return(noopDialect{})
				db.On("CopyFrom", mock.Anything, mock.Anything, true).Return(nil)
				return db
			},
			Table: &schema.Table{
				Name:                 "simple",
				Resolver:             returnValueResolver,
				Columns:              commonColumns,
				PostResourceResolver: postResourceResolver,
			},
			ExpectedResourceCount: 1,
		},
		{
			Name: "post_resource_resolver_fail",
			SetupStorage: func(t *testing.T) Storage {
				db := new(DatabaseMock)
				db.On("RemoveStaleData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				db.On("Dialect").Return(noopDialect{})
				db.On("CopyFrom", mock.Anything, mock.Anything, true).Return(nil)
				return db
			},
			Table: &schema.Table{
				Name:                 "simple",
				Resolver:             returnValueResolver,
				Columns:              commonColumns,
				PostResourceResolver: postResourceResolverError,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "some error",
					Resource: "simple",
					Severity: diag.ERROR,
					Type:     diag.RESOLVING,
					Summary:  "some error",
				},
				{
					Err:      "some error 2",
					Resource: "simple",
					Severity: diag.ERROR,
					Type:     diag.RESOLVING,
					Summary:  "some error 2",
				},
			},
		},
		{
			Name: "post_resource_resolver_warning",
			SetupStorage: func(t *testing.T) Storage {
				db := new(DatabaseMock)
				db.On("RemoveStaleData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				db.On("Dialect").Return(schema.PostgresDialect{})
				db.On("CopyFrom", mock.Anything, mock.Anything, true).Return(nil).Run(
					func(args mock.Arguments) {
						resources := args.Get(1).(schema.Resources)
						if !assert.Greater(t, len(resources), 0) {
							return
						}

						assert.NotNil(t, resources[0].Get("cq_id"))
					})
				return db
			},
			Table: &schema.Table{
				Name:                 "post_resource_resolver_warning_table",
				Resolver:             returnValueResolver,
				Columns:              commonColumns,
				PostResourceResolver: postResourceResolverWarning,
			},
			ExpectedResourceCount: 1,
			ErrorExpected:         true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "some warning",
					Resource: "post_resource_resolver_warning_table",
					Severity: diag.WARNING,
					Type:     diag.RESOLVING,
					Summary:  "some warning",
				},
			},
		},
		{
			Name: "failing_column",
			SetupStorage: func(t *testing.T) Storage {
				db := new(DatabaseMock)
				db.On("RemoveStaleData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				db.On("Dialect").Return(noopDialect{})
				db.On("CopyFrom", mock.Anything, mock.Anything, true).Return(nil)
				return db
			},
			Table: &schema.Table{
				Name:     "column",
				Resolver: returnValueResolver,
				Columns: schema.ColumnList{
					{
						Name: "name",
						Resolver: func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
							return fmt.Errorf("failed column")
						},
					},
				},
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "failed column",
					Resource: "failing_column",
					Severity: diag.ERROR,
					Type:     diag.RESOLVING,
					Summary:  `column resolver "name" failed for table "column": failed column`,
				},
			},
		},
		{
			Name: "internal_column",
			Table: &schema.Table{
				Name:     "column",
				Resolver: returnValueResolver,
				Columns:  commonColumns,
			},
			ExpectedResourceCount: 1,
		},
		{
			Name: "return_wrap_error",
			Table: &schema.Table{
				Name:     "simple",
				Resolver: returnWrapErrorResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      `error at github.com/cloudquery/cq-provider-sdk/provider/execution.glob..func4[execution_test.go:74] some error`,
					Resource: "return_wrap_error",
					Severity: diag.ERROR,
					Summary:  `failed to resolve table "simple": error at github.com/cloudquery/cq-provider-sdk/provider/execution.glob..func4[execution_test.go:74] some error`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "timeout_resolver",
			Table: &schema.Table{
				Name:     "timeout_resolver",
				Resolver: timeoutResolver,
				Columns:  commonColumns,
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "context deadline exceeded",
					Resource: "timeout_resolver",
					Severity: diag.ERROR,
					Summary:  `failed to resolve table "timeout_resolver": context deadline exceeded`,
					Type:     diag.RESOLVING,
				},
			},
		},
		{
			Name: "panic_column",
			Table: &schema.Table{
				Name:     "panic_column_table",
				Resolver: returnValueResolver,
				Columns: schema.ColumnList{
					{
						Name: "name",
						Resolver: func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
							return resource.Set(c.Name, "name_value")
						},
					},
					{
						Name: "tags",
						Resolver: func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
							panic("oops")
						},
					},
				},
			},
			ErrorExpected: true,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "column resolve panic: oops",
					Resource: "panic_column",
					Severity: diag.PANIC,
					Type:     diag.RESOLVING,
					Summary:  `resolve column "tags" in table "panic_column_table" recovered from panic: column resolve panic: oops`,
				},
			},
		},
		{
			Name: "ignore_error_recursive",
			Table: &schema.Table{
				Name:        "simple",
				Resolver:    returnValueResolver,
				IgnoreError: func(err error) bool { return true },
				Columns:     commonColumns,
				Relations: []*schema.Table{
					{
						Name:     "simple",
						Resolver: returnErrorResolver,
						Columns:  commonColumns,
					},
				},
			},
			ErrorExpected:         true,
			ExpectedResourceCount: 1,
			ExpectedDiags: []diag.FlatDiag{
				{
					Err:      "some error",
					Resource: "ignore_error_recursive",
					Severity: diag.IGNORE,
					Summary:  `table "simple" resolver ignored error: some error`,
					Type:     diag.RESOLVING,
				},
			},
		},
	}

	executionClient := executionClient{testlog.New(t)}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var storage Storage = noopStorage{}
			if tc.SetupStorage != nil {
				storage = tc.SetupStorage(t)
			}
			if tc.Name == "ignore_error_recursive" {
				fmt.Println("debug")
			}
			limiter := semaphore.NewWeighted(int64(limit.GetMaxGoRoutines()))
			exec := NewTableExecutor(tc.Name, storage, testlog.New(t), tc.Table, nil, nil, limiter, 10*time.Second)
			count, diags := exec.Resolve(context.Background(), executionClient)
			assert.Equal(t, tc.ExpectedResourceCount, count)
			if tc.ErrorExpected {
				require.True(t, diags.HasDiags())
				if tc.ExpectedDiags != nil {
					assert.EqualValues(t, tc.ExpectedDiags, diag.FlattenDiags(diags, true))
				}
			} else {
				require.Empty(t, diags)
			}
		})
	}
}

func TestTableExecutor_resolveResourceValues(t *testing.T) {
	testCases := []resolveColumnsTestCase{
		{
			Name:  "resolve all zeroed columns",
			Table: testZeroTable,
			ResourceData: func() interface{} {
				object := zeroValuedStruct{}
				_ = defaults.Set(&object)
				return object
			}(),
			MetaData:       nil,
			SetupStorage:   nil,
			ExpectedValues: []interface{}{false, 0, true, 5, ptr.Int(0), ptr.Int(5), ""},
			ExpectedDiags:  nil,
		},
		{
			Name:  "resolve_columns with dialect",
			Table: testZeroTable,
			ResourceData: func() interface{} {
				object := zeroValuedStruct{}
				_ = defaults.Set(&object)
				return object
			}(),
			MetaData: nil,
			SetupStorage: func(t *testing.T) Storage {
				return &noopStorage{schema.PostgresDialect{}}
			},
			ExpectedValues: []interface{}{nil, nil, false, 0, true, 5, ptr.Int(0), ptr.Int(5), ""},
			CompareValues: func(t *testing.T, r *schema.Resource, want []interface{}) {
				v, err := r.Values()
				// skip cq_id, cq_meta
				assert.Equal(t, want[2:], v[2:])
				assert.NotNil(t, v[0])
				assert.NotNil(t, v[1])
				assert.Nil(t, err)
			},
			ExpectedDiags: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var storage Storage = noopStorage{}
			if tc.SetupStorage != nil {
				storage = tc.SetupStorage(t)
			}
			limiter := semaphore.NewWeighted(int64(limit.GetMaxGoRoutines()))
			exec := NewTableExecutor(tc.Name, storage, testlog.New(t), tc.Table, nil, nil, limiter, 0)

			r := schema.NewResourceData(storage.Dialect(), tc.Table, nil, tc.ResourceData, tc.MetaData, exec.executionStart)
			// columns should be resolved from ColumnResolver functions or default functions
			cl := executionClient{testlog.New(t)}
			diags := exec.resolveResourceValues(context.Background(), cl, r)
			if tc.ExpectedDiags != nil {
				require.True(t, diags.HasDiags())
				if tc.ExpectedDiags != nil {
					assert.Equal(t, tc.ExpectedDiags, diag.FlattenDiags(diags, true))
				}
			} else {
				require.Nil(t, diags)
			}
			if tc.CompareValues != nil {
				tc.CompareValues(t, r, tc.ExpectedValues)
			} else {
				v, err := r.Values()
				assert.Equal(t, tc.ExpectedValues, v)
				assert.Nil(t, err)
			}
		})
	}
}
