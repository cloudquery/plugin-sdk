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
var CqFetchTime = Column{Name: "_cq_fetch_time", Type: TypeTimestamp, Description: "Internal CQ row of when fetch was started (this will be the same for all rows in a single fetch)"}

var CqColumns = []Column{
	CqIDColumn,
	CqFetchTime,
}

func cqUUIDResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		uuidGen := uuid.New()
		return r.Set(c.Name, uuidGen)
	}
}
