package plugin

import (
	"context"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type WriterTestSuite struct {
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

// SafeMigrations defines which migrations are supported by the plugin in safe migrate mode
type SafeMigrations struct {
	AddColumn           bool
	AddColumnNotNull    bool
	RemoveColumn        bool
	RemoveColumnNotNull bool
	ChangeColumn        bool
}

type PluginTestSuiteTests struct {
	// SkipUpsert skips testing with MessageInsert and Upsert=true.
	// Usually when a destination is not supporting primary keys
	SkipUpsert bool

	// SkipDeleteStale skips testing MessageDelete events.
	SkipDeleteStale bool

	// SkipAppend skips testing MessageInsert and Upsert=false.
	SkipInsert bool

	// SkipMigrate skips testing migration
	SkipMigrate bool

	// SafeMigrations defines which tests should work with force migration
	// and which should pass with safe migration
	SafeMigrations SafeMigrations
}

type NewPluginFunc func() *Plugin

func WithTestSourceAllowNull(allowNull func(arrow.DataType) bool) func(o *WriterTestSuite) {
	return func(o *WriterTestSuite) {
		o.allowNull = allowNull
	}
}

func WithTestIgnoreNullsInLists() func(o *WriterTestSuite) {
	return func(o *WriterTestSuite) {
		o.ignoreNullsInLists = true
	}
}

func WithTestDataOptions(opts schema.TestSourceOptions) func(o *WriterTestSuite) {
	return func(o *WriterTestSuite) {
		o.genDatOptions = opts
	}
}

func TestWriterSuiteRunner(t *testing.T, p *Plugin, tests PluginTestSuiteTests, opts ...func(o *WriterTestSuite)) {
	t.Helper()
	suite := &WriterTestSuite{
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

	t.Run("TestDeleteStale", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipDeleteStale {
			t.Skip("skipping " + t.Name())
		}
		if err := suite.testDeleteStale(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestMigrate", func(t *testing.T) {
		t.Helper()
		if suite.tests.SkipMigrate {
			t.Skip("skipping " + t.Name())
		}
		suite.testMigrate(ctx, t, false)
		suite.testMigrate(ctx, t, true)
	})
}
