package transformers_test

import (
	"reflect"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/stretchr/testify/require"
)

type myItem struct {
	ID   string
	Name string
}

func TestTransformWithStruct_PopulatesItemSample(t *testing.T) {
	tbl := &schema.Table{Name: "t1", Transform: transformers.TransformWithStruct(&myItem{})}
	require.NoError(t, tbl.Transform(tbl))

	got := tbl.ItemSampleType()
	require.NotNil(t, got, "TransformWithStruct should populate itemSample")
	require.Equal(t, reflect.TypeOf(myItem{}), got)
}
