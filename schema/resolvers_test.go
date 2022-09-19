package schema

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var resolverTestTable = &Table{
	Name: "testTable",
	Columns: []Column{
		{
			Name:     "string_column",
			Type:     TypeString,
			Resolver: PathResolver("PathResolver"),
		},
		{
			Name:     "at",
			Type:     TypeTimestamp,
			Resolver: PathResolver("At"),
		},
	},
}

var (
	timestamp        = time.Now().UTC()
	tsProto          = timestamppb.New(timestamp)
	resolverTestItem = map[string]interface{}{
		"PathResolver": "test",
		"At":           tsProto,
	}
)

var resolverTestCases = []struct {
	Name                 string
	Columns              ColumnList
	Resource             *Resource
	ExpectedResourceData map[string]interface{}
}{
	{
		Name:     "PathResolver",
		Columns:  resolverTestTable.Columns,
		Resource: NewResourceData(resolverTestTable, nil, time.Now(), resolverTestItem),
		ExpectedResourceData: map[string]interface{}{
			"string_column": "test",
			"at":            timestamp,
		},
	},
}

func TestResolvers(t *testing.T) {
	for _, tc := range resolverTestCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			for _, column := range tc.Columns {
				err := column.Resolver(context.Background(), nil, tc.Resource, column)
				require.NoError(t, err)
			}
			delete(tc.Resource.Data, "_cq_fetch_time")
			require.Equal(t, tc.ExpectedResourceData, tc.Resource.Data)
		})
	}
}
