package serve

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

func TestPluginPublish(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current file path")
	}
	dir := filepath.Dir(filepath.Dir(filename))
	simplePluginPath := filepath.Join(dir, "examples/simple_plugin")
	p := plugin.NewPlugin(
		"testPlugin",
		"v1.0.0",
		memdb.NewMemDBClient)
	srv := Plugin(p)
	cmd := srv.newCmdPluginRoot()
	cmd.SetArgs([]string{"publish", simplePluginPath, "--dist-dir", t.TempDir()})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}
