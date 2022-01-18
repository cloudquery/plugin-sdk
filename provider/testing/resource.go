package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/cloudquery/cq-provider-sdk/database"
	"github.com/cloudquery/cq-provider-sdk/migration"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/cq-provider-sdk/testlog"
	"github.com/cloudquery/faker/v3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type ResourceTestCase struct {
	Provider       *provider.Provider
	Table          *schema.Table
	Config         string
	SkipEmptyJsonB bool
	// SkipEmptyColumn will skip checking results for empty columns
	SkipEmptyColumn bool
	// SkipEmptyRows will skip checking that results were returned and will just check that fetch worked
	SkipEmptyRows bool
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

	conn, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}

	l := testlog.New(t)
	l.SetLevel(hclog.Debug)
	resource.Provider.Logger = l
	tableCreator := migration.NewTableCreator(l, schema.PostgresDialect{})
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
}

// fetch - fetches resources from the cloud and puts them into database. database config can be specified via DATABASE_URL env variable
func fetch(t *testing.T, resource *ResourceTestCase) error {
	t.Helper()
	t.Logf("%s fetch resources", resource.Table.Name)

	if _, err := resource.Provider.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "",
		Connection: cqproto.ConnectionDetails{DSN: getEnv("DATABASE_URL",
			"host=localhost user=postgres password=pass DB.name=postgres port=5432")},
		Config: []byte(resource.Config),
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

func deleteTables(conn schema.QueryExecer, table *schema.Table) error {
	s := sq.Delete(table.Name)
	sql, args, err := s.ToSql()
	if err != nil {
		return err
	}

	return conn.Exec(context.TODO(), sql, args...)
}

func verifyNoEmptyColumns(t *testing.T, tc ResourceTestCase, conn pgxscan.Querier) {
	t.Helper()
	if tc.SkipEmptyRows {
		t.Logf("table %s marked with SkipEmptyRows. Skipping...", tc.Table.Name)
		return
	}
	// Test that we don't have missing columns and have exactly one entry for each table
	for _, table := range getTablesFromMainTable(tc.Table) {
		if table.IgnoreInTests {
			t.Logf("table %s marked as IgnoreInTest. Skipping...", table.Name)
			continue
		}
		s := sq.StatementBuilder.
			PlaceholderFormat(sq.Dollar).
			Select(fmt.Sprintf("json_agg(%s)", table.Name)).
			From(table.Name)
		query, args, err := s.ToSql()
		if err != nil {
			t.Fatal(err)
		}
		var data []map[string]interface{}
		if err := pgxscan.Get(context.Background(), conn, &data, query, args...); err != nil {
			t.Fatal(err)
		}

		if len(data) == 0 {
			t.Errorf("expected to have at least 1 entry at table %s got zero", table.Name)
			return
		}

		nilColumns := map[string]bool{}
		// mark all columns as nil
		for _, c := range table.Columns {
			if !c.IgnoreInTests {
				nilColumns[c.Name] = true
			}
		}

		for _, row := range data {
			for c, v := range row {
				if v != nil {
					// as long as we had one row or result with this column not nil it means the resolver worked
					nilColumns[c] = false
				}
			}
		}

		var nilColumnsArr []string
		for c, v := range nilColumns {
			if v {
				nilColumnsArr = append(nilColumnsArr, c)
			}
		}

		if len(nilColumnsArr) != 0 {
			b, err := json.MarshalIndent(data, "", "\t")
			if err != nil {
				t.Fatal(err)
			}
			t.Errorf("found nil column in table %s. rows=\n%s\ncolumns=%s\n", table.Name, string(b), strings.Join(nilColumnsArr, ","))
		}
		// if tc.SkipEmptyJsonB {
		// 	continue
		// }

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
	pool       schema.QueryExecer
	dbErr      error
)

func setupDatabase() (schema.QueryExecer, error) {
	dbConnOnce.Do(func() {
		pool, dbErr = database.New(context.Background(), hclog.NewNullLogger(), getEnv("DATABASE_URL", "host=localhost user=postgres password=pass DB.name=postgres port=5432"))
		if dbErr != nil {
			return
		}
	})
	return pool, dbErr
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTablesFromMainTable(table *schema.Table) []*schema.Table {
	var res []*schema.Table
	res = append(res, table)
	for _, t := range table.Relations {
		res = append(res, getTablesFromMainTable(t)...)
	}
	return res
}
