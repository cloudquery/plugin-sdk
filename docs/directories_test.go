package docs

import (
	"os"
	"path"
	"testing"
)

func TestPrepareDirectory(t *testing.T) {
	tmpdir, tmpErr := os.MkdirTemp("", "docs_test_*")
	if tmpErr != nil {
		t.Fatalf("failed to create temporary directory: %v", tmpErr)
	}
	defer os.RemoveAll(tmpdir)

	p := path.Join(tmpdir, "testdir")
	err := PrepareDirectory(p)
	if err != nil {
		t.Errorf("PrepareDirectory returned error: %v", err)
	}

	_, err = os.Create(path.Join(p, "testfile"))
	if err != nil {
		t.Fatalf("failed to create file for test: %v", err)
	}
	err = PrepareDirectory(p)
	if err != nil {
		t.Errorf("PrepareDirectory with existing directory returned error: %v", err)
	}

	_, err = os.Open(path.Join(p, "testfile"))
	if err == nil {
		t.Errorf("file was not deleted by PrepareDirectory")
	}
}
