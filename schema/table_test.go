package schema

import (
	"context"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

type tableTestCase struct {
	Table     Table
	Resources []*Resource
}

var tableTestCases = []tableTestCase{
	{
		Table: Table{
			Name: "testResolver",
			Columns: []Column{
				{
					Name: "test",
					Type: TypeInt,
				},
			},
			Resolver: func(ctx context.Context, meta ClientMeta, parent *Resource, res chan<- interface{}) error {
				res <- []map[string]interface{}{
					{
						"test": 1,
					},
				}
				return nil
			},
		},
		Resources: []*Resource{
			{
				data: map[string]interface{}{
					"test": 1,
				},
			},
		},
	},
}

type testClient struct {
}

func (c testClient) Logger() *zerolog.Logger {
	return &zerolog.Logger{}
}

func TestTableExecution(t *testing.T) {
	ctx := context.Background()
	for _, tc := range tableTestCases {
		tc := tc
		t.Run(tc.Table.Name, func(t *testing.T) {
			m := testClient{}
			resources := make(chan *Resource)
			go func() {
				defer close(resources)
				tc.Table.Resolve(ctx, m, nil, resources)
			}()
			var i = 0
			for resource := range resources {
				if reflect.DeepEqual(resource.data, tc.Resources[i].data) {
					t.Errorf("expected %v, got %v", tc.Resources[i].data, resource)
				}
				i += 1
			}
		})
	}
}
