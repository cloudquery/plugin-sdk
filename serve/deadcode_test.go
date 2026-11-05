package serve

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPluginBinaryKeepsMethodDeadCodeElimination(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a plugin binary")
	}
	out, err := exec.Command("go", "build", "-tags=grpcnotrace", "-ldflags=-dumpdep", "-o", filepath.Join(t.TempDir(), "plugin"), "./testdata/deadcode").CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %v\n%s", err, out[max(0, len(out)-4096):])
	}
	for line := range strings.Lines(string(out)) {
		if strings.Contains(line, "<ReflectMethod>") {
			t.Fatalf("a reflect method call is reachable, so the linker keeps every exported method of every type: %s", line)
		}
	}
}
