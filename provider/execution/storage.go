package execution

import (
	"context"
	"io"
	"time"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/georgysavva/scany/pgxscan"
)

//go:generate mockgen -package=mock -destination=../schema/mock/mock_storage.go . Storage
type Storage interface {
	QueryExecer
	Copier
	TXer
	Insert(ctx context.Context, t *schema.Table, instance schema.Resources, shouldCascade bool) error
	Delete(ctx context.Context, t *schema.Table, kvFilters []interface{}) error
	RemoveStaleData(ctx context.Context, t *schema.Table, executionStart time.Time, kvFilters []interface{}) error
	CopyFrom(ctx context.Context, resources schema.Resources, shouldCascade bool) error
	Close()
	Dialect() schema.Dialect
}

type QueryExecer interface {
	pgxscan.Querier
	Exec(ctx context.Context, query string, args ...interface{}) error
}

type Copier interface {
	RawCopyTo(ctx context.Context, w io.Writer, sql string) error
	RawCopyFrom(ctx context.Context, r io.Reader, sql string) error
}

type TXer interface {
	Begin(context.Context) (TXQueryExecer, error)
}

type TXQueryExecer interface {
	QueryExecer
	TXer
	Rollback(context.Context) error
	Commit(context.Context) error
}
