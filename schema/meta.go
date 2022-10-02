package schema

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ClientMeta interface {
	Logger() *zerolog.Logger
}

type Meta struct {
	LastUpdate time.Time `json:"last_updated"`
	FetchID    string    `json:"fetch_id,omitempty"`
}

const FetchIDMetaKey = "cq_fetch_id"

var CqIDColumn = Column{Name: "_cq_id", Type: TypeUUID, Description: "Internal CQ ID of the row", Resolver: cqUUIDResolver()}
var CqParentIDColumn = Column{Name: "_cq_parent_id", Type: TypeUUID, Description: "Internal CQ ID of the parent row", Resolver: parentCqUUIDResolver()}
var CqSyncTime = Column{Name: "_cq_sync_time", Type: TypeTimestamp, Description: "Internal CQ row of when sync was started (this will be the same for all rows in a single fetch)"}
var CqSourceName = Column{Name: "_cq_source_name", Type: TypeString, Description: "Internal CQ row that references the source plugin name data was retrieved"}

func cqUUIDResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		uuidGen := uuid.New()
		return r.Set(c.Name, uuidGen)
	}
}

func parentCqUUIDResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		if r.Parent == nil {
			return nil
		}
		parentCqID := r.Parent.Get(CqIDColumn.Name)
		return r.Set(c.Name, parentCqID)
	}
}
