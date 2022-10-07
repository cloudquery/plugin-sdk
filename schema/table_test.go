package schema

import (
	"context"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

type tableTestCase struct {
	Table     *Table
	Resources []*Resource
	Summary   *SyncSummary
}

type testClient struct {
}

func testTableSuccess() *Table {
	return &Table{
		Name: "testTableSuccess",
		Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []Column{
			{
				Name: "test_column",
				Type: TypeInt,
			},
		},
	}
}

func testTableRelationSuccess() *Table {
	return &Table{
		Name: "testTableRelationSuccess",
		Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []Column{
			{
				Name: "test_column",
				Type: TypeInt,
			},
		},
		Relations: []*Table{
			testTableSuccess(),
		},
	}
}

func testTableRelationPanic() *Table {
	return &Table{
		Name: "testTableRelationSuccess",
		Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []Column{
			{
				Name: "test_column",
				Type: TypeInt,
			},
		},
		Relations: []*Table{
			testTableResolverPanic(),
		},
	}
}

func testTableResolverPanic() *Table {
	return &Table{
		Name: "testTableResolverPanic",
		Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
			panic("Resolver")
		},
		Columns: []Column{
			{
				Name: "test_column",
				Type: TypeInt,
			},
		},
	}
}

func testPreResourceResolverPanic() *Table {
	return &Table{
		Name: "testPreResourceResolverPanic",
		PreResourceResolver: func(ctx context.Context, meta ClientMeta, resource *Resource) error {
			panic("PreResourceResolver")
		},
		Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []Column{
			{
				Name: "test_column",
				Type: TypeInt,
			},
		},
	}
}

func testColumnResolverPanic() *Table {
	return &Table{
		Name: "testColumnResolverPanic",
		Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []Column{
			{
				Name: "test_column",
				Type: TypeInt,
				Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
					panic("ColumnResolver")
				},
			},
		},
	}
}

var tableTestCases = []tableTestCase{
	{
		Table: testTableSuccess(),
		Resources: []*Resource{
			{
				Data: map[string]interface{}{
					"test": 1,
				},
			},
		},
		Summary: &SyncSummary{
			Resources: 1,
		},
	},
	{
		Table:     testTableResolverPanic(),
		Resources: nil,
		Summary: &SyncSummary{
			Panics: 1,
		},
	},
	{
		Table:     testPreResourceResolverPanic(),
		Resources: []*Resource{},
		Summary: &SyncSummary{
			Panics: 1,
		},
	},
	{
		Table:     testColumnResolverPanic(),
		Resources: []*Resource{},
		Summary: &SyncSummary{
			Panics: 1,
		},
	},
	{
		Table: testTableRelationSuccess(),
		Resources: []*Resource{
			{
				Data: map[string]interface{}{
					"test": 1,
				},
			},
			{
				Data: map[string]interface{}{
					"test": 1,
				},
			},
		},
		Summary: &SyncSummary{
			Resources: 2,
		},
	},
	{
		Table: testTableRelationPanic(),
		Resources: []*Resource{
			{
				Data: map[string]interface{}{
					"test": 1,
				},
			},
		},
		Summary: &SyncSummary{
			Panics:    1,
			Resources: 1,
		},
	},
}

func (testClient) Logger() *zerolog.Logger {
	return &zerolog.Logger{}
}

func TestTableExecution(t *testing.T) {
	ctx := context.Background()
	for _, tc := range tableTestCases {
		tc := tc
		t.Run(tc.Table.Name, func(t *testing.T) {
			m := testClient{}
			resources := make(chan *Resource)
			var summary *SyncSummary
			go func() {
				defer close(resources)
				summary = tc.Table.Resolve(ctx, m, nil, resources)
			}()
			var i = uint64(0)
			for resource := range resources {
				if reflect.DeepEqual(resource.Data, tc.Resources[i].Data) {
					t.Errorf("expected %v, got %v", tc.Resources[i].Data, resource)
				}
				i++
			}
			if tc.Summary.Resources != i {
				t.Errorf("expected %d resources, got %d", tc.Summary.Resources, i)
			}
			if tc.Summary.Resources != summary.Resources {
				t.Errorf("expected %d summary resources, got %d", tc.Summary.Resources, summary.Resources)
			}
			if tc.Summary.Errors != summary.Errors {
				t.Errorf("expected %d errors, got %d", tc.Summary.Errors, summary.Errors)
			}
			if tc.Summary.Panics != summary.Panics {
				t.Errorf("expected %d panics, got %d", tc.Summary.Panics, summary.Panics)
			}
		})
	}
}
