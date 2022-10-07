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

var CqSyncTime = Column{
	Name:          "_cq_sync_time",
	Type:          TypeTimestamp,
	Description:   "Internal CQ row of when sync was started (this will be the same for all rows in a single fetch)",
	Resolver:      noopResolver(), // _cq_sync_time is set later in the SDK, so we use noopResolver.
	IgnoreInTests: true,
}

var CqSourceName = Column{
	Name:          "_cq_source_name",
	Type:          TypeString,
	Description:   "Internal CQ row that references the source plugin name data was retrieved",
	Resolver:      noopResolver(), // _cq_source_name is set later in the SDK, so we use noopResolver.
	IgnoreInTests: true,
}

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

// A placeholder resolver that doesn't do anything.
// Use this when you need a resolver that doesn't do anything (rather than `nil“), because the default
// resolver actually uses reflection to look for fields in the upstream resource.
func noopResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return nil
	}
}
