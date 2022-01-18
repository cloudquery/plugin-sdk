package postgres

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect connects to the given DSN and returns a pgxpool
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		UUIDType := pgtype.DataType{
			Value: &UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		}

		conn.ConnInfo().RegisterDataType(UUIDType)
		return nil
	}
	poolCfg.LazyConnect = true
	return pgxpool.ConnectConfig(ctx, poolCfg)
}
