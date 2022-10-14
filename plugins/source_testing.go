package plugins

import (
	"context"
	"encoding/json"
	"testing"

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

func TestSourcePluginSync(t *testing.T, plugin *SourcePlugin, logger zerolog.Logger, spec specs.Source, opts ...TestSourcePluginOption) {
	t.Helper()

	o := &testSourcePluginOptions{}
	for _, opt := range opts {
		opt(o)
	}
	if !o.NoParallel {
		t.Parallel()
	}

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

type TestSourcePluginOption func(*testSourcePluginOptions)

func TestSourcePluginWithoutParallel() TestSourcePluginOption {
	return func(f *testSourcePluginOptions) {
		f.NoParallel = true
	}
}

type testSourcePluginOptions struct {
	NoParallel bool
}

func getTableResources(t *testing.T, table *schema.Table, resources []*schema.Resource) []*schema.Resource {
	t.Helper()

	tableResources := make([]*schema.Resource, 0)

	for _, resource := range resources {
		if resource.Table.Name == table.Name {
			tableResources = append(tableResources, resource)
		}
	}

	return tableResources
}

func validateTable(t *testing.T, table *schema.Table, resources []*schema.Resource) {
	t.Helper()
	tableResources := getTableResources(t, table, resources)
	if len(tableResources) == 0 {
		t.Errorf("Expected table %s to be synced but it was not found", table.Name)
		return
	}
	validateResources(t, tableResources)
}

func validateTables(t *testing.T, tables schema.Tables, resources []*schema.Resource) {
	t.Helper()
	for _, table := range tables {
		validateTable(t, table, resources)
		validateTables(t, table.Relations, resources)
	}
}

// Validates that every column has at least one non-nil value.
// Also does some additional validations.
func validateResources(t *testing.T, resources []*schema.Resource) {
	t.Helper()

	table := resources[0].Table

	// A set of column-names that have values in at least one of the resources.
	columnsWithValues := make(map[string]bool)

	for _, resource := range resources {
		// we want to marshal and unmarshal to mimic over-the-wire behavior
		b, err := json.Marshal(resource.Data)
		if err != nil {
			t.Fatalf("failed to marshal resource data: %v", err)
		}
		var data map[string]interface{}
		if err := json.Unmarshal(b, &data); err != nil {
			t.Fatalf("failed to unmarshal resource data: %v", err)
		}

		for columnName, value := range data {
			if value == nil {
				continue
			}

			columnsWithValues[columnName] = true

			switch resource.Table.Columns.Get(columnName).Type {
			case schema.TypeJSON:
				switch value.(type) {
				case string, []byte:
					t.Errorf("table: %s JSON column %s is being set with a string or byte slice. Either the unmarhsalled object should be passed in, or the column type should be changed to string", resource.Table.Name, columnName)
					continue
				}
				if _, err := json.Marshal(value); err != nil {
					t.Errorf("table: %s with invalid json column %s", table.Name, columnName)
				}
			default:
				// todo
			}
		}

		// check that every key in the returned object exist as a column in the table
		for key := range data {
			if col := resource.Table.Columns.Get(key); col == nil {
				t.Errorf("table: %s with unknown column %s", table.Name, key)
			}
		}
	}

	// Make sure every column has at least one value.
	for _, columnName := range table.Columns.Names() {
		if _, ok := columnsWithValues[columnName]; !ok && !table.Columns.Get(columnName).IgnoreInTests {
			t.Errorf("Expected column %s to have at least one non-nil value but it was not found", columnName)
		}
	}
}
