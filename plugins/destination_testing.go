package plugins

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/testdata"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type destinationTestSuite struct {
	tests DestinationTestSuiteTests
}

type DestinationTestSuiteTests struct {
	// SkipOverwrite skips testing for "overwrite" mode. Use if the destination
	//	// plugin doesn't support this feature.
	SkipOverwrite bool

	// SkipDeleteStale skips testing "delete-stale" mode. Use if the destination
	// plugin doesn't support this feature.
	SkipDeleteStale bool

	// SkipAppend skips testing for "append" mode. Use if the destination
	// plugin doesn't support this feature.
	SkipAppend bool
}

func getTestLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func (s *destinationTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *DestinationPlugin, logger zerolog.Logger, spec specs.Destination) error {
	// ----------------------- Write two resources -----------------------
	spec.WriteMode = specs.WriteModeOverwrite
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := "cq_test_write_overwrite"
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC()
	tables := []*schema.Table{
		table,
	}

	// These are the indexes for fields in the test-table (without _cq_sync_time and _cq_source_name)
	// To get the indices in structs that include those fields, add +2
	// We need to calculate this before calling Migrate (because Migrate adds 2 columns to the table)
	uuidFieldIndex := getColumnIndex(table, "uuid")
	intFieldIndex := getColumnIndex(table, "int")
	const cqSyncTimeFieldIndex = -1

	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "source_name"
	resource1 := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.TestData(),
	}
	resource2 := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.TestData(),
	}

	if err := resource2.Data[uuidFieldIndex].Set(uuid.New().String()); err != nil {
		return err
	}

	resources := []schema.DestinationResource{
		resource1,
		resource2,
	}

	if err := p.writeAll(ctx, tables, sourceName, syncTime, resources); err != nil {
		return fmt.Errorf("failed to write all: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resource, got %d", len(resourcesRead))
	}

	expectedResourceData1 := appendSourceNameAndSyncTimeToCqTypes(sourceName, syncTime, resource1.Data)
	expectedResourceData2 := appendSourceNameAndSyncTimeToCqTypes(sourceName, syncTime, resource2.Data)

	expectedResourceData := []schema.CQTypes{
		expectedResourceData1,
		expectedResourceData2,
	}

	// Because we have 2 resources, we need to sort (by UUID) to compare them effectively.
	sort.Slice(expectedResourceData, func(i, j int) bool {
		// We add '+2' to uuidFieldIndex because we add 2 columns (_cq_sync_time and source_name)
		return expectedResourceData[i][uuidFieldIndex+2].String() < expectedResourceData[j][uuidFieldIndex+2].String()
	})
	sort.Slice(resourcesRead, func(i, j int) bool {
		// We add '+2' to uuidFieldIndex because we add 2 columns (_cq_sync_time and source_name)
		return resourcesRead[i][uuidFieldIndex+2].String() < resourcesRead[j][uuidFieldIndex+2].String()
	})

	if !expectedResourceData[0].Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data[0] to be %v, got %v", expectedResourceData[0], resourcesRead[0])
	}
	if !expectedResourceData[1].Equal(resourcesRead[1]) {
		return fmt.Errorf("expected data[1] to be %v, got %v", expectedResourceData[1], resourcesRead[1])
	}

	// ----------------------- Overwrite the first resource -----------------------
	secondSyncTime := syncTime.Add(time.Minute).UTC()
	if err := resource1.Data[intFieldIndex].Set(22); err != nil {
		return err
	}
	if err := resources[0].Data[intFieldIndex].Set(22); err != nil {
		return err
	}
	if err := expectedResourceData[0][intFieldIndex+2].Set(22); err != nil {
		return err
	}
	if err := expectedResourceData[0][cqSyncTimeFieldIndex+2].Set(secondSyncTime); err != nil {
		return err
	}

	// write second time
	if err := p.writeOne(ctx, tables, sourceName, secondSyncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}

	sort.Slice(resourcesRead, func(i, j int) bool {
		// We add '+2' to uuidFieldIndex because we add 2 columns (_cq_sync_time and source_name)
		return resourcesRead[i][uuidFieldIndex+2].String() < resourcesRead[j][uuidFieldIndex+2].String()
	})

	if !expectedResourceData[0].Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data[0] to be %v, got %v", expectedResourceData[0], resourcesRead[0])
	}
	if !expectedResourceData[1].Equal(resourcesRead[1]) {
		return fmt.Errorf("expected data[1] to be %v, got %v", expectedResourceData[1], resourcesRead[1])
	}

	// ----------------------- test delete-stale -----------------------

	if !s.tests.SkipDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, secondSyncTime); err != nil {
			return fmt.Errorf("failed to delete stale data second time: %w", err)
		}

		resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
		if err != nil {
			return fmt.Errorf("failed to read all second time: %w", err)
		}
		if len(resourcesRead) != 1 {
			return fmt.Errorf("expected 1 resource, got %d", len(resourcesRead))
		}

		if !expectedResourceData[0].Equal(resourcesRead[0]) {
			return fmt.Errorf("expected data to be %v, got %v", expectedResourceData[0], resourcesRead[0])
		}
	}

	return nil
}

