package destination

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/types"
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

func RecordDiff(l arrow.Record, r arrow.Record) string {
	var sb strings.Builder
	if l.NumCols() != r.NumCols() {
		return fmt.Sprintf("different number of columns: %d vs %d", l.NumCols(), r.NumCols())
	}
	if l.NumRows() != r.NumRows() {
		return fmt.Sprintf("different number of rows: %d vs %d", l.NumRows(), r.NumRows())
	}
	for i := 0; i < int(l.NumCols()); i++ {
		edits, err := array.Diff(l.Column(i), r.Column(i))
		if err != nil {
			panic(fmt.Sprintf("left: %v, right: %v, error: %v", l.Column(i).DataType(), r.Column(i).DataType(), err))
		}
		diff := edits.UnifiedDiff(l.Column(i), r.Column(i))
		if diff != "" {
			sb.WriteString(l.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			sb.WriteString("\n")
		}
	}
	return sb.String()
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
		destSpec.MigrateMode = specs.MigrateModeSafe
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
		destSpec.MigrateMode = specs.MigrateModeSafe
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

func sortRecordsBySyncTime(table *schema.Table, records []arrow.Record) {
	syncTimeIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	cqIDIndex := table.Columns.Index(schema.CqIDColumn.Name)
	sort.Slice(records, func(i, j int) bool {
		// sort by sync time, then UUID
		first := records[i].Column(syncTimeIndex).(*array.Timestamp).Value(0).ToTime(arrow.Millisecond)
		second := records[j].Column(syncTimeIndex).(*array.Timestamp).Value(0).ToTime(arrow.Millisecond)
		if first.Equal(second) {
			firstUUID := records[i].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			secondUUID := records[j].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			return strings.Compare(firstUUID, secondUUID) < 0
		}
		return first.Before(second)
	})
}
