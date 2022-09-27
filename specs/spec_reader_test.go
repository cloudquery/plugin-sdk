package specs

import (
	"testing"
)

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader([]string{"testdata/gcp.yml", "testdata/dir"})
	if err != nil {
		t.Fatal(err)
	}
	if len(specReader.Sources) != 2 {
		t.Fatalf("got: %d expected: 1", len(specReader.Sources))
	}
	if len(specReader.Destinations) != 2 {
		t.Fatalf("got: %d expected: 2", len(specReader.Destinations))
	}

	_, err = NewSpecReader([]string{"testdata/gcp.yml", "testdata/gcp.yml"})
	if err != nil && err.Error() != "duplicate source name gcp" {
		t.Fatalf("got: %s expected: duplicate source name error", err)
	}
}
