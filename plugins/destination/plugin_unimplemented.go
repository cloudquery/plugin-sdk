package destination

import (
	"context"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
)

type UnimplementedUnmanagedWriter struct{}

func (UnimplementedUnmanagedWriter) Write(context.Context, specs.Source, schema.Tables, time.Time, <-chan arrow.Record) error {
	panic("Write not implemented")
}

func (UnimplementedUnmanagedWriter) Metrics() Metrics {
	panic("Metrics not implemented")
}

type UnimplementedManagedWriter struct{}

func (UnimplementedManagedWriter) WriteTableBatch(context.Context, specs.Source, *schema.Table, time.Time, []arrow.Record) error {
	panic("WriteTableBatch not implemented")
}

type UnimplementedRead struct{}

func (UnimplementedRead) Read(context.Context, *schema.Table, string, chan<- arrow.Record) error {
	panic("Read not implemented")
}

type UnimplementedMigrate struct{}

func (UnimplementedMigrate) Migrate(context.Context, schema.Tables) error {
	return nil // Special case, we don't want to error here
}

type UnimplementedDeleteStale struct{}

func (UnimplementedDeleteStale) DeleteStale(context.Context, schema.Tables, string, time.Time) error {
	return nil // Special case, we don't want to error here
}

type UnimplementedClose struct{}

func (UnimplementedClose) Close(context.Context) error {
	return nil // Special case, we don't want to error here
}

type UnimplementedClient struct {
	UnimplementedMigrate
	UnimplementedRead
	UnimplementedDeleteStale
	UnimplementedClose
	UnimplementedManagedWriter
	UnimplementedUnmanagedWriter
}

var (
	_ UnmanagedWriter = UnimplementedUnmanagedWriter{}
	_ ManagedWriter   = UnimplementedManagedWriter{}
	_ Client          = UnimplementedClient{}
)
