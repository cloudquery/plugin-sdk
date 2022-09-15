package specs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var sources = map[string]Source{
	"aws.cq.yml": {
		Name:          "aws",
		Path:          "aws",
		Version:       "v1.0.0",
		MaxGoRoutines: 10,
		Registry:      RegistryLocal,
	},
}

var destinations = map[string]Destination{
	"pg.cq.yml": {
		Name:      "postgresql",
		Path:      "postgresql",
		Version:   "v1.0.0",
		Registry:  RegistryGrpc,
		WriteMode: WriteModeOverwrite,
	},
}

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader("testdata")
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, sources, specReader.sources)
	require.Equal(t, destinations, specReader.destinations)
}
