package specs

import "testing"

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
