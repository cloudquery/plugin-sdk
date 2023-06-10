package serve

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

func TestPluginDocs(t *testing.T) {
	tmpDir := t.TempDir()
	p := plugin.NewPlugin(
		"testPlugin",
		"v1.0.0",
		memdb.NewMemDBClient)
	srv := Plugin(p, WithArgs("doc", tmpDir), WithTestListener())
	if err := srv.newCmdPluginDoc().Execute(); err != nil {
		t.Fatal(err)
	}
}
