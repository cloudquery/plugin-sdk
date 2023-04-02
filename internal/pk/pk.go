package pk

import (
	"strings"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/schema"
)

func String(table *schema.Table, resource arrow.Record) string {
	parts := make([]string, 0, len(table.PrimaryKeys()))
	for i, col := range table.Columns {
		if !col.CreationOptions.PrimaryKey {
			continue
		}
		parts = append(parts, resource.Column(i).String())
	}

	return "(" + strings.Join(parts, ",") + ")"
}
