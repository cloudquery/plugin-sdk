package types

import (
	"net"
	"testing"

	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/stretchr/testify/require"
)

func mustParseInet(s string) *net.IPNet {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	ipnet.IP = ip
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
	b.Append(mustParseInet("192.168.0.1/24"))

	require.Equal(t, 7, b.Len(), "unexpected Len()")

	a := b.NewArray()

	// check state of builder after NewInetBuilder
	require.Zero(t, b.Len(), "unexpected ArrayBuilder.Len(), did not reset state")
	require.Zero(t, b.Cap(), "unexpected ArrayBuilder.Cap(), did not reset state")
	require.Zero(t, b.NullN(), "unexpected ArrayBuilder.NullN(), did not reset state")

	require.Equal(t, `["192.168.0.0/24" (null) "192.168.0.0/25" (null) "192.168.0.0/26" "192.168.0.0/27" "192.168.0.1/24"]`, a.String())
	st, err := a.MarshalJSON()
	require.NoError(t, err)

	b.Release()
	a.Release()

	b = NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
	err = b.UnmarshalJSON(st)
	require.NoError(t, err)

	a = b.NewArray()
	require.Equal(t, `["192.168.0.0/24" (null) "192.168.0.0/25" (null) "192.168.0.0/26" "192.168.0.0/27" "192.168.0.1/24"]`, a.String())
	b.Release()
	a.Release()
}
