package spec

import (
	"reflect"
	"testing"
)

// This test is testing both unmarshalling and LoadSpecs together
// so to add a new test case add a new file to testdata and an expected unmarshaled objectSpec
var expectedSpecs = map[string]interface{}{
	"aws.cq.yml": SourceSpec{
		Name:          "aws",
		Version:       "1.0.0",
		MaxGoRoutines: 10,
	},
}

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for fileName, spec := range specReader.sources {
		t.Run(fileName, func(t *testing.T) {
			if !reflect.DeepEqual(spec, expectedSpecs[fileName]) {
				t.Fatalf("expected %v, got %v", expectedSpecs[fileName], spec)
			}
		})
	}
}
