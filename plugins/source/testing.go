package source

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
)

type Validator func(t *testing.T, plugin *Plugin, resources []*schema.Resource)

func TestPluginSync(t *testing.T, plugin *Plugin, spec specs.Source, opts ...TestPluginOption) {
	t.Helper()

	o := &testPluginOptions{
		parallel:   true,
		validators: []Validator{validatePlugin},
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
	for _, validator := range o.validators {
		validator(t, plugin, syncedResources)
	}
}

type TestPluginOption func(*testPluginOptions)

func WithTestPluginNoParallel() TestPluginOption {
	return func(f *testPluginOptions) {
		f.parallel = false
	}
}

func WithTestPluginAdditionalValidators(v Validator) TestPluginOption {
	return func(f *testPluginOptions) {
		f.validators = append(f.validators, v)
	}
}

type testPluginOptions struct {
	parallel   bool
	validators []Validator
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

func validatePlugin(t *testing.T, plugin *Plugin, resources []*schema.Resource) {
	t.Helper()
	tables := extractTables(plugin.tables)
	for _, table := range tables {
		validateTable(t, table, resources)
	}
}

func extractTables(tables schema.Tables) []*schema.Table {
	result := make([]*schema.Table, 0)
	for _, table := range tables {
		result = append(result, table)
		result = append(result, extractTables(table.Relations)...)
	}
	return result
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
			if value.IsValid() {
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
