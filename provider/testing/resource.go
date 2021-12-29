package testing

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/cq-provider-sdk/testlog"
	"github.com/cloudquery/faker/v3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

type ResourceTestCase struct {
	Provider       *provider.Provider
	Table          *schema.Table
	Config         string
	SnapshotsDir   string
	SkipEmptyJsonB bool
}

// IntegrationTest - creates resources using terraform, fetches them to db and compares with expected values
func TestResource(t *testing.T, resource ResourceTestCase) {
	t.Parallel()
	t.Helper()
	if err := faker.SetRandomMapAndSliceMinSize(1); err != nil {
		t.Fatal(err)
	}
	if err := faker.SetRandomMapAndSliceMaxSize(1); err != nil {
		t.Fatal(err)
	}

	// No need for configuration or db connection, get it out of the way first
	// testTableIdentifiersForProvider(t, resource.Provider)

	pool, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Release()

	l := testlog.New(t)
	l.SetLevel(hclog.Debug)
	resource.Provider.Logger = l
	tableCreator := provider.NewTableCreator(l)
	if err := tableCreator.CreateTable(context.Background(), conn, resource.Table, nil); err != nil {
		assert.FailNow(t, fmt.Sprintf("failed to create tables %s", resource.Table.Name), err)
	}

	if err := deleteTables(conn, resource.Table); err != nil {
		t.Fatal(err)
	}

	if err = fetch(t, &resource); err != nil {
		t.Fatal(err)
	}

	verifyNoEmptyColumns(t, resource, conn)

	if err := conn.Conn().Close(ctx); err != nil {
		t.Fatal(err)
	}

}

// fetch - fetches resources from the cloud and puts them into database. database config can be specified via DATABASE_URL env variable
func fetch(t *testing.T, resource *ResourceTestCase) error {
	t.Logf("%s fetch resources", resource.Table.Name)

	if _, err := resource.Provider.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "",
		Connection: cqproto.ConnectionDetails{DSN: getEnv("DATABASE_URL",
			"host=localhost user=postgres password=pass DB.name=postgres port=5432")},
		Config:        []byte(resource.Config),
		DisableDelete: true,
	}); err != nil {
		return err
	}

	var resourceSender = &fakeResourceSender{
		Errors: []string{},
	}

	if err := resource.Provider.FetchResources(context.Background(),
		&cqproto.FetchResourcesRequest{
			Resources: []string{findResourceFromTableName(resource.Table, resource.Provider.ResourceMap)},
		},
		resourceSender,
	); err != nil {
		return err
	}

	if len(resourceSender.Errors) > 0 {
		return fmt.Errorf("error/s occur during test, %s", strings.Join(resourceSender.Errors, ", "))
	}

	return nil
}

func deleteTables(conn *pgxpool.Conn, table *schema.Table) error {
	s := sq.Delete(table.Name)
	sql, args, err := s.ToSql()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.TODO(), sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func verifyNoEmptyColumns(t *testing.T, tc ResourceTestCase, conn pgxscan.Querier) {
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

var (
	dbConnOnce sync.Once
	pool       *pgxpool.Pool
	dbErr      error
)

func setupDatabase() (*pgxpool.Pool, error) {
	dbConnOnce.Do(func() {
		var dbCfg *pgxpool.Config
		dbCfg, dbErr = pgxpool.ParseConfig(getEnv("DATABASE_URL", "host=localhost user=postgres password=pass DB.name=postgres port=5432"))
		if dbErr != nil {
			return
		}
		ctx := context.Background()
		dbCfg.MaxConns = 15
		dbCfg.LazyConnect = true
		pool, dbErr = pgxpool.ConnectConfig(ctx, dbCfg)
	})
	return pool, dbErr

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTablesFromMainTable(table *schema.Table) []string {
	var res []string
	res = append(res, table.Name)
	for _, t := range table.Relations {
		res = append(res, getTablesFromMainTable(t)...)
	}
	return res
}
