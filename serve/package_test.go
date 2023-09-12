package serve

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/google/go-cmp/cmp"
)

func TestPluginPackage(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current file path")
	}
	dir := filepath.Dir(filepath.Dir(filename))
	simplePluginPath := filepath.Join(dir, "examples/simple_plugin")
	packageVersion := "v1.2.3"
	p := plugin.NewPlugin(
		"testPlugin",
		"development",
		memdb.NewMemDBClient,
		plugin.WithBuildTargets([]plugin.BuildTarget{
			{OS: plugin.GoOSLinux, Arch: plugin.GoArchAmd64},
			{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64},
			{OS: plugin.GoOSDarwin, Arch: plugin.GoArchAmd64},
		}),
	)
	srv := Plugin(p)
	cmd := srv.newCmdPluginRoot()
	distDir := t.TempDir()
	cmd.SetArgs([]string{"package", "--dist-dir", distDir, simplePluginPath, packageVersion})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(distDir)
	if err != nil {
		t.Fatal(err)
	}
	expect := []string{
		"package.json",
		"plugin-testPlugin-v1.2.3-darwin-amd64.zip",
		"plugin-testPlugin-v1.2.3-linux-amd64.zip",
		"plugin-testPlugin-v1.2.3-windows-amd64.zip",
		"tables.json",
	}
	if diff := cmp.Diff(expect, fileNames(files)); diff != "" {
		t.Fatalf("unexpected files in dist directory (-want +got):\n%s", diff)
	}

	expectPackage := PackageJSON{
		Name:      "testPlugin",
		Version:   "v1.2.3",
		Protocols: []int{3},
		SupportedTargets: []TargetBuild{
			{OS: plugin.GoOSLinux, Arch: plugin.GoArchAmd64, Path: "plugin-testPlugin-v1.2.3-linux-amd64.zip"},
			{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64, Path: "plugin-testPlugin-v1.2.3-windows-amd64.zip"},
			{OS: plugin.GoOSDarwin, Arch: plugin.GoArchAmd64, Path: "plugin-testPlugin-v1.2.3-darwin-amd64.zip"},
		},
		PackageType: plugin.PackageTypeNative,
	}
	checkPackageJSONContents(t, filepath.Join(distDir, "package.json"), expectPackage)
}

func checkPackageJSONContents(t *testing.T, filename string, expect PackageJSON) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open package.json: %v", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read package.json: %v", err)
	}
	j := PackageJSON{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		t.Fatalf("failed to unmarshal package.json: %v", err)
	}
	if diff := cmp.Diff(expect, j); diff != "" {
		t.Fatalf("package.json contents mismatch (-want +got):\n%s", diff)
	}
}

func fileNames(files []os.DirEntry) []string {
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}
