package database

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/database/postgres"
	"github.com/cloudquery/cq-provider-sdk/provider/execution"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/hashicorp/go-hclog"
)

// DB encapsulates a schema.Storage and the (auto-detected) dialect it was configured with
type DB struct {
	execution.Storage

	dialectType schema.DialectType
}

// New creates a new DB using the provided DSN. It will auto-detect the dialect based on the DSN and pass that info to NewPgDatabase
func New(ctx context.Context, logger hclog.Logger, dsn string) (*DB, error) {
	dType, newDSN, err := ParseDialectDSN(dsn)
	if err != nil {
		return nil, err
	}

	dialect, err := schema.GetDialect(dType)
	if err != nil {
		return nil, err
	}

	db, err := postgres.NewPgDatabase(ctx, logger, newDSN, dialect)
	if err != nil {
		return nil, err
	}

	return &DB{
		Storage:     db,
		dialectType: dType,
	}, nil
}

// DialectType returns the dialect type the DB was configured with
func (d *DB) DialectType() schema.DialectType {
	return d.dialectType
}
