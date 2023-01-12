package specs

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestIncrementalJsonMarshalUnmarshal(t *testing.T) {
	b, err := json.Marshal(IncrementalBoth)
	if err != nil {
		t.Fatal("failed to marshal:", err)
	}
	var incremental Incremental
	if err := json.Unmarshal(b, &incremental); err != nil {
		t.Fatal("failed to unmarshal:", err)
	}
	if incremental != IncrementalBoth {
		t.Fatal("expected incremental to be both, but got:", incremental)
	}
}

func TestIncrementalYamlMarshalUnmarsahl(t *testing.T) {
	b, err := yaml.Marshal(IncrementalBoth)
	if err != nil {
		t.Fatal("failed to marshal:", err)
	}
	var incremental Incremental
	if err := yaml.Unmarshal(b, &incremental); err != nil {
		t.Fatal("failed to unmarshal:", err)
	}
	if incremental != IncrementalBoth {
		t.Fatal("expected registry to be `both`, but got:", IncrementalBoth)
	}
}
