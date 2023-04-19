package destination

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/cloudquery/plugin-sdk/v2/types"
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
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		destSpec.Name = "test_write_overwrite"
		p := newPlugin()
		if err := suite.destinationPluginTestWriteOverwrite(ctx, mem, p, logger, destSpec); err != nil {
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
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		destSpec.Name = "test_write_overwrite_delete_stale"
		p := newPlugin()
		if err := suite.destinationPluginTestWriteOverwriteDeleteStale(ctx, mem, p, logger, destSpec); err != nil {
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
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		destSpec.WriteMode = specs.WriteModeOverwrite
		destSpec.Name = "test_migrate_overwrite"
		suite.destinationPluginTestMigrate(ctx, mem, t, newPlugin, logger, destSpec, tests.MigrateStrategyOverwrite)
	})

	t.Run("TestMigrateOverwriteForce", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrateOverwriteForce {
			t.Skip("skipping " + t.Name())
		}
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		destSpec.WriteMode = specs.WriteModeOverwrite
		destSpec.MigrateMode = specs.MigrateModeForced
		destSpec.Name = "test_migrate_overwrite_force"
		suite.destinationPluginTestMigrate(ctx, mem, t, newPlugin, logger, destSpec, tests.MigrateStrategyOverwrite)
	})

	t.Run("TestWriteAppend", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipAppend {
			t.Skip("skipping " + t.Name())
		}
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		destSpec.Name = "test_write_append"
		p := newPlugin()
		if err := suite.destinationPluginTestWriteAppend(ctx, mem, p, logger, destSpec); err != nil {
			t.Fatal(err)
		}
		if err := p.Close(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestMigrateAppend", func(t *testing.T) {
		t.Helper()
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		if suite.tests.SkipMigrateAppend {
			t.Skip("skipping " + t.Name())
		}
		destSpec.WriteMode = specs.WriteModeAppend
		destSpec.Name = "test_migrate_append"
		suite.destinationPluginTestMigrate(ctx, mem, t, newPlugin, logger, destSpec, tests.MigrateStrategyAppend)
	})

	t.Run("TestMigrateAppendForce", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrateAppendForce {
			t.Skip("skipping " + t.Name())
		}
		mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
		defer mem.AssertSize(t, 0)
		destSpec.WriteMode = specs.WriteModeAppend
		destSpec.MigrateMode = specs.MigrateModeForced
		destSpec.Name = "test_migrate_append_force"
		suite.destinationPluginTestMigrate(ctx, mem, t, newPlugin, logger, destSpec, tests.MigrateStrategyAppend)
	})
}

func sortRecordsBySyncTime(table *schema.Table, records []arrow.Record) {
	syncTimeIndex := table.Columns.Index(schema.CqSyncTimeColumn.Name)
	cqIDIndex := table.Columns.Index(schema.CqIDColumn.Name)
	sortRecordsBySyncTimeIndex(syncTimeIndex, cqIDIndex, records)
}

func sortRecordsBySyncTimeArrow(s *arrow.Schema, records []arrow.Record) {
	syncTimeIndex := s.FieldIndices(schema.CqSyncTimeColumn.Name)
	if len(syncTimeIndex) != 1 {
		panic("no CqSyncTimeColumn in schema")
	}
	cqIDIndex := s.FieldIndices(schema.CqIDColumn.Name)
	if len(syncTimeIndex) != 1 {
		panic("no CqIDColumn in schema")
	}
	sortRecordsBySyncTimeIndex(syncTimeIndex[0], cqIDIndex[0], records)
}

func sortRecordsBySyncTimeIndex(syncTimeIndex, cqIDIndex int, records []arrow.Record) {
	sort.Slice(records, func(i, j int) bool {
		// sort by sync time, then UUID
		t1 := records[i].Column(syncTimeIndex).(*array.Timestamp)
		t2 := records[j].Column(syncTimeIndex).(*array.Timestamp)
		first := t1.Value(0).ToTime(t1.DataType().(*arrow.TimestampType).Unit)
		second := t2.Value(0).ToTime(t2.DataType().(*arrow.TimestampType).Unit)
		if first.Equal(second) {
			// Since our cq_id UUIDs are version 4 (completely random, no time-component UUID) this is only a stable tie-breaker
			firstUUID := records[i].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			secondUUID := records[j].Column(cqIDIndex).(*types.UUIDArray).Value(0).String()
			return strings.Compare(firstUUID, secondUUID) < 0
		}
		return first.Before(second)
	})
}
