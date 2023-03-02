package source

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
	pluralize "github.com/gertd/go-pluralize"
)

type Validator func(plugin *Plugin, resources []*schema.Resource) error

func validateColumnsHaveData(plugin *Plugin, resources []*schema.Resource) error {
	tables := extractTables(plugin.tables)
	for _, table := range tables {
		err := validateTable(table, resources)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateTableNamePlural(plugin *Plugin, _ []*schema.Resource) error {
	pluralizeClient := pluralize.NewClient()
	tables := extractTables(plugin.tables)
	for _, table := range tables {
		if !pluralizeClient.IsPlural(table.Name) {
			return fmt.Errorf("invalid table name: %s. must be plural", table.Name)
		}
	}
	return nil
}

func getTableResources(table *schema.Table, resources []*schema.Resource) []*schema.Resource {
	tableResources := make([]*schema.Resource, 0)

	for _, resource := range resources {
		if resource.Table.Name == table.Name {
			tableResources = append(tableResources, resource)
		}
	}

	return tableResources
}

func validateTable(table *schema.Table, resources []*schema.Resource) error {
	tableResources := getTableResources(table, resources)
	if len(tableResources) == 0 {
		return fmt.Errorf("expected table %s to be synced but it was not found", table.Name)
	}
	return validateResources(tableResources)
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
func validateResources(resources []*schema.Resource) error {
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
			return fmt.Errorf("table: %s column %s has no values", table.Name, table.Columns[i].Name)
		}
	}
	return nil
}
