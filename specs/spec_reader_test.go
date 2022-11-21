package specs

import (
	"bytes"
	"os"
	"path"
	"runtime"
	"testing"
)

type specLoaderTestCase struct {
	name         string
	path         []string
	err          func() string
	sources      int
	destinations int
}

func getPath(pathParts ...string) string {
	return path.Join("testdata", path.Join(pathParts...))
}

var specLoaderTestCases = []specLoaderTestCase{
	{
		name: "success",
		path: []string{getPath("gcp.yml"), getPath("dir")},
		err: func() string {
			return ""
		},
		sources:      2,
		destinations: 2,
	},
	{
		name: "duplicate_source",
		path: []string{getPath("gcp.yml"), getPath("gcp.yml")},
		err: func() string {
			return "duplicate source name gcp"
		},
	},
	{
		name: "no_such_file",
		path: []string{getPath("dir", "no_such_file.yml"), getPath("dir", "postgresql.yml")},
		err: func() string {
			if runtime.GOOS == "windows" {
				return "open testdata/dir/no_such_file.yml: The system cannot find the file specified."
			}
			return "open testdata/dir/no_such_file.yml: no such file or directory"
		},
	},
	{
		name: "duplicate_destination",
		path: []string{getPath("dir", "postgresql.yml"), getPath("dir", "postgresql.yml")},
		err: func() string {
			return "duplicate destination name postgresql"
		},
	},
	{
		name: "different_versions_for_destinations",
		path: []string{getPath("gcp.yml"), getPath("gcpv2.yml")},
		err: func() string {
			return "destination postgresqlv2 is used by multiple sources cloudquery/gcp with different versions"
		},
	},
	{
		name: "multiple sources success",
		path: []string{getPath("multiple_sources.yml")},
		err: func() string {
			return ""
		},
		sources:      2,
		destinations: 1,
	},
}

func TestLoadSpecs(t *testing.T) {
	for _, tc := range specLoaderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			specReader, err := NewSpecReader(tc.path)
			expectedErr := tc.err()
			if err != nil {
				if err.Error() != expectedErr {
					t.Fatalf("expected error: '%s', got: '%s'", expectedErr, err)
				}
				return
			}
			if expectedErr != "" {
				t.Fatalf("expected error: %s, got nil", expectedErr)
			}
			if len(specReader.Sources) != tc.sources {
				t.Fatalf("got: %d expected: %d", len(specReader.Sources), tc.sources)
			}
			if len(specReader.Destinations) != tc.destinations {
				t.Fatalf("got: %d expected: %d", len(specReader.Destinations), tc.destinations)
			}
		})
	}
}

func TestExpandFile(t *testing.T) {
	cfg := []byte(`
kind: source
spec:
	name: test
	version: v1.0.0
	spec:
		credentials: ${file:./testdata/creds.txt}
		otherstuff: 2
		credentials1: [${file:./testdata/creds.txt}, ${file:./testdata/creds1.txt}]
	`)
	expectedCfg := []byte(`
kind: source
spec:
	name: test
	version: v1.0.0
	spec:
		credentials: mytestcreds
		otherstuff: 2
		credentials1: [mytestcreds, anothercredtest]
	`)
	expandedCfg, err := expandFileConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expandedCfg, expectedCfg) {
		t.Fatalf("got: %s expected: %s", expandedCfg, expectedCfg)
	}

	badCfg := []byte(`
kind: source
spec:
	name: test
	version: v1.0.0
	spec:
		credentials: ${file:./testdata/creds2.txt}
		otherstuff: 2
	`)
	_, err = expandFileConfig(badCfg)
	if !os.IsNotExist(err) {
		t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
	}
}
