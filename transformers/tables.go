package transformers

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// Set parent links on relational tables
func setParents(tables schema.Tables, parent *schema.Table) {
	for _, table := range tables {
		table.Parent = parent
		setParents(table.Relations, table)
	}
}

// Apply transformations to tables
func TransformTables(tables schema.Tables) error {
	for _, table := range tables {
		if table.Transform != nil {
			if err := table.Transform(table); err != nil {
				return fmt.Errorf("failed to transform table %s: %w", table.Name, err)
			}
		}
		if err := TransformTables(table.Relations); err != nil {
			return err
		}
	}
	return nil
}
