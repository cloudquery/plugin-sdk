package schema

import (
	"testing"
)

type Foo struct {
	Num int
}

func TestJSONSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result JSON
	}{
		{source: "", result: JSON{Bytes: []byte(""), Status: Null}},
		{source: "{}", result: JSON{Bytes: []byte("{}"), Status: Present}},
		{source: `"test"`, result: JSON{Bytes: []byte(`"test"`), Status: Present}},
		{source: "1", result: JSON{Bytes: []byte("1"), Status: Present}},
		{source: "[1, 2, 3]", result: JSON{Bytes: []byte("[1, 2, 3]"), Status: Present}},
		{source: []byte("{}"), result: JSON{Bytes: []byte("{}"), Status: Present}},
		{source: []byte(`"test"`), result: JSON{Bytes: []byte(`"test"`), Status: Present}},
		{source: []byte("1"), result: JSON{Bytes: []byte("1"), Status: Present}},
		{source: []byte("[1, 2, 3]"), result: JSON{Bytes: []byte("[1, 2, 3]"), Status: Present}},
		{source: ([]byte)(nil), result: JSON{Status: Null}},
		{source: (*string)(nil), result: JSON{Status: Null}},

		{source: []int{1, 2, 3}, result: JSON{Bytes: []byte("[1,2,3]"), Status: Present}},
		{source: []int(nil), result: JSON{Bytes: []byte(`[]`), Status: Present}},
		{source: []int{}, result: JSON{Bytes: []byte(`[]`), Status: Present}},
		{source: []Foo(nil), result: JSON{Bytes: []byte(`[]`), Status: Present}},
		{source: []Foo{}, result: JSON{Bytes: []byte(`[]`), Status: Present}},
		{source: []Foo{{1}}, result: JSON{Bytes: []byte(`[{"Num":1}]`), Status: Present}},

		{source: map[string]any{"foo": "bar"}, result: JSON{Bytes: []byte(`{"foo":"bar"}`), Status: Present}},
		{source: map[string]any(nil), result: JSON{Bytes: []byte(`{}`), Status: Present}},
		{source: map[string]any{}, result: JSON{Bytes: []byte(`{}`), Status: Present}},
		{source: map[string]string{"foo": "bar"}, result: JSON{Bytes: []byte(`{"foo":"bar"}`), Status: Present}},
		{source: map[string]string(nil), result: JSON{Bytes: []byte(`{}`), Status: Present}},
		{source: map[string]string{}, result: JSON{Bytes: []byte(`{}`), Status: Present}},
		{source: map[string]Foo{"foo": {1}}, result: JSON{Bytes: []byte(`{"foo":{"Num":1}}`), Status: Present}},
		{source: map[string]Foo(nil), result: JSON{Bytes: []byte(`{}`), Status: Present}},
		{source: map[string]Foo{}, result: JSON{Bytes: []byte(`{}`), Status: Present}},

		{source: nil, result: JSON{Status: Null}},

		{source: map[string]any{"test1": "a&b", "test2": "😀"}, result: JSON{Bytes: []byte(`{"test1": "a&b", "test2": "😀"}`), Status: Present}},
	}

	for i, tt := range successfulTests {
		var d JSON
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !d.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, d, tt.result)
		}
	}
}

func TestJSON_Size(t *testing.T) {
	tests := []struct {
		name string
		j    JSON
		want int
	}{
		{
			name: "empty",
			j:    JSON{Status: Null},
			want: 0,
		},
		{
			name: "present",
			j:    JSON{Bytes: []byte(`{"foo":"bar"}`), Status: Present},
			want: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.Size(); got != tt.want {
				t.Errorf("JSON.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
