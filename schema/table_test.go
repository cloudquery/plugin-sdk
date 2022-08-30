package schema

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type tableTestCase struct {
	Table     Table
	Resources []*Resource
}

type testClient struct {
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
				Data: map[string]interface{}{
					"test": 1,
				},
			},
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
			go func() {
				defer close(resources)
				tc.Table.Resolve(ctx, m, time.Now(), nil, resources)
			}()
			var i = 0
			for resource := range resources {
				if reflect.DeepEqual(resource.Data, tc.Resources[i].Data) {
					t.Errorf("expected %v, got %v", tc.Resources[i].Data, resource)
				}
				i++
			}
		})
	}
}
