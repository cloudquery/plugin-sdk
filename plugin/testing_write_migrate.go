package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func tableUUIDSuffix() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "_")
}

func testMigration(ctx context.Context, _ *testing.T, p *Plugin, logger zerolog.Logger, spec pbPlugin.Spec, target *schema.Table, source *schema.Table, mode pbPlugin.WriteSpec_MIGRATE_MODE, testOpts PluginTestSuiteRunnerOptions) error {
	if err := p.Init(ctx, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}

	if err := p.Migrate(ctx, schema.Tables{source}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := target.Name
	sourceSpec := pbPlugin.Spec{
		Name: sourceName,
	}
	syncTime := time.Now().UTC().Round(1 * time.Second)
	opts := schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      syncTime,
		MaxRows:       1,
		TimePrecision: testOpts.TimePrecision,
	}
	resource1 := schema.GenTestData(source, opts)[0]
	if err := p.writeOne(ctx, sourceSpec, syncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one: %w", err)
	}

	if err := p.Migrate(ctx, schema.Tables{target}); err != nil {
		return fmt.Errorf("failed to migrate existing table: %w", err)
	}
	opts.SyncTime = syncTime.Add(time.Second).UTC()
	resource2 := schema.GenTestData(target, opts)
	if err := p.writeAll(ctx, sourceSpec, syncTime, resource2); err != nil {
		return fmt.Errorf("failed to write one after migration: %w", err)
	}

	testOpts.AllowNull.replaceNullsByEmpty(resource2)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(resource2)
	}

	resourcesRead, err := p.readAll(ctx, target, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortRecordsBySyncTime(target, resourcesRead)
	if mode == pbPlugin.WriteSpec_SAFE {
		if len(resourcesRead) != 2 {
			return fmt.Errorf("expected 2 resources after write, got %d", len(resourcesRead))
		}
		if !array.RecordApproxEqual(resourcesRead[1], resource2[0]) {
			diff := RecordDiff(resourcesRead[1], resource2[0])
			return fmt.Errorf("resource1 and resource2 are not equal. diff: %s", diff)
		}
	} else {
		if len(resourcesRead) != 1 {
			return fmt.Errorf("expected 1 resource after write, got %d", len(resourcesRead))
		}
		if !array.RecordApproxEqual(resourcesRead[0], resource2[0]) {
			diff := RecordDiff(resourcesRead[0], resource2[0])
			return fmt.Errorf("resource1 and resource2 are not equal. diff: %s", diff)
		}
	}

	return nil
}

func (*PluginTestSuite) destinationPluginTestMigrate(
	ctx context.Context,
	t *testing.T,
	newPlugin NewPluginFunc,
	logger zerolog.Logger,
	spec pbPlugin.Spec,
	strategy MigrateStrategy,
	testOpts PluginTestSuiteRunnerOptions,
) {
	spec.WriteSpec.BatchSize = 1

	t.Run("add_column", func(t *testing.T) {
		if strategy.AddColumn == pbPlugin.WriteSpec_FORCE && spec.WriteSpec.MigrateMode == pbPlugin.WriteSpec_SAFE {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "add_column_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.AddColumn, testOpts); err != nil {
			t.Fatalf("failed to migrate %s: %v", tableName, err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("add_column_not_null", func(t *testing.T) {
		if strategy.AddColumnNotNull == pbPlugin.WriteSpec_FORCE && spec.WriteSpec.MigrateMode == pbPlugin.WriteSpec_SAFE {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "add_column_not_null_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.AddColumnNotNull, testOpts); err != nil {
			t.Fatalf("failed to migrate add_column_not_null: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove_column", func(t *testing.T) {
		if strategy.RemoveColumn == pbPlugin.WriteSpec_FORCE && spec.WriteSpec.MigrateMode == pbPlugin.WriteSpec_SAFE {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "remove_column_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
			}}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.RemoveColumn, testOpts); err != nil {
			t.Fatalf("failed to migrate remove_column: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove_column_not_null", func(t *testing.T) {
		if strategy.RemoveColumnNotNull == pbPlugin.WriteSpec_FORCE && spec.WriteSpec.MigrateMode == pbPlugin.WriteSpec_SAFE {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "remove_column_not_null_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
			}}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.RemoveColumnNotNull, testOpts); err != nil {
			t.Fatalf("failed to migrate remove_column_not_null: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("change_column", func(t *testing.T) {
		if strategy.ChangeColumn == pbPlugin.WriteSpec_FORCE && spec.WriteSpec.MigrateMode == pbPlugin.WriteSpec_SAFE {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "change_column_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.BinaryTypes.String, NotNull: true},
			}}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.ChangeColumn, testOpts); err != nil {
			t.Fatalf("failed to migrate change_column: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("double_migration", func(t *testing.T) {
		tableName := "double_migration_" + tableUUIDSuffix()
		table := schema.TestTable(tableName, testOpts.TestSourceOptions)

		p := newPlugin()
		require.NoError(t, p.Init(ctx, spec))
		require.NoError(t, p.Migrate(ctx, schema.Tables{table}))

		nonForced := spec
		nonForced.WriteSpec.MigrateMode = pbPlugin.WriteSpec_SAFE
		require.NoError(t, p.Init(ctx, nonForced))
		require.NoError(t, p.Migrate(ctx, schema.Tables{table}))
	})
}
