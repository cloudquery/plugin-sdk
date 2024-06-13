package scalar

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructEncodeDecode(t *testing.T) {
	oddBytes := make([]byte, 4)
	_, err := base64.StdEncoding.Decode(oddBytes, []byte("7049Ug=="))
	require.NoError(t, err)
	tl := []struct {
		name     string
		dt       *arrow.StructType
		input    any
		expected string
	}{
		{
			name:     "binary",
			dt:       arrow.StructOf(arrow.Field{Name: "binary", Type: arrow.BinaryTypes.Binary}),
			input:    `{"binary":"7049Ug=="}`,
			expected: `{"binary":"7049Ug=="}`,
		},
		{
			name: "binary_from_struct",
			dt:   arrow.StructOf(arrow.Field{Name: "binary", Type: arrow.BinaryTypes.Binary}),
			input: struct {
				Binary []byte `json:"binary"`
			}{Binary: oddBytes},
			expected: `{"binary":"7049Ug=="}`,
		},
		{
			name:     "uuid",
			dt:       arrow.StructOf(arrow.Field{Name: "uuid", Type: types.ExtensionTypes.UUID}),
			input:    `{"uuid":"f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}`,
			expected: `{"uuid":"f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}`,
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
			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(data))
		})
	}
}

func TestStructMissingKeys(t *testing.T) {
	tl := []struct {
		dt          *arrow.StructType
		input       any
		expectPanic bool
	}{
		{dt: arrow.StructOf(arrow.Field{Name: "i64", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "s", Type: arrow.BinaryTypes.String}), input: `{"i64": 1, "s": "foo"}`, expectPanic: false},
		{dt: arrow.StructOf(arrow.Field{Name: "i64", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "s", Type: arrow.BinaryTypes.String}), input: `{"i64": 1, "s": "foo", "extra":"bar"}`, expectPanic: true},
	}

	for idx, tc := range tl {
		tc := tc
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			panicked := false
			defer func() {
				if t.Failed() {
					return
				}

				if panicked && !tc.expectPanic {
					t.Errorf("unexpected panic")
				}
				if !panicked && tc.expectPanic {
					t.Errorf("expected panic")
				}
			}()

			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()

			bldr := array.NewBuilder(memory.DefaultAllocator, tc.dt)
			defer bldr.Release()

			s := NewScalar(tc.dt)
			if !arrow.TypeEqual(s.DataType(), tc.dt) {
				t.Fatalf("expected %v, got %v", tc.dt, s.DataType())
			}

			assert.NoError(t, s.Set(tc.input))
			if t.Failed() {
				return
			}

			assert.True(t, s.IsValid())
			AppendToBuilder(bldr, s)
		})
	}
}

func TestStructSet(t *testing.T) {
	type testType any

	tl := []struct {
		schema *arrow.StructType
		input  any
	}{
		{
			schema: arrow.StructOf(arrow.Field{Name: "a", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "b", Type: arrow.PrimitiveTypes.Int64}),
			input:  map[string]any{"a": 1, "b": 2},
		},
		{
			schema: arrow.StructOf(arrow.Field{Name: "a", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "b", Type: arrow.PrimitiveTypes.Int64}),
			input:  map[string]testType{"a": 1, "b": 2},
		},
		{
			schema: arrow.StructOf(arrow.Field{Name: "nested", Type: arrow.StructOf(arrow.Field{Name: "a", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "b", Type: arrow.PrimitiveTypes.Int64})}),
			input:  map[string]any{"x": map[string]any{"a": 1, "b": 2}},
		},
		{
			schema: arrow.StructOf(arrow.Field{Name: "nested", Type: arrow.StructOf(arrow.Field{Name: "a", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "b", Type: arrow.PrimitiveTypes.Int64})}),
			input:  map[string]testType{"x": map[string]testType{"a": 1, "b": 2}},
		},
		{
			schema: arrow.StructOf(arrow.Field{Name: "nested", Type: arrow.StructOf(arrow.Field{Name: "a", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "b", Type: arrow.PrimitiveTypes.Int64})}),
			input:  map[string]testType{"x": map[string]any{"a": 1, "b": 2}},
		},
		{
			schema: arrow.StructOf(arrow.Field{Name: "nested", Type: arrow.StructOf(arrow.Field{Name: "a", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "b", Type: arrow.PrimitiveTypes.Int64})}),
			input:  map[string]any{"x": map[string]testType{"a": 1, "b": 2}},
		},
	}

	for _, tc := range tl {
		t.Run(tc.schema.String(), func(t *testing.T) {
			s := NewScalar(tc.schema)
			assert.IsType(t, new(Struct), s)
			require.NoError(t, s.Set(tc.input))
		})
	}
}
