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

func TestTable_SetItemSample_IdempotentSameType(t *testing.T) {
	tbl := &Table{Name: "t1"}
	tbl.SetItemSample(sampleItem{})
	// Second call with the SAME type is a no-op.
	tbl.SetItemSample(sampleItem{})
	got := tbl.ItemSampleType()
	require.Equal(t, reflect.TypeOf(sampleItem{}), got)
}

func TestTable_SetItemSample_PanicsOnConflict(t *testing.T) {
	tbl := &Table{Name: "t1"}
	tbl.SetItemSample(sampleItem{})
	require.PanicsWithValue(t,
		`schema.Table "t1": itemSample already set to schema.sampleItem, got conflicting int`,
		func() {
			tbl.SetItemSample(42)
		},
	)
}

func TestTable_SetItemSample_IdempotentValueVsPointer(t *testing.T) {
	tbl := &Table{Name: "t1"}
	tbl.SetItemSample(sampleItem{})
	// Pointer-to-same-type is idempotent since pointer is unwrapped.
	tbl.SetItemSample(&sampleItem{})
	got := tbl.ItemSampleType()
	require.Equal(t, reflect.TypeOf(sampleItem{}), got)
}
