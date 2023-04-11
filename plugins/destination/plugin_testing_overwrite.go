package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/cloudquery/plugin-sdk/v2/testdata"
	"github.com/cloudquery/plugin-sdk/v2/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, mem memory.Allocator, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeOverwrite
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_%s_%d", spec.Name, time.Now().Unix())
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*arrow.Schema{
		table.ToArrowSchema(),
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testOverwriteSource" + uuid.NewString()
	sourceSpec := specs.Source{
		Name: sourceName,
	}

	opts := testdata.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   syncTime,
		MaxRows:    2,
	}
	resources := testdata.GenTestData(mem, schema.CQSchemaToArrow(table), opts)
	defer func() {
		for _, r := range resources {
			r.Release()
		}
	}()
	if err := p.writeAll(ctx, sourceSpec, syncTime, resources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	sortRecordsBySyncTime(table, resources)

	resourcesRead, err := p.readAll(ctx, table.ToArrowSchema(), sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}

	if !array.RecordEqual(resources[0], resourcesRead[0]) {
		diff := RecordDiff(resources[0], resourcesRead[0])
		return fmt.Errorf("expected first resource to be equal. diff=%s", diff)
	}

	if !array.RecordEqual(resources[1], resourcesRead[1]) {
		diff := RecordDiff(resources[1], resourcesRead[1])
		return fmt.Errorf("expected second resource to be equal. diff=%s", diff)
	}

	secondSyncTime := time.Now().UTC().Round(1 * time.Second).Add(time.Hour).UTC()

	// copy first resource but update the sync time
	u := resources[0].Column(2).(*types.UUIDArray).Value(0)
	opts = testdata.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   secondSyncTime,
		MaxRows:    1,
		StableUUID: *u,
	}
	updatedResource := testdata.GenTestData(mem, schema.CQSchemaToArrow(table), opts)[0]
	defer updatedResource.Release()
	// write second time
	if err := p.writeOne(ctx, sourceSpec, secondSyncTime, updatedResource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, table.ToArrowSchema(), sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)
	if len(resourcesRead) != 2 {
		return fmt.Errorf("after overwrite expected 2 resources, got %d", len(resourcesRead))
	}

	if !array.RecordEqual(resources[1], resourcesRead[0]) {
		diff := RecordDiff(resources[1], resourcesRead[0])
		return fmt.Errorf("after overwrite expected first resource to be equal. diff=%s", diff)
	}
	if !array.RecordEqual(updatedResource, resourcesRead[1]) {
		diff := RecordDiff(updatedResource, resourcesRead[1])
		return fmt.Errorf("after overwrite expected second resource to be equal. diff=%s", diff)
	}

	return nil
}
