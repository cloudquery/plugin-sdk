package pk

import (
	"strings"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func String(resource arrow.Record) string {
	sc := resource.Schema()
	table, err := schema.NewTableFromArrowSchema(sc)
	if err != nil {
		panic(err)
	}
	pkIndices := table.PrimaryKeysIndexes()
	parts := make([]string, 0, len(pkIndices))
	for _, i := range pkIndices {
		parts = append(parts, resource.Column(i).String())
	}

	return "(" + strings.Join(parts, ",") + ")"
}
