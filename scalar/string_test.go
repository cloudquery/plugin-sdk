package scalar

import "testing"

func TestStringSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result String
	}{
		{source: "foo", result: String{Value: "foo", Valid: true}},
		{source: _string("bar"), result: String{Value: "bar", Valid: true}},
		{source: (*string)(nil), result: String{}},
	}

	for i, tt := range successfulTests {
		var d String
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}
