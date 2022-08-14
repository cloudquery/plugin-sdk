package specs

import (
	"embed"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

//go:embed testdata/*.cq.yml
var testSpecsFS embed.FS

var testSpecs = map[string]Spec{
	"testdata/pg.cq.yml": {
		Kind: "destination",
		Spec: &DestinationSpec{
			Name:      "postgresql",
			Version:   "v1.0.0",
			Registry:  RegistryGrpc,
			WriteMode: WriteModeOverwrite,
		},
	},
	"testdata/aws.cq.yml": {
		Kind: "source",
		Spec: &SourceSpec{
			Name:          "aws",
			Path:          "aws",
			Version:       "v1.0.0",
			MaxGoRoutines: 10,
			Registry:      RegistryLocal,
		},
	},
}

func TestSpecYamlMarshal(t *testing.T) {
	for fileName, expectedSpec := range testSpecs {
		t.Run(fileName, func(t *testing.T) {
			b, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatal(err)
			}
			var spec Spec
			if err := yaml.Unmarshal(b, &spec); err != nil {
				t.Fatal(err)
			}
			b, err = yaml.Marshal(spec)
			if err != nil {
				t.Fatal(err)
			}
			if err := yaml.Unmarshal(b, &spec); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(spec, expectedSpec) {
				t.Errorf("expected spec %s to be:\n%v\nbut got:\n%v", fileName, expectedSpec.Spec, spec.Spec)
			}
		})
	}
}
