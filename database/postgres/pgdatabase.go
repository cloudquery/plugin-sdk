package postgres

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cloudquery/cq-provider-sdk/provider/diag"

	sq "github.com/Masterminds/squirrel"
	"github.com/cloudquery/cq-provider-sdk/provider/execution"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cast"
)

type PgDatabase struct {
	pool *pgxpool.Pool
	log  hclog.Logger
	sd   schema.Dialect
}

func NewPgDatabase(ctx context.Context, logger hclog.Logger, dsn string, sd schema.Dialect) (*PgDatabase, error) {
	pool, err := Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &PgDatabase{
		pool: pool,
		log:  logger,
		sd:   sd,
	}, nil
}

var _ execution.Storage = (*PgDatabase)(nil)

// Insert inserts all resources to given table, table and resources are assumed from same table.
func (p PgDatabase) Insert(ctx context.Context, t *schema.Table, resources schema.Resources, shouldCascade bool, cascadeDeleteFilters map[string]interface{}) error {
	if len(resources) == 0 {
		return nil
	}

	// It is safe to assume that all resources have the same columns
	cols := quoteColumns(resources.ColumnNames())
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlStmt := psql.Insert(t.Name).Columns(cols...)
	for _, res := range resources {
		if res.TableName() != t.Name {
			return fmt.Errorf("resource table expected %s got %s", t.Name, res.TableName())
		}
		values, err := res.Values()
		if err != nil {
			return fmt.Errorf("table %s insert failed %w", t.Name, err)
		}
		sqlStmt = sqlStmt.Values(values...)
	}
	if t.Global {
		updateColumns := make([]string, len(cols))
		for i, c := range cols {
			updateColumns[i] = fmt.Sprintf("%[1]s = excluded.%[1]s", c)
		}
		sqlStmt = sqlStmt.Suffix(fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s",
			strings.Join(p.sd.PrimaryKeys(t), ","), strings.Join(updateColumns, ",")))
	}

	s, args, err := sqlStmt.ToSql()
	if err != nil {
		return diag.NewBaseError(err, diag.DATABASE, diag.WithResourceName(t.Name), diag.WithSummary("bad insert SQL statement created"), diag.WithDetails("SQL statement %q is invalid", s))
	}

	err = p.pool.BeginTxFunc(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	}, func(tx pgx.Tx) error {
		if shouldCascade {
			if err := deleteResourceByCQId(ctx, tx, resources, cascadeDeleteFilters); err != nil {
				return err
			}
		}

		_, err := tx.Exec(ctx, s, args...)
		return err
	})
	if err == nil {
		return nil
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		// This should rarely occur, but if it occurs we want to print the SQL to debug it further
		if pgerrcode.IsSyntaxErrororAccessRuleViolation(pgErr.Code) {
			p.log.Debug("insert syntax error", "sql", s)
		}
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			p.log.Debug("insert integrity violation error", "constraint", pgErr.ConstraintName, "errMsg", pgErr.Message)
		}
		return diag.NewBaseError(err, diag.DATABASE, diag.WithResourceName(t.Name), diag.WithSummary("failed to insert to table %q", t.Name), diag.WithDetails("%s", pgErr.Message))
	}
	return diag.NewBaseError(err, diag.DATABASE, diag.WithResourceName(t.Name))
}

// CopyFrom copies all resources from []*Resource
func (p PgDatabase) CopyFrom(ctx context.Context, resources schema.Resources, shouldCascade bool, cascadeDeleteFilters map[string]interface{}) error {
	if len(resources) == 0 {
		return nil
	}
	err := p.pool.BeginTxFunc(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	}, func(tx pgx.Tx) error {
		if shouldCascade {
			if err := deleteResourceByCQId(ctx, tx, resources, cascadeDeleteFilters); err != nil {
				return err
			}
		}
		copied, err := tx.CopyFrom(
			ctx, pgx.Identifier{resources.TableName()}, resources.ColumnNames(),
			pgx.CopyFromSlice(len(resources), func(i int) ([]interface{}, error) {
				// use getResourceValues instead of Resource.Values since values require some special encoding for CopyFrom
				return p.sd.GetResourceValues(resources[i])
			}))
		if err != nil {
			return err
		}
		if copied != int64(len(resources)) {
			return fmt.Errorf("not all resources copied %d != %d to %s", copied, len(resources), resources.TableName())
		}
		return nil
	})
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

// QueryOne  allows execution of postgres queries with given args returning data result
func (p PgDatabase) QueryOne(ctx context.Context, query string, args ...interface{}) pgx.Row {
	row := p.pool.QueryRow(ctx, query, args...)
	return row
}

func (p PgDatabase) Delete(ctx context.Context, t *schema.Table, kvFilters []interface{}) error {
	nc := len(kvFilters)
	if nc%2 != 0 {
		return fmt.Errorf("number of args to delete should be even. Got %d", nc)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	ds := psql.Delete(t.Name)
	for i := 0; i < nc; i += 2 {
		ds = ds.Where(sq.Eq{kvFilters[i].(string): kvFilters[i+1]})
	}
	sql, args, err := ds.ToSql()
	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, sql, args...)
	return err
}

func (p PgDatabase) RemoveStaleData(ctx context.Context, t *schema.Table, executionStart time.Time, kvFilters []interface{}) error {
	q := goqu.Delete(t.Name).WithDialect("postgres").Where(goqu.L(`extract(epoch from (cq_meta->>'last_updated')::timestamp)`).Lt(executionStart.Unix()))
	if len(kvFilters)%2 != 0 {
		return fmt.Errorf("expected even number of k,v delete filters received %s", kvFilters)
	}
	for i := 0; i < len(kvFilters); i += 2 {
		q = q.Where(goqu.Ex{cast.ToString(kvFilters[i]): goqu.Op{"eq": kvFilters[i+1]}})
	}
	sql, args, err := q.Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("failed building query: %w", err)
	}
	_, err = p.pool.Exec(ctx, sql, args...)
	return err
}

func (p PgDatabase) Close() {
	p.pool.Close()
}

func (p PgDatabase) RawCopyTo(ctx context.Context, w io.Writer, sql string) error {
	c, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()
	_, err = c.Conn().PgConn().CopyTo(ctx, w, sql)
	return err
}
func (p PgDatabase) RawCopyFrom(ctx context.Context, r io.Reader, sql string) error {
	c, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()
	_, err = c.Conn().PgConn().CopyFrom(ctx, r, sql)
	return err
}

func (p PgDatabase) Dialect() schema.Dialect {
	return p.sd
}

func (p PgDatabase) Begin(ctx context.Context) (execution.TXQueryExecer, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgTx{tx}, nil
}

type PgTx struct {
	pgx.Tx
}

func (p PgTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, v := p.Tx.Exec(ctx, query, args...)
	return v
}

func (p PgTx) Begin(ctx context.Context) (execution.TXQueryExecer, error) {
	v, err := p.Tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgTx{v}, nil
}

func quoteColumns(columns []string) []string {
	ret := make([]string, len(columns))
	for i, v := range columns {
		ret[i] = strconv.Quote(v)
	}
	return ret
}

func deleteResourceByCQId(ctx context.Context, tx pgx.Tx, resources schema.Resources, cascadeDeleteFilters map[string]interface{}) error {
	q := goqu.Dialect("postgres").Delete(resources.TableName()).Where(goqu.Ex{"cq_id": resources.GetIds()})
	for k, v := range cascadeDeleteFilters {
		q = q.Where(goqu.Ex{k: goqu.Op{"eq": v}})
	}
	sql, args, err := q.Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}
