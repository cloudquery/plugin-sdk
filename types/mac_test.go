package types

import (
	"net"
	"testing"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
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
