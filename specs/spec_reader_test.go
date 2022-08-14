package specs

import (
	"reflect"
	"testing"
)

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for fileName, spec := range specReader.sources {
		t.Run(fileName, func(t *testing.T) {
			if !reflect.DeepEqual(spec, testSpecs[fileName]) {
				t.Errorf("expected spec %s to be:\n%v\nbut got:\n%v", fileName, testSpecs[fileName].Spec, spec.Spec)
			}
		})
	}
}
