package scalar

import "testing"

type Foo struct {
	Num int
}

func TestJSONSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result JSON
	}{
		{source: "", result: JSON{Value: []byte("")}},
		{source: "{}", result: JSON{Value: []byte("{}"), Valid: true}},
		{source: `"test"`, result: JSON{Value: []byte(`"test"`), Valid: true}},
		{source: "1", result: JSON{Value: []byte("1"), Valid: true}},
		{source: "[1, 2, 3]", result: JSON{Value: []byte("[1, 2, 3]"), Valid: true}},
		{source: []byte("{}"), result: JSON{Value: []byte("{}"), Valid: true}},
		{source: []byte(`"test"`), result: JSON{Value: []byte(`"test"`), Valid: true}},
		{source: []byte("1"), result: JSON{Value: []byte("1"), Valid: true}},
		{source: []byte("[1, 2, 3]"), result: JSON{Value: []byte("[1, 2, 3]"), Valid: true}},
		{source: ([]byte)(nil), result: JSON{}},
		{source: (*string)(nil), result: JSON{}},

		{source: []int{1, 2, 3}, result: JSON{Value: []byte("[1,2,3]"), Valid: true}},
		{source: []int(nil), result: JSON{Value: []byte(`[]`), Valid: true}},
		{source: []int{}, result: JSON{Value: []byte(`[]`), Valid: true}},
		{source: []Foo(nil), result: JSON{Value: []byte(`[]`), Valid: true}},
		{source: []Foo{}, result: JSON{Value: []byte(`[]`), Valid: true}},
		{source: []Foo{{1}}, result: JSON{Value: []byte(`[{"Num":1}]`), Valid: true}},

		{source: map[string]any{"foo": "bar"}, result: JSON{Value: []byte(`{"foo":"bar"}`), Valid: true}},
		{source: map[string]any(nil), result: JSON{Value: []byte(`{}`), Valid: true}},
		{source: map[string]any{}, result: JSON{Value: []byte(`{}`), Valid: true}},
		{source: map[string]string{"foo": "bar"}, result: JSON{Value: []byte(`{"foo":"bar"}`), Valid: true}},
		{source: map[string]string(nil), result: JSON{Value: []byte(`{}`), Valid: true}},
		{source: map[string]string{}, result: JSON{Value: []byte(`{}`), Valid: true}},
		{source: map[string]Foo{"foo": {1}}, result: JSON{Value: []byte(`{"foo":{"Num":1}}`), Valid: true}},
		{source: map[string]Foo(nil), result: JSON{Value: []byte(`{}`), Valid: true}},
		{source: map[string]Foo{}, result: JSON{Value: []byte(`{}`), Valid: true}},

		{source: nil, result: JSON{}},

		{source: map[string]any{"test1": "a&b", "test2": "😀"}, result: JSON{Value: []byte(`{"test1": "a&b", "test2": "😀"}`), Valid: true}},
		{source: &JSON{Value: []byte(`{"test1": "a&b", "test2": "😀"}`), Valid: true}, result: JSON{Value: []byte(`{"test1": "a&b", "test2": "😀"}`), Valid: true}},
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
