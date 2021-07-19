package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/modern-go/reflect2"

	"github.com/doug-martin/goqu/v9"

	"github.com/jackc/pgx/v4"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// MaxTableLength in postgres is 63 when building _fk or _pk we want to truncate the name to 60 chars max
	maxTableNamePKConstraint = 60
)

//go:generate mockgen -package=mock -destination=./mocks/mock_database.go . Database
type Database interface {
	Insert(ctx context.Context, t *Table, instance []*Resource) error
	Exec(ctx context.Context, query string, args ...interface{}) error
	Delete(ctx context.Context, t *Table, args []interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	CopyFrom(ctx context.Context, resources Resources, shouldCascade bool, CascadeDeleteFilters map[string]interface{}) error
}

type PgDatabase struct {
	pool *pgxpool.Pool
}

func NewPgDatabase(ctx context.Context, dsn string) (*PgDatabase, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &PgDatabase{pool: pool}, nil
}

// Insert inserts all resources to given table, table and resources are assumed from same table.
func (p PgDatabase) Insert(ctx context.Context, t *Table, resources []*Resource) error {
	if len(resources) == 0 {
		return nil
	}
	// It is safe to assume that all resources have the same columns
	cols := quoteColumns(resources[0].columns)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlStmt := psql.Insert(t.Name).Columns(cols...).Suffix(fmt.Sprintf("ON CONFLICT ON CONSTRAINT %s_pk DO UPDATE SET %s", TruncateTableConstraint(t.Name), buildReplaceColumns(cols)))
	for _, res := range resources {
		if res.table != t {
			return fmt.Errorf("resource table expected %s got %s", t.Name, res.table.Name)
		}
		values, err := res.Values()
		if err != nil {
			return fmt.Errorf("table %s insert failed %w", t.Name, err)
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

// CopyFrom copies all resources from []*Resource
func (p PgDatabase) CopyFrom(ctx context.Context, resources Resources, shouldCascade bool, cascadeDeleteFilters map[string]interface{}) error {
	if len(resources) == 0 {
		return nil
	}
	err := p.pool.BeginTxFunc(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	}, func(tx pgx.Tx) error {
		if shouldCascade {
			q := goqu.Dialect("postgres").Delete(resources.TableName()).Where(goqu.Ex{"cq_id": resources.GetIds()})
			for k, v := range cascadeDeleteFilters {
				q = q.Where(goqu.Ex{k: goqu.Op{"eq": v}})
			}
			sql, args, err := q.Prepared(true).ToSQL()
			if err != nil {
				return err
			}
			_, err = tx.Exec(ctx, sql, args...)
			if err != nil {
				return err
			}
		}
		copied, err := tx.CopyFrom(
			ctx, pgx.Identifier{resources.TableName()}, resources.ColumnNames(),
			pgx.CopyFromSlice(len(resources), func(i int) ([]interface{}, error) {
				// use getResourceValues instead of Resource.Values since values require some special encoding for CopyFrom
				return getResourceValues(resources[i])
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
	case TypeInet:
		return "inet"
	case TypeMacAddr:
		return "mac"
	case TypeInetArray:
		return "inet[]"
	case TypeMacAddrArray:
		return "mac[]"
	case TypeCIDR:
		return "cidr"
	case TypeCIDRArray:
		return "cidr[]"
	default:
		panic("invalid type")
	}
}

func quoteColumns(columns []string) []string {
	for i, v := range columns {
		columns[i] = strconv.Quote(v)
	}
	return columns
}

func buildReplaceColumns(columns []string) string {
	replaceColumns := make([]string, len(columns))
	for i, c := range columns {
		replaceColumns[i] = fmt.Sprintf("%[1]s = EXCLUDED.%[1]s", c)
	}
	return strings.Join(replaceColumns, ",")
}

func TruncateTableConstraint(name string) string {
	if len(name) > maxTableNamePKConstraint {
		return name[:maxTableNamePKConstraint]
	}
	return name
}

func getResourceValues(r *Resource) ([]interface{}, error) {
	values := make([]interface{}, 0)
	for _, c := range append(r.table.Columns, GetDefaultSDKColumns()...) {
		v := r.Get(c.Name)
		if err := c.ValidateType(v); err != nil {
			return nil, err
		}
		if c.Type == TypeJSON {
			if v == nil {
				values = append(values, v)
				continue
			}
			if reflect2.TypeOf(v).Kind() == reflect.Map {
				values = append(values, v)
				continue
			}
			switch data := v.(type) {
			case map[string]interface{}:
				values = append(values, data)
			case string:
				newV := make(map[string]interface{})
				err := json.Unmarshal([]byte(data), &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			case *string:
				var newV interface{}
				err := json.Unmarshal([]byte(*data), &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			case []byte:
				var newV interface{}
				err := json.Unmarshal(data, &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			}
		} else {
			values = append(values, v)
		}
	}
	for _, v := range r.extraFields {
		values = append(values, v)
	}
	return values, nil
}
