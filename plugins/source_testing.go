package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/faker/v3"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/stretchr/testify/require"
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

func TestSourcePluginSync(t *testing.T, plugin *SourcePlugin, spec specs.Source) {
	t.Helper()

	resourcesChannel := make(chan *schema.Resource)
	var fetchErr error

	go func() {
		defer close(resourcesChannel)
		fetchErr = plugin.Sync(context.Background(), spec, resourcesChannel)
	}()

	syncedResources := make([]*schema.Resource, 0)
	for resource := range resourcesChannel {
		syncedResources = append(syncedResources, resource)
	}
	require.NoError(t, fetchErr)

	validateTables(t, plugin.Tables(), syncedResources)
}

func getTableResource(t *testing.T, table *schema.Table, resources []*schema.Resource) *schema.Resource {
	t.Helper()
	for _, resource := range resources {
		if resource.Table.Name == table.Name {
			return resource
		}
	}

	return nil
}

func validateTable(t *testing.T, table *schema.Table, resources []*schema.Resource) {
	t.Helper()
	resource := getTableResource(t, table, resources)
	if resource != nil {
		validateResource(t, resource)
		return
	}
	t.Errorf("Expected table %s to be synced but it was not found", table.Name)
}

func validateTables(t *testing.T, tables schema.Tables, resources []*schema.Resource) {
	t.Helper()
	for _, table := range tables {
		validateTable(t, table, resources)
		validateTables(t, table.Relations, resources)
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
