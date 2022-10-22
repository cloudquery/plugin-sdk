package schema

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudquery/plugin-sdk/cqtypes"
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
	ExpectedResourceData []interface{}
}{
	{
		Name:                 "PathResolver",
		Column:               resolverTestTable.Columns[0],
		ColumnResolver:       PathResolver("PathResolver"),
		Resource:             NewResourceData(resolverTestTable, nil, resolverTestItem),
		ExpectedResourceData: []interface{}{&cqtypes.Text{String: "test", Status: cqtypes.Present}},
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
			// delete(tc.Resource.data, "_cq_fetch_time")
			if !reflect.DeepEqual(tc.ExpectedResourceData, tc.Resource.data) {
				t.Errorf("Expected %v, got %v", tc.ExpectedResourceData, tc.Resource.data)
			}
		})
	}
}