func (*destinationTestSuite) destinationPluginTestWriteAppend(ctx context.Context, p *DestinationPlugin, logger zerolog.Logger, spec specs.Destination) error {
	// -----------------------------------------------

	spec.WriteMode = specs.WriteModeAppend
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := "cq_test_write_append"
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC()
	tables := []*schema.Table{
		table,
	}

	// These are the indexes for fields in the test-table (without _cq_sync_time and _cq_source_name)
	// To get the indices in structs that include those fields, add +2
	// We need to calculate this before calling Migrate (because Migrate adds 2 columns to the table)
	const cqSyncTimeFieldIndex = -1
	uuidFieldIndex := getColumnIndex(table, "uuid")

	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "source_name"
	resource1 := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.TestData(),
	}
	resource2 := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.TestData(),
	}

	if err := resource2.Data[uuidFieldIndex].Set(uuid.New().String()); err != nil {
		return err
	}

	if err := p.writeOne(ctx, tables, sourceName, syncTime, resource1); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	// we dont use time.now because looks like there is some strange
	// issue on windows machine on github actions where it returns the same thing
	// for all calls.
	secondSyncTime := syncTime.Add(time.Second).UTC()
	// write second time
	if err := p.writeOne(ctx, tables, sourceName, secondSyncTime, resource2); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resources, got %d", len(resourcesRead))
	}

	expectedResourceData := []schema.CQTypes{
		appendSourceNameAndSyncTimeToCqTypes(sourceName, syncTime, resource1.Data),
		appendSourceNameAndSyncTimeToCqTypes(sourceName, secondSyncTime, resource2.Data),
	}

	sort.Slice(expectedResourceData, func(i, j int) bool {
		// We add '+2' to uuidFieldIndex because we add 2 columns (_cq_sync_time and source_name)
		return expectedResourceData[i][cqSyncTimeFieldIndex+2].String() < expectedResourceData[j][cqSyncTimeFieldIndex+2].String()
	})
	sort.Slice(resourcesRead, func(i, j int) bool {
		// We add '+2' to uuidFieldIndex because we add 2 columns (_cq_sync_time and source_name)
		return resourcesRead[i][cqSyncTimeFieldIndex+2].String() < resourcesRead[j][cqSyncTimeFieldIndex+2].String()
	})

	if !expectedResourceData[0].Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data[0] to be %v, got %v", expectedResourceData[0], resourcesRead[0])
	}
	if !expectedResourceData[1].Equal(resourcesRead[1]) {
		return fmt.Errorf("expected data[1] to be %v, got %v", expectedResourceData[1], resourcesRead[1])
	}

	return nil
}

func DestinationPluginTestSuiteRunner(t *testing.T, p *DestinationPlugin, spec interface{}, tests DestinationTestSuiteTests) {
	t.Helper()
	destSpec := specs.Destination{
		Name: "testsuite",
		Spec: spec,
	}
	suite := &destinationTestSuite{
		tests: tests,
	}
	ctx := context.Background()
	logger := getTestLogger(t)

	t.Run("TestWriteOverwrite", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipOverwrite {
			t.Skip("skipping TestWriteOverwrite")
			return
		}
		if err := suite.destinationPluginTestWriteOverwrite(ctx, p, logger, destSpec); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestWriteAppend", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipAppend {
			t.Skip("skipping TestWriteAppend")
			return
		}
		if err := suite.destinationPluginTestWriteAppend(ctx, p, logger, destSpec); err != nil {
			t.Fatal(err)
		}
	})
}

// appendSourceNameAndSyncTimeToCqTypes appends {source_name, sync_time} to the beginning of the table CQTypes.
// nolint: unparam
func appendSourceNameAndSyncTimeToCqTypes(sourceName string, syncTime time.Time, cqtypes schema.CQTypes) schema.CQTypes {
	return append(schema.CQTypes{
		&schema.Text{
			Str:    sourceName,
			Status: schema.Present,
		},
		&schema.Timestamptz{
			Time:   syncTime,
			Status: schema.Present,
		},
	}, cqtypes...)
}

// Returns the index of a column in a table.
// Note that the table doesn't contain the _cq_sync_time and source_name columns - so if your data contains
// them you probably need to add `2` to this result.
func getColumnIndex(table *schema.Table, column string) int {
	for i, c := range table.Columns {
		if c.Name == column {
			return i
		}
	}

	panic("Failed to get index of column" + column)
}
