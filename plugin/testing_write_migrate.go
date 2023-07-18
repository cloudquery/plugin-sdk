package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func tableUUIDSuffix() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "_")[:8] // use only first 8 chars
}

// nolint:revive
func (s *WriterTestSuite) migrate(ctx context.Context, target *schema.Table, source *schema.Table, supportsSafeMigrate bool, writeOptionMigrateForce bool) error {
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table:        source,
		MigrateForce: writeOptionMigrateForce,
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
	tg := schema.NewTestDataGenerator()
	resource1 := tg.Generate(source, opts)[0]
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: resource1,
	}); err != nil {
		return fmt.Errorf("failed to insert first record: %w", err)
	}
	resource1 = s.handleNulls(resource1) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, source)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems := TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}
	if diff := RecordDiff(records[0], resource1); diff != "" {
		return fmt.Errorf("first record differs from expectation: %s", diff)
	}

	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table:        target,
		MigrateForce: writeOptionMigrateForce,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	resource2 := tg.Generate(target, opts)[0]
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: resource2,
	}); err != nil {
		return fmt.Errorf("failed to insert second record: %w", err)
	}
	resource2 = s.handleNulls(resource2) // we process nulls after writing

	records, err = s.plugin.readAll(ctx, target)
	if err != nil {
		return fmt.Errorf("failed to readAll: %w", err)
	}
	sortRecords(target, records, "id")

	// if force migration is not required, we don't expect any items to be dropped (so there should be 2 items)
	if !writeOptionMigrateForce || supportsSafeMigrate {
		if err := expectRows(records, 2, resource2); err != nil {
			if writeOptionMigrateForce && TotalRows(records) == 1 {
				// if force migration is required, we can also expect 1 item to be dropped
				return expectRows(records, 1, resource2)
			}

			return err
		}

		return nil
	}

	return expectRows(records, 1, resource2)
}

// nolint:revive
func (s *WriterTestSuite) testMigrate(
	ctx context.Context,
	t *testing.T,
	forceMigrate bool,
) {
	suffix := "_safe"
	if forceMigrate {
		suffix = "_force"
	}
	t.Run("add_column"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.AddColumn {
			t.Skip("skipping test: add_column")
		}
		tableName := "add_column" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			},
		}
		if err := s.migrate(ctx, target, source, s.tests.SafeMigrations.AddColumn, forceMigrate); err != nil {
			t.Fatalf("failed to migrate %s: %v", tableName, err)
		}
	})

	t.Run("add_column_not_null"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.AddColumnNotNull {
			t.Skip("skipping test: add_column_not_null")
		}
		tableName := "add_column_not_null" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
			},
		}

		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		if err := s.migrate(ctx, target, source, s.tests.SafeMigrations.AddColumnNotNull, forceMigrate); err != nil {
			t.Fatalf("failed to migrate add_column_not_null: %v", err)
		}
	})

	t.Run("remove_column"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.RemoveColumn {
			t.Skip("skipping test: remove_column")
		}
		tableName := "remove_column" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
			}}
		if err := s.migrate(ctx, target, source, s.tests.SafeMigrations.RemoveColumn, forceMigrate); err != nil {
			t.Fatalf("failed to migrate remove_column: %v", err)
		}
	})

	t.Run("remove_column_not_null"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.RemoveColumnNotNull {
			t.Skip("skipping test: remove_column_not_null")
		}
		tableName := "remove_column_not_null" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
			}}
		if err := s.migrate(ctx, target, source, s.tests.SafeMigrations.RemoveColumnNotNull, forceMigrate); err != nil {
			t.Fatalf("failed to migrate remove_column_not_null: %v", err)
		}
	})

	t.Run("change_column"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.ChangeColumn {
			t.Skip("skipping test: change_column")
		}
		tableName := "change_column" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.BinaryTypes.String, NotNull: true},
			}}
		if err := s.migrate(ctx, target, source, s.tests.SafeMigrations.ChangeColumn, forceMigrate); err != nil {
			t.Fatalf("failed to migrate change_column: %v", err)
		}
	})

	t.Run("double_migration", func(t *testing.T) {
		tableName := "double_migration_" + tableUUIDSuffix()
		table := schema.TestTable(tableName, s.genDatOptions)
		// s.migrate will perform create->write->migrate->write
		require.NoError(t, s.migrate(ctx, table, table, true, false))
	})
}

func expectRows(records []arrow.Record, expectTotal int64, expectedLast arrow.Record) error {
	totalItems := TotalRows(records)
	if totalItems != expectTotal {
		return fmt.Errorf("expected %d items, got %d", expectTotal, totalItems)
	}
	if diff := RecordDiff(records[totalItems-1], expectedLast); diff != "" {
		return fmt.Errorf("record #%d differs from expectation: %s", totalItems-1, diff)
	}
	return nil
}
