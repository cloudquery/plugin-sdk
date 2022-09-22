package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToSnake(t *testing.T) {
	type test struct {
		Camel string
		Snake string
	}

	generatorTests := []test{
		{Camel: "TestCamelCase", Snake: "test_camel_case"},
		{Camel: "TestCamelCase", Snake: "test_camel_case"},
		{Camel: "AccountID", Snake: "account_id"},
		{Camel: "IDs", Snake: "ids"},
		{Camel: "PostgreSQL", Snake: "postgre_sql"},
		{Camel: "QueryStoreRetention", Snake: "query_store_retention"},
		{Camel: "TestCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Camel: "testCamelCaseLongString", Snake: "test_camel_case_long_string"},
	}
	t.Parallel()
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Snake, ToSnake(tc.Camel))
		})
	}
}
