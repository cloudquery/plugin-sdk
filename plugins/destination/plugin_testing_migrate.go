package destination

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudquery/plugin-sdk/caser"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/testdata"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestMigrate(
	ctx context.Context,
	p *Plugin,
	logger zerolog.Logger,
	spec specs.Destination,
) error {
	spec.BatchSize = 1
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	suffix := strings.ToLower(strings.ReplaceAll(spec.WriteMode.String(), "-", "_"))
	tableName := fmt.Sprintf("cq_test_migrate_%s_%d", suffix, time.Now().Unix())
	table := testdata.TestTable(tableName)
	if err := p.Migrate(ctx, []*schema.Table{table}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testMigrate" + caser.New().ToPascal(suffix) + "Source" + uuid.NewString()
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	syncTime := time.Now().UTC().Round(1 * time.Second)
	resource1 := createTestResources(table, sourceName, syncTime, 1)[0]
	if err := p.writeOne(ctx, sourceSpec, []*schema.Table{table}, syncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one: %w", err)
	}

	// check that migrations and writes still succeed when column ordering is changed
	a := table.Columns.Index("uuid")
	b := table.Columns.Index("float")
	table.Columns[a], table.Columns[b] = table.Columns[b], table.Columns[a]
	if err := p.Migrate(ctx, []*schema.Table{table}); err != nil {
		return fmt.Errorf("failed to migrate table with changed column ordering: %w", err)
	}
	resource2 := createTestResources(table, sourceName, syncTime, 1)[0]
	if err := p.writeOne(ctx, sourceSpec, []*schema.Table{table}, syncTime, resource2); err != nil {
		return fmt.Errorf("failed to write one after column order change: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources after second write, got %d", len(resourcesRead))
	}

	// check that migrations succeed when a new column is added
	table.Columns = append(table.Columns, schema.Column{
		Name: "new_column",
		Type: schema.TypeInt,
	})
	if err := p.Migrate(ctx, []*schema.Table{table}); err != nil {
		return fmt.Errorf("failed to migrate table with new column: %w", err)
	}
	resource3 := createTestResources(table, sourceName, syncTime, 1)[0]
	if err := p.writeOne(ctx, sourceSpec, []*schema.Table{table}, syncTime, resource3); err != nil {
		return fmt.Errorf("failed to write one after column order change: %w", err)
	}
	resourcesRead, err = p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	if len(resourcesRead) != 3 {
		return fmt.Errorf("expected 3 resources after third write, got %d", len(resourcesRead))
	}

	// check that migration still succeeds when there is an extra column in the destination table,
	// which should be ignored
	oldTable := testdata.TestTable(tableName)
	if err := p.Migrate(ctx, []*schema.Table{oldTable}); err != nil {
		return fmt.Errorf("failed to migrate table with extra column in destination: %w", err)
	}
	resource4 := createTestResources(oldTable, sourceName, syncTime, 1)[0]
	if err := p.writeOne(ctx, sourceSpec, []*schema.Table{oldTable}, syncTime, resource4); err != nil {
		return fmt.Errorf("failed to write one after column order change: %w", err)
	}
	totalExpectedResources := 4
	if spec.MigrateMode == specs.MigrateModeForced {
		table.Columns[len(table.Columns)-1].Type = schema.TypeString
		if err := p.Migrate(ctx, []*schema.Table{table}); err != nil {
			return fmt.Errorf("failed to migrate table with changed column type: %w", err)
		}
		resource5 := createTestResources(table, sourceName, syncTime, 1)[0]
		if err := p.writeOne(ctx, sourceSpec, []*schema.Table{table}, syncTime, resource5); err != nil {
			return fmt.Errorf("failed to write one after column type change: %w", err)
		}
		totalExpectedResources++
	}

	resourcesRead, err = p.readAll(ctx, oldTable, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	if len(resourcesRead) != totalExpectedResources {
		return fmt.Errorf("expected %d resources after fourth write, got %d", totalExpectedResources, len(resourcesRead))
	}
	cqIDIndex := table.Columns.Index(schema.CqIDColumn.Name)
	found := false
	for _, r := range resourcesRead {
		if !r[cqIDIndex].Equal(resource4.Data[cqIDIndex]) {
			continue
		}
		found = true
		if !r.Equal(resource4.Data) {
			return fmt.Errorf("expected resource to be equal to original resource, but got diff: %s", r.Diff(resource4.Data))
		}
	}
	if !found {
		return fmt.Errorf("expected to find resource with cq_id %s, but none matched", resource4.Data[cqIDIndex])
	}

	return nil
}
