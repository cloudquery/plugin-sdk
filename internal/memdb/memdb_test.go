package memdb

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/google/uuid"
)

var migrateStrategyOverwrite = plugin.MigrateStrategy{
	AddColumn:           plugin.MigrateModeForce,
	AddColumnNotNull:    plugin.MigrateModeForce,
	RemoveColumn:        plugin.MigrateModeForce,
	RemoveColumnNotNull: plugin.MigrateModeForce,
	ChangeColumn:        plugin.MigrateModeForce,
}

var migrateStrategyAppend = plugin.MigrateStrategy{
	AddColumn:           plugin.MigrateModeForce,
	AddColumnNotNull:    plugin.MigrateModeForce,
	RemoveColumn:        plugin.MigrateModeForce,
	RemoveColumnNotNull: plugin.MigrateModeForce,
	ChangeColumn:        plugin.MigrateModeForce,
}

func TestPluginUnmanagedClient(t *testing.T) {
	plugin.PluginTestSuiteRunner(
		t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewMemDBClient)
		},
		nil,
		plugin.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		},
	)
}

func TestPluginManagedClientWithCQPKs(t *testing.T) {
	plugin.PluginTestSuiteRunner(t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewMemDBClient)
		},
		pbPlugin.Spec{
			WriteSpec: &pbPlugin.WriteSpec{
				PkMode: pbPlugin.WriteSpec_CQ_ID_ONLY,
			},
		},
		plugin.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginOnNewError(t *testing.T) {
	ctx := context.Background()
	p := plugin.NewPlugin("test", "development", NewMemDBClientErrOnNew)
	err := p.Init(ctx, nil)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithErrOnWrite())
	p := plugin.NewPlugin("test", "development", newClientFunc)
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
	p := plugin.NewPlugin("test", "development", newClientFunc)
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
