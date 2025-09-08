package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
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
	var rowsPerRecord = int(10)
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table:        source,
		MigrateForce: writeOptionMigrateForce,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	records, err := s.plugin.readAll(ctx, source)
	if err != nil {
		return fmt.Errorf("failed to read initial records: %w", err)
	}
	initialItems := int(TotalRows(records))

	sourceName := target.Name
	syncTime := time.Now().UTC().Round(1 * time.Second)
	opts := schema.GenTestDataOptions{
		SourceName:         sourceName,
		SyncTime:           syncTime,
		MaxRows:            rowsPerRecord,
		TimePrecision:      s.genDatOptions.TimePrecision,
		UseHomogeneousType: s.useHomogeneousTypes,
	}
	// Test Generator should be initialized with the current number of items in the destination
	// this allows us to have multi-pass tests that ensure the migrations are stable
	// create--> write --> migrate --> write -->migrate -->write-->migrate -->write
	tg := schema.NewTestDataGenerator(uint64(initialItems))
	resource1 := tg.Generate(source, opts)
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: resource1,
	}); err != nil {
		return fmt.Errorf("failed to insert first record: %w", err)
	}
	resource1 = s.handleNulls(resource1) // we process nulls after writing

	records, err = s.plugin.readAll(ctx, source)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	sortRecords(source, records, "id")
	records = records[initialItems:]

	totalItems := TotalRows(records)
	if totalItems != int64(rowsPerRecord) {
		return fmt.Errorf("expected items: %d, got: %d", rowsPerRecord, totalItems)
	}

	if diff := RecordsDiff(source.ToArrowSchema(), records, []arrow.RecordBatch{resource1}); diff != "" {
		return fmt.Errorf("first record differs from expectation: %s", diff)
	}

	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table:        target,
		MigrateForce: writeOptionMigrateForce,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	resource2 := tg.Generate(target, opts)
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
	records = records[initialItems:]
	lastRow := resource2.NewSlice(resource2.NumRows()-1, resource2.NumRows())
	// if force migration is not required, we don't expect any items to be dropped (so there should be 2 items)
	if !writeOptionMigrateForce || supportsSafeMigrate {
		if err := expectRows(target.ToArrowSchema(), records, 2*int64(rowsPerRecord), lastRow); err != nil {
			if writeOptionMigrateForce && TotalRows(records) == int64(rowsPerRecord) {
				// if force migration is required, we can also expect 1 item to be dropped
				return expectRows(target.ToArrowSchema(), records, int64(rowsPerRecord), lastRow)
			}

			return err
		}

		return nil
	}

	return expectRows(target.ToArrowSchema(), records, int64(rowsPerRecord), lastRow)
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
		tableName := "cq_add_column" + suffix + "_" + tableUUIDSuffix()
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
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.AddColumn, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})

	t.Run("add_column_not_null"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.AddColumnNotNull {
			t.Skip("skipping test: add_column_not_null")
		}
		tableName := "cq_add_column_not_null" + suffix + "_" + tableUUIDSuffix()
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
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.AddColumnNotNull, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}

	})

	t.Run("remove_column"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.RemoveColumn {
			t.Skip("skipping test: remove_column")
		}
		tableName := "cq_remove_column" + suffix + "_" + tableUUIDSuffix()
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
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.RemoveColumn, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})

	t.Run("remove_column_not_null"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.RemoveColumnNotNull {
			t.Skip("skipping test: remove_column_not_null")
		}
		tableName := "cq_remove_column_not_null" + suffix + "_" + tableUUIDSuffix()
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
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.RemoveColumnNotNull, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})

	t.Run("change_column"+suffix, func(t *testing.T) {
		if !forceMigrate && !s.tests.SafeMigrations.ChangeColumn {
			t.Skip("skipping test: change_column")
		}
		tableName := "cq_change_column" + suffix + "_" + tableUUIDSuffix()
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
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.ChangeColumn, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})

	t.Run("remove_unique_constraint_only"+suffix, func(t *testing.T) {
		if s.tests.SkipSpecificMigrations.RemoveUniqueConstraint {
			t.Skip("skipping test completely: remove_unique_constraint_only")
		}
		if !forceMigrate && !s.tests.SafeMigrations.RemoveUniqueConstraint {
			t.Skip("skipping test: remove_unique_constraint_only")
		}
		tableName := "remove_unique_constraint_only" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "_cq_id", Type: types.ExtensionTypes.UUID, Unique: true},
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "_cq_id", Type: types.ExtensionTypes.UUID},
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.RemoveUniqueConstraint, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})

	t.Run("move_to_cq_id_only"+suffix, func(t *testing.T) {
		if s.tests.SkipSpecificMigrations.MovePKToCQOnly {
			t.Skip("skipping test completely: move_to_cq_id_only")
		}
		if !forceMigrate && !s.tests.SafeMigrations.MovePKToCQOnly {
			t.Skip("skipping test: move_to_cq_id_only")
		}
		tableName := "cq_move_to_cq_id_only" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "_cq_id", Type: types.ExtensionTypes.UUID, NotNull: true, Unique: true},
				{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "_cq_id", Type: types.ExtensionTypes.UUID, NotNull: true, Unique: true, PrimaryKey: true},
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.MovePKToCQOnly, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})
	t.Run("move_to_cq_id_only_adding_pkc"+suffix, func(t *testing.T) {
		if s.tests.SkipSpecificMigrations.MovePKToCQOnly {
			t.Skip("skipping test completely: move_to_cq_id_only_adding_pkc")
		}
		if !forceMigrate && !s.tests.SafeMigrations.MovePKToCQOnly {
			t.Skip("skipping test: move_to_cq_id_only_adding_pkc")
		}
		tableName := "cq_move_to_cq_id_only_adding_pkc" + suffix + "_" + tableUUIDSuffix()
		source := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "_cq_id", Type: types.ExtensionTypes.UUID, NotNull: true, Unique: true},
				{Name: "id", Type: arrow.PrimitiveTypes.Int64, PrimaryKey: true},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true, PrimaryKeyComponent: true},
			}}
		target := &schema.Table{
			Name: tableName,
			Columns: schema.ColumnList{
				{Name: "_cq_id", Type: types.ExtensionTypes.UUID, NotNull: true, Unique: true, PrimaryKey: true},
				{Name: "id", Type: arrow.PrimitiveTypes.Int64},
				{Name: "uuid", Type: types.ExtensionTypes.UUID},
				{Name: "bool", Type: arrow.FixedWidthTypes.Boolean, NotNull: true},
			}}
		require.NoError(t, s.migrate(ctx, target, source, s.tests.SafeMigrations.MovePKToCQOnly, forceMigrate))
		if !forceMigrate {
			require.NoError(t, s.migrate(ctx, target, target, true, false))
		}
	})

	t.Run("double_migration", func(t *testing.T) {
		if forceMigrate {
			t.Skip("double migration test has sense only for safe migrations")
		}
		tableName := "cq_double_migration_" + tableUUIDSuffix()
		table := schema.TestTable(tableName, s.genDatOptions)
		// s.migrate will perform create->write->migrate->write
		require.NoError(t, s.migrate(ctx, table, table, true, false))
	})
}

func expectRows(sc *arrow.Schema, records []arrow.RecordBatch, expectTotal int64, expectedLast arrow.RecordBatch) error {
	totalItems := TotalRows(records)
	if totalItems != expectTotal {
		return fmt.Errorf("expected %d items, got %d", expectTotal, totalItems)
	}
	lastRecord := records[len(records)-1]
	lastRow := lastRecord.NewSlice(lastRecord.NumRows()-1, lastRecord.NumRows())
	if diff := RecordsDiff(sc, []arrow.RecordBatch{lastRow}, []arrow.RecordBatch{expectedLast}); diff != "" {
		return fmt.Errorf("record #%d differs from expectation: %s", totalItems, diff)
	}
	return nil
}
