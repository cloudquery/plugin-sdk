package types

import (
	"net"
	"testing"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseMac(s string) net.HardwareAddr {
	mac, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return mac
}

func TestMacBuilder(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewMacBuilder(array.NewExtensionBuilder(mem, NewMacType()))

	b.Append(mustParseMac("00:00:00:00:00:01"))
	b.AppendNull()
	b.Append(mustParseMac("00:00:00:00:00:02"))
	b.AppendNull()

	require.Equal(t, 4, b.Len(), "unexpected Len()")
	require.Equal(t, 2, b.NullN(), "unexpected NullN()")

	values := []net.HardwareAddr{
		mustParseMac("00:00:00:00:00:03"),
		mustParseMac("00:00:00:00:00:04"),
	}
	b.AppendValues(values, nil)

	require.Equal(t, 6, b.Len(), "unexpected Len()")

	a := b.NewArray()

	// check state of builder after NewArray
	require.Zero(t, b.Len(), "unexpected ArrayBuilder.Len(), did not reset state")
	require.Zero(t, b.Cap(), "unexpected ArrayBuilder.Cap(), did not reset state")
	require.Zero(t, b.NullN(), "unexpected ArrayBuilder.NullN(), did not reset state")

	require.Equal(t, `["00:00:00:00:00:01" (null) "00:00:00:00:00:02" (null) "00:00:00:00:00:03" "00:00:00:00:00:04"]`, a.String())
	st, err := a.MarshalJSON()
	require.NoError(t, err)

	b.Release()
	a.Release()

	b = NewMacBuilder(array.NewExtensionBuilder(mem, NewMacType()))
	err = b.UnmarshalJSON(st)
	require.NoError(t, err)

	a = b.NewArray()
	require.Equal(t, `["00:00:00:00:00:01" (null) "00:00:00:00:00:02" (null) "00:00:00:00:00:03" "00:00:00:00:00:04"]`, a.String())
	b.Release()
	a.Release()
}

func TestMacArray_ValueStr(t *testing.T) {
	// 1. create array
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewMacBuilder(array.NewExtensionBuilder(mem, NewMacType()))
	defer b.Release()

	b.AppendNull()
	b.Append(mustParseMac("00:00:00:00:00:01"))
	b.AppendNull()
	b.Append(mustParseMac("00:00:00:00:00:02"))
	b.AppendNull()

	arr := b.NewMacArray()
	defer arr.Release()

	// 2. create array via AppendValueFromString
	b1 := NewMacBuilder(array.NewExtensionBuilder(mem, NewMacType()))
	defer b1.Release()

	for i := 0; i < arr.Len(); i++ {
		assert.NoError(t, b1.AppendValueFromString(arr.ValueStr(i)))
	}

	arr1 := b1.NewMacArray()
	defer arr1.Release()

	assert.Equal(t, arr.Len(), arr1.Len())
	for i := 0; i < arr.Len(); i++ {
		assert.Equal(t, arr.IsValid(i), arr1.IsValid(i))
		assert.Equal(t, arr.ValueStr(i), arr1.ValueStr(i))
	}
}

func TestMacBuilder_AppendValueFromString(t *testing.T) {
	// 1. create array
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewMacBuilder(array.NewExtensionBuilder(mem, NewMacType()))
	defer b.Release()

	b.AppendNull()
	b.Append(mustParseMac("00:00:00:00:00:01"))
	b.AppendNull()
	b.Append(mustParseMac("00:00:00:00:00:02"))
	b.AppendNull()

	arr := b.NewMacArray()
	defer arr.Release()

	// 2. create array via AppendValueFromString
	b1 := NewMacBuilder(array.NewExtensionBuilder(mem, NewMacType()))
	defer b1.Release()

	for i := 0; i < arr.Len(); i++ {
		assert.NoError(t, b1.AppendValueFromString(arr.ValueStr(i)))
	}

	arr1 := b1.NewMacArray()
	defer arr1.Release()

	assert.Equal(t, arr.Len(), arr1.Len())
	for i := 0; i < arr.Len(); i++ {
		assert.Equal(t, arr.IsValid(i), arr1.IsValid(i))
		if arr.IsValid(i) {
			assert.Exactly(t, arr.Value(i), arr1.Value(i))
		}
	}
}
