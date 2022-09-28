package specs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testDestinationSpec struct {
	ConnectionString string `json:"connection_string"`
}

func TestWriteModeFromString(t *testing.T) {
	var writeMode WriteMode
	if err := writeMode.UnmarshalJSON([]byte(`"append"`)); err != nil {
		t.Fatal(err)
	}
	if writeMode != WriteModeAppend {
		t.Fatalf("expected WriteModeAppend, got %v", writeMode)
	}
	if err := writeMode.UnmarshalJSON([]byte(`"overwrite"`)); err != nil {
		t.Fatal(err)
	}
	if writeMode != WriteModeOverwrite {
		t.Fatalf("expected WriteModeOverwrite, got %v", writeMode)
	}
}

func TestDestinationSpecUnmarshalSpec(t *testing.T) {
	destination := Destination{
		Spec: map[string]interface{}{
			"connection_string": "postgres://user:pass@host:port/db",
		},
	}
	var spec testDestinationSpec
	if err := destination.UnmarshalSpec(&spec); err != nil {
		t.Fatal(err)
	}
	if spec.ConnectionString != "postgres://user:pass@host:port/db" {
		t.Fatalf("expected postgres://user:pass@host:port/db, got %s", spec.ConnectionString)
	}
}

var destinationUnmarshalSpecTestCases = []struct {
	name   string
	spec   string
	err    string
	source *Source
}{
	{
		"invalid_kind",
		`kind: nice`,
		"failed to decode spec: unknown kind nice",
		nil,
	},
	{
		"invalid_type",
		`kind: source
spec:
  name: 3
`,
		"failed to decode spec: json: cannot unmarshal number into Go struct field Source.name of type string",
		&Source{
			Name:   "test",
			Tables: []string{"*"},
		},
	},
	{
		"unknown_field",
		`kind: source
spec:
  namea: 3
`,
		`failed to decode spec: json: unknown field "namea"`,
		&Source{
			Name:   "test",
			Tables: []string{"*"},
		},
	},
}

func TestDestinationUnmarshalSpec(t *testing.T) {
	for _, tc := range destinationUnmarshalSpecTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var spec Spec
			err = SpecUnmarshalYamlStrict([]byte(tc.spec), &spec)
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("expected:%s got:%s", tc.err, err.Error())
				}
				return
			}

			source := spec.Spec.(*Source)
			if cmp.Diff(source, tc.source) != "" {
				t.Fatalf("expected:%v got:%v", tc.source, source)
			}
		})
	}
}

var destinationUnmarshalSpecValidateTestCases = []struct {
	name        string
	spec        string
	err         string
	destination *Destination
}{
	{
		"required_name",
		`kind: destination
spec:`,
		"name is required",
		nil,
	},
	{
		"required_version",
		`kind: destination
spec:
  name: test
`,
		"version is required",
		nil,
	},
	{
		"required_version_format",
		`kind: destination
spec:
  name: test
  version: 1.1.0
`,
		"version must start with v",
		nil,
	},
	{
		"version_is_not_required_for_grpc_registry",
		`kind: destination
spec:
  name: test
  registry: grpc
  path: "localhost:9999"
`,
		"",
		&Destination{
			Name:     "test",
			Registry: RegistryGrpc,
			Path:     "localhost:9999",
		},
	},
	{
		"version_is_not_required_for_local_registry",
		`kind: destination
spec:
  name: test
  registry: local
  path: "/home/user/some_executable"
`,
		"",
		&Destination{
			Name:     "test",
			Registry: RegistryLocal,
			Path:     "/home/user/some_executable",
		},
	},
	{
		"success",
		`kind: destination
spec:
  name: test
  version: v1.1.0
`,
		"",
		&Destination{
			Name:     "test",
			Registry: RegistryGithub,
			Path:     "cloudquery/test",
			Version:  "v1.1.0",
		},
	},
}

func TestDestinationUnmarshalSpecValidate(t *testing.T) {
	for _, tc := range destinationUnmarshalSpecValidateTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var spec Spec
			err = SpecUnmarshalYamlStrict([]byte(tc.spec), &spec)
			if err != nil {
				t.Fatal(err)
			}
			destination := spec.Spec.(*Destination)
			destination.SetDefaults()
			err = destination.Validate()
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("expected:%s got:%s", tc.err, err.Error())
				}
				return
			}

			if cmp.Diff(destination, tc.destination) != "" {
				t.Fatalf("expected:%v got:%v", tc.destination, destination)
			}
		})
	}
}
