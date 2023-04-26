package types

import (
	"bytes"
	"testing"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestJSONBuilder(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer mem.AssertSize(t, 0)

	b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
	b.Append(map[string]any{"a": 1, "b": 2})
	b.AppendNull()
	b.Append(map[string]any{"c": 3, "d": 4})
	b.AppendNull()

	require.Equal(t, 4, b.Len(), "unexpected Len()")
	require.Equal(t, 2, b.NullN(), "unexpected NullN()")

	values := []any{
		map[string]any{"e": 5, "f": 6},
		map[string]any{"g": 7, "h": 8},
	}
	b.AppendValues(values, []bool{true, true})

	require.Equal(t, 6, b.Len(), "unexpected Len()")

	a := b.NewArray()

	// check state of builder after NewJSONBuilder
	require.Zero(t, b.Len(), "unexpected ArrayBuilder.Len(), NewJSONBuilder did not reset state")
	require.Zero(t, b.Cap(), "unexpected ArrayBuilder.Cap(), NewJSONBuilder did not reset state")
	require.Zero(t, b.NullN(), "unexpected ArrayBuilder.NullN(), NewJSONBuilder did not reset state")
	require.Equal(t, `["{\"a\":1,\"b\":2}" (null) "{\"c\":3,\"d\":4}" (null) "{\"e\":5,\"f\":6}" "{\"g\":7,\"h\":8}"]`, a.String())
	st, err := a.MarshalJSON()
	require.NoError(t, err)

	b.Release()
	a.Release()

	b = NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
	err = b.UnmarshalJSON(st)
	require.NoError(t, err)

	a = b.NewArray()
	require.Equal(t, `["{\"a\":1,\"b\":2}" (null) "{\"c\":3,\"d\":4}" (null) "{\"e\":5,\"f\":6}" "{\"g\":7,\"h\":8}"]`, a.String())
	b.Release()
	a.Release()
}

func TestJSONBuilder_UnmarshalOne(t *testing.T) {
	cases := []struct {
		name string
		data string
		want string
	}{
		{
			name: `map`,
			data: `{"a": 1, "b": 2}`,
			want: `["{\"a\":1,\"b\":2}"]`,
		},
		{
			name: `two maps`,
			data: `{"a": 1, "b": 2}{"c": 3, "d": 4}`,
			want: `["{\"a\":1,\"b\":2}"]`,
		},
		{
			name: `array`,
			data: `[1, 2, 3]`,
			want: `["[1,2,3]"]`,
		},
		{
			name: `two arrays`,
			data: `[1, 2, 3][4, 5, 6]`,
			want: `["[1,2,3]"]`,
		},
		{
			name: `null`,
			data: `null`,
			want: `[(null)]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
			defer mem.AssertSize(t, 0)
			b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
			defer b.Release()
			dec := json.NewDecoder(bytes.NewReader([]byte(tc.data)))
			err := b.UnmarshalOne(dec)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			a := b.NewArray()
			defer a.Release()
			require.Equal(t, tc.want, a.String())
		})
	}
}

func TestJSONArray_GetOneForMarshal(t *testing.T) {
	cases := []struct {
		name string
		data string
		want json.RawMessage
		nil  bool
	}{
		{
			name: `map`,
			data: `{"a": 1, "b": 2}`,
			want: json.RawMessage(`{"a":1,"b":2}`),
		},
		{
			name: `two maps`,
			data: `{"a": 1, "b": 2}{"c": 3, "d": 4}`,
			want: json.RawMessage(`{"a":1,"b":2}`),
		},
		{
			name: `array`,
			data: `[1, 2, 3]`,
			want: json.RawMessage(`[1,2,3]`),
		},
		{
			name: `two arrays`,
			data: `[1, 2, 3][4, 5, 6]`,
			want: json.RawMessage(`[1,2,3]`),
		},
		{
			name: `null`,
			data: `null`,
			nil:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
			defer mem.AssertSize(t, 0)
			b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
			defer b.Release()
			dec := json.NewDecoder(bytes.NewReader([]byte(tc.data)))
			err := b.UnmarshalOne(dec)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			a := b.NewArray()
			defer a.Release()
			if tc.nil {
				require.Nil(t, a.GetOneForMarshal(0))
			} else {
				require.Exactly(t, tc.want, a.GetOneForMarshal(0))
			}
		})
	}
}

func TestJSONArray_ValueStr(t *testing.T) {
	cases := []struct {
		name string
		data string
		want string
		nil  bool
	}{
		{
			name: `map`,
			data: `{"a": 1, "b": 2}`,
			want: `{"a":1,"b":2}`,
		},
		{
			name: `two maps`,
			data: `{"a": 1, "b": 2}{"c": 3, "d": 4}`,
			want: `{"a":1,"b":2}`,
		},
		{
			name: `array`,
			data: `[1, 2, 3]`,
			want: `[1,2,3]`,
		},
		{
			name: `two arrays`,
			data: `[1, 2, 3][4, 5, 6]`,
			want: `[1,2,3]`,
		},
		{
			name: `null`,
			data: `null`,
			want: array.NullValueStr,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
			defer mem.AssertSize(t, 0)
			b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
			defer b.Release()
			dec := json.NewDecoder(bytes.NewReader([]byte(tc.data)))
			err := b.UnmarshalOne(dec)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			a := b.NewArray()
			defer a.Release()
			require.Exactly(t, tc.want, a.ValueStr(0))
		})
	}
}

func TestJSONArray_Value(t *testing.T) {
	cases := []struct {
		name string
		data string
		want any
	}{
		{
			name: `map`,
			data: `{"a": 1, "b": 2}`,
			want: map[string]any{"a": float64(1), "b": float64(2)},
		},
		{
			name: `two maps`,
			data: `{"a": 1, "b": 2}{"c": 3, "d": 4}`,
			want: map[string]any{"a": float64(1), "b": float64(2)},
		},
		{
			name: `array`,
			data: `[1, 2, 3]`,
			want: []any{float64(1), float64(2), float64(3)},
		},
		{
			name: `two arrays`,
			data: `[1, 2, 3][4, 5, 6]`,
			want: []any{float64(1), float64(2), float64(3)},
		},
		{
			name: `null`,
			data: `null`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
			defer mem.AssertSize(t, 0)
			b := NewJSONBuilder(array.NewExtensionBuilder(mem, NewJSONType()))
			defer b.Release()
			dec := json.NewDecoder(bytes.NewReader([]byte(tc.data)))
			err := b.UnmarshalOne(dec)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			a := b.NewArray().(*JSONArray)
			defer a.Release()
			require.Equal(t, tc.want, a.Value(0))
		})
	}
}
