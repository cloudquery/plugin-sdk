package migration

import (
	"context"
	"strconv"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

// CreateTableDefinitions reads schema.Table and builds the CREATE TABLE statement for it, also processing and returning subrelation tables
func CreateTableDefinitions(ctx context.Context, dialect schema.Dialect, t *schema.Table, parent *schema.Table) ([]string, error) {
	b := &strings.Builder{}

	// Build a SQL to create a table
	b.WriteString("CREATE TABLE IF NOT EXISTS " + strconv.Quote(t.Name) + " (\n")

	for _, c := range dialect.Columns(t) {
		b.WriteByte('\t')
		b.WriteString(strconv.Quote(c.Name) + " " + dialect.DBTypeFromType(c.Type))
		if c.CreationOptions.NotNull {
			b.WriteString(" NOT NULL")
		}
		// c.CreationOptions.Unique is handled in the Constraints() call below
		b.WriteString(",\n")
	}

	cons := dialect.Constraints(t, parent)
	for i, cn := range cons {
		b.WriteByte('\t')
		b.WriteString(cn)

		if i < len(cons)-1 {
			b.WriteByte(',')
		}

		b.WriteByte('\n')
	}

	b.WriteString(");")

	up := make([]string, 0, 1+len(t.Relations))
	up = append(up, b.String())
	up = append(up, dialect.Extra(t, parent)...)

	// Create relation tables
	for _, r := range t.Relations {
		if cr, err := CreateTableDefinitions(ctx, dialect, r, t); err != nil {
			return nil, err
		} else {
			up = append(up, cr...)
		}
	}

	return up, nil
}
