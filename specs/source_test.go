package specs

import (
	"strings"
	"testing"
)

type testSourceSpec struct {
	Accounts []string `json:"accounts"`
}

func TestSourceSetDefaults(t *testing.T) {
	source := Source{
		Name: "testSource",
	}
	source.SetDefaults()
	if source.Registry != RegistryGithub {
		t.Fatalf("expected RegistryGithub, got %v", source.Registry)
	}
	if source.Path != "cloudquery/testSource" {
		t.Fatalf("expected source.Path (%s), got %s", source.Name, source.Path)
	}
	if source.Version != "latest" {
		t.Fatalf("expected latest, got %s", source.Version)
	}
}

const expectedSourceExample = `kind: "source"
spec:
  # Name of the plugin.
  name: "testSource"

  # Version of the plugin to use.
  version: "v0.1.0"

  # Registry to use (one of "github", "local" or "grpc").
  registry: "github"

  # Path to plugin. Required format depends on the registry.
  path: "cloudquery/testSource"

  # List of tables to sync.
  tables: ["*"]

  ## Tables to skip during sync. Optional.
  # skip_tables: []

  # Names of destination plugins to sync to.
  destinations: ["postgresql"]

  ## Approximate cap on number of requests to perform concurrently. Optional.
  # concurrency: 1000

  # Plugin-specific configuration.
  spec:
    # Check documentation here: https://github.com/cloudquery/cloudquery/tree/main/plugins/source/testSource`

func TestSourceWriteExample(t *testing.T) {
	spec := Source{
		Name: "testSource",
		Version: "v0.1.0",
		Path: "cloudquery/testSource",
		Registry: RegistryGithub,
	}
	var sb strings.Builder
	if err := spec.WriteExample(&sb); err != nil {
		t.Fatalf("failed to write example: %v", err)
	}
	if sb.String() != expectedSourceExample {
		t.Fatalf("expected example:\n%s\ngot\n%s\n", expectedSourceExample, sb.String())
	}
}

func TestSourceUnmarshalSpec(t *testing.T) {
	source := Source{
		Spec: map[string]interface{}{
			"accounts": []string{"test_account1", "test_account2"},
		},
	}
	var spec testSourceSpec
	if err := source.UnmarshalSpec(&spec); err != nil {
		t.Fatal(err)
	}
	if len(spec.Accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(spec.Accounts))
	}
	if spec.Accounts[0] != "test_account1" {
		t.Fatalf("expected test_account1, got %s", spec.Accounts[0])
	}
	if spec.Accounts[1] != "test_account2" {
		t.Fatalf("expected test_account2, got %s", spec.Accounts[1])
	}
}
