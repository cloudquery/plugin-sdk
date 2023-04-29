package destination

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/cloudquery/plugin-sdk/v2/schemav2"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/cloudquery/plugin-sdk/v2/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func tableUUIDSuffix() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "_")
}

func testMigration(ctx context.Context, _ *testing.T, p *Plugin, logger zerolog.Logger, spec specs.Destination, target *schemav2.Table, source *schemav2.Table, mode specs.MigrateMode) error {
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}

	if err := p.Migrate(ctx, []*schemav2.Table{source}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := target.Name
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	syncTime := time.Now().UTC().Round(1 * time.Second)
	opts := schemav2.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   syncTime,
		MaxRows:    1,
	}
	resource1 := schemav2.GenTestData(source, opts)[0]
	if err := p.writeOne(ctx, sourceSpec, syncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one: %w", err)
	}

	if err := p.Migrate(ctx, schemav2.Tables{target}); err != nil {
		return fmt.Errorf("failed to migrate existing table: %w", err)
	}
	opts.SyncTime = syncTime.Add(time.Second).UTC()
	resource2 := schemav2.GenTestData(target, opts)[0]
	if err := p.writeOne(ctx, sourceSpec, syncTime, resource2); err != nil {
		return fmt.Errorf("failed to write one after migration: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, target, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortRecordsBySyncTime(target, resourcesRead)
	if mode == specs.MigrateModeSafe {
		if len(resourcesRead) != 2 {
			return fmt.Errorf("expected 2 resources after write, got %d", len(resourcesRead))
		}
		if !array.RecordApproxEqual(resourcesRead[1], resource2) {
			diff := RecordDiff(resourcesRead[1], resource2)
			return fmt.Errorf("resource1 and resource2 are not equal. diff: %s", diff)
		}
	} else {
		if len(resourcesRead) != 1 {
			return fmt.Errorf("expected 1 resource after write, got %d", len(resourcesRead))
		}
		if !array.RecordApproxEqual(resourcesRead[0], resource2) {
			diff := RecordDiff(resourcesRead[0], resource2)
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
	spec specs.Destination,
	strategy MigrateStrategy,
) {
	spec.BatchSize = 1

	t.Run("add_column", func(t *testing.T) {
		if strategy.AddColumn == specs.MigrateModeForced && spec.MigrateMode == specs.MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "add_column_" + tableUUIDSuffix()
		source := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID, CreationOptions: schemav2.ColumnCreationOptions{NotNull: false}},
			},
		}
		target := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID, CreationOptions: schemav2.ColumnCreationOptions{NotNull: false}},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, CreationOptions: schemav2.ColumnCreationOptions{NotNull: false}},
			},
		}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.AddColumn); err != nil {
			t.Fatalf("failed to migrate %s: %v", tableName, err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("add_column_not_null", func(t *testing.T) {
		if strategy.AddColumnNotNull == specs.MigrateModeForced && spec.MigrateMode == specs.MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "add_column_not_null_" + tableUUIDSuffix()
		source := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID, CreationOptions: schemav2.ColumnCreationOptions{NotNull: false}},
			},
		}
		target := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID, CreationOptions: schemav2.ColumnCreationOptions{NotNull: false}},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, CreationOptions: schemav2.ColumnCreationOptions{NotNull: true}},
			},
		}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.AddColumnNotNull); err != nil {
			t.Fatalf("failed to migrate add_column_not_null: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove_column", func(t *testing.T) {
		if strategy.RemoveColumn == specs.MigrateModeForced && spec.MigrateMode == specs.MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "remove_column_" + tableUUIDSuffix()
		source := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		}
		target := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
			},
		}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.RemoveColumn); err != nil {
			t.Fatalf("failed to migrate remove_column: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove_column_not_null", func(t *testing.T) {
		if strategy.RemoveColumnNotNull == specs.MigrateModeForced && spec.MigrateMode == specs.MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "remove_column_not_null_" + tableUUIDSuffix()
		source := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, CreationOptions: schemav2.ColumnCreationOptions{NotNull: true}},
			},
		}
		target := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
			},
		}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.RemoveColumnNotNull); err != nil {
			t.Fatalf("failed to migrate remove_column_not_null: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("change_column", func(t *testing.T) {
		if strategy.ChangeColumn == specs.MigrateModeForced && spec.MigrateMode == specs.MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "change_column_" + tableUUIDSuffix()
		source := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, CreationOptions: schemav2.ColumnCreationOptions{NotNull: true}},
			},
		}
		target := &schemav2.Table{
			Name: tableName,
			Columns: []schemav2.Column{
				schemav2.CqSourceNameColumn,
				schemav2.CqSyncTimeColumn,
				schemav2.CqIDColumn,
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.BinaryTypes.String, CreationOptions: schemav2.ColumnCreationOptions{NotNull: true}},
			},
		}

		p := newPlugin()
		if err := testMigration(ctx, t, p, logger, spec, target, source, strategy.ChangeColumn); err != nil {
			t.Fatalf("failed to migrate change_column: %v", err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("double_migration", func(t *testing.T) {
		tableName := "double_migration_" + tableUUIDSuffix()
		table := schemav2.TestTable(tableName)

		p := newPlugin()
		require.NoError(t, p.Init(ctx, logger, spec))
		require.NoError(t, p.Migrate(ctx, schemav2.Tables{table}))

		nonForced := spec
		nonForced.MigrateMode = specs.MigrateModeSafe
		require.NoError(t, p.Init(ctx, logger, nonForced))
		require.NoError(t, p.Migrate(ctx, schemav2.Tables{table}))
	})
}
