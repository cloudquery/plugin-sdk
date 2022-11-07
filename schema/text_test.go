package schema

import "testing"

func TestTextSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result Text
	}{
		{source: "foo", result: Text{Str: "foo", Status: Present}},
		{source: _string("bar"), result: Text{Str: "bar", Status: Present}},
		{source: (*string)(nil), result: Text{Status: Null}},
	}

	for i, tt := range successfulTests {
		var d Text
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}
