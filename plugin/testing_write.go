package plugin

import (
	"context"
	"math/rand"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type WriterTestSuite struct {
	tests WriterTestSuiteTests

	plugin *Plugin

	// AllowNull is a custom func to determine whether a data type may be correctly represented as null.
	// Destinations that have problems representing some data types should provide a custom implementation here.
	// If this param is empty, the default is to allow all data types to be nullable.
	// When the value returned by this func is `true` the comparison is made with the empty value instead of null.
	allowNull AllowNullFunc

	// IgnoreNullsInLists allows stripping null values from lists before comparison.
	// Destination setups that don't support nulls in lists should set this to true.
	ignoreNullsInLists bool

	// useHomogeneousTypes use the same type for all items within an array - useful for destinations that don't support mixed types.
	useHomogeneousTypes bool

	// genDataOptions define how to generate test data and which data types to skip
	genDatOptions schema.TestSourceOptions

	// random seed to use
	randSeed int64

	// rand.Rand
	rand *rand.Rand
}

// SafeMigrations defines which migrations are supported by the plugin in safe migrate mode
type SafeMigrations struct {
	AddColumn              bool
	AddColumnNotNull       bool
	RemoveColumn           bool
	RemoveColumnNotNull    bool
	ChangeColumn           bool
	RemoveUniqueConstraint bool
	MovePKToCQOnly         bool
}

// Migrations defines which migrations should be skipped completely
type Migrations struct {
	RemoveUniqueConstraint bool
	MovePKToCQOnly         bool
}

// WriteTests defines which tests should be skipped in the write test suite
type WriteTests struct {
	DuplicatePK bool
}

type WriterTestSuiteTests struct {
	// SkipUpsert skips testing with message.Insert and Upsert=true.
	// Usually when a destination is not supporting primary keys
	SkipUpsert bool

	// SkipDeleteStale skips testing message.Delete events.
	SkipDeleteStale bool

	// SkipDeleteRecord skips testing message.DeleteRecord events.
	SkipDeleteRecord bool

	// SkipAppend skips testing message.Insert and Upsert=false.
	SkipInsert bool

	// SkipMigrate skips testing migration
	SkipMigrate bool

	// SafeMigrations defines which tests should work with force migration
	// and which should pass with safe migration
	SafeMigrations SafeMigrations

	SkipSpecificMigrations Migrations

	SkipSpecificWriteTests WriteTests
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

func WithRandomSeed(seed int64) func(o *WriterTestSuite) {
	return func(o *WriterTestSuite) {
		o.randSeed = seed
	}
}

func WithHomogeneousTypes() func(o *WriterTestSuite) {
	return func(o *WriterTestSuite) {
		o.useHomogeneousTypes = true
	}
}

func TestWriterSuiteRunner(t *testing.T, p *Plugin, tests WriterTestSuiteTests, opts ...func(o *WriterTestSuite)) {
	suite := &WriterTestSuite{
		tests:  tests,
		plugin: p,
	}

	for _, opt := range opts {
		opt(suite)
	}

	suite.rand = rand.New(rand.NewSource(suite.randSeed))

	ctx := context.Background()

	t.Run("TestUpsert", func(t *testing.T) {
		if suite.tests.SkipUpsert {
			t.Skip("skipping " + t.Name())
		}
		t.Run("Basic", func(t *testing.T) {
			if err := suite.testUpsertBasic(ctx); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("All", func(t *testing.T) {
			if err := suite.testUpsertAll(ctx); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Duplicate PKs", func(t *testing.T) {
			if tests.SkipSpecificWriteTests.DuplicatePK {
				t.Skip("skipping " + t.Name())
			}
			if err := suite.testInsertDuplicatePK(ctx); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("TestInsert", func(t *testing.T) {
		if suite.tests.SkipInsert {
			t.Skip("skipping " + t.Name())
		}
		t.Run("Basic", func(t *testing.T) {
			if err := suite.testInsertBasic(ctx); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("All", func(t *testing.T) {
			if err := suite.testInsertAll(ctx); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("TestDeleteStale", func(t *testing.T) {
		if suite.tests.SkipDeleteStale {
			t.Skip("skipping " + t.Name())
		}
		t.Run("Basic", func(t *testing.T) {
			suite.testDeleteStaleBasic(ctx, t)
		})
		t.Run("All", func(t *testing.T) {
			suite.testDeleteStaleAll(ctx, t)
		})
	})

	t.Run("TestDeleteRecord", func(t *testing.T) {
		if suite.tests.SkipDeleteRecord {
			t.Skip("skipping " + t.Name())
		}
		t.Run("Basic", func(t *testing.T) {
			suite.testDeleteRecordBasic(ctx, t)
		})
		t.Run("DeleteAll", func(t *testing.T) {
			suite.testDeleteAllRecords(ctx, t)
		})
	})

	t.Run("TestMigrate", func(t *testing.T) {
		if suite.tests.SkipMigrate {
			t.Skip("skipping " + t.Name())
		}
		suite.testMigrate(ctx, t, false)
		suite.testMigrate(ctx, t, true)
	})
}
