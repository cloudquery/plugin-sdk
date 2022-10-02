package plugins

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudquery/faker/v3"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
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

func TestSourcePluginSync(t *testing.T, plugin *SourcePlugin, logger zerolog.Logger, spec specs.Source) {
	t.Helper()

	resourcesChannel := make(chan *schema.Resource)
	var fetchErr error

	go func() {
		defer close(resourcesChannel)
		_, fetchErr = plugin.Sync(context.Background(), logger, spec, resourcesChannel)
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
	if resource == nil {
		t.Errorf("Expected table %s to be synced but it was not found", table.Name)
		return
	}
	validateResource(t, resource)
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
	// we want to marshal and unmarshal to mimic over-the-wire behavior
	b, err := json.Marshal(resource.Data)
	if err != nil {
		t.Fatalf("failed to marshal resource data: %v", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatalf("failed to unmarshal resource data: %v", err)
	}
	for _, columnName := range resource.Table.Columns.Names() {
		if data[columnName] == nil && !resource.Table.Columns.Get(columnName).IgnoreInTests {
			t.Errorf("table: %s with unset column %s", resource.Table.Name, columnName)
		}
		val := data[columnName]
		if val != nil {
			switch resource.Table.Columns.Get(columnName).Type {
			case schema.TypeJSON:
				switch val.(type) {
				case string, []byte:
					t.Errorf("table: %s JSON column %s is being set with a string or byte slice. Either the unmarhsalled object should be passed in, or the column type should be changed to string", resource.Table.Name, columnName)
					continue
				}
				if _, err := json.Marshal(val); err != nil {
					t.Errorf("table: %s with invalid json column %s", resource.Table.Name, columnName)
				}
			default:
				// todo
			}
		}
	}

	// check that every key in the returned object exist as a column in the table
	for key := range data {
		if col := resource.Table.Columns.Get(key); col == nil {
			t.Errorf("table: %s with unknown column %s", resource.Table.Name, key)
		}
	}
}
