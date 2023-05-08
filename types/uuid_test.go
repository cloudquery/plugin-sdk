package types

import (
	"testing"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestUUIDArray_ValueStr(t *testing.T) {
	// 1. create array
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))
	defer b.Release()

	b.AppendNull()
	b.Append(uuid.NameSpaceURL)
	b.AppendNull()
	b.Append(uuid.NameSpaceDNS)
	b.AppendNull()

	arr := b.NewUUIDArray()
	defer arr.Release()

	// 2. create array via AppendValueFromString
	b1 := NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))
	defer b1.Release()

	for i := 0; i < arr.Len(); i++ {
		assert.NoError(t, b1.AppendValueFromString(arr.ValueStr(i)))
	}

	arr1 := b1.NewUUIDArray()
	defer arr1.Release()

	assert.Equal(t, arr.Len(), arr1.Len())
	for i := 0; i < arr.Len(); i++ {
		assert.Equal(t, arr.IsValid(i), arr1.IsValid(i))
		assert.Equal(t, arr.ValueStr(i), arr1.ValueStr(i))
	}
}

func TestUUIDBuilder_AppendValueFromString(t *testing.T) {
	// 1. create array
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))
	defer b.Release()

	b.AppendNull()
	b.Append(uuid.NameSpaceURL)
	b.AppendNull()
	b.Append(uuid.NameSpaceDNS)
	b.AppendNull()

	arr := b.NewUUIDArray()
	defer arr.Release()

	// 2. create array via AppendValueFromString
	b1 := NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))
	defer b1.Release()

	for i := 0; i < arr.Len(); i++ {
		assert.NoError(t, b1.AppendValueFromString(arr.ValueStr(i)))
	}

	arr1 := b1.NewUUIDArray()
	defer arr1.Release()

	assert.Equal(t, arr.Len(), arr1.Len())
	for i := 0; i < arr.Len(); i++ {
		assert.Equal(t, arr.IsValid(i), arr1.IsValid(i))
		if arr.IsValid(i) {
			assert.Exactly(t, arr.Value(i), arr1.Value(i))
		}
	}
}
