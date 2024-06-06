package types

import (
	"testing"

	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUUIDBuilder(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))

	b.Append(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	b.AppendNull()
	b.Append(uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	b.AppendNull()

	require.Equal(t, 4, b.Len(), "unexpected Len()")
	require.Equal(t, 2, b.NullN(), "unexpected NullN()")

	values := []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		uuid.MustParse("00000000-0000-0000-0000-000000000004"),
	}
	b.AppendValues(values, nil)

	require.Equal(t, 6, b.Len(), "unexpected Len()")

	a := b.NewArray()

	// check state of builder after NewUUIDBuilder
	require.Zero(t, b.Len(), "unexpected ArrayBuilder.Len(), NewUUIDBuilder did not reset state")
	require.Zero(t, b.Cap(), "unexpected ArrayBuilder.Cap(), NewUUIDBuilder did not reset state")
	require.Zero(t, b.NullN(), "unexpected ArrayBuilder.NullN(), NewUUIDBuilder did not reset state")

	require.Equal(t, `["00000000-0000-0000-0000-000000000001" (null) "00000000-0000-0000-0000-000000000002" (null) "00000000-0000-0000-0000-000000000003" "00000000-0000-0000-0000-000000000004"]`, a.String())
	st, err := a.MarshalJSON()
	require.NoError(t, err)

	b.Release()
	a.Release()

	b = NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))
	err = b.UnmarshalJSON(st)
	require.NoError(t, err)

	a = b.NewArray()
	require.Equal(t, `["00000000-0000-0000-0000-000000000001" (null) "00000000-0000-0000-0000-000000000002" (null) "00000000-0000-0000-0000-000000000003" "00000000-0000-0000-0000-000000000004"]`, a.String())
	b.Release()
	a.Release()
}
