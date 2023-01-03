package local

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/specs"
)

func TestLocal(t *testing.T) {
	tmpDir := t.TempDir()
	ss := specs.Source{
		Name:    "test",
		Version: "vtest",
		Path:    "test",
		Backend: specs.BackendLocal,
		BackendSpec: Spec{
			Path: tmpDir,
		},
	}
	local, err := New(ss)
	if err != nil {
		t.Fatalf("failed to create local backend: %v", err)
	}
	if local.spec.Path != tmpDir {
		t.Fatalf("expected path to be %s, but got %s", tmpDir, local.spec.Path)
	}

	got, err := local.Get("test_table", "test_key")
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value, but got %s", got)
	}

	err = local.Set("test_table", "test_key", "test_value")
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	got, err = local.Get("test_table", "test_key")
	if err != nil {
		t.Fatalf("failed to get value after setting it: %v", err)
	}
	if got != "test_value" {
		t.Fatalf("expected value to be test_value, but got %s", got)
	}

	err = local.Close()
	if err != nil {
		t.Fatalf("failed to close local backend: %v", err)
	}

	local, err = New(ss)
	if err != nil {
		t.Fatalf("failed to open local backend after closing it: %v", err)
	}

	got, err = local.Get("test_table", "test_key")
	if err != nil {
		t.Fatalf("failed to get value after closing and reopening local backend: %v", err)
	}
	if got != "test_value" {
		t.Fatalf("expected value to be test_value, but got %s", got)
	}

	got, err = local.Get("some_other_table", "test_key")
	if err != nil {
		t.Fatalf("failed to get value after closing and reopening local backend: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value for some_other_table -> test_key, but got %s", got)
	}

	// check that state is namespaced by source name
	ss.Name = "test2"
	local, err = New(ss)
	if err != nil {
		t.Fatalf("failed to create local backend for test2: %v", err)
	}
	got, err = local.Get("test_table", "test_key")
	if err != nil {
		t.Fatalf("failed to get value for local backend test2: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value for test2 -> test_table -> test_key, but got %s", got)
	}
}
