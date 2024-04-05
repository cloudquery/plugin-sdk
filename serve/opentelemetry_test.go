package serve

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

func TestNewResource(t *testing.T) {
	p := plugin.NewPlugin(
		"testPluginV3",
		"v1.0.0",
		memdb.NewMemDBClient)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("failed to instantiate resource: %v", r)
		}
	}()
	newResource(p)
}
