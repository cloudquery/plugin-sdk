package specs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSpecs(t *testing.T) {
	var tests = []struct {
		files              []string
		expectSources      int
		expectDestinations int
		expectError        bool
		expectWarning      bool
	}{
		{
			files:              []string{"testdata/gcp.yml", "testdata/dir"},
			expectSources:      2,
			expectDestinations: 2,
		},
		{
			files:              []string{"testdata/gcp_deprecated.yml"},
			expectWarning:      true,
			expectSources:      1,
			expectDestinations: 1,
		},
	}
	for _, tc := range tests {
		specReader, err := NewSpecReader(tc.files)
		if tc.expectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		if !tc.expectWarning {
			for _, w := range specReader.Warnings {
				assert.Empty(t, w)
			}
		} else {
			for _, w := range specReader.Warnings {
				assert.NotEmpty(t, w)
			}
		}

		assert.Equal(t, tc.expectSources, len(specReader.Sources))
		assert.Equal(t, tc.expectDestinations, len(specReader.Destinations))
	}

	_, err := NewSpecReader([]string{"testdata/gcp.yml", "testdata/gcp.yml"})
	if err != nil && err.Error() != "duplicate source name gcp" {
		t.Fatalf("got: %s expected: duplicate source name error", err)
	}
}
