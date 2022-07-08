package postgres

import (
	"context"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/database/dsn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect connects to the given DSN and returns a pgxpool
func Connect(ctx context.Context, dsnURI string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsnURI)
	if err != nil {
		return nil, dsn.RedactParseError(err)
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// force a known search_path if DSN doesn't specify one
		if !strings.Contains(dsnURI, "&search_path=") && !strings.Contains(dsnURI, "?search_path=") {
			_, err := conn.Exec(ctx, "SET search_path=public")
			if err != nil {
				return err
			}
		}

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
