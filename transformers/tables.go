package transformers

import (
	"context"
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

// Add internal columns
func AddInternalColumns(tables []*schema.Table) error {
	for _, table := range tables {
		if c := table.Column("_cq_id"); c != nil {
			return fmt.Errorf("table %s already has column _cq_id", table.Name)
		}
		cqID := schema.CqIDColumn
		if len(table.PrimaryKeys()) == 0 {
			cqID.PrimaryKey = true
		}
		cqSourceName := schema.CqSourceNameColumn
		cqSyncTime := schema.CqSyncTimeColumn
		cqSourceName.Resolver = func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			return resource.Set(c.Name, p.sourceName)
		}
		cqSyncTime.Resolver = func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			return resource.Set(c.Name, p.syncTime)
		}

		table.Columns = append([]schema.Column{cqSourceName, cqSyncTime, cqID, schema.CqParentIDColumn}, table.Columns...)
		if err := AddInternalColumns(table.Relations); err != nil {
			return err
		}
	}
	return nil
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
