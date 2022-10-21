package schema

import "testing"

func TestInt64(t *testing.T) {
	v := &Int64{}
	if err := v.Scan(1); err != nil {
		t.Fatal(err)
	}
}
