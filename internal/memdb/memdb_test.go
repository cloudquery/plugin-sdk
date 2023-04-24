package memdb

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/v2/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/cloudquery/plugin-sdk/v2/testdata"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var migrateStrategyOverwrite = destination.MigrateStrategy{
	AddColumn:           specs.MigrateModeForced,
	AddColumnNotNull:    specs.MigrateModeForced,
	RemoveColumn:        specs.MigrateModeForced,
	RemoveColumnNotNull: specs.MigrateModeForced,
	ChangeColumn:        specs.MigrateModeForced,
}

var migrateStrategyAppend = destination.MigrateStrategy{
	AddColumn:           specs.MigrateModeForced,
	AddColumnNotNull:    specs.MigrateModeForced,
	RemoveColumn:        specs.MigrateModeForced,
	RemoveColumnNotNull: specs.MigrateModeForced,
	ChangeColumn:        specs.MigrateModeForced,
}

func TestPluginUnmanagedClient(t *testing.T) {
	destination.PluginTestSuiteRunner(
		t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", NewClient)
		},
		specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		},
	)
}

func TestPluginManagedClient(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", NewClient, destination.WithManagedWriter())
		},
		specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithSmallBatchSize(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", NewClient, destination.WithManagedWriter(),
				destination.WithDefaultBatchSize(1),
				destination.WithDefaultBatchSizeBytes(1))
		}, specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithLargeBatchSize(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", NewClient, destination.WithManagedWriter(),
				destination.WithDefaultBatchSize(100000000),
				destination.WithDefaultBatchSizeBytes(100000000))
		},
		specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithCQPKs(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", NewClient)
		},
		specs.Destination{PKMode: specs.PKModeCQID},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginOnNewError(t *testing.T) {
	ctx := context.Background()
	p := destination.NewPlugin("test", "development", NewClientErrOnNew)
	err := p.Init(ctx, getTestLogger(t), specs.Destination{})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithErrOnWrite())
	p := destination.NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, getTestLogger(t), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	table := testdata.TestTable("test")
	tables := []*arrow.Schema{
		table.ToArrowSchema(),
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	ch := make(chan arrow.Record, 1)
	opts := testdata.GenTestDataOptions{
		SourceName: "test",
		SyncTime:   time.Now(),
		MaxRows:    1,
		StableUUID: uuid.Nil,
	}
	record := testdata.GenTestData(table.ToArrowSchema(), opts)[0]
	ch <- record
	close(ch)
	err := p.Write(ctx, sourceSpec, tables, syncTime, ch)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "errOnWrite" {
		t.Fatalf("expected errOnWrite, got %s", err.Error())
	}
}

func TestOnWriteCtxCancelled(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithBlockingWrite())
	p := destination.NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, getTestLogger(t), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	table := testdata.TestTable("test")
	tables := []*arrow.Schema{
		table.ToArrowSchema(),
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	ch := make(chan arrow.Record, 1)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	opts := testdata.GenTestDataOptions{
		SourceName: "test",
		SyncTime:   time.Now(),
		MaxRows:    1,
		StableUUID: uuid.Nil,
	}
	record := testdata.GenTestData(table.ToArrowSchema(), opts)[0]
	ch <- record
	defer cancel()
	err := p.Write(ctx, sourceSpec, tables, syncTime, ch)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPluginInit(t *testing.T) {
	const (
		batchSize      = 100
		batchSizeBytes = 1000
	)

	var (
		batchSizeObserved      int
		batchSizeBytesObserved int
	)
	p := destination.NewPlugin(
		"test",
		"development",
		func(ctx context.Context, logger zerolog.Logger, s specs.Destination) (destination.Client, error) {
			batchSizeObserved = s.BatchSize
			batchSizeBytesObserved = s.BatchSizeBytes
			return NewClient(ctx, logger, s)
		},
		destination.WithDefaultBatchSize(batchSize),
		destination.WithDefaultBatchSizeBytes(batchSizeBytes),
	)
	require.NoError(t, p.Init(context.TODO(), getTestLogger(t), specs.Destination{}))

	require.Equal(t, batchSize, batchSizeObserved)
	require.Equal(t, batchSizeBytes, batchSizeBytesObserved)
}
