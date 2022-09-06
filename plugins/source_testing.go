package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/faker/v3"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

type ResourceTestCase struct {
	Plugin *SourcePlugin
	Spec   specs.Source
	// ParallelFetchingLimit limits parallel resources fetch at a time
	ParallelFetchingLimit uint64
	// SkipIgnoreInTest flag which detects if schema.Table or schema.Column should be ignored
	SkipIgnoreInTest bool
}

func init() {
	_ = faker.SetRandomMapAndSliceMinSize(1)
	_ = faker.SetRandomMapAndSliceMaxSize(1)
}

// type

func TestSourcePluginSync(t *testing.T, plugin *SourcePlugin, spec specs.Source) {
	// t.Parallel()
	t.Helper()
	// No need for configuration or db connection, get it out of the way first
	// testTableIdentifiersForProvider(t, resource.Provider)

	// l := testlog.New(t)
	// l.SetLevel(hclog.Info)
	// resource.Plugin.Logger = l
	resources := make(chan *schema.Resource)
	var fetchErr error

	go func() {
		defer close(resources)
		fetchErr = plugin.Sync(context.Background(), spec, resources)
	}()
	totalResources := 0
	for resource := range resources {
		totalResources++
		validateResource(t, resource)
	}
	if fetchErr != nil {
		t.Fatal(fetchErr)
	}
	if totalResources == 0 {
		t.Fatal("no resources fetched")
	}
}

func validateResource(t *testing.T, resource *schema.Resource) {
	t.Helper()
	for _, columnName := range resource.Table.Columns.Names() {
		if resource.Get(columnName) == nil && !resource.Table.Columns.Get(columnName).IgnoreInTests {
			t.Errorf("table: %s with unset column %s", resource.Table.Name, columnName)
		}
	}
}
