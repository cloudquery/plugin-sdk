package schema

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database interface {
	Insert(ctx context.Context, t *Table, instance []*Resource) error
	Exec(ctx context.Context, query string, args ...interface{}) error
	Delete(ctx context.Context, t *Table, args []interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type PgDatabase struct {
	pool *pgxpool.Pool
}

func NewPgDatabase(dsn string) (*PgDatabase, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return &PgDatabase{pool: pool}, nil
}

// Insert inserts all resources to given table, table and resources are assumed from same table.
func (p PgDatabase) Insert(ctx context.Context, t *Table, resources []*Resource) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlStmt := psql.Insert(t.Name).Columns(t.ColumnNames()...)
	for _, res := range resources {
		if res.table != t {
			return fmt.Errorf("resource table expected %s got %s", t.Name, res.table.Name)
		}
		values, err := res.Values()
		if err != nil {
			return err
		}
		sqlStmt = sqlStmt.Values(values...)
	}

	s, args, err := sqlStmt.ToSql()
	if err != nil {
		return err
	}
	_, err = p.pool.Exec(ctx, s, args...)
	return err
}

// Exec allows executions of postgres queries with given args returning error of execution
func (p PgDatabase) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.pool.Exec(ctx, query, args...)
	return err
}

// Query  allows execution of postgres queries with given args returning data result
func (p PgDatabase) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := p.pool.Query(ctx, query, args...)
	return rows, err
}

// Query  allows execution of postgres queries with given args returning data result
func (p PgDatabase) QueryOne(ctx context.Context, query string, args ...interface{}) pgx.Row {
	row := p.pool.QueryRow(ctx, query, args...)
	return row
}

func (p PgDatabase) Delete(ctx context.Context, t *Table, args []interface{}) error {
	nc := len(args)
	if nc%2 != 0 {
		return fmt.Errorf("number of args to delete should be even. Got %d", nc)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	ds := psql.Delete(t.Name)
	for i := 0; i < nc; i += 2 {
		ds = ds.Where(sq.Eq{args[i].(string): args[i+1]})
	}
	sql, args, err := ds.ToSql()
	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, sql, args...)
	return err
}

func GetPgTypeFromType(v ValueType) string {
	switch v {
	case TypeBool:
		return "boolean"
	case TypeInt:
		return "integer"
	case TypeBigInt:
		return "bigint"
	case TypeSmallInt:
		return "smallint"
	case TypeFloat:
		return "float"
	case TypeUUID:
		return "uuid"
	case TypeString:
		return "text"
	case TypeJSON:
		return "jsonb"
	case TypeIntArray:
		return "integer[]"
	case TypeStringArray:
		return "text[]"
	case TypeTimestamp:
		return "timestamp without time zone"
	case TypeByteArray:
		return "bytea"
	case TypeInvalid:
		fallthrough
	default:
		panic("invalid type")
	}
}
