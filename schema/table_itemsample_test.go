package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type sampleItem struct{ Name string }

func TestTable_SetItemSample(t *testing.T) {
	tbl := &Table{Name: "t1"}
	require.Nil(t, tbl.ItemSampleType())

	tbl.SetItemSample(sampleItem{})
	got := tbl.ItemSampleType()
	require.NotNil(t, got)
	require.Equal(t, reflect.TypeOf(sampleItem{}), got)
}

func TestTable_SetItemSample_PointerUnwrapped(t *testing.T) {
	tbl := &Table{Name: "t1"}
	tbl.SetItemSample(&sampleItem{})
	got := tbl.ItemSampleType()
	require.NotNil(t, got)
	require.Equal(t, reflect.TypeOf(sampleItem{}), got, "pointer should be unwrapped to the element type")
}

func TestTable_SetItemSample_Idempotent(t *testing.T) {
	tbl := &Table{Name: "t1"}
	tbl.SetItemSample(sampleItem{})
	// Second call with a different type is a no-op — first-write-wins.
	tbl.SetItemSample(42)
	got := tbl.ItemSampleType()
	require.Equal(t, reflect.TypeOf(sampleItem{}), got)
}
