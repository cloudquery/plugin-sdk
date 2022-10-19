package schema

import (
	"context"
	"reflect"
	"testing"
)

var resolverTestTable = &Table{
	Name: "test_table",
	Columns: []Column{
		{
			Name: "string_column",
			Type: TypeString,
		},
	},
}

var resolverTestItem = map[string]interface{}{
	"PathResolver": "test",
}

var resolverTestCases = []struct {
	Name                 string
	Column               Column
	ColumnResolver       ColumnResolver
	Resource             *Resource
	ExpectedResourceData map[string]interface{}
}{
	{
		Name:                 "PathResolver",
		Column:               resolverTestTable.Columns[0],
		ColumnResolver:       PathResolver("PathResolver"),
		Resource:             NewResourceData(resolverTestTable, nil, resolverTestItem),
		ExpectedResourceData: map[string]interface{}{"string_column": "test"},
	},
}

func TestResolvers(t *testing.T) {
	for _, tc := range resolverTestCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.ColumnResolver(context.Background(), nil, tc.Resource, tc.Column)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			delete(tc.Resource.Data, "_cq_fetch_time")
			if !reflect.DeepEqual(tc.ExpectedResourceData, tc.Resource.Data) {
				t.Errorf("Expected %v, got %v", tc.ExpectedResourceData, tc.Resource.Data)
			}
		})
	}
}
