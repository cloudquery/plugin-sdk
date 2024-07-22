package plugin

import (
	"testing"
	"time"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTestIgnoreNullsInLists(t *testing.T) {
	s := &WriterTestSuite{ignoreNullsInLists: true}

	tg := schema.NewTestDataGenerator(0)
	source := schema.TestTable("ignore_nulls_in_lists", schema.TestSourceOptions{})
	resource := s.handleNulls(tg.Generate(source, schema.GenTestDataOptions{
		SourceName: "allow_null",
		SyncTime:   time.Now(),
		MaxRows:    100,
		NullRows:   false,
	}))
	for _, c := range resource.Columns() {
		assertNoNullsInLists(t, c)
	}

	resource = s.handleNulls(tg.Generate(source, schema.GenTestDataOptions{
		SourceName: "ignore_nulls_in_lists",
		SyncTime:   time.Now(),
		MaxRows:    100,
		NullRows:   true,
	}))
	for _, c := range resource.Columns() {
		assertNoNullsInLists(t, c)
	}
}

func assertNoNullsInLists(t *testing.T, arr arrow.Array) {
	// traverse
	switch arr := arr.(type) {
	case array.ListLike:
		assert.Zero(t, arr.ListValues().NullN())
	case *array.Struct:
		for i := 0; i < arr.NumField(); i++ {
			assertNoNullsInLists(t, arr.Field(i))
		}
	}
}

func TestWithTestSourceAllowNull(t *testing.T) {
	s := &WriterTestSuite{allowNull: func(dt arrow.DataType) bool {
		switch dt.(type) {
		case *arrow.StructType, arrow.ListLikeType:
			return false
		default:
			return true
		}
	}}

	tg := schema.NewTestDataGenerator(0)
	source := schema.TestTable("allow_null", schema.TestSourceOptions{})
	resource := s.handleNulls(tg.Generate(source, schema.GenTestDataOptions{
		SourceName: "allow_null",
		SyncTime:   time.Now(),
		MaxRows:    100,
		NullRows:   false,
	}))
	for _, c := range resource.Columns() {
		assertNoNulls(t, s.allowNull, c)
	}

	resource = s.handleNulls(tg.Generate(source, schema.GenTestDataOptions{
		SourceName: "allow_null",
		SyncTime:   time.Now(),
		MaxRows:    100,
		NullRows:   true,
	}))
	for _, c := range resource.Columns() {
		assertNoNulls(t, s.allowNull, c)
	}
}

func assertNoNulls(t *testing.T, allowNull AllowNullFunc, arr arrow.Array) {
	require.NotNil(t, allowNull)

	if !allowNull(arr.DataType()) {
		assert.Zero(t, arr.NullN())
	}

	// traverse
	switch arr := arr.(type) {
	case array.ListLike:
		assertNoNulls(t, allowNull, arr.ListValues())
	case *array.Struct:
		for i := 0; i < arr.NumField(); i++ {
			assertNoNulls(t, allowNull, arr.Field(i))
		}
	}
}
