package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *Plugin, logger zerolog.Logger, spec pbPlugin.Spec, testOpts PluginTestSuiteRunnerOptions) error {
	spec.WriteSpec.WriteMode = pbPlugin.WRITE_MODE_WRITE_MODE_OVERWRITE
	if err := p.Init(ctx, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_%s_%d", spec.Name, time.Now().Unix())
	table := schema.TestTable(tableName, testOpts.TestSourceOptions)
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := schema.Tables{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testOverwriteSource" + uuid.NewString()
	sourceSpec := pbPlugin.Spec{
		Name: sourceName,
	}

	opts := schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      syncTime,
		MaxRows:       2,
		TimePrecision: testOpts.TimePrecision,
	}
	resources := schema.GenTestData(table, opts)
	if err := p.writeAll(ctx, sourceSpec, syncTime, resources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	sortRecordsBySyncTime(table, resources)
	testOpts.AllowNull.replaceNullsByEmpty(resources)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(resources)
	}
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
	cqIDInds := resources[0].Schema().FieldIndices(schema.CqIDColumn.Name)
	u := resources[0].Column(cqIDInds[0]).(*types.UUIDArray).Value(0)
	opts = schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      secondSyncTime,
		MaxRows:       1,
		StableUUID:    u,
		TimePrecision: testOpts.TimePrecision,
	}
	updatedResource := schema.GenTestData(table, opts)
	// write second time
	if err := p.writeAll(ctx, sourceSpec, secondSyncTime, updatedResource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	testOpts.AllowNull.replaceNullsByEmpty(updatedResource)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(updatedResource)
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
	if !array.RecordApproxEqual(updatedResource[0], resourcesRead[1]) {
		diff := RecordDiff(updatedResource[0], resourcesRead[1])
		return fmt.Errorf("after overwrite expected second resource to be equal. diff=%s", diff)
	}

	return nil
}