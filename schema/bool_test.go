package schema

import "testing"

type typeTestCase struct {
	value       any
	expected    any
	expectedErr error
}

func BoolTest(t *testing.T) {
	b := Bool{}
	if err := b.Scan(true); err != nil {
		t.Fatal(err)
	}
	if b.Bool {
		t.Fatalf("expected %t, got %t", true, b.Bool)
	}
}
