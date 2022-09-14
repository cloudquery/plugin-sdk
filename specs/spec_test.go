package specs

import (
	"os"
	"reflect"
	"testing"
)

var testSpecs = map[string]Spec{
	"testdata/pg.cq.yml": {
		Kind: KindDestination,
		Spec: &Destination{
			Name:      "postgresql",
			Path:      "postgresql",
			Version:   "v1.0.0",
			Registry:  RegistryGrpc,
			WriteMode: WriteModeOverwrite,
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
			if err := SpecUnmarshalYamlStrict(b, &spec); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(spec, expectedSpec) {
				t.Errorf("expected spec %s to be:\n%+v\nbut got:\n%+v", fileName, expectedSpec.Spec, spec.Spec)
			}
		})
	}
}
