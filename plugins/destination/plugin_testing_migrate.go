package destination

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func tableUUIDSuffix() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "_")
}

func testMigration(ctx context.Context, t *testing.T, p *Plugin, logger zerolog.Logger, spec specs.Destination, target *schema.Table, source *schema.Table, mode specs.MigrateMode) error {
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}

	source.Columns = append(schema.ColumnList{
		schema.CqSourceNameColumn,
		schema.CqSyncTimeColumn,
		schema.CqIDColumn,
	}, source.Columns...)
	target.Columns = append(schema.ColumnList{
		schema.CqSourceNameColumn,
		schema.CqSyncTimeColumn,
		schema.CqIDColumn,
	}, target.Columns...)

	if err := p.Migrate(ctx, []*schema.Table{source}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := target.Name
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	syncTime := time.Now().UTC().Round(1 * time.Second)
	resource1 := createTestResources(source, sourceName, syncTime, 1)[0]
	if err := p.writeOne(ctx, sourceSpec, []*schema.Table{source}, syncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one: %w", err)
	}

	if err := p.Migrate(ctx, []*schema.Table{target}); err != nil {
		return fmt.Errorf("failed to migrate existing table: %w", err)
	}
	resource2 := createTestResources(target, sourceName, syncTime, 1)[0]
	if err := p.writeOne(ctx, sourceSpec, []*schema.Table{target}, syncTime, resource2); err != nil {
		return fmt.Errorf("failed to write one after migration: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, target, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	if mode == specs.MigrateModeSafe {
		if len(resourcesRead) != 2 {
			return fmt.Errorf("expected 2 resources after write, got %d", len(resourcesRead))
		}
		require.Contains(t, resourcesRead, resource2.Data)
	} else {
		if len(resourcesRead) != 1 {
			return fmt.Errorf("expected 1 resource after write, got %d", len(resourcesRead))
		}
		if diff := resourcesRead[0].Diff(resource2.Data); diff != "" {
			return fmt.Errorf("resource1 diff: %s", diff)
		}
	}

	if p.spec.PKMode == specs.PKModeCQID {
		for _, tColumn := range target.Columns {
			if tColumn.Name != schema.CqIDColumn.Name && tColumn.CreationOptions.PrimaryKey {
				return fmt.Errorf("unexpected primary key on %s", tColumn.Name)
			}
			if tColumn.Name == schema.CqIDColumn.Name && !tColumn.CreationOptions.PrimaryKey {
				return fmt.Errorf("expected primary key on %s", tColumn.Name)
			}
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
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
				{
					Name: "bool",
					Type: schema.TypeBool,
				},
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
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
				{
					Name: "bool",
					Type: schema.TypeBool,
					CreationOptions: schema.ColumnCreationOptions{
						NotNull: true,
					},
				},
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
		source := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
				{
					Name: "bool",
					Type: schema.TypeBool,
				},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
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
		source := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
				{
					Name: "bool",
					Type: schema.TypeBool,
					CreationOptions: schema.ColumnCreationOptions{
						NotNull: true,
					},
				},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
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
		source := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
				{
					Name: "bool",
					Type: schema.TypeBool,
				},
			},
		}
		target := &schema.Table{
			Name: tableName,
			Columns: []schema.Column{
				{
					Name: "id",
					Type: schema.TypeUUID,
				},
				{
					Name: "bool",
					Type: schema.TypeString,
				},
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
}
