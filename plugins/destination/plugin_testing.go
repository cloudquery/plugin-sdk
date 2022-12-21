package destination

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/testdata"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PluginTestSuite struct {
	tests PluginTestSuiteTests
}

type PluginTestSuiteTests struct {
	// SkipOverwrite skips testing for "overwrite" mode. Use if the destination
	//	// plugin doesn't support this feature.
	SkipOverwrite bool

	// SkipDeleteStale skips testing "delete-stale" mode. Use if the destination
	// plugin doesn't support this feature.
	SkipDeleteStale bool

	// SkipAppend skips testing for "append" mode. Use if the destination
	// plugin doesn't support this feature.
	SkipAppend bool

	// SkipSecondAppend skips the second append step in the test.
	// This is useful in cases like cloud storage where you can't append to an
	// existing object after the file has been closed.
	SkipSecondAppend bool
}

func (s *PluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeOverwrite
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := "cq_test_write_overwrite"
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*schema.Table{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testOverwriteSource" + uuid.NewString()
	resource := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.GenTestData(table),
	}
	_ = resource.Data[0].Set(sourceName)
	_ = resource.Data[1].Set(syncTime)

	resource2 := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.GenTestData(table),
	}
	_ = resource2.Data[0].Set(sourceName)
	_ = resource2.Data[1].Set(syncTime)

	resources := []schema.DestinationResource{
		resource,
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

	if diff := resource.Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("expected first resource diff: %s", diff)
	}

	if diff := resource2.Data.Diff(resourcesRead[1]); diff != "" {
		return fmt.Errorf("expected second resource diff: %s", diff)
		// return fmt.Errorf("expected second resource to be:\n%v\ngot:\n%v", resource.Data, resourcesRead[1])
	}

	secondSyncTime := syncTime.Add(time.Second).UTC()
	_ = resource.Data[1].Set(secondSyncTime)
	// write second time
	if err := p.writeOne(ctx, tables, sourceName, secondSyncTime, resource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("after overwrite expected 2 resources, got %d", len(resourcesRead))
	}

	if diff := resource2.Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("after overwrite expected first resource diff: %s", diff)
	}

	if diff := resource.Data.Diff(resourcesRead[1]); diff != "" {
		return fmt.Errorf("after overwrite expected second resource diff: %s", diff)
	}

	if !s.tests.SkipDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, secondSyncTime); err != nil {
			return fmt.Errorf("failed to delete stale data second time: %w", err)
		}
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	if len(resourcesRead) != 1 {
		return fmt.Errorf("expected 1 resource after delete stale, got %d", len(resourcesRead))
	}

	if diff := resource.Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("after delete stale expected second resource diff: %s", diff)
	}

	return nil
}

func (s *PluginTestSuite) destinationPluginTestWriteAppend(ctx context.Context, p *Plugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeAppend
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := "cq_test_write_append"
	table := testdata.TestTable(tableName)
	syncTime := time.Now().UTC().Round(1 * time.Second)
	tables := []*schema.Table{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "testAppendSource" + uuid.NewString()
	resource := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.GenTestData(table),
	}
	_ = resource.Data[0].Set(sourceName)
	_ = resource.Data[1].Set(syncTime)

	if err := p.writeOne(ctx, tables, sourceName, syncTime, resource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resource2 := schema.DestinationResource{
		TableName: table.Name,
		Data:      testdata.GenTestData(table),
	}

	if !s.tests.SkipSecondAppend {
		// we dont use time.now because looks like there is some strange
		// issue on windows machine on github actions where it returns the same thing
		// for all calls.
		secondSyncTime := syncTime.Add(10 * time.Second).UTC()
		_ = resource2.Data[0].Set(sourceName)
		_ = resource2.Data[1].Set(secondSyncTime)
		// write second time
		if err := p.writeOne(ctx, tables, sourceName, secondSyncTime, resource2); err != nil {
			return fmt.Errorf("failed to write one second time: %w", err)
		}
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}

	expectedResource := 2
	if s.tests.SkipSecondAppend {
		expectedResource = 1
	}

	if len(resourcesRead) != expectedResource {
		return fmt.Errorf("expected %d resources, got %d", expectedResource, len(resourcesRead))
	}

	if diff := resource.Data.Diff(resourcesRead[0]); diff != "" {
		return fmt.Errorf("first expected resource diff: %s", diff)
	}

	if !s.tests.SkipSecondAppend {
		if diff := resource2.Data.Diff(resourcesRead[1]); diff != "" {
			return fmt.Errorf("second expected resource diff: %s", diff)
		}
	}

	return nil
}

func getTestLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.TraceLevel).With().Timestamp().Logger()
}

func PluginTestSuiteRunner(t *testing.T, p *Plugin, spec interface{}, tests PluginTestSuiteTests) {
	t.Helper()
	destSpec := specs.Destination{
		Name: "testsuite",
		Spec: spec,
	}
	suite := &PluginTestSuite{
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
