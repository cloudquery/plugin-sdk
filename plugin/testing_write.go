package plugin

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type PluginTestSuite struct {
	tests PluginTestSuiteTests

	plugin *Plugin

	// AllowNull is a custom func to determine whether a data type may be correctly represented as null.
	// Destinations that have problems representing some data types should provide a custom implementation here.
	// If this param is empty, the default is to allow all data types to be nullable.
	// When the value returned by this func is `true` the comparison is made with the empty value instead of null.
	allowNull AllowNullFunc

	// IgnoreNullsInLists allows stripping null values from lists before comparison.
	// Destination setups that don't support nulls in lists should set this to true.
	ignoreNullsInLists bool

	// genDataOptions define how to generate test data and which data types to skip
	genDatOptions schema.TestSourceOptions
}

// MigrateStrategy defines which tests we should include
type MigrateStrategy struct {
	AddColumn           MigrateMode
	AddColumnNotNull    MigrateMode
	RemoveColumn        MigrateMode
	RemoveColumnNotNull MigrateMode
	ChangeColumn        MigrateMode
}

type PluginTestSuiteTests struct {
	// SkipUpsert skips testing with MessageInsert and Upsert=true.
	// Usually when a destination is not supporting primary keys
	SkipUpsert bool

	// SkipDelete skips testing MessageDelete events.
	SkipDelete bool

	// SkipAppend skips testing MessageInsert and Upsert=false.
	SkipInsert bool

	// SkipMigrate skips testing migration
	SkipMigrate bool

	// MigrateStrategy defines which tests should work with force migration
	// and which should pass with safe migration
	MigrateStrategy MigrateStrategy
}

type NewPluginFunc func() *Plugin

func WithTestSourceAllowNull(allowNull func(arrow.DataType) bool) func(o *PluginTestSuite) {
	return func(o *PluginTestSuite) {
		o.allowNull = allowNull
	}
}

func WithTestIgnoreNullsInLists() func(o *PluginTestSuite) {
	return func(o *PluginTestSuite) {
		o.ignoreNullsInLists = true
	}
}

func WithTestDataOptions(opts schema.TestSourceOptions) func(o *PluginTestSuite) {
	return func(o *PluginTestSuite) {
		o.genDatOptions = opts
	}
}

func PluginTestSuiteRunner(t *testing.T, p *Plugin, tests PluginTestSuiteTests, opts ...func(o *PluginTestSuite)) {
	t.Helper()
	suite := &PluginTestSuite{
		tests:  tests,
		plugin: p,
	}

	for _, opt := range opts {
		opt(suite)
	}

	ctx := context.Background()

	t.Run("TestUpsert", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipUpsert {
			t.Skip("skipping " + t.Name())
		}
		if err := suite.testUpsert(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestInsert", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipInsert {
			t.Skip("skipping " + t.Name())
		}
		if err := suite.testInsert(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestDelete", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipDelete {
			t.Skip("skipping " + t.Name())
		}
		if err := suite.testDelete(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestMigrate", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrate {
			t.Skip("skipping " + t.Name())
		}
		migrateMode := MigrateModeSafe
		writeMode := WriteModeOverwrite
		suite.destinationPluginTestMigrate(ctx, t, p, migrateMode, writeMode, tests.MigrateStrategyOverwrite, opts)
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
