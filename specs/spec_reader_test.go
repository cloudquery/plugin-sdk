package specs

import (
	"testing"
)

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader("testdata/valid")
	if err != nil {
		t.Fatal(err)
	}
	if len(specReader.sources) != 1 {
		t.Fatalf("expected 1 source got %d", len(specReader.sources))
	}
	if len(specReader.destinations) != 1 {
		t.Fatalf("expected 1 destination got %d", len(specReader.destinations))
	}
}

func TestWrongKind(t *testing.T) {
	_, err := NewSpecReader("testdata/wrong_source")
	require.Equal(t, err.Error(), "failed to unmarshal file invalid.yml: failed to decode json: unknown kind test")
}

func TestNoSpecs(t *testing.T) {
	_, err := NewSpecReader("testdata")
	require.Equal(t, err.Error(), "no valid config files found in directory testdata")
}
