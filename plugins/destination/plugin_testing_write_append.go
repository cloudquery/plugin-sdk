package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (s *PluginTestSuite) destinationPluginTestWriteAppend(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination, testOpts PluginTestSuiteRunnerOptions) error {
	spec.WriteMode = specs.WriteModeAppend
	if err := p.Init(ctx, logger, spec); err != nil {
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

	sourceName := "testAppendSource" + uuid.NewString()
	specSource := specs.Source{
		Name: sourceName,
	}

	opts := schema.GenTestDataOptions{
		SourceName:    sourceName,
		SyncTime:      syncTime,
		MaxRows:       2,
		TimePrecision: testOpts.TimePrecision,
	}
	record1 := schema.GenTestData(table, opts)
	if err := p.writeAll(ctx, specSource, syncTime, record1); err != nil {
		return fmt.Errorf("failed to write record first time: %w", err)
	}

	secondSyncTime := syncTime.Add(10 * time.Second).UTC()
	opts.SyncTime = secondSyncTime
	opts.MaxRows = 1
	record2 := schema.GenTestData(table, opts)

	if !s.tests.SkipSecondAppend {
		// write second time
		if err := p.writeAll(ctx, specSource, secondSyncTime, record2); err != nil {
			return fmt.Errorf("failed to write one second time: %w", err)
		}
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	sortRecordsBySyncTime(table, resourcesRead)

	expectedResource := 3
	if s.tests.SkipSecondAppend {
		expectedResource = 2
	}

	if len(resourcesRead) != expectedResource {
		return fmt.Errorf("expected %d resources, got %d", expectedResource, len(resourcesRead))
	}

	testOpts.AllowNull.replaceNullsByEmpty(record1)
	testOpts.AllowNull.replaceNullsByEmpty(record2)
	if testOpts.IgnoreNullsInLists {
		stripNullsFromLists(record1)
		stripNullsFromLists(record2)
	}
	if !array.RecordApproxEqual(record1[0], resourcesRead[0]) {
		diff := RecordDiff(record1[0], resourcesRead[0])
		return fmt.Errorf("first expected resource diff at row 0: %s", diff)
	}
	if !array.RecordApproxEqual(record1[1], resourcesRead[1]) {
		diff := RecordDiff(record1[1], resourcesRead[1])
		return fmt.Errorf("first expected resource diff at row 1: %s", diff)
	}

	if !s.tests.SkipSecondAppend {
		if !array.RecordApproxEqual(record2[0], resourcesRead[2]) {
			diff := RecordDiff(record2[0], resourcesRead[2])
			return fmt.Errorf("second expected resource diff: %s", diff)
		}
	}

	return nil
}
