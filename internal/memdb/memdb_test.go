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
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var migrateStrategyOverwrite = plugin.MigrateStrategy{
	AddColumn:           pbPlugin.WriteSpec_FORCE,
	AddColumnNotNull:    pbPlugin.WriteSpec_FORCE,
	RemoveColumn:        pbPlugin.WriteSpec_FORCE,
	RemoveColumnNotNull: pbPlugin.WriteSpec_FORCE,
	ChangeColumn:        pbPlugin.WriteSpec_FORCE,
}

var migrateStrategyAppend = plugin.MigrateStrategy{
	AddColumn:           pbPlugin.WriteSpec_FORCE,
	AddColumnNotNull:    pbPlugin.WriteSpec_FORCE,
	RemoveColumn:        pbPlugin.WriteSpec_FORCE,
	RemoveColumnNotNull: pbPlugin.WriteSpec_FORCE,
	ChangeColumn:        pbPlugin.WriteSpec_FORCE,
}

func TestPluginUnmanagedClient(t *testing.T) {
	plugin.PluginTestSuiteRunner(
		t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewClient)
		},
		pbPlugin.Spec{},
		plugin.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		},
	)
}

func TestPluginManagedClient(t *testing.T) {
	plugin.PluginTestSuiteRunner(t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewClient, plugin.WithManagedWriter())
		},
		pbPlugin.Spec{},
		plugin.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithSmallBatchSize(t *testing.T) {
	plugin.PluginTestSuiteRunner(t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewClient, plugin.WithManagedWriter(),
				plugin.WithDefaultBatchSize(1),
				plugin.WithDefaultBatchSizeBytes(1))
		}, pbPlugin.Spec{},
		plugin.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithLargeBatchSize(t *testing.T) {
	plugin.PluginTestSuiteRunner(t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewClient, plugin.WithManagedWriter(),
				plugin.WithDefaultBatchSize(100000000),
				plugin.WithDefaultBatchSizeBytes(100000000))
		},
		pbPlugin.Spec{},
		plugin.PluginTestSuiteTests{
			MigrateStrategyOverwrite: migrateStrategyOverwrite,
			MigrateStrategyAppend:    migrateStrategyAppend,
		})
}

func TestPluginManagedClientWithCQPKs(t *testing.T) {
	plugin.PluginTestSuiteRunner(t,
		func() *plugin.Plugin {
			return plugin.NewPlugin("test", "development", NewClient)
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
	p := plugin.NewPlugin("test", "development", NewClientErrOnNew)
	err := p.Init(ctx, pbPlugin.Spec{})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithErrOnWrite())
	p := plugin.NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, pbPlugin.Spec{}); err != nil {
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
	if err := p.Init(ctx, pbPlugin.Spec{}); err != nil {
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

func TestPluginInit(t *testing.T) {
	const (
		batchSize      = 100
		batchSizeBytes = 1000
	)

	var (
		batchSizeObserved      uint64
		batchSizeBytesObserved uint64
	)
	p := plugin.NewPlugin(
		"test",
		"development",
		func(ctx context.Context, logger zerolog.Logger, s pbPlugin.Spec) (plugin.Client, error) {
			batchSizeObserved = s.WriteSpec.BatchSize
			batchSizeBytesObserved = s.WriteSpec.BatchSizeBytes
			return NewClient(ctx, logger, s)
		},
		plugin.WithDefaultBatchSize(batchSize),
		plugin.WithDefaultBatchSizeBytes(batchSizeBytes),
	)
	require.NoError(t, p.Init(context.TODO(), pbPlugin.Spec{}))

	require.Equal(t, batchSize, batchSizeObserved)
	require.Equal(t, batchSizeBytes, batchSizeBytesObserved)
}