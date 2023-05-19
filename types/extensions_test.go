package types

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueStrRoundTrip(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	cases := []struct {
		arr     arrow.Array
		builder array.Builder
	}{
		{
			arr: func() arrow.Array {
				b := NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType()))
				defer b.Release()

				b.AppendNull()
				b.Append(mustParseInet("192.168.0.0/24"))
				b.AppendNull()
				b.Append(mustParseInet("192.168.0.0/25"))
				b.AppendNull()

				return b.NewInetArray()
			}(),
			builder: NewInetBuilder(array.NewExtensionBuilder(mem, NewInetType())),
		},
		{
			arr: func() arrow.Array {
				b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
				defer b.Release()

				b.AppendNull()
				b.Append(map[string]any{"a": 1, "b": 2})
				b.AppendNull()
				b.Append([]any{1, 2, 3})
				b.AppendNull()
				b.Append(map[string]any{"MyKey": "A\u0026B"})
				b.AppendNull()

				return b.NewJSONArray()
			}(),
			builder: NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType())),
		},
		{
			arr: func() arrow.Array {
				b := NewMACBuilder(array.NewExtensionBuilder(mem, NewMACType()))
				defer b.Release()

				b.AppendNull()
				b.Append(mustParseMAC("00:00:00:00:00:01"))
				b.AppendNull()
				b.Append(mustParseMAC("00:00:00:00:00:02"))
				b.AppendNull()

				return b.NewMACArray()
			}(),
			builder: NewMACBuilder(array.NewExtensionBuilder(mem, NewMACType())),
		},
		{
			arr: func() arrow.Array {
				b := NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType()))
				defer b.Release()

				b.AppendNull()
				b.Append(uuid.NameSpaceURL)
				b.AppendNull()
				b.Append(uuid.NameSpaceDNS)
				b.AppendNull()

				return b.NewUUIDArray()
			}(),
			builder: NewUUIDBuilder(array.NewExtensionBuilder(mem, NewUUIDType())),
		},
	}

	for _, tc := range cases {
		t.Run(tc.arr.DataType().(arrow.ExtensionType).ExtensionName(), func(t *testing.T) {
			defer tc.arr.Release()
			defer tc.builder.Release()
			t.Helper()

			for i := 0; i < tc.arr.Len(); i++ {
				assert.NoError(t, tc.builder.AppendValueFromString(tc.arr.ValueStr(i)))
			}

			arr := tc.builder.NewArray()
			defer arr.Release()

			require.True(t, array.Equal(tc.arr, arr))
		})
	}
}
