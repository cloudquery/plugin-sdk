package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
)

func tableUUIDSuffix() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "_")
}

func (s *PluginTestSuite) migrate(ctx context.Context, target *schema.Table, source *schema.Table, strategy MigrateMode, mode MigrateMode) error {
	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageCreateTable{
		Table: source,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	sourceName := target.Name
	syncTime := time.Now().UTC().Round(1 * time.Second)
	opts := schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      syncTime,
		MaxRows:       1,
		TimePrecision: s.genDatOptions.TimePrecision,
	}

	resource1 := schema.GenTestData(source, opts)[0]

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageInsert{
		Record: resource1,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	messages, err := s.plugin.syncAll(ctx, SyncOptions{
		Tables: []string{source.Name},
	})
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems := messages.InsertItems()
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageCreateTable{
		Table: target,
		Force: strategy == MigrateModeForce,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	if err := s.plugin.writeOne(ctx, WriteOptions{}, &MessageInsert{
		Record: resource1,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	messages, err = s.plugin.syncAll(ctx, SyncOptions{
		Tables: []string{source.Name},
	})
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	if strategy == MigrateModeSafe || mode == MigrateModeSafe {
		totalItems = messages.InsertItems()
		if totalItems != 2 {
			return fmt.Errorf("expected 2 item, got %d", totalItems)
		}
	} else {
		totalItems = messages.InsertItems()
		if totalItems != 1 {
			return fmt.Errorf("expected 1 item, got %d", totalItems)
		}
	}

	return nil
}

func (s *PluginTestSuite) testMigrate(
	ctx context.Context,
	t *testing.T,
	mode MigrateMode,
) {
	t.Run("add_column", func(t *testing.T) {
		if s.tests.MigrateStrategy.AddColumn == MigrateModeForce && mode == MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "add_column_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		}
		if err := s.migrate(ctx, target, source, s.tests.MigrateStrategy.AddColumn, mode); err != nil {
			t.Fatalf("failed to migrate %s: %v", tableName, err)
		}
	})

	t.Run("add_column_not_null", func(t *testing.T) {
		if s.tests.MigrateStrategy.AddColumnNotNull == MigrateModeForce && mode == MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "add_column_not_null_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		if err := s.migrate(ctx, target, source, s.tests.MigrateStrategy.AddColumnNotNull, mode); err != nil {
			t.Fatalf("failed to migrate add_column_not_null: %v", err)
		}
	})

	t.Run("remove_column", func(t *testing.T) {
		if s.tests.MigrateStrategy.RemoveColumn == MigrateModeForce && mode == MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "remove_column_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
			}}
		if err := s.migrate(ctx, target, source, s.tests.MigrateStrategy.RemoveColumn, mode); err != nil {
			t.Fatalf("failed to migrate remove_column: %v", err)
		}
	})

	t.Run("remove_column_not_null", func(t *testing.T) {
		if s.tests.MigrateStrategy.RemoveColumnNotNull == MigrateModeForce && mode == MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "remove_column_not_null_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
			}}
		if err := s.migrate(ctx, target, source, s.tests.MigrateStrategy.RemoveColumnNotNull, mode); err != nil {
			t.Fatalf("failed to migrate remove_column_not_null: %v", err)
		}
	})

	t.Run("change_column", func(t *testing.T) {
		if s.tests.MigrateStrategy.ChangeColumn == MigrateModeForce && mode == MigrateModeSafe {
			t.Skip("skipping as migrate mode is safe")
			return
		}
		tableName := "change_column_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.BinaryTypes.String, NotNull: true},
			}}
		if err := s.migrate(ctx, target, source, s.tests.MigrateStrategy.ChangeColumn, mode); err != nil {
			t.Fatalf("failed to migrate change_column: %v", err)
		}
	})

	t.Run("double_migration", func(t *testing.T) {
		// tableName := "double_migration_" + tableUUIDSuffix()
		// table := schema.TestTable(tableName, testOpts.TestSourceOptions)
		// require.NoError(t, p.Migrate(ctx, schema.Tables{table}, MigrateOptions{MigrateMode: MigrateModeForce}))
		// require.NoError(t, p.Migrate(ctx, schema.Tables{table}, MigrateOptions{MigrateMode: MigrateModeForce}))
	})
}
