package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/testdata"
	"github.com/cloudquery/plugin-sdk/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwriteDeleteStale(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeOverwriteDeleteStale
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_%s_%d", spec.Name, time.Now().Unix())
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

	resources := testdata.GenTestData(table.ToArrowSchema(), sourceName, syncTime, uuid.Nil, 2)
	incResources := testdata.GenTestData(incTable.ToArrowSchema(), sourceName, syncTime, uuid.Nil, 2)
	allResources := resources
	allResources = append(allResources, incResources...)
	if err := p.writeAll(ctx, sourceSpec, tables, syncTime, allResources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	sortRecordsBySyncTime(table, resources)

	resourcesRead, err := p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}
	if !array.RecordEqual(resources[0], resourcesRead[0]) {
		diff := RecordDiff(resources[0], resourcesRead[0])
		return fmt.Errorf("expected first resource to be equal. diff: %s", diff)
	}

	if !array.RecordEqual(resources[1], resourcesRead[1]) {
		diff := RecordDiff(resources[1], resourcesRead[1])
		return fmt.Errorf("expected second resource to be equal. diff: %s", diff)
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
	u := resources[0].Column(2).(*types.UUIDArray).Value(0)
	updatedResources := testdata.GenTestData(table.ToArrowSchema(), sourceName, secondSyncTime, *u, 1)[0]

	if err := p.writeOne(ctx, sourceSpec, tables, secondSyncTime, updatedResources); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)
	if len(resourcesRead) != 1 {
		return fmt.Errorf("after overwrite expected 1 resource, got %d", len(resourcesRead))
	}
	if array.RecordEqual(resources[0], resourcesRead[0]) {
		diff := RecordDiff(resources[0], resourcesRead[0])
		return fmt.Errorf("after overwrite expected first resource to be different. diff: %s", diff)
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	if len(resourcesRead) != 1 {
		return fmt.Errorf("expected 1 resource after delete stale, got %d", len(resourcesRead))
	}

	// we expect the only resource returned to match the updated resource we wrote
	if !array.RecordEqual(updatedResources, resourcesRead[0]) {
		diff := RecordDiff(updatedResources, resourcesRead[0])
		return fmt.Errorf("after delete stale expected resource to be equal. diff: %s", diff)
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
