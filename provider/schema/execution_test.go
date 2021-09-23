package schema

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var alwaysDeleteTable = &Table{
	Name:         "always_delete_test_table",
	AlwaysDelete: true,
	Columns:      []Column{{Name: "name", Type: TypeString}},
}

var testMultiplexTable = &Table{
	Name: "test_multiplex_table",
	Multiplex: func(meta ClientMeta) []ClientMeta {
		return []ClientMeta{meta}
	},
	Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
		return nil
	},
	Columns: []Column{
		{
			Name: "name",
			Type: TypeString,
		},
	},
	Relations: []*Table{
		{
			Name: "test_relation_multiplex_table",
			Multiplex: func(meta ClientMeta) []ClientMeta {
				return []ClientMeta{meta}
			},
			Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
				return nil
			},
			Columns: []Column{
				{
					Name: "name",
					Type: TypeString,
				},
			},
		},
	},
}

var testTable = &Table{
	Name: "test_table",
	Columns: []Column{
		{
			Name: "name",
			Type: TypeString,
		},
		{
			Name:     "name_no_prefix",
			Type:     TypeString,
			Resolver: PathResolver("Inner.NameNoPrefix"),
		},
		{
			Name:     "prefix_name",
			Type:     TypeString,
			Resolver: PathResolver("Prefix.Name"),
		},
	},
}

type testTableStruct struct {
	Name  string `default:"test"`
	Inner struct {
		NameNoPrefix string `default:"name_no_prefix"`
	}
	Prefix struct {
		Name string `default:"prefix_name"`
	}
}

var testDefaultsTable = &Table{
	Name: "test_table",
	Columns: []Column{
		{
			Name:    "name",
			Type:    TypeString,
			Default: "defaultValue",
		},
	},
}

type testDefaultsTableData struct {
	Name         *string
	DefaultValue string
}

var testBadColumnResolverTable = &Table{
	Name: "test_table",
	Columns: []Column{
		{
			Name: "name",
			Type: TypeString,
			Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
				data := resource.Item.(testDefaultsTableData)
				if data.Name != nil && *data.Name == "noError" {
					return nil
				}
				return errors.New("ERROR")
			},
		},
	},
	Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
		res <- testDefaultsTableData{Name: nil}
		return nil
	},
}

var testIgnoreErrorColumnResolverTable = &Table{
	Name: "test_table",
	Columns: []Column{
		{
			Name: "name",
			Type: TypeString,
			IgnoreError: func(err error) bool {
				return true
			},
			Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
				return errors.New("ERROR")
			},
		},
		{
			Name: "default_value",
			Type: TypeString,
			IgnoreError: func(err error) bool {
				return true
			},
			Default: "TestValue",
			Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
				return errors.New("ERROR")
			},
		},
	},
	Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
		res <- testDefaultsTableData{Name: nil}
		return nil
	},
}

func failingTableResolver(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
	return fmt.Errorf("table resolve failed")
}

func doNothingResolver(_ context.Context, _ ClientMeta, _ *Resource, _ chan interface{}) error {
	return nil
}

func dataReturningResolver(_ context.Context, _ ClientMeta, _ *Resource, res chan interface{}) error {
	object := testTableStruct{}
	_ = defaults.Set(&object)
	res <- []testTableStruct{object, object, object}
	return nil
}

func dataReturningSingleResolver(_ context.Context, _ ClientMeta, _ *Resource, res chan interface{}) error {
	object := testTableStruct{}
	_ = defaults.Set(&object)
	res <- object
	return nil
}

func passingNilResolver(_ context.Context, _ ClientMeta, _ *Resource, res chan interface{}) error {
	res <- nil
	return nil
}

