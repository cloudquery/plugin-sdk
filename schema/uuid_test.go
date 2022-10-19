package schema

import "testing"

func TestUUID(t *testing.T) {
	v := &UUID{}

	if err := v.Scan("10000000-0000-0000-0000-000000000000"); err != nil {
		t.Fatal(err)
	}
	
}