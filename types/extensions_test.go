package types

import (
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
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

func TestStorageArrayConv(t *testing.T) {
	const amount = 100
	cases := []arrow.ExtensionType{ExtensionTypes.UUID, ExtensionTypes.MAC, ExtensionTypes.JSON, ExtensionTypes.Inet}
	for _, dt := range cases {
		t.Run(dt.String(), func(t *testing.T) {
			storageBuilder := array.NewBuilder(memory.DefaultAllocator, dt.StorageType())
			defer storageBuilder.Release()
			builder := array.NewBuilder(memory.DefaultAllocator, dt)
			defer builder.Release()

			for i := 0; i < amount; i++ {
				if i%2 == 0 {
					storageBuilder.AppendNull()
					builder.AppendNull()
					continue
				}
				storageBuilder.AppendEmptyValue()
				builder.AppendEmptyValue()
			}

			storage := storageBuilder.NewArray()
			defer storage.Release()
			arr := builder.NewArray().(array.ExtensionArray)
			defer arr.Release()

			// check matching
			assert.True(t, array.Equal(storage, arr.Storage()))

			// check that creating extension from storage matches
			fromStorage := array.NewExtensionArrayWithStorage(dt, storage)
			defer fromStorage.Release()

			assert.True(t, array.Equal(arr, fromStorage))

			// assert that no issues are in the fromStorage array
			for i := 0; i < fromStorage.Len(); i++ {
				assert.NotPanics(t, func() { fromStorage.ValueStr(i) })
				assert.Equal(t, arr.ValueStr(i), fromStorage.ValueStr(i))
			}
		})
	}
}
