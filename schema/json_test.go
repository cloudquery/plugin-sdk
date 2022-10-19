package schema

import "testing"

func TestJson(t *testing.T) {
	v := &Json{}
	if err := v.Scan([]byte(`{"foo": "bar"}`)); err != nil {
		t.Fatal(err)
	}
}