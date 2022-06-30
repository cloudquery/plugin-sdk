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

type noopDialect struct {
}

var _ Storage = (*noopStorage)(nil)

func (noopStorage) Begin(ctx context.Context) (TXQueryExecer, error) {
	return nil, fmt.Errorf("not implemented")
}

func (noopStorage) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (noopStorage) Exec(ctx context.Context, query string, args ...interface{}) error {
	return nil
}

func (noopStorage) Insert(ctx context.Context, t *schema.Table, instance schema.Resources, shouldCascade bool) error {
	return nil
}

func (noopStorage) Delete(ctx context.Context, t *schema.Table, kvFilters []interface{}) error {
	return nil
}

func (noopStorage) RemoveStaleData(ctx context.Context, t *schema.Table, executionStart time.Time, kvFilters []interface{}) error {
	return nil
}

func (noopStorage) CopyFrom(ctx context.Context, resources schema.Resources, shouldCascade bool) error {
	return nil
}

func (noopStorage) RawCopyFrom(ctx context.Context, r io.Reader, sql string) error {
	return nil
}

func (noopStorage) RawCopyTo(ctx context.Context, w io.Writer, sql string) error {
	return nil
}

func (noopStorage) Close() {}

func (f noopStorage) Dialect() schema.Dialect {
	if f.D != nil {
		return f.D
	}
	return noopDialect{}
}

func (noopDialect) PrimaryKeys(t *schema.Table) []string {
	return t.Options.PrimaryKeys
}

func (noopDialect) Columns(t *schema.Table) schema.ColumnList {
	return t.Columns
}

func (noopDialect) Constraints(t, parent *schema.Table) []string {
	return []string{}
}

func (noopDialect) Extra(t, parent *schema.Table) []string {
	return []string{}
}

func (noopDialect) DBTypeFromType(v schema.ValueType) string {
	return v.String()
}

func (noopDialect) GetResourceValues(r *schema.Resource) ([]interface{}, error) {
	return r.Values()
}
