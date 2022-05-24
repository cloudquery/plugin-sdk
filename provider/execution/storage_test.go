package execution

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/jackc/pgx/v4"
)

type noopStorage struct {
	D schema.Dialect
}

func (f noopStorage) Begin(ctx context.Context) (TXQueryExecer, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f noopStorage) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (f noopStorage) Exec(ctx context.Context, query string, args ...interface{}) error {
	return nil
}

func (f noopStorage) Insert(ctx context.Context, t *schema.Table, instance schema.Resources, shouldCascade bool, cascadeDeleteFilters map[string]interface{}) error {
	return nil
}

func (f noopStorage) Delete(ctx context.Context, t *schema.Table, kvFilters []interface{}) error {
	return nil
}

func (f noopStorage) RemoveStaleData(ctx context.Context, t *schema.Table, executionStart time.Time, kvFilters []interface{}) error {
	return nil
}

func (f noopStorage) CopyFrom(ctx context.Context, resources schema.Resources, shouldCascade bool, cascadeDeleteFilters map[string]interface{}) error {
	return nil
}

func (f noopStorage) RawCopyFrom(ctx context.Context, r io.Reader, sql string) error {
	return nil
}

func (f noopStorage) RawCopyTo(ctx context.Context, w io.Writer, sql string) error {
	return nil
}

func (f noopStorage) Close() {}

func (f noopStorage) Dialect() schema.Dialect {
	if f.D != nil {
		return f.D
	}
	return noopDialect{}
}

type noopDialect struct {
}

func (d noopDialect) PrimaryKeys(t *schema.Table) []string {
	return t.Options.PrimaryKeys
}

func (d noopDialect) Columns(t *schema.Table) schema.ColumnList {
	return t.Columns
}

func (d noopDialect) Constraints(t, parent *schema.Table) []string {
	return []string{}
}

func (d noopDialect) Extra(t, parent *schema.Table) []string {
	return []string{}
}

func (d noopDialect) DBTypeFromType(v schema.ValueType) string {
	return v.String()
}

func (d noopDialect) GetResourceValues(r *schema.Resource) ([]interface{}, error) {
	return r.Values()
}

var _ Storage = (*noopStorage)(nil)
