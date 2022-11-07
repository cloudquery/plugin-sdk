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
	ExpectedResourceData CQTypes
}{
	{
		Name:                 "PathResolver",
		Column:               resolverTestTable.Columns[0],
		ColumnResolver:       PathResolver("PathResolver"),
		Resource:             NewResourceData(resolverTestTable, nil, resolverTestItem),
		ExpectedResourceData: CQTypes{&Text{Str: "test", Status: Present}},
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

			if len(tc.ExpectedResourceData) != len(tc.Resource.data) {
				t.Errorf("expected %d columns, got %d", len(tc.ExpectedResourceData), len(tc.Resource.data))
				return
			}
			for i, expected := range tc.ExpectedResourceData {
				if !reflect.DeepEqual(expected, tc.Resource.data[i]) {
					t.Errorf("expected %v, got %v", expected, tc.Resource.data[i])
				}
			}
		})
	}
}
