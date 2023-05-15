package schema

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v3/types"
	"github.com/google/uuid"
)

type ClientMeta interface {
	ID() string
}

// These columns are managed and populated by the source plugins
var CqIDColumn = Column{
	Name:            CqIDField.Name,
	Type:            CqIDField.Type,
	Description:     "Internal CQ ID of the row",
	CreationOptions: CreationOptions{NotNull: true, Unique: true},
}
var CqParentIDColumn = Column{
	Name:          CqParentIDField.Name,
	Type:          CqParentIDField.Type,
	Description:   "Internal CQ ID of the parent row",
	Resolver:      parentCqUUIDResolver(),
	IgnoreInTests: true,
}

// These columns are managed and populated by the destination plugin.
var CqSyncTimeColumn = Column{
	Name:        CqSyncTimeField.Name,
	Type:        CqSyncTimeField.Type,
	Description: "Internal CQ row of when sync was started (this will be the same for all rows in a single fetch)",
}
var CqSourceNameColumn = Column{
	Name:        CqSourceNameField.Name,
	Type:        CqSourceNameField.Type,
	Description: "Internal CQ row that references the source plugin name data was retrieved",
}

var CqIDField = arrow.Field{
	Name:     "_cq_id",
	Type:     types.ExtensionTypes.UUID,
	Metadata: arrow.MetadataFrom(map[string]string{MetadataUnique: MetadataTrue}),
}

var CqParentIDField = arrow.Field{
	Name:     "_cq_parent_id",
	Type:     types.ExtensionTypes.UUID,
	Metadata: arrow.MetadataFrom(map[string]string{MetadataUnique: MetadataTrue}),
}

var CqSyncTimeField = arrow.Field{
	Name: "_cq_sync_time",
	Type: arrow.FixedWidthTypes.Timestamp_us,
}

var CqSourceNameField = arrow.Field{
	Name: "_cq_source_name",
	Type: arrow.BinaryTypes.String,
}

func parentCqUUIDResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		if r.Parent == nil {
			return r.Set(c.Name, nil)
		}
		parentCqID := r.Parent.Get(CqIDColumn.Name)
		if parentCqID == nil {
			return r.Set(c.Name, nil)
		}
		pUUID, ok := parentCqID.(*uuid.UUID)
		if !ok {
			return r.Set(c.Name, nil)
		}
		return r.Set(c.Name, pUUID)
	}
}
