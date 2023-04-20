package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/cloudquery/plugin-sdk/v2/testdata"
	"github.com/cloudquery/plugin-sdk/v2/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeOverwrite
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_%s_%d", spec.Name, time.Now().Unix())
	table := testdata.TestTable(tableName).ToArrowSchema()
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*arrow.Schema{
		table,
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
	resources := testdata.GenTestData(table, opts)
	if err := p.writeAll(ctx, sourceSpec, syncTime, resources); err != nil {
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

	if !array.RecordApproxEqual(resources[0], resourcesRead[0]) {
		diff := RecordDiff(resources[0], resourcesRead[0])
		return fmt.Errorf("expected first resource to be equal. diff=%s", diff)
	}

	if !array.RecordApproxEqual(resources[1], resourcesRead[1]) {
		diff := RecordDiff(resources[1], resourcesRead[1])
		return fmt.Errorf("expected second resource to be equal. diff=%s", diff)
	}

	secondSyncTime := syncTime.Add(time.Second).UTC()

	// copy first resource but update the sync time
	u := resources[0].Column(2).(*types.UUIDArray).Value(0)
	opts = testdata.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   secondSyncTime,
		MaxRows:    1,
		StableUUID: *u,
	}
	updatedResource := testdata.GenTestData(table, opts)[0]
	// write second time
	if err := p.writeOne(ctx, sourceSpec, secondSyncTime, updatedResource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, table, sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)
	if len(resourcesRead) != 2 {
		return fmt.Errorf("after overwrite expected 2 resources, got %d", len(resourcesRead))
	}

	if !array.RecordApproxEqual(resources[1], resourcesRead[0]) {
		diff := RecordDiff(resources[1], resourcesRead[0])
		return fmt.Errorf("after overwrite expected first resource to be equal. diff=%s", diff)
	}
	if !array.RecordApproxEqual(updatedResource, resourcesRead[1]) {
		diff := RecordDiff(updatedResource, resourcesRead[1])
		return fmt.Errorf("after overwrite expected second resource to be equal. diff=%s", diff)
	}

	return nil
}