func TestExecutionData_ResolveTable(t *testing.T) {

	mockedClient := new(mockedClientMeta)
	logger := logging.New(&hclog.LoggerOptions{
		Name:   "test_log",
		Level:  hclog.Error,
		Output: nil,
	})
	mockedClient.On("Logger", mock.Anything).Return(logger)

	t.Run("failing table column resolver", func(t *testing.T) {
		testTable.Resolver = failingTableResolver
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Error(t, err)
		execFailing := NewExecutionData(mockDb, logger, testBadColumnResolverTable, false, nil, false)
		_, err = execFailing.ResolveTable(context.Background(), mockedClient, nil)
		assert.Error(t, err)
	})

	t.Run("ignore error table column resolver w/partialFetch", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		exec := NewExecutionData(mockDb, logger, testIgnoreErrorColumnResolverTable, false, nil, true)
		var expectedResource *Resource
		testIgnoreErrorColumnResolverTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			expectedResource = parent
			return nil
		}
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Len(t, exec.PartialFetchFailureResult, 0)
		assert.Equal(t, "TestValue", expectedResource.Get("default_value"))
		assert.Nil(t, expectedResource.Get("name"))
	})

	t.Run("error table column resolver w/partialFetch", func(t *testing.T) {
		testBadColumnResolverTable.Resolver = func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
			someString := "noError"
			res <- []testDefaultsTableData{{Name: &someString}, {Name: nil}, {Name: &someString}}
			return nil
		}
		mockDb := new(DatabaseMock)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		exec := NewExecutionData(mockDb, logger, testBadColumnResolverTable, false, nil, true)
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Len(t, exec.PartialFetchFailureResult, 1)
	})

	t.Run("doing nothing resolver", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		testTable.Resolver = doNothingResolver
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
	})

	t.Run("simple returning resources insert", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		testTable.Resolver = dataReturningResolver
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		mockDb.AssertNumberOfCalls(t, "CopyFrom", 1)
	})
	t.Run("simple returning resources insert w/disable_delete", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, true, mock.Anything).Return(nil)
		mockDb.On("RemoveStaleData", mock.Anything, testTable, exec.executionStart, mock.Anything).Return(nil)
		testTable.Resolver = dataReturningResolver
		exec.disableDelete = true
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		exec.disableDelete = false
		mockDb.AssertNumberOfCalls(t, "CopyFrom", 1)
		assert.Nil(t, err)
	})
	t.Run("simple returning single resources insert", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		testTable.Resolver = dataReturningSingleResolver
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
	})
	t.Run("simple returning nil resources insert", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		testTable.Resolver = passingNilResolver
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		mockDb.AssertNumberOfCalls(t, "CopyFrom", 0)
	})
	t.Run("check post row resolver", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, nil, false)
		testTable.Resolver = dataReturningSingleResolver
		var expectedResource *Resource
		testTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			err := parent.Set("name", "other")
			assert.Nil(t, err)
			expectedResource = parent
			return nil
		}
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Equal(t, expectedResource.data["name"], "other")
		assert.Nil(t, err)
		testTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			return errors.New("error")
		}
		_, err = exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Error(t, err)
	})

	t.Run("test resolving with default column values", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		execDefault := NewExecutionData(mockDb, logger, testDefaultsTable, false, nil, false)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		testDefaultsTable.Resolver = func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
			res <- testDefaultsTableData{Name: nil}
			return nil
		}
		var expectedResource *Resource
		testDefaultsTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			expectedResource = parent
			return nil
		}
		_, err := execDefault.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Equal(t, expectedResource.data["name"], "defaultValue")
	})

	t.Run("disable delete", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, true, nil, false)
		//mockDb.On("CopyFrom", mock.Anything, mock.Anything, true, mock.Anything).Return(nil)
		testTable.Resolver = dataReturningSingleResolver
		testTable.DeleteFilter = func(meta ClientMeta, r *Resource) []interface{} {
			return nil
		}
		var expectedResource *Resource
		testTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			err := parent.Set("name", "other")
			assert.Nil(t, err)
			expectedResource = parent
			return nil
		}
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, true, mock.Anything).Return(nil)
		mockDb.On("Delete", mock.Anything, testTable, mock.Anything).Return(nil)
		mockDb.On("RemoveStaleData", mock.Anything, testTable, exec.executionStart, mock.Anything).Return(nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 0)
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 0)
		mockDb.AssertNumberOfCalls(t, "CopyFrom", 1)
		assert.Equal(t, expectedResource.data["name"], "other")
		assert.Nil(t, err)
		exec = NewExecutionData(mockDb, logger, testTable, false, nil, false)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		_, err = exec.ResolveTable(context.Background(), mockedClient, nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 1)
		mockDb.AssertNumberOfCalls(t, "CopyFrom", 2)
		assert.Nil(t, err)
	})
	t.Run("disable delete failed copy from", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, true, nil, false)
		testTable.Resolver = dataReturningSingleResolver
		testTable.DeleteFilter = func(meta ClientMeta, r *Resource) []interface{} {
			return nil
		}
		var expectedResource *Resource
		testTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			err := parent.Set("name", "other")
			assert.Nil(t, err)
			expectedResource = parent
			return nil
		}
		mockDb.On("RemoveStaleData", mock.Anything, testTable, exec.executionStart, mock.Anything).Return(nil)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, true, mock.Anything).Return(fmt.Errorf("some error"))
		mockDb.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockDb.On("Delete", mock.Anything, testTable, mock.Anything).Return(nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 0)
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 0)
		mockDb.AssertNumberOfCalls(t, "CopyFrom", 1)
		mockDb.AssertNumberOfCalls(t, "Insert", 1)
		assert.Equal(t, expectedResource.data["name"], "other")
		assert.Nil(t, err)
	})

	t.Run("always delete with disable delete", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, alwaysDeleteTable, true, nil, false)
		alwaysDeleteTable.Resolver = dataReturningSingleResolver
		alwaysDeleteTable.DeleteFilter = func(meta ClientMeta, r *Resource) []interface{} {
			return nil
		}
		var expectedResource *Resource
		alwaysDeleteTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			err := parent.Set("name", "other")
			assert.Nil(t, err)
			expectedResource = parent
			return nil
		}
		mockDb.On("RemoveStaleData", mock.Anything, alwaysDeleteTable, exec.executionStart, mock.Anything).Return(nil)
		mockDb.On("Delete", mock.Anything, alwaysDeleteTable, mock.Anything).Return(nil)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, true, mock.Anything).Return(nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 0)
		_, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		mockDb.AssertNumberOfCalls(t, "Delete", 1)
		assert.Equal(t, expectedResource.data["name"], "other")
		assert.Nil(t, err)
	})

	t.Run("inject fields into execution", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		exec := NewExecutionData(mockDb, logger, testTable, false, map[string]interface{}{"injected_field": 1}, false)
		testTable.Resolver = dataReturningSingleResolver
		testTable.DeleteFilter = nil
		var expectedResource *Resource
		testTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			err := parent.Set("name", "other")
			assert.Nil(t, err)
			expectedResource = parent
			return nil
		}
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, map[string]interface{}{"injected_field": 1}).Return(nil)
		count, err := exec.ResolveTable(context.Background(), mockedClient, nil)
		assert.Equal(t, count, uint64(1))
		assert.Nil(t, err)
		assert.Equal(t, expectedResource.data["name"], "other")
		assert.Equal(t, 1, expectedResource.extraFields["injected_field"])
		values, err := expectedResource.Values()
		assert.Nil(t, err)
		assert.Equal(t, []string{"name", "name_no_prefix", "prefix_name", "cq_id", "meta", "injected_field"}, expectedResource.columns)
		assert.Equal(t, []interface{}{"other", "name_no_prefix", "prefix_name", expectedResource.cqId, expectedResource.Get("meta"), 1}, values)
	})

	t.Run("test partial fetch post resource resolver", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		execDefault := NewExecutionData(mockDb, logger, testDefaultsTable, false, nil, true)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		testDefaultsTable.Resolver = func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
			res <- testDefaultsTableData{Name: nil}
			return nil
		}
		var expectedResource *Resource
		testDefaultsTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			expectedResource = parent
			return fmt.Errorf("random failure")
		}
		_, err := execDefault.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Equal(t, expectedResource.data["name"], "defaultValue")
		assert.Len(t, execDefault.PartialFetchFailureResult, 1)
		assert.Equal(t, "failed to resolve resource: post resource resolver failed: random failure", execDefault.PartialFetchFailureResult[0].Error)
	})

	t.Run("test partial fetch resolver", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		execDefault := NewExecutionData(mockDb, logger, testDefaultsTable, false, nil, true)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		testDefaultsTable.Resolver = func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
			res <- testDefaultsTableData{Name: nil}
			return fmt.Errorf("random failure")
		}
		var expectedResource *Resource
		testDefaultsTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			expectedResource = parent
			return nil
		}
		_, err := execDefault.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Equal(t, expectedResource.data["name"], "defaultValue")
		assert.Len(t, execDefault.PartialFetchFailureResult, 1)
		assert.Equal(t, "table resolve error: random failure", execDefault.PartialFetchFailureResult[0].Error)
	})

	t.Run("test partial fetch resolver panic", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		execDefault := NewExecutionData(mockDb, logger, testDefaultsTable, false, nil, true)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		testDefaultsTable.Resolver = func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
			res <- testDefaultsTableData{Name: nil}
			panic("test panic")
		}
		var expectedResource *Resource
		testDefaultsTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			expectedResource = parent
			return nil
		}
		_, err := execDefault.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Equal(t, expectedResource.data["name"], "defaultValue")
		assert.Len(t, execDefault.PartialFetchFailureResult, 1)
		assert.Equal(t, "table resolve error: failed table test_table fetch. Error: test panic", execDefault.PartialFetchFailureResult[0].Error)
	})

	t.Run("test partial fetch post resource resolver panic", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		execDefault := NewExecutionData(mockDb, logger, testDefaultsTable, false, nil, true)
		mockDb.On("CopyFrom", mock.Anything, mock.Anything, false, mock.Anything).Return(nil)
		testDefaultsTable.Resolver = func(ctx context.Context, meta ClientMeta, parent *Resource, res chan interface{}) error {
			res <- testDefaultsTableData{Name: nil}
			return nil
		}
		var expectedResource *Resource
		testDefaultsTable.PostResourceResolver = func(ctx context.Context, meta ClientMeta, parent *Resource) error {
			expectedResource = parent
			panic("test panic")
		}
		_, err := execDefault.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.Equal(t, expectedResource.data["name"], "defaultValue")
		assert.Len(t, execDefault.PartialFetchFailureResult, 1)
		assert.Equal(t, "failed to resolve resource: recovered from panic: test panic", execDefault.PartialFetchFailureResult[0].Error)
	})

	t.Run("test table with multiplex", func(t *testing.T) {
		mockDb := new(DatabaseMock)
		execDefault := NewExecutionData(mockDb, logger, testMultiplexTable, false, nil, true)
		var parentMultiplexCalled, relationMultiplexCalled = false, false
		testMultiplexTable.Multiplex = func(meta ClientMeta) []ClientMeta {
			parentMultiplexCalled = true
			return []ClientMeta{meta}
		}
		testMultiplexTable.Relations[0].Multiplex = func(meta ClientMeta) []ClientMeta {
			relationMultiplexCalled = true
			return []ClientMeta{meta}
		}
		_, err := execDefault.ResolveTable(context.Background(), mockedClient, nil)
		assert.Nil(t, err)
		assert.True(t, parentMultiplexCalled)
		assert.False(t, relationMultiplexCalled)
	})
}

// ClientMeta is an autogenerated mock type for the ClientMeta type
type mockedClientMeta struct {
	mock.Mock
}

// Logger provides a mock function with given fields:
func (_m *mockedClientMeta) Logger() hclog.Logger {
	ret := _m.Called()

	var r0 hclog.Logger
	if rf, ok := ret.Get(0).(func() hclog.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(hclog.Logger)
		}
	}

	return r0
}
