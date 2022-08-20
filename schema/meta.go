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

var (
	// cqMeta = Column{
	// 	Name:        "cq_meta",
	// 	Type:        TypeJSON,
	// 	Description: "Meta column holds fetch information",
	// 	Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
	// 		mi := Meta{
	// 			LastUpdate: time.Now().UTC(),
	// 		}
	// 		if val, ok := resource.GetMeta(FetchIdMetaKey); ok {
	// 			if s, ok := val.(string); ok {
	// 				mi.FetchId = s
	// 			}
	// 		}
	// 		b, _ := json.Marshal(mi)
	// 		return resource.Set(c.Name, b)
	// 	},
	// 	internal: true,
	// }
	cqIdColumn = Column{
		Name:        "cq_id",
		Type:        TypeUUID,
		Description: "Unique CloudQuery Id added to every resource",
		// Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
		// 	if err := resource.GenerateCQId(); err != nil {
		// 		if resource.Parent == nil {
		// 			return err
		// 		}

		// 		meta.Logger().Debug("one of the table pk is nil", "table", resource.table.Name)
		// 	}
		// 	return resource.Set(c.Name, resource.Id())
		// },
		CreationOptions: ColumnCreationOptions{
			Unique:  true,
			NotNull: true,
		},
		internal: true,
	}
	// cqFetchDateColumn = Column{
	// 	Name:        "cq_fetch_date",
	// 	Type:        TypeTimestamp,
	// 	Description: "Time of fetch for this resource",
	// 	// Resolver: func(ctx context.Context, meta ClientMeta, resource *Resource, c Column) error {
	// 	// 	val, ok := resource.GetMeta("cq_fetch_date")
	// 	// 	if !ok && !resource.executionStart.IsZero() {
	// 	// 		val = resource.executionStart
	// 	// 	}
	// 	// 	if val == nil {
	// 	// 		return fmt.Errorf("zero cq_fetch date")
	// 	// 	}
	// 	// 	return resource.Set(c.Name, val)
	// 	// },
	// 	CreationOptions: ColumnCreationOptions{
	// 		NotNull: true,
	// 	},
	// 	internal: true,
	// }
)
