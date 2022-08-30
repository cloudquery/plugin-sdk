package schema

import (
	"time"

	"github.com/rs/zerolog"
)

type ClientMeta interface {
	Logger() *zerolog.Logger
}

type Meta struct {
	LastUpdate time.Time `json:"last_updated"`
	FetchId    string    `json:"fetch_id,omitempty"`
}

const FetchIdMetaKey = "cq_fetch_id"

var CqIdColumn = Column{Name: "_cq_id", Type: TypeUUID, Description: "Internal CQ ID of the row", CreationOptions: ColumnCreationOptions{Unique: true}, Resolver: CQUUIDResolver()}
var CqFetchTime = Column{Name: "_cq_fetch_time", Type: TypeTimestamp, Description: "Internal CQ row of when fetch was started (this will be the same for all rows in a single fetch)"}

var CqColumns = []Column{
	CqIdColumn,
	CqFetchTime,
}
