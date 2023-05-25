package memdb

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination/batchingwriter"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/google/uuid"
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
			return destination.NewPlugin("test", "development", NewClient, destination.WithManagedWriter(batchingwriter.New()))
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
			return destination.NewPlugin("test", "development", NewClient,
				destination.WithManagedWriter(batchingwriter.New(batchingwriter.WithDefaultBatchSize(1, 1))),
			)
		}, specs.Destination{},
		destination.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithLargeBatchSize(t *testing.T) {
	destination.PluginTestSuiteRunner(t,
		func() *destination.Plugin {
			return destination.NewPlugin("test", "development", NewClient,
				destination.WithManagedWriter(batchingwriter.New(batchingwriter.WithDefaultBatchSize(100000000, 100000000))),
			)
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
	table := schema.TestTable("test", schema.TestSourceOptions{})
	tables := schema.Tables{
		table,
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	ch := make(chan arrow.Record, 1)
	opts := schema.GenTestDataOptions{
		SourceName: "test",
		SyncTime:   time.Now(),
		MaxRows:    1,
		StableUUID: uuid.Nil,
	}
	record := schema.GenTestData(table, opts)[0]
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
	table := schema.TestTable("test", schema.TestSourceOptions{})
	tables := schema.Tables{
		table,
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	ch := make(chan arrow.Record, 1)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	opts := schema.GenTestDataOptions{
		SourceName: "test",
		SyncTime:   time.Now(),
		MaxRows:    1,
		StableUUID: uuid.Nil,
	}
	record := schema.GenTestData(table, opts)[0]
	ch <- record
	defer cancel()
	err := p.Write(ctx, sourceSpec, tables, syncTime, ch)
	if err != nil {
		t.Fatal(err)
	}
}
