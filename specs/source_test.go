package specs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var sourceUnmarshalSpecTestCases = []struct {
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

func TestSourceUnmarshalSpec(t *testing.T) {
	for _, tc := range sourceUnmarshalSpecTestCases {
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

var sourceUnmarshalSpecValidateTestCases = []struct {
	name   string
	spec   string
	err    string
	source *Source
}{
	{
		"required_name",
		`kind: source
spec:`,
		"name is required",
		nil,
	},
	{
		"required_version",
		`kind: source
spec:
  name: test
`,
		"version is required",
		nil,
	},
	{
		"required_version_format",
		`kind: source
spec:
  name: test
  version: 1.1.0
`,
		"version must start with v",
		nil,
	},
	{
		"destination_required",
		`kind: source
spec:
  name: test
  version: v1.1.0
`,
		"at least one destination is required",
		nil,
	},
	{
		"success",
		`kind: source
spec:
  name: test
  version: v1.1.0
  destinations: ["test"]
`,
		"",
		&Source{
			Name:         "test",
			Registry:     RegistryGithub,
			Path:         "cloudquery/test",
			Concurrency:  defaultConcurrency,
			Version:      "v1.1.0",
			Tables:       []string{"*"},
			Destinations: []string{"test"},
		},
	},
}

func TestSourceUnmarshalSpecValidate(t *testing.T) {
	for _, tc := range sourceUnmarshalSpecValidateTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var spec Spec
			err = SpecUnmarshalYamlStrict([]byte(tc.spec), &spec)
			if err != nil {
				t.Fatal(err)
			}
			source := spec.Spec.(*Source)
			source.SetDefaults()
			err = source.Validate()
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("expected:%s got:%s", tc.err, err.Error())
				}
				return
			}

			if cmp.Diff(source, tc.source) != "" {
				t.Fatalf("expected:%v got:%v", tc.source, source)
			}
		})
	}
}
