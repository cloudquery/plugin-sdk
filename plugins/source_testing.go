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

func getResourcesCountPerTable(table *schema.Table) int {
	resourcesCount := 1
	for _, relation := range table.Relations {
		resourcesCount += getResourcesCountPerTable(relation)
	}

	return resourcesCount
}

func getResourcesCount(tables schema.Tables) int {
	total := 0
	for _, table := range tables {
		total += getResourcesCountPerTable(table)
	}
	return total
}

func TestSourcePluginSync(t *testing.T, plugin *SourcePlugin, spec specs.Source) {
	t.Helper()

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
	require.NoError(t, fetchErr)

	resourcesCount := getResourcesCount(plugin.Tables())
	require.Equal(t, resourcesCount, totalResources, "expected %d resources, got %d", resourcesCount, totalResources)
}

func validateResource(t *testing.T, resource *schema.Resource) {
	t.Helper()
	for _, columnName := range resource.Table.Columns.Names() {
		if resource.Get(columnName) == nil && !resource.Table.Columns.Get(columnName).IgnoreInTests {
			t.Errorf("table: %s with unset column %s", resource.Table.Name, columnName)
		}
	}
}
