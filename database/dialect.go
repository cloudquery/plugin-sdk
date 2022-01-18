package database

import (
	"strings"

	"github.com/cloudquery/cq-provider-sdk/database/dsn"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

// ParseDialectDSN parses a DSN and returns the suggested DialectType, as well as a new version of the DSN if applicable.
// The DSN change is done to support protocol-compatible databases without needing to add support for custom URL schemes to 3rd party packages.
func ParseDialectDSN(inputDSN string) (d schema.DialectType, newDSN string, err error) {
	u, err := dsn.ParseConnectionString(inputDSN)
	if err != nil {
		return schema.Postgres, inputDSN, err
	}

	switch u.Scheme {
	case "timescaledb", "tsdb", "timescale":
		// Replace tsdb schemes to look like postgres, so that postgres-protocol compatible tools (like go-migrate) work
		// Keep/return the DialectType separately from the DSN so we can refer to it later
		fixedDSN := strings.Replace(u.String(), u.Scheme+"://", "postgres://", 1)
		return schema.TSDB, fixedDSN, nil
	default:
		return schema.Postgres, inputDSN, nil
	}
}
