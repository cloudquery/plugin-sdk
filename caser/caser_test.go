package caser

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

func Test_ToCamel(t *testing.T) {
	type test struct {
		Camel string
		Snake string
	}

	generatorTests := []test{
		{Camel: "testCamelCase", Snake: "test_camel_case"},
		{Camel: "testCamelCase", Snake: "test_camel_case"},
		{Camel: "accountID", Snake: "account_id"},
		{Camel: "arns", Snake: "arns"},
		{Camel: "postgreSQL", Snake: "postgre_sql"},
		{Camel: "queryStoreRetention", Snake: "query_store_retention"},
		{Camel: "testCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Camel: "testCamelCaseLongString", Snake: "test_camel_case_long_string"},
	}
	t.Parallel()
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Camel, ToCamel(tc.Snake))
		})
	}
}

func Test_ToPascal(t *testing.T) {
	type test struct {
		Camel string
		Snake string
	}

	generatorTests := []test{
		{Camel: "TestCamelCase", Snake: "test_camel_case"},
		{Camel: "TestCamelCase", Snake: "test_camel_case"},
		{Camel: "AccountID", Snake: "account_id"},
		{Camel: "Arns", Snake: "arns"},
		{Camel: "PostgreSQL", Snake: "postgre_sql"},
		{Camel: "QueryStoreRetention", Snake: "query_store_retention"},
		{Camel: "TestCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Camel: "TestCamelCaseLongString", Snake: "test_camel_case_long_string"},
	}
	t.Parallel()
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Camel, ToPascal(tc.Snake))
		})
	}
}

func Test_Configure(t *testing.T) {
	type test struct {
		Camel string
		Snake string
	}

	generatorTests := []test{
		{Camel: "CDNs", Snake: "cdns"},
		{Camel: "ARNs", Snake: "arns"},
	}
	ConfigureInitialisms(map[string]bool{"CDN": true, "ARN": true})
	t.Parallel()
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Snake, ToSnake(tc.Camel))
		})
	}
}
