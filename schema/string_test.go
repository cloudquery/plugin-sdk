package schema

import "testing"

func TestString(t *testing.T) {
	v := &String{}
	if err := v.Scan("hello"); err != nil {
		t.Fatal(err)
	}
}
