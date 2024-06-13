package scalar

import (
	"encoding/json"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapEncodeDecode(t *testing.T) {
	tl := []struct {
		name     string
		dt       *arrow.MapType
		input    any
		expected any
	}{
		{
			name:     "binary_to_string",
			dt:       arrow.MapOf(arrow.BinaryTypes.Binary, arrow.BinaryTypes.String),
			input:    `{"7049Ug==":"binary"}`,
			expected: `[{"key":"NzA0OVVnPT0=","value":"binary"}]`,
		},
		{
			name:     "binary_to_string_from_map",
			dt:       arrow.MapOf(arrow.BinaryTypes.Binary, arrow.BinaryTypes.String),
			input:    map[string]string{"7049Ug==": "binary"},
			expected: `[{"key":"NzA0OVVnPT0=","value":"binary"}]`,
		},
		{
			name:     "uuid_to_int",
			dt:       arrow.MapOf(types.ExtensionTypes.UUID, arrow.PrimitiveTypes.Int64),
			input:    `{"f81d4fae-7dec-11d0-a765-00a0c91e6bf6":123}`,
			expected: `[{"key":"f81d4fae-7dec-11d0-a765-00a0c91e6bf6","value":123}]`,
		},
		{
			name: "uuid_to_int_from_map",
			dt:   arrow.MapOf(types.ExtensionTypes.UUID, arrow.PrimitiveTypes.Int64),
			input: map[uuid.UUID]int64{
				uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"): 123,
			},
			expected: `[{"key":"f81d4fae-7dec-11d0-a765-00a0c91e6bf6","value":123}]`,
		},
	}

	for _, tc := range tl {
		t.Run(tc.name, func(t *testing.T) {
			bldr := array.NewBuilder(memory.DefaultAllocator, tc.dt)
			defer bldr.Release()

			s := NewScalar(tc.dt)
			require.Truef(t, arrow.TypeEqual(s.DataType(), tc.dt), "expected %v, got %v", tc.dt, s.DataType())

			require.NoError(t, s.Set(tc.input))
			assert.True(t, s.IsValid())
			AppendToBuilder(bldr, s)
			arr := bldr.NewArray()
			one := arr.GetOneForMarshal(0)
			data, err := json.Marshal(one)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(data))
		})
	}
}

func TestMap_Set(t *testing.T) {
	type testType any

	tl := []struct {
		schema *arrow.MapType
		input  any
	}{
		{
			schema: arrow.MapOf(arrow.BinaryTypes.String, arrow.PrimitiveTypes.Int64),
			input:  map[string]int64{"a": 1, "b": 2},
		},
		{
			schema: arrow.MapOf(arrow.BinaryTypes.String, arrow.PrimitiveTypes.Int64),
			input:  `{"a": 1, "b": 2}`,
		},
		{
			schema: arrow.MapOf(arrow.FixedWidthTypes.Float16, arrow.FixedWidthTypes.Boolean),
			input:  map[float32]bool{32: true, 47: false},
		},
		{
			schema: arrow.MapOf(arrow.FixedWidthTypes.Float16, arrow.FixedWidthTypes.Boolean),
			input:  `{"32": "true", "47": "false"}`,
		},
	}

	for _, tc := range tl {
		t.Run(tc.schema.String(), func(t *testing.T) {
			s := NewScalar(tc.schema)
			require.IsType(t, new(Map), s)
			require.NoError(t, s.Set(tc.input))
		})
	}
}
