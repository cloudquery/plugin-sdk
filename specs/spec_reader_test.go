package specs

import (
	"testing"
)

type specLoaderTestCase struct {
	name     string
	path []string
	err string
	sources int
	destinations int
}

var specLoaderTestCases = []specLoaderTestCase{
	{
		name: "sucess",
		path: []string{"testdata/gcp.yml", "testdata/dir"},
		err: "",
		sources: 2,
		destinations: 2,
	},
	{
		name: "duplicate_source",
		path: []string{"testdata/gcp.yml", "testdata/gcp.yml"},
		err: "duplicate source name gcp",
	},
	{
		name: "no_such_file",
		path: []string{"testdata/dir/no_such_file.yml", "testdata/dir/postgresql.yml"},
		err: "open testdata/dir/no_such_file.yml: no such file or directory",
	},
	{
		name: "duplicate_destination",
		path: []string{"testdata/dir/postgresql.yml", "testdata/dir/postgresql.yml"},
		err: "duplicate destination name postgresql",
	},
	{
		name: "different_versions_for_destinations",
		path: []string{"testdata/gcp.yml", "testdata/gcpv2.yml"},
		err: "destination postgresqlv2 is used by multiple sources cloudquery/gcp with different versions",
	},
}

func TestLoadSpecs(t *testing.T) {
	for _, tc := range specLoaderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			specReader, err := NewSpecReader(tc.path)
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("expected error: '%s', got: '%s'", tc.err, err)
				}
				return
			}
			if tc.err != "" {
				t.Fatalf("expected error: %s, got nil", tc.err)
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
