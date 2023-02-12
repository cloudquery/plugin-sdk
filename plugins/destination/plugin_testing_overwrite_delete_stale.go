package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/testdata"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwriteDeleteStale(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeOverwriteDeleteStale
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_test_write_overwrite_delete_stale_%d", time.Now().Unix())
	table := testdata.TestTable(tableName)
	incTable := testdata.TestTable(tableName + "_incremental")
	incTable.IsIncremental = true
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*schema.Table{
		table,
		incTable,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testOverwriteSource" + uuid.NewString()
	sourceSpec := specs.Source{
		Name:    sourceName,
		Backend: specs.BackendLocal,
	}

	resources := createTestResources(table, sourceName, syncTime, 2)
	incResources := createTestResources(incTable, sourceName, syncTime, 2)
	if err := p.writeAll(ctx, sourceSpec, tables, syncTime, append(resources, incResources...)); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	sortResources(table, resources)

	resourcesRead, err := p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortCQTypes(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}

	if diff := resources[0].Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("expected first resource diff: %s", diff)
	}

	if diff := resources[1].Data.Diff(resourcesRead[1]); diff != "" {
		return fmt.Errorf("expected second resource diff: %s", diff)
	}

	// read from incremental table
	resourcesRead, err = p.readAll(ctx, incTable, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources in incremental table, got %d", len(resourcesRead))
	}

	secondSyncTime := syncTime.Add(time.Second).UTC()

	// copy first resource but update the sync time
	updatedResource := schema.DestinationResource{
		TableName: table.Name,
		Data:      make(schema.CQTypes, len(resources[0].Data)),
	}
	copy(updatedResource.Data, resources[0].Data)
	_ = updatedResource.Data[1].Set(secondSyncTime)

	// write second time
	if err := p.writeOne(ctx, sourceSpec, tables, secondSyncTime, updatedResource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortCQTypes(table, resourcesRead)
	if len(resourcesRead) != 1 {
		return fmt.Errorf("after overwrite expected 1 resource, got %d", len(resourcesRead))
	}

	if diff := resources[0].Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("after overwrite expected first resource diff: %s", diff)
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	if len(resourcesRead) != 1 {
		return fmt.Errorf("expected 1 resource after delete stale, got %d", len(resourcesRead))
	}

	// we expect the only resource returned to match the updated resource we wrote
	if diff := updatedResource.Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("after delete stale expected resource diff: %s", diff)
	}

	// we expect the incremental table to still have 2 resources, because delete-stale should
	// not apply there
	resourcesRead, err = p.readAll(ctx, tables[1], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all from incremental table: %w", err)
	}
	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources in incremental table after delete-stale, got %d", len(resourcesRead))
	}

	return nil
}
