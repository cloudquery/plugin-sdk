package pk

import (
	"strings"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/schema"
)

func String(resource arrow.Record) string {
	sc := resource.Schema()
	pkIndices := schema.PrimaryKeyIndices(sc)
	parts := make([]string, 0, len(pkIndices))
	for _, i := range pkIndices {
		parts = append(parts, resource.Column(i).String())
	}

	return "(" + strings.Join(parts, ",") + ")"
}
