package types

import (
	"testing"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/stretchr/testify/require"
)

func TestJSONBuilder(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
	b.Append(map[string]any{"a": 1, "b": 2})
	b.AppendNull()
	b.Append(map[string]any{"c": 3, "d": 4})
	b.AppendNull()

	require.Equal(t, 4, b.Len(), "unexpected Len()")
	require.Equal(t, 2, b.NullN(), "unexpected NullN()")

	values := []any{
		map[string]any{"e": 5, "f": 6},
		map[string]any{"g": 7, "h": 8},
	}
	b.AppendValues(values, []bool{true, true})

	require.Equal(t, 6, b.Len(), "unexpected Len()")

	a := b.NewArray()

	// check state of builder after NewJSONBuilder
	require.Zero(t, b.Len(), "unexpected ArrayBuilder.Len(), NewJSONBuilder did not reset state")
	require.Zero(t, b.Cap(), "unexpected ArrayBuilder.Cap(), NewJSONBuilder did not reset state")
	require.Zero(t, b.NullN(), "unexpected ArrayBuilder.NullN(), NewJSONBuilder did not reset state")
	require.Equal(t, `["{"a":1,"b":2}" (null) "{"c":3,"d":4}" (null) "{"e":5,"f":6}" "{"g":7,"h":8}"]`, a.String())
	st, err := a.MarshalJSON()
	require.NoError(t, err)

	b.Release()
	a.Release()

	b = NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
	err = b.UnmarshalJSON(st)
	require.NoError(t, err)

	a = b.NewArray()
	require.Equal(t, `["{"a":1,"b":2}" (null) "{"c":3,"d":4}" (null) "{"e":5,"f":6}" "{"g":7,"h":8}"]`, a.String())
	b.Release()
	a.Release()
}