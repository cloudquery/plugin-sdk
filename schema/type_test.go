package schema

import "testing"

func TestCQTypeFromSchema(t *testing.T) {
	// check that all values are registered in CQTypeFromSchema function
	for i := TypeBool; i < TypeEnd; i++ {
		if CQTypeFromSchema(i) == nil {
			t.Fatalf("missing CQType for schema type %d", i)
		}
	}
}
