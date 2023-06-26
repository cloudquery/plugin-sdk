package serve

import (
	"context"
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
	if err := p.Init(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	srv := Plugin(p)
	cmd := srv.newCmdPluginRoot()
	cmd.SetArgs([]string{"doc", tmpDir})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}
