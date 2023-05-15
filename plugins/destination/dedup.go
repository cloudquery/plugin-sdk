package destination

import (
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v3/internal/pk"
	"github.com/cloudquery/plugin-sdk/v3/schema"
)

func RemoveDuplicatesByPK(table *schema.Table, resources []arrow.Record) []arrow.Record {
	pkIndices := table.PrimaryKeysIndexes()
	// special case where there's no PK at all
	if len(pkIndices) == 0 {
		return resources
	}

	pks := make(map[string]struct{}, len(resources))
	res := make([]arrow.Record, 0, len(resources))
	for _, r := range resources {
		if r.NumRows() > 1 {
			panic(fmt.Sprintf("record with more than 1 row: %d", r.NumRows()))
		}
		key := pk.String(r)
		_, ok := pks[key]
		if !ok {
			pks[key] = struct{}{}
			res = append(res, r)
			continue
		}
		// duplicate, release
		r.Release()
	}

	return res
}
