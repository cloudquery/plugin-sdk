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
		{Camel: "testIPv4", Snake: "test_ipv4"},
		{Camel: "CoreIPs", Snake: "core_ips"},
		{Camel: "CoreIps", Snake: "core_ips"},
		{Camel: "CoreV1", Snake: "core_v1"},
		{Camel: "APIVersion", Snake: "api_version"},
		{Camel: "TTLSecondsAfterFinished", Snake: "ttl_seconds_after_finished"},
		{Camel: "PodCIDRs", Snake: "pod_cidrs"},
		{Camel: "IAMRoles", Snake: "iam_roles"},
		{Camel: "testIAM", Snake: "test_iam"},
		{Camel: "TestAWSMode", Snake: "test_aws_mode"},
	}
	t.Parallel()
	c := New()
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Snake, c.ToSnake(tc.Camel))
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
		{Camel: "accountID", Snake: "account_id"},
		{Camel: "arns", Snake: "arns"},
		{Camel: "postgreSQL", Snake: "postgre_sql"},
		{Camel: "queryStoreRetention", Snake: "query_store_retention"},
		{Camel: "testCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Camel: "testCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Camel: "testIPv4", Snake: "test_ipv4"},
	}
	t.Parallel()
	c := New()
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Camel, c.ToCamel(tc.Snake))
		})
	}
}

func Test_ToPascal(t *testing.T) {
	type test struct {
		Pascal string
		Snake  string
	}

	generatorTests := []test{
		{Pascal: "TestCamelCase", Snake: "test_camel_case"},
		{Pascal: "AccountID", Snake: "account_id"},
		{Pascal: "Arns", Snake: "arns"},
		{Pascal: "PostgreSQL", Snake: "postgre_sql"},
		{Pascal: "QueryStoreRetention", Snake: "query_store_retention"},
		{Pascal: "TestCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Pascal: "TestCamelCaseLongString", Snake: "test_camel_case_long_string"},
		{Pascal: "TestV1", Snake: "test_v1"},
		{Pascal: "TestIPv4", Snake: "test_ipv4"},
		{Pascal: "Ec2", Snake: "ec2"},
		{Pascal: "S3", Snake: "s3"},
	}
	t.Parallel()
	c := New()
	for _, tc := range generatorTests {
		t.Run(tc.Pascal, func(t *testing.T) {
			assert.Equal(t, tc.Pascal, c.ToPascal(tc.Snake))
		})
	}
}

func TestInversion(t *testing.T) {
	type test struct {
		Pascal string
	}

	generatorTests := []test{
		{Pascal: "TestCamelCase"},
		{Pascal: "AccountID"},
		{Pascal: "Arns"},
		{Pascal: "PostgreSQL"},
		{Pascal: "QueryStoreRetention"},
		{Pascal: "TestCamelCaseLongString"},
		{Pascal: "TestCamelCaseLongString"},
		{Pascal: "TestV1"},
		{Pascal: "TestIPv4"},
		{Pascal: "TestIPv4"},
		{Pascal: "S3"},
	}
	t.Parallel()
	c := New()
	for _, tc := range generatorTests {
		t.Run(tc.Pascal, func(t *testing.T) {
			assert.Equal(t, tc.Pascal, c.ToPascal(c.ToSnake(tc.Pascal)))
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
		{Camel: "EC2", Snake: "ec2"},
		{Camel: "S3", Snake: "s3"},
	}
	t.Parallel()
	c := New(WithCustomInitialisms(map[string]bool{"CDN": true, "ARN": true, "EC2": true}))
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Snake, c.ToSnake(tc.Camel))
		})
	}
}

func Test_Exceptions(t *testing.T) {
	type test struct {
		Camel string
		Snake string
	}

	generatorTests := []test{
		{Camel: "TEst", Snake: "test"},
		{Camel: "TTv2", Snake: "ttv2"},
	}
	t.Parallel()
	c := New(WithCustomExceptions(map[string]string{"test": "TEst", "ttv2": "TTv2"}))
	for _, tc := range generatorTests {
		t.Run(tc.Camel, func(t *testing.T) {
			assert.Equal(t, tc.Camel, c.ToCamel(tc.Snake))
		})
	}
}
