package types

import (
	"net"
	"testing"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseInet(s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ipnet
}

func TestInetBuilder(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))

	b.Append(mustParseInet("192.168.0.0/24"))
	b.AppendNull()
	b.Append(mustParseInet("192.168.0.0/25"))
	b.AppendNull()

	require.Equal(t, 4, b.Len(), "unexpected Len()")
	require.Equal(t, 2, b.NullN(), "unexpected NullN()")

	values := []*net.IPNet{
		mustParseInet("192.168.0.0/26"),
		mustParseInet("192.168.0.0/27"),
	}
	b.AppendValues(values, nil)

	require.Equal(t, 6, b.Len(), "unexpected Len()")

	a := b.NewArray()

	// check state of builder after NewInetBuilder
	require.Zero(t, b.Len(), "unexpected ArrayBuilder.Len(), did not reset state")
	require.Zero(t, b.Cap(), "unexpected ArrayBuilder.Cap(), did not reset state")
	require.Zero(t, b.NullN(), "unexpected ArrayBuilder.NullN(), did not reset state")

	require.Equal(t, `["192.168.0.0/24" (null) "192.168.0.0/25" (null) "192.168.0.0/26" "192.168.0.0/27"]`, a.String())
	st, err := a.MarshalJSON()
	require.NoError(t, err)

	b.Release()
	a.Release()

	b = NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
	err = b.UnmarshalJSON(st)
	require.NoError(t, err)

	a = b.NewArray()
	require.Equal(t, `["192.168.0.0/24" (null) "192.168.0.0/25" (null) "192.168.0.0/26" "192.168.0.0/27"]`, a.String())
	b.Release()
	a.Release()
}

func TestInetArray_ValueStr(t *testing.T) {
	// 1. create array
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
	defer b.Release()

	b.AppendNull()
	b.Append(mustParseInet("192.168.0.0/24"))
	b.AppendNull()
	b.Append(mustParseInet("192.168.0.0/25"))
	b.AppendNull()

	arr := b.NewInetArray()
	defer arr.Release()

	// 2. create array via AppendValueFromString
	b1 := NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
	defer b1.Release()

	for i := 0; i < arr.Len(); i++ {
		assert.NoError(t, b1.AppendValueFromString(arr.ValueStr(i)))
	}

	arr1 := b1.NewInetArray()
	defer arr1.Release()

	assert.Equal(t, arr.Len(), arr1.Len())
	for i := 0; i < arr.Len(); i++ {
		assert.Equal(t, arr.IsValid(i), arr1.IsValid(i))
		assert.Equal(t, arr.ValueStr(i), arr1.ValueStr(i))
	}
}

func TestInetBuilder_AppendValueFromString(t *testing.T) {
	// 1. create array
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
	defer b.Release()

	b.AppendNull()
	b.Append(mustParseInet("192.168.0.0/24"))
	b.AppendNull()
	b.Append(mustParseInet("192.168.0.0/25"))
	b.AppendNull()

	arr := b.NewInetArray()
	defer arr.Release()

	// 2. create array via AppendValueFromString
	b1 := NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
	defer b1.Release()

	for i := 0; i < arr.Len(); i++ {
		assert.NoError(t, b1.AppendValueFromString(arr.ValueStr(i)))
	}

	arr1 := b1.NewInetArray()
	defer arr1.Release()

	assert.Equal(t, arr.Len(), arr1.Len())
	for i := 0; i < arr.Len(); i++ {
		assert.Equal(t, arr.IsValid(i), arr1.IsValid(i))
		if arr.IsValid(i) {
			assert.Exactly(t, arr.Value(i), arr1.Value(i))
		}
	}
}
