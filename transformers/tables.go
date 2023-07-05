package transformers

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// Set parent links on relational tables
func SetParents(tables schema.Tables, parent *schema.Table) {
	for _, table := range tables {
		table.Parent = parent
		SetParents(table.Relations, table)
	}
}

// TransformTables runs given Tables' transformers as defined in the table definitions, and recursively on the relations
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

// ApplyTransformers applies the given transformers to the given Tables, and recursively on the relations
func ApplyTransformers(tables schema.Tables, extraTransformers ...schema.Transform) error {
	for _, table := range tables {
		for _, tf := range extraTransformers {
			if err := tf(table); err != nil {
				return fmt.Errorf("failed to transform table %s: %w", table.Name, err)
			}
		}
		if err := ApplyTransformers(table.Relations, extraTransformers...); err != nil {
			return err
		}
	}
	return nil
}
