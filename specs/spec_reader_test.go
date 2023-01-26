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
	envVariables map[string]string
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
		name: "success_yaml_extension",
		path: []string{getPath("gcp.yml"), getPath("dir_yaml")},
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
	{
		name: "environment variables",
		path: []string{getPath("env_variables.yml")},
		err: func() string {
			return ""
		},
		sources:      2,
		destinations: 1,
		envVariables: map[string]string{
			"VERSION":           "v1",
			"DESTINATIONS":      "postgresql",
			"CONNECTION_STRING": "postgresql://localhost:5432/cloudquery?sslmode=disable",
		},
	},
	{
		name: "environment variables with error",
		path: []string{getPath("env_variables.yml")},
		err: func() string {
			return "failed to expand environment variable in file testdata/env_variables.yml (section 3): env variable CONNECTION_STRING not found"
		},
		sources:      2,
		destinations: 1,
		envVariables: map[string]string{
			"VERSION":      "v1",
			"DESTINATIONS": "postgresql",
		},
	},
	{
		name: "environment variables in string without error",
		path: []string{getPath("env_variable_in_string.yml")},
		err: func() string {
			return ""
		},
		sources:      1,
		destinations: 1,
		envVariables: map[string]string{
			"VERSION": "v1",
		},
	},
	{
		name: "environment variables in string with error",
		path: []string{getPath("env_variable_in_string.yml")},
		err: func() string {
			return "failed to expand environment variable in file testdata/env_variable_in_string.yml (section 2): env variable VERSION not found"
		},
		sources:      1,
		destinations: 1,
		envVariables: map[string]string{},
	},
}

func TestLoadSpecs(t *testing.T) {
	for _, tc := range specLoaderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envVariables {
				t.Setenv(k, v)
			}
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

func TestExpandEnv(t *testing.T) {
	os.Setenv("TEST_ENV_CREDS", "mytestcreds")
	os.Setenv("TEST_ENV_CREDS2", "anothercredtest")
	cfg := []byte(`
kind: source
spec:
	name: test
	version: v1.0.0
	spec:
		credentials: ${TEST_ENV_CREDS}
		otherstuff: 2
		credentials1: [${TEST_ENV_CREDS}, ${TEST_ENV_CREDS2}]
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
	expandedCfg, err := expandEnv(cfg)
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
		credentials: ${TEST_ENV_CREDS1}
		otherstuff: 2
	`)
	_, err = expandEnv(badCfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
