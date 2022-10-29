package cqtypes

import (
	"testing"
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

		if !d.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, d, tt.result)
		}
	}
}
