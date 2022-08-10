package specs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// This test is testing both unmarshalling and LoadSpecs together
// so to add a new test case add a new file to testdata and an expected unmarshaled objectSpec
var (
	expectedSpecs = map[string]interface{}{
		"aws.cq.yml": SourceSpec{
			Name:          "aws",
			Path:          "cloudquery/aws",
			Registry:      "github",
			Version:       "1.0.0",
			MaxGoRoutines: 10,
		},
	}
	expectedAWSInnerSpec = map[string]interface{}{
		"regions": []interface{}{"us-east-1"},
	}
)

func TestLoadSpecs(t *testing.T) {
	specReader, err := NewSpecReader("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for fileName, spec := range specReader.sources {
		t.Run(fileName, func(t *testing.T) {
			innerSpec, _ := yaml.Marshal(spec.Spec)
			spec.Spec = yaml.Node{}
			assert.Equal(t, expectedSpecs[fileName], spec)

			if !t.Failed() {
				var innerSpecMap map[string]interface{}
				_ = yaml.Unmarshal(innerSpec, &innerSpecMap)
				assert.Equal(t, innerSpecMap, expectedAWSInnerSpec)
			}
		})
	}
}
