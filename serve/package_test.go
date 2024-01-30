package serve

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/memdbtables.json
var memDBPackageJSON string

//go:embed testdata/source_spec_schema.json
var sourceSpecSchema string

//go:embed testdata/destination_spec_schema.json
var destinationSpecSchema string

func TestPluginPackage_Source(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current file path")
	}
	dir := filepath.Dir(filepath.Dir(filename))
	simplePluginPath := filepath.Join(dir, "examples/simple_plugin")
	packageVersion := "v1.2.3"
	p := plugin.NewPlugin(
		"test-plugin",
		"development",
		memdb.NewMemDBClient,
		plugin.WithBuildTargets([]plugin.BuildTarget{
			{OS: plugin.GoOSLinux, Arch: plugin.GoArchAmd64},
			{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64},
		}),
		plugin.WithKind("source"),
		plugin.WithTeam("test-team"),
		plugin.WithJSONSchema(sourceSpecSchema),
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
			t.Setenv("CGO_ENABLED", "0") // disable CGO to ensure we environmental differences don't interfere with the test
			srv := Plugin(p)
			cmd := srv.newCmdPluginRoot()
			distDir := t.TempDir()
			cmd.SetArgs([]string{"package", "--dist-dir", distDir, "-m", tc.message, packageVersion, simplePluginPath})
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
				"plugin-test-plugin-v1.2.3-linux-amd64.zip",
				"plugin-test-plugin-v1.2.3-windows-amd64.zip",
				"spec_json_schema.json",
				"tables.json",
			}
			if diff := cmp.Diff(expect, fileNames(files)); diff != "" {
				t.Fatalf("unexpected files in dist directory (-want +got):\n%s", diff)
			}
			// expect SHA-256 for the zip files to differ
			sha1 := sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-linux-amd64.zip"))
			sha2 := sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-windows-amd64.zip"))
			if sha1 == sha2 {
				t.Fatalf("expected SHA-256 for linux and windows zip files to differ, but they are the same: %s", sha1)
			}

			expectPackage := PackageJSON{
				SchemaVersion: 1,
				Name:          "test-plugin",
				Team:          "test-team",
				Kind:          "source",
				Message:       msg,
				Version:       packageVersion,
				Protocols:     []int{3},
				SupportedTargets: []TargetBuild{
					{OS: plugin.GoOSLinux, Arch: plugin.GoArchAmd64, Path: "plugin-test-plugin-v1.2.3-linux-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-linux-amd64.zip"))},
					{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64, Path: "plugin-test-plugin-v1.2.3-windows-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-windows-amd64.zip"))},
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
			checkFileContent(t, filepath.Join(distDir, "spec_json_schema.json"), sourceSpecSchema)
		})
	}
}

func TestPluginPackage_Destination(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current file path")
	}
	dir := filepath.Dir(filepath.Dir(filename))
	simplePluginPath := filepath.Join(dir, "examples/simple_plugin")
	packageVersion := "v1.2.3"
	p := plugin.NewPlugin(
		"test-plugin",
		"development",
		memdb.NewMemDBClient,
		plugin.WithBuildTargets([]plugin.BuildTarget{
			{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64},
			{OS: plugin.GoOSDarwin, Arch: plugin.GoArchAmd64},
		}),
		plugin.WithKind("destination"),
		plugin.WithTeam("test-team"),
		plugin.WithJSONSchema(destinationSpecSchema),
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
			t.Setenv("CGO_ENABLED", "0") // disable CGO to ensure we environmental differences don't interfere with the test
			srv := Plugin(p)
			cmd := srv.newCmdPluginRoot()
			distDir := t.TempDir()
			cmd.SetArgs([]string{"package", "--dist-dir", distDir, "-m", tc.message, packageVersion, simplePluginPath})
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
				"plugin-test-plugin-v1.2.3-darwin-amd64.zip",
				"plugin-test-plugin-v1.2.3-windows-amd64.zip",
				"spec_json_schema.json",
			}
			if diff := cmp.Diff(expect, fileNames(files)); diff != "" {
				t.Fatalf("unexpected files in dist directory (-want +got):\n%s", diff)
			}
			// expect SHA-256 for the zip files to differ
			sha1 := sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-windows-amd64.zip"))
			sha2 := sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-darwin-amd64.zip"))
			if sha1 == sha2 {
				t.Fatalf("expected SHA-256 for windows and darwin zip files to differ, but they are the same: %s", sha1)
			}

			expectPackage := PackageJSON{
				SchemaVersion: 1,
				Team:          "test-team",
				Kind:          "destination",
				Name:          "test-plugin",
				Message:       msg,
				Version:       "v1.2.3",
				Protocols:     []int{3},
				SupportedTargets: []TargetBuild{
					{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64, Path: "plugin-test-plugin-v1.2.3-windows-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-windows-amd64.zip"))},
					{OS: plugin.GoOSDarwin, Arch: plugin.GoArchAmd64, Path: "plugin-test-plugin-v1.2.3-darwin-amd64.zip", Checksum: "sha256:" + sha256sum(filepath.Join(distDir, "plugin-test-plugin-v1.2.3-darwin-amd64.zip"))},
				},
				PackageType: plugin.PackageTypeNative,
			}
			checkPackageJSONContents(t, filepath.Join(distDir, "package.json"), expectPackage)

			expectDocs := []string{
				"configuration.md",
				"overview.md",
			}
			checkDocs(t, filepath.Join(distDir, "docs"), expectDocs)
			checkFileContent(t, filepath.Join(distDir, "spec_json_schema.json"), destinationSpecSchema)
		})
	}
}

func TestVersionRegex(t *testing.T) {
	t.Parallel()
	re := Plugin(&plugin.Plugin{}).versionRegex()
	testCases := []struct {
		input string
		want  bool
	}{
		{
			input: `var Version = ""`,
			want:  true,
		},
		{
			input: `var (
Version = ""
)
`,
			want: true,
		},
		{
			input: ` var Version = ""`,
			want:  false,
		},
		{
			input: `var (
  Version = ""
)
`,
		},
		{
			input: `var (
	Version = ""
)
`,
			want: true,
		},
	}
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			t.Parallel()
			// only match the line with "Version"
			lines := strings.Split(tc.input, "\n")
			realInput := ""
			for _, line := range lines {
				if strings.Contains(line, "Version") {
					realInput = line
					break
				}
			}
			if realInput == "" {
				t.Fatalf("failed to find line with Version: %q", tc.input)
			}
			got := re.MatchString(realInput)
			require.Equalf(t, tc.want, got, "input: %q", realInput)
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
	var actual any
	err = json.Unmarshal(content, &actual)
	require.NoError(t, err)

	var expected any
	err = json.Unmarshal([]byte(memDBPackageJSON), &expected)
	require.NoError(t, err)

	require.Equal(t, expected, actual, "tables.json content mismatch")
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

func checkFileContent(t *testing.T, filename string, expect string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open %s: %v", filename, err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read %s: %v", filename, err)
	}
	if diff := cmp.Diff(expect, string(b)); diff != "" {
		t.Fatalf("%s contents mismatch (-want +got):\n%s", filename, diff)
	}
}

func fileNames(files []os.DirEntry) []string {
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}
