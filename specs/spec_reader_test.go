package specs

import (
	"testing"
)

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if len(specReader.sources) != 1 {
		t.Fatalf("got: %d expected: 1", len(specReader.sources))
	}
	if len(specReader.destinations) != 1 {
		t.Fatalf("got: %d expected: 1", len(specReader.destinations))
	}
}
