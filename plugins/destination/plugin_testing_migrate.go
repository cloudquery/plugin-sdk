package destination

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func tableUUIDSuffix() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "_")
}

func testMigration(ctx context.Context, _ *testing.T, p *Plugin, logger zerolog.Logger, spec specs.Destination, target *schema.Table, source *schema.Table, mode specs.MigrateMode) error {
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}

	if err := p.Migrate(ctx, schema.Tables{source}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := target.Name
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	syncTime := time.Now().UTC().Round(1 * time.Second)
	opts := schema.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   syncTime,
		MaxRows:    1,
	}
	resource1 := schema.GenTestData(source, opts)[0]
	if err := p.writeOne(ctx, sourceSpec, syncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one: %w", err)
	}

	if err := p.Migrate(ctx, schema.Tables{target}); err != nil {
		return fmt.Errorf("failed to migrate existing table: %w", err)
	}
	opts.SyncTime = syncTime.Add(time.Second).UTC()
	resource2 := schema.GenTestData(target, opts)[0]
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
		source := &schema.Table{
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
				{Field: arrow.Field{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, Nullable: true}},
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
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
			},
		}

		target := &schema.Table{
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
				{Field: arrow.Field{Name: "bool", Type: arrow.FixedWidthTypes.Boolean}},
			}}
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
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID}},
				{Field: arrow.Field{Name: "bool", Type: arrow.FixedWidthTypes.Boolean}},
			}}
		target := &schema.Table{
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
			}}

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
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
				{Field: arrow.Field{Name: "bool", Type: arrow.FixedWidthTypes.Boolean}},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID}},
			}}

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
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
				{Field: arrow.Field{Name: "bool", Type: arrow.FixedWidthTypes.Boolean}},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				schema.CqSourceNameColumn,
				schema.CqSyncTimeColumn,
				schema.CqIDColumn,
				{Field: arrow.Field{Name: "id", Type: types.ExtensionTypes.UUID, Nullable: true}},
				{Field: arrow.Field{Name: "bool", Type: arrow.BinaryTypes.String}},
			}}

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
		table := schema.TestTable(tableName)

		p := newPlugin()
		require.NoError(t, p.Init(ctx, logger, spec))
		require.NoError(t, p.Migrate(ctx, schema.Tables{table}))

		nonForced := spec
		nonForced.MigrateMode = specs.MigrateModeSafe
		require.NoError(t, p.Init(ctx, logger, nonForced))
		require.NoError(t, p.Migrate(ctx, schema.Tables{table}))
	})
}
