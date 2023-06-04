package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/google/uuid"
)

var migrateStrategyOverwrite = MigrateStrategy{
	AddColumn:           MigrateModeForce,
	AddColumnNotNull:    MigrateModeForce,
	RemoveColumn:        MigrateModeForce,
	RemoveColumnNotNull: MigrateModeForce,
	ChangeColumn:        MigrateModeForce,
}

var migrateStrategyAppend = MigrateStrategy{
	AddColumn:           MigrateModeForce,
	AddColumnNotNull:    MigrateModeForce,
	RemoveColumn:        MigrateModeForce,
	RemoveColumnNotNull: MigrateModeForce,
	ChangeColumn:        MigrateModeForce,
}

func TestPluginUnmanagedClient(t *testing.T) {
	PluginTestSuiteRunner(
		t,
		func() *Plugin {
			return NewPlugin("test", "development", NewMemDBClient)
		},
		nil,
		PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		},
	)
}

func TestPluginManagedClient(t *testing.T) {
	PluginTestSuiteRunner(t,
		func() *Plugin {
			return NewPlugin("test", "development", NewMemDBClient, WithManagedWriter())
		},
		nil,
		PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithSmallBatchSize(t *testing.T) {
	PluginTestSuiteRunner(t,
		func() *Plugin {
			return NewPlugin("test", "development", NewMemDBClient, WithManagedWriter(),
				WithDefaultBatchSize(1),
				WithDefaultBatchSizeBytes(1))
		}, nil,
		PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithLargeBatchSize(t *testing.T) {
	PluginTestSuiteRunner(t,
		func() *Plugin {
			return NewPlugin("test", "development", NewMemDBClient, WithManagedWriter(),
				WithDefaultBatchSize(100000000),
				WithDefaultBatchSizeBytes(100000000))
		},
		nil,
		PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithCQPKs(t *testing.T) {
	PluginTestSuiteRunner(t,
		func() *Plugin {
			return NewPlugin("test", "development", NewMemDBClient)
		},
		pbPlugin.Spec{
			WriteSpec: &pbPlugin.WriteSpec{
				PkMode: pbPlugin.WriteSpec_CQ_ID_ONLY,
			},
		},
		PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginOnNewError(t *testing.T) {
	ctx := context.Background()
	p := NewPlugin("test", "development", NewMemDBClientErrOnNew)
	err := p.Init(ctx, nil)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithErrOnWrite())
	p := NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, nil); err != nil {
		t.Fatal(err)
	}
	table := schema.TestTable("test", schema.TestSourceOptions{})
	tables := schema.Tables{
		table,
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	sourceSpec := pbPlugin.Spec{
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
	p := NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, pbPlugin.Spec{
		WriteSpec: &pbPlugin.WriteSpec{},
	}); err != nil {
		t.Fatal(err)
	}
	table := schema.TestTable("test", schema.TestSourceOptions{})
	tables := schema.Tables{
		table,
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	sourceSpec := pbPlugin.Spec{
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
