package cqtypes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJSONSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result JSON
	}{
		{source: "{}", result: JSON{Bytes: []byte("{}"), Status: Present}},
		{source: []byte("{}"), result: JSON{Bytes: []byte("{}"), Status: Present}},
		{source: ([]byte)(nil), result: JSON{Status: Null}},
		{source: (*string)(nil), result: JSON{Status: Null}},
		{source: []int{1, 2, 3}, result: JSON{Bytes: []byte("[1,2,3]"), Status: Present}},
		{source: map[string]interface{}{"foo": "bar"}, result: JSON{Bytes: []byte(`{"foo":"bar"}`), Status: Present}},
	}

	for i, tt := range successfulTests {
		var d JSON
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if diff := cmp.Diff(d, tt.result); diff != "" {
			t.Errorf("%d: got diff:\n%s", i, diff)
		}
	}
}
