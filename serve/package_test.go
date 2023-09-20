package serve

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
	msg := `Test message
with multiple lines and **markdown**`
	testCases := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "inline message",
			message: msg,
		},
		{
			name:    "message from file",
			message: "@testdata/message.txt",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := Plugin(p)
			cmd := srv.newCmdPluginRoot()
			distDir := t.TempDir()
			cmd.SetArgs([]string{"package", "--dist-dir", distDir, "-m", tc.message, simplePluginPath, packageVersion})
			err := cmd.Execute()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			files, err := os.ReadDir(distDir)
			if err != nil {
				t.Fatal(err)
			}
			expect := []string{
				"docs",
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
				SchemaVersion: 1,
				Name:          "testPlugin",
				Message:       msg,
				Version:       "v1.2.3",
				Protocols:     []int{3},
				SupportedTargets: []TargetBuild{
					{OS: plugin.GoOSLinux, Arch: plugin.GoArchAmd64, Path: "plugin-testPlugin-v1.2.3-linux-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-testPlugin-v1.2.3-linux-amd64.zip"))},
					{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64, Path: "plugin-testPlugin-v1.2.3-windows-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-testPlugin-v1.2.3-windows-amd64.zip"))},
					{OS: plugin.GoOSDarwin, Arch: plugin.GoArchAmd64, Path: "plugin-testPlugin-v1.2.3-darwin-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-testPlugin-v1.2.3-darwin-amd64.zip"))},
				},
				PackageType: plugin.PackageTypeNative,
			}
			checkPackageJSONContents(t, filepath.Join(distDir, "package.json"), expectPackage)

			expectDocs := []string{
				"configuration.md",
				"overview.md",
			}
			checkDocs(t, filepath.Join(distDir, "docs"), expectDocs)
			checkTables(t, distDir)
		})
	}
}

func sha256sum(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func checkDocs(t *testing.T, dir string, expect []string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expect, fileNames(files)); diff != "" {
		t.Fatalf("unexpected files in docs directory (-want +got):\n%s", diff)
	}
}

func checkTables(t *testing.T, distDir string) {
	content, err := os.ReadFile(filepath.Join(distDir, "tables.json"))
	if err != nil {
		t.Fatal(err)
	}
	tablesString := string(content)

	if diff := cmp.Diff(tablesString, "[{\"name\":\"table1\",\"relations\":[\"table2\"]},{\"name\":\"table2\"}]\n"); diff != "" {
		t.Fatalf("unexpected content in tables.json (-want +got):\n%s", diff)
	}
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
