package schema

import (
	"context"
	"time"
)

var (
	meta = Column{
		Name:        "meta",
		Type:        TypeJSON,
		Description: "Meta column holds fetch information",
		Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
			return resource.Set(c.Name, map[string]interface{}{"last_updated": time.Now().UTC().Format(time.RFC3339)})
		},
	}
	cqIdColumn = Column{
		Name:        "cq_id",
		Type:        TypeUUID,
		Description: "Unique CloudQuery Id added to every resource",
		Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
			if err := resource.GenerateCQId(); err != nil {
				if resource.Parent == nil {
					return err
				} else {
					meta.Logger().Debug("one of the table pk is nil", "table", resource.table.Name)
				}
			}
			return resource.Set(c.Name, resource.Id())
		},
		CreationOptions: ColumnCreationOptions{
			Unique: true,
		},
	}
)

// GetDefaultSDKColumns Default columns of the SDK, these columns are added to each table by default
func GetDefaultSDKColumns() []Column {
	return []Column{cqIdColumn, meta}
}
