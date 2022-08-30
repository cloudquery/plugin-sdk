package specs

import (
	"testing"
)

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

func TestDestinationSetDefaults(t *testing.T) {
	destination := Destination{
		Name: "testDestination",
	}
	destination.SetDefaults()
	if destination.Registry != RegistryGithub {
		t.Fatalf("expected RegistryGithub, got %v", destination.Registry)
	}
	if destination.Path != "cloudquery/testDestination" {
		t.Fatalf("expected destination.Path (%s), got %s", destination.Name, destination.Path)
	}
	if destination.Version != "latest" {
		t.Fatalf("expected latest, got %s", destination.Version)
	}
}

type testDestinationSpec struct {
	ConnectionString string `json:"connection_string"`
}

func TestDestinationUnmarshalSpec(t *testing.T) {
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
