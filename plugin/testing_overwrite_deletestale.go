package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (*PluginTestSuite) destinationPluginTestWriteOverwriteDeleteStale(ctx context.Context, p *Plugin, logger zerolog.Logger, spec any, testOpts PluginTestSuiteRunnerOptions) error {
	writeMode := WriteModeOverwriteDeleteStale
	if err := p.Init(ctx, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := fmt.Sprintf("cq_overwrite_delete_stale_%d", time.Now().Unix())
	table := schema.TestTable(tableName, testOpts.TestSourceOptions)
	incTable := schema.TestTable(tableName+"_incremental", testOpts.TestSourceOptions)
	incTable.IsIncremental = true
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := schema.Tables{
		table,
		incTable,
	}
	if err := p.Migrate(ctx, tables, MigrateModeSafe); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testOverwriteSource" + uuid.NewString()

	opts := schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      syncTime,
		MaxRows:       2,
		TimePrecision: testOpts.TimePrecision,
	}
	resources := schema.GenTestData(table, opts)
	incResources := schema.GenTestData(incTable, opts)
	allResources := resources
	allResources = append(allResources, incResources...)
	if err := p.writeAll(ctx, sourceName, syncTime, writeMode, allResources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}
	sortRecordsBySyncTime(table, resources)

	resourcesRead, err := p.syncAll(ctx, SyncOptions{
		Tables:     []string{tableName},
		SyncTime:   syncTime,
		SourceName: sourceName,
	})
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}
	testOpts.AllowNull.replaceNullsByEmpty(resources)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(resources)
	}
	if !recordApproxEqual(resources[0], resourcesRead[0]) {
		diff := RecordDiff(resources[0], resourcesRead[0])
		return fmt.Errorf("expected first resource to be equal. diff: %s", diff)
	}

	if !recordApproxEqual(resources[1], resourcesRead[1]) {
		diff := RecordDiff(resources[1], resourcesRead[1])
		return fmt.Errorf("expected second resource to be equal. diff: %s", diff)
	}

	// read from incremental table
	resourcesRead, err = p.syncAll(ctx, SyncOptions{
		Tables:     []string{tableName},
		SyncTime:   syncTime,
		SourceName: sourceName,
	})
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}
	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources in incremental table, got %d", len(resourcesRead))
	}

	secondSyncTime := syncTime.Add(time.Second).UTC()
	// copy first resource but update the sync time
	cqIDInds := resources[0].Schema().FieldIndices(schema.CqIDColumn.Name)
	u := resources[0].Column(cqIDInds[0]).(*types.UUIDArray).Value(0)
	opts = schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      secondSyncTime,
		StableUUID:    u,
		MaxRows:       1,
		TimePrecision: testOpts.TimePrecision,
	}
	updatedResources := schema.GenTestData(table, opts)
	updatedIncResources := schema.GenTestData(incTable, opts)
	allUpdatedResources := updatedResources
	allUpdatedResources = append(allUpdatedResources, updatedIncResources...)

	if err := p.writeAll(ctx, sourceName, secondSyncTime, writeMode, allUpdatedResources); err != nil {
		return fmt.Errorf("failed to write all second time: %w", err)
	}

	resourcesRead, err = p.syncAll(ctx, SyncOptions{
		Tables:     []string{tableName},
		SyncTime:   secondSyncTime,
		SourceName: sourceName,
	})
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)
	if len(resourcesRead) != 1 {
		return fmt.Errorf("after overwrite expected 1 resource, got %d", len(resourcesRead))
	}
	testOpts.AllowNull.replaceNullsByEmpty(resources)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(resources)
	}
	if recordApproxEqual(resources[0], resourcesRead[0]) {
		diff := RecordDiff(resources[0], resourcesRead[0])
		return fmt.Errorf("after overwrite expected first resource to be different. diff: %s", diff)
	}

	resourcesRead, err = p.syncAll(ctx, SyncOptions{
		Tables:     []string{tableName},
		SyncTime:   syncTime,
		SourceName: sourceName,
	})
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	if len(resourcesRead) != 1 {
		return fmt.Errorf("expected 1 resource after delete stale, got %d", len(resourcesRead))
	}

	// we expect the only resource returned to match the updated resource we wrote
	testOpts.AllowNull.replaceNullsByEmpty(updatedResources)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(updatedResources)
	}
	if !recordApproxEqual(updatedResources[0], resourcesRead[0]) {
		diff := RecordDiff(updatedResources[0], resourcesRead[0])
		return fmt.Errorf("after delete stale expected resource to be equal. diff: %s", diff)
	}

	// we expect the incremental table to still have 3 resources, because delete-stale should
	// not apply there
	resourcesRead, err = p.syncAll(ctx, SyncOptions{
		Tables:     []string{incTable.Name},
		SyncTime:   secondSyncTime,
		SourceName: sourceName,
	})
	if err != nil {
		return fmt.Errorf("failed to read all from incremental table: %w", err)
	}
	if len(resourcesRead) != 3 {
		return fmt.Errorf("expected 3 resources in incremental table after delete-stale, got %d", len(resourcesRead))
	}

	return nil
}
