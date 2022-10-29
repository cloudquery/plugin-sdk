package schema

import (
	"context"
	"time"

	// "github.com/cloudquery/plugin-sdk/cqtypes"
	"github.com/cloudquery/plugin-sdk/cqtypes"
	"github.com/google/uuid"
)

type ClientMeta interface {
	Name() string
}

type Meta struct {
	LastUpdate time.Time `json:"last_updated"`
	FetchID    string    `json:"fetch_id,omitempty"`
}

// These columns are managed and populated by the source plugins
var CqIDColumn = Column{
	Name:        "_cq_id",
	Type:        TypeUUID,
	Description: "Internal CQ ID of the row",
	Resolver:    cqUUIDResolver(),
}
var CqParentIDColumn = Column{
	Name:          "_cq_parent_id",
	Type:          TypeUUID,
	Description:   "Internal CQ ID of the parent row",
	Resolver:      parentCqUUIDResolver(),
	IgnoreInTests: true,
}

// These columns are managed and populated by the destination plugin.
var CqSyncTimeColumn = Column{
	Name:        "_cq_sync_time",
	Type:        TypeTimestamp,
	Description: "Internal CQ row of when sync was started (this will be the same for all rows in a single fetch)",
}
var CqSourceNameColumn = Column{
	Name:        "_cq_source_name",
	Type:        TypeString,
	Description: "Internal CQ row that references the source plugin name data was retrieved",
}

var CqSyncTimeColumnD = Column{
	Name: "_cq_sync_time",
	Type: TypeTimestamp,
}
var CqSourceNameColumnD = Column{
	Name: "_cq_source_name",
	Type: TypeString,
}

func cqUUIDResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		uuidGen := uuid.New()
		b, _ := uuidGen.MarshalBinary()
		return r.Set(c.Name, b)
	}
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
		pUUID, ok := parentCqID.(*cqtypes.UUID)
		if !ok {
			return r.Set(c.Name, nil)
		}
		return r.Set(c.Name, pUUID.Bytes)
	}
}
