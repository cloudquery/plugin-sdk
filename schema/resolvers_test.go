package schema

import (
	"context"
	"reflect"
	"testing"
)

type innerStruct struct {
	Value string
}

type testStruct struct {
	Inner      innerStruct
	Value      int
	unexported bool
}

type testDateStruct struct {
	Date string
}

type testNetStruct struct {
	IP  string
	MAC string
	Net string
	IPS []string
}

type testTransformersStruct struct {
	Int      int
	String   string
	Float    float64
	BadFloat string
}

type testUUIDStruct struct {
	UUID    string
	BadUUID string
}

var resolverTestTable = &Table{
	Name: "testTable",
	Columns: []Column{
		{
			Name: "stringColumn",
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
		ExpectedResourceData: map[string]interface{}{"stringColumn": "test"},
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
			if !reflect.DeepEqual(tc.Resource.data, tc.ExpectedResourceData) {
				t.Errorf("Expected %v, got %v", tc.ExpectedResourceData, tc.Resource.data)
			}
		})
	}
}
