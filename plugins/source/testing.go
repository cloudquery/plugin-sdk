package source

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

func TestPluginSync(t *testing.T, plugin *Plugin, spec specs.Source, opts ...TestPluginOption) {
	t.Helper()

	o := &testPluginOptions{
		parallel: true,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.parallel {
		t.Parallel()
	}

	resourcesChannel := make(chan *schema.Resource)
	var syncErr error

	if err := plugin.Init(context.Background(), spec); err != nil {
		t.Fatal(err)
	}

	go func() {
		defer close(resourcesChannel)
		syncErr = plugin.Sync(context.Background(), resourcesChannel)
	}()

	syncedResources := make([]*schema.Resource, 0)
	for resource := range resourcesChannel {
		syncedResources = append(syncedResources, resource)
	}
	if syncErr != nil {
		t.Fatal(syncErr)
	}

	validateTables(t, plugin.Tables(), syncedResources)
}

type TestPluginOption func(*testPluginOptions)

func WithTestPluginNoParallel() TestPluginOption {
	return func(f *testPluginOptions) {
		f.parallel = false
	}
}

type testPluginOptions struct {
	parallel bool
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
	columnsWithValues := make([]bool, len(table.Columns))

	for _, resource := range resources {
		for i, value := range resource.GetValues() {
			if value == nil {
				continue
			}
			if value.Get() != nil && value.Get() != schema.Undefined {
				columnsWithValues[i] = true
			}
		}
	}

	// Make sure every column has at least one value.
	for i, hasValue := range columnsWithValues {
		col := table.Columns[i]
		emptyExpected := col.Name == "_cq_parent_id" && table.Parent == nil
		if !hasValue && !emptyExpected && !col.IgnoreInTests {
			t.Errorf("table: %s column %s has no values", table.Name, table.Columns[i].Name)
		}
	}
}
