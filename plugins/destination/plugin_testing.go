package destination

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/cloudquery/plugin-sdk/testdata"
	"github.com/rs/zerolog"
)

type PluginTestSuite struct {
	tests PluginTestSuiteTests
}

// MigrateStrategy defines which tests we should include
type MigrateStrategy struct {
	AddColumn           specs.MigrateMode
	AddColumnNotNull    specs.MigrateMode
	RemoveColumn        specs.MigrateMode
	RemoveColumnNotNull specs.MigrateMode
	ChangeColumn        specs.MigrateMode
}

type PluginTestSuiteTests struct {
	// SkipOverwrite skips testing for "overwrite" mode. Use if the destination
	// plugin doesn't support this feature.
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

	// SkipMigrateAppend skips a test for the migrate function where a column is added,
	// data is appended, then the column is removed and more data appended, checking that the migrations handle
	// this correctly.
	SkipMigrateAppend bool
	// SkipMigrateAppendForce skips a test for the migrate function where a column is changed in force mode
	SkipMigrateAppendForce bool

	// SkipMigrateOverwrite skips a test for the migrate function where a column is added,
	// data is appended, then the column is removed and more data overwritten, checking that the migrations handle
	// this correctly.
	SkipMigrateOverwrite bool
	// SkipMigrateOverwriteForce skips a test for the migrate function where a column is changed in force mode
	SkipMigrateOverwriteForce bool

	MigrateStrategyOverwrite MigrateStrategy
	MigrateStrategyAppend    MigrateStrategy
}

func getTestLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.TraceLevel).With().Timestamp().Logger()
}

type NewPluginFunc func() *Plugin

func PluginTestSuiteRunner(t *testing.T, newPlugin NewPluginFunc, destSpec specs.Destination, tests PluginTestSuiteTests) {
	t.Helper()
	destSpec.Name = "testsuite"

	suite := &PluginTestSuite{
		tests: tests,
	}
	ctx := context.Background()
	logger := getTestLogger(t)

	t.Run("TestWriteOverwrite", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipOverwrite {
			t.Skip("skipping " + t.Name())
		}
		destSpec.Name = "test_write_overwrite"
		p := newPlugin()
		if err := suite.destinationPluginTestWriteOverwrite(ctx, p, logger, destSpec); err != nil {
			t.Fatal(err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestWriteOverwriteDeleteStale", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipOverwrite || suite.tests.SkipDeleteStale {
			t.Skip("skipping " + t.Name())
		}
		destSpec.Name = "test_write_overwrite_delete_stale"
		p := newPlugin()
		if err := suite.destinationPluginTestWriteOverwriteDeleteStale(ctx, p, logger, destSpec); err != nil {
			t.Fatal(err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestMigrateOverwrite", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrateOverwrite {
			t.Skip("skipping " + t.Name())
		}
		destSpec.WriteMode = specs.WriteModeOverwrite
		destSpec.Name = "test_migrate_overwrite"
		suite.destinationPluginTestMigrate(ctx, t, newPlugin, logger, destSpec, tests.MigrateStrategyOverwrite)
	})

	t.Run("TestMigrateOverwriteForce", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrateOverwriteForce {
			t.Skip("skipping " + t.Name())
		}
		destSpec.WriteMode = specs.WriteModeOverwrite
		destSpec.MigrateMode = specs.MigrateModeForced
		destSpec.Name = "test_migrate_overwrite_force"
		suite.destinationPluginTestMigrate(ctx, t, newPlugin, logger, destSpec, tests.MigrateStrategyOverwrite)
	})

	t.Run("TestWriteAppend", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipAppend {
			t.Skip("skipping " + t.Name())
		}
		destSpec.Name = "test_write_append"
		p := newPlugin()
		if err := suite.destinationPluginTestWriteAppend(ctx, p, logger, destSpec); err != nil {
			t.Fatal(err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestMigrateAppend", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrateAppend {
			t.Skip("skipping " + t.Name())
		}
		destSpec.WriteMode = specs.WriteModeAppend
		destSpec.Name = "test_migrate_append"
		suite.destinationPluginTestMigrate(ctx, t, newPlugin, logger, destSpec, tests.MigrateStrategyAppend)
	})

	t.Run("TestMigrateAppendForce", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrateAppendForce {
			t.Skip("skipping " + t.Name())
		}
		destSpec.WriteMode = specs.WriteModeAppend
		destSpec.MigrateMode = specs.MigrateModeForced
		destSpec.Name = "test_migrate_append_force"
		suite.destinationPluginTestMigrate(ctx, t, newPlugin, logger, destSpec, tests.MigrateStrategyAppend)
	})
}

func createTestResources(table *schema.Table, sourceName string, syncTime time.Time, count int) []schema.DestinationResource {
	resources := make([]schema.DestinationResource, count)
	for i := 0; i < count; i++ {
		resource := schema.DestinationResource{
			TableName: table.Name,
			Data:      testdata.GenTestData(table),
		}
		_ = resource.Data[0].Set(sourceName)
		_ = resource.Data[1].Set(syncTime)
		resources[i] = resource
	}
	return resources
}

func sortResources(table *schema.Table, resources []schema.DestinationResource) {
	cqIDIndex := table.Columns.Index(schema.CqIDColumn.Name)
	syncTimeIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	sort.Slice(resources, func(i, j int) bool {
		// sort by sync time, then UUID
		if !resources[i].Data[syncTimeIndex].Equal(resources[j].Data[syncTimeIndex]) {
			return resources[i].Data[syncTimeIndex].Get().(time.Time).Before(resources[j].Data[syncTimeIndex].Get().(time.Time))
		}
		return resources[i].Data[cqIDIndex].String() < resources[j].Data[cqIDIndex].String()
	})
}

func sortCQTypes(table *schema.Table, resources []schema.CQTypes) {
	cqIDIndex := table.Columns.Index(schema.CqIDColumn.Name)
	syncTimeIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	sort.Slice(resources, func(i, j int) bool {
		// sort by sync time, then UUID
		if !resources[i][syncTimeIndex].Equal(resources[j][syncTimeIndex]) {
			return resources[i][syncTimeIndex].Get().(time.Time).Before(resources[j][syncTimeIndex].Get().(time.Time))
		}
		return resources[i][cqIDIndex].String() < resources[j][cqIDIndex].String()
	})
}
