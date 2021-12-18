package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/faker/v3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmccombs/hcl2json/convert"
)

type ResourceTestData struct {
	Table          *schema.Table
	Config         interface{}
	Resources      []string
	Configure      func(logger hclog.Logger, data interface{}) (schema.ClientMeta, error)
	SkipEmptyJsonB bool
}

func TestResource(t *testing.T, providerCreator func() *provider.Provider, resource ResourceTestData) {
	if err := faker.SetRandomMapAndSliceMinSize(1); err != nil {
		t.Fatal(err)
	}
	if err := faker.SetRandomMapAndSliceMaxSize(1); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	pool, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Release()
	l := logging.New(hclog.DefaultOptions)
	migrator := provider.NewTableCreator(l)
	if err := migrator.CreateTable(ctx, conn, resource.Table, nil); err != nil {
		assert.FailNow(t, fmt.Sprintf("failed to create tables %s", resource.Table.Name), err)
	}
	// Write configuration as a block and extract it out passing that specific block data as part of the configure provider
	f := hclwrite.NewFile()
	f.Body().AppendBlock(gohcl.EncodeAsBlock(resource.Config, "configuration"))
	data, err := convert.Bytes(f.Bytes(), "config.json", convert.Options{})
	require.Nil(t, err)
	hack := map[string]interface{}{}
	require.Nil(t, json.Unmarshal(data, &hack))
	data, err = json.Marshal(hack["configuration"].([]interface{})[0])
	require.Nil(t, err)

	testProvider := providerCreator()

	// No need for configuration or db connection, get it out of the way first
	testTableIdentifiersForProvider(t, testProvider)

	testProvider.Logger = l
	testProvider.Configure = resource.Configure
	_, err = testProvider.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "",
		Connection: cqproto.ConnectionDetails{DSN: getEnv("DATABASE_URL",
			"host=localhost user=postgres password=pass DB.name=postgres port=5432")},
		Config: data,
	})
	assert.Nil(t, err)

	err = testProvider.FetchResources(context.Background(), &cqproto.FetchResourcesRequest{Resources: []string{findResourceFromTableName(resource.Table, testProvider.ResourceMap)}}, &fakeResourceSender{Errors: []string{}})
	assert.Nil(t, err)
	verifyNoEmptyColumns(t, resource, conn)
}

func findResourceFromTableName(table *schema.Table, tables map[string]*schema.Table) string {
	for resource, t := range tables {
		if table.Name == t.Name {
			return resource
		}
	}
	return ""
}

type fakeResourceSender struct {
	Errors []string
}

func (f *fakeResourceSender) Send(r *cqproto.FetchResourcesResponse) error {
	if r.Error != "" {
		fmt.Printf(r.Error)
		f.Errors = append(f.Errors, r.Error)
	}
	return nil
}

func setupDatabase() (*pgxpool.Pool, error) {
	dbCfg, err := pgxpool.ParseConfig(getEnv("DATABASE_URL",
		"host=localhost user=postgres password=pass DB.name=postgres port=5432"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse config. %w", err)
	}
	ctx := context.Background()
	dbCfg.MaxConns = 1
	dbCfg.LazyConnect = true
	pool, err := pgxpool.ConnectConfig(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database. %w", err)
	}
	return pool, nil

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func verifyNoEmptyColumns(t *testing.T, tc ResourceTestData, conn pgxscan.Querier) {
	// Test that we don't have missing columns and have exactly one entry for each table
	for _, table := range getTablesFromMainTable(tc.Table) {
		query := fmt.Sprintf("select * FROM %s ", table)
		rows, err := conn.Query(context.Background(), query)
		if err != nil {
			t.Fatal(err)
		}
		count := 0
		for rows.Next() {
			count += 1
		}
		if count < 1 {
			t.Fatalf("expected to have at least 1 entry at table %s got %d", table, count)
		}
		if tc.SkipEmptyJsonB {
			continue
		}
		query = fmt.Sprintf("select t.* FROM %s as t WHERE to_jsonb(t) = jsonb_strip_nulls(to_jsonb(t))", table)
		rows, err = conn.Query(context.Background(), query)
		if err != nil {
			t.Fatal(err)
		}
		count = 0
		for rows.Next() {
			count += 1
		}
		if count < 1 {
			t.Fatalf("row at table %s has an empty column", table)
		}
	}
}

func getTablesFromMainTable(table *schema.Table) []string {
	var res []string
	res = append(res, table.Name)
	for _, t := range table.Relations {
		res = append(res, getTablesFromMainTable(t)...)
	}
	return res
}

func testTableIdentifiersForProvider(t *testing.T, prov *provider.Provider) {
	t.Run("testTableIdentifiersForProvider", func(t *testing.T) {
		t.Parallel()
		for _, res := range prov.ResourceMap {
			res := res
			t.Run(res.Name, func(t *testing.T) {
				testTableIdentifiers(t, res)
			})
		}
	})
}

func testTableIdentifiers(t *testing.T, table *schema.Table) {
	t.Parallel()
	assert.NoError(t, schema.ValidateTable(table))

	for _, res := range table.Relations {
		res := res
		t.Run(res.Name, func(t *testing.T) {
			testTableIdentifiers(t, res)
		})
	}
}
