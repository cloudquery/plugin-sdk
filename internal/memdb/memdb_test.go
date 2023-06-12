package memdb

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

func TestPlugin(t *testing.T) {
	ctx := context.Background()
	p := plugin.NewPlugin("test", "development", NewMemDBClient)
	if err := p.Init(ctx, nil); err != nil {
		t.Fatal(err)
	}
	plugin.PluginTestSuiteRunner(
		t,
		p,
		plugin.PluginTestSuiteTests{
			MigrateStrategy: plugin.MigrateStrategy{
				AddColumn:           plugin.MigrateModeForce,
				AddColumnNotNull:    plugin.MigrateModeForce,
				RemoveColumn:        plugin.MigrateModeForce,
				RemoveColumnNotNull: plugin.MigrateModeForce,
				ChangeColumn:        plugin.MigrateModeForce,
			},
		},
	)
}

// func TestPluginOnNewError(t *testing.T) {
// 	ctx := context.Background()
// 	p := plugin.NewPlugin("test", "development", NewMemDBClientErrOnNew)
// 	err := p.Init(ctx, nil)

// 	if err == nil {
// 		t.Fatal("expected error")
// 	}
// }

// func TestOnWriteError(t *testing.T) {
// 	ctx := context.Background()
// 	newClientFunc := GetNewClient(WithErrOnWrite())
// 	p := plugin.NewPlugin("test", "development", newClientFunc)
// 	if err := p.Init(ctx, nil); err != nil {
// 		t.Fatal(err)
// 	}
// 	table := schema.TestTable("test", schema.TestSourceOptions{})
// 	tables := schema.Tables{
// 		table,
// 	}
// 	sourceName := "TestDestinationOnWriteError"
// 	syncTime := time.Now()
// 	sourceSpec := pbPlugin.Spec{
// 		Name: sourceName,
// 	}
// 	ch := make(chan arrow.Record, 1)
// 	opts := schema.GenTestDataOptions{
// 		SourceName: "test",
// 		SyncTime:   time.Now(),
// 		MaxRows:    1,
// 		StableUUID: uuid.Nil,
// 	}
// 	record := schema.GenTestData(table, opts)[0]
// 	ch <- record
// 	close(ch)
// 	err := p.Write(ctx, sourceSpec, tables, syncTime, ch)
// 	if err == nil {
// 		t.Fatal("expected error")
// 	}
// 	if err.Error() != "errOnWrite" {
// 		t.Fatalf("expected errOnWrite, got %s", err.Error())
// 	}
// }

// func TestOnWriteCtxCancelled(t *testing.T) {
// 	ctx := context.Background()
// 	newClientFunc := GetNewClient(WithBlockingWrite())
// 	p := plugin.NewPlugin("test", "development", newClientFunc)
// 	if err := p.Init(ctx, pbPlugin.Spec{
// 		WriteSpec: &pbPlugin.WriteSpec{},
// 	}); err != nil {
// 		t.Fatal(err)
// 	}
// 	table := schema.TestTable("test", schema.TestSourceOptions{})
// 	tables := schema.Tables{
// 		table,
// 	}
// 	sourceName := "TestDestinationOnWriteError"
// 	syncTime := time.Now()
// 	sourceSpec := pbPlugin.Spec{
// 		Name: sourceName,
// 	}
// 	ch := make(chan arrow.Record, 1)
// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	opts := schema.GenTestDataOptions{
// 		SourceName: "test",
// 		SyncTime:   time.Now(),
// 		MaxRows:    1,
// 		StableUUID: uuid.Nil,
// 	}
// 	record := schema.GenTestData(table, opts)[0]
// 	ch <- record
// 	defer cancel()
// 	err := p.Write(ctx, sourceSpec, tables, syncTime, ch)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
