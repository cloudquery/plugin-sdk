package testing

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/cloudquery/cq-provider-sdk/database"
	"github.com/cloudquery/cq-provider-sdk/migration"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/execution"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/cq-provider-sdk/testlog"
	"github.com/cloudquery/faker/v3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type ResourceTestCase struct {
	Provider *provider.Provider
	Config   string
	// we want it to be parallel by default
	NotParallel bool
	// ParallelFetchingLimit limits parallel resources fetch at a time
	ParallelFetchingLimit uint64
	// SkipIgnoreInTest flag which detects if schema.Table or schema.Column should be ignored
	SkipIgnoreInTest bool
	// Verifiers are map from resource name to its verifiers.
	// If no verifiers specified for resource (resource name is not in key set of map),
	// non emptiness check of all columns in table and its relations will be performed.
	Verifiers map[string][]Verifier
}

// Verifier verifies tables specified by table schema (main table and its relations).
type Verifier func(t *testing.T, table *schema.Table, conn pgxscan.Querier, shouldSkipIgnoreInTest bool)

type testResourceSender struct {
	Errors []string
}

var (
	dbConnOnce sync.Once
	pool       execution.QueryExecer
	dbErr      error
)

func init() {
	_ = faker.SetRandomMapAndSliceMinSize(1)
	_ = faker.SetRandomMapAndSliceMaxSize(1)
}

func TestResource(t *testing.T, resource ResourceTestCase) {
	if !resource.NotParallel {
		t.Parallel()
	}
	t.Helper()

	// No need for configuration or db connection, get it out of the way first
	// testTableIdentifiersForProvider(t, resource.Provider)

	conn, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}

	l := testlog.New(t)
	l.SetLevel(hclog.Info)
	resource.Provider.Logger = l

	// configure provider
	if resp, configErr := resource.Provider.ConfigureProvider(context.Background(), &cqproto.ConfigureProviderRequest{
		CloudQueryVersion: "",
		Connection: cqproto.ConnectionDetails{DSN: getEnv("DATABASE_URL",
			"host=localhost user=postgres password=pass DB.name=postgres port=5432")},
		Config: []byte(resource.Config),
	}); configErr != nil {
		t.Fatal("failed to configure provider", configErr)
	} else if resp != nil && resp.Diagnostics.HasErrors() {
		t.Fatal("errors while configuring provider", configErr)
	}

	for resourceName, table := range resource.Provider.ResourceMap {
		t.Run(resourceName, func(t *testing.T) {
			testErr := testResource(t, resource, resourceName, table, conn)
			if testErr != nil {
				t.Errorf("Error testing %v: %v", table.Name, testErr)
			}
		})
	}
}

func testResource(t *testing.T, resource ResourceTestCase, name string, table *schema.Table, conn execution.QueryExecer) error {
	t.Helper()

	if createErr := dropAndCreateTable(context.Background(), conn, table); createErr != nil {
		assert.FailNow(t, fmt.Sprintf("failed to create table %s", table.Name), createErr)
	}

	if !resource.SkipIgnoreInTest && table.IgnoreInTests {
		t.Logf("skipping fetch of resource: %s in tests", name)
	} else {
		if err := fetchResource(t, &resource, name); err != nil {
			return err
		}
	}

	if verifiers, ok := resource.Verifiers[name]; ok {
		for _, verifier := range verifiers {
			verifier(t, table, conn, resource.SkipIgnoreInTest)
		}
	} else {
		// fallback to default verification
		verifyNoEmptyColumns(t, table, conn, resource.SkipIgnoreInTest)
	}

	return nil
}

// fetchResource - fetches a resource from the cloud and puts them into database. database config can be specified via DATABASE_URL env variable
func fetchResource(t *testing.T, resource *ResourceTestCase, resourceName string) error {
	t.Helper()

	t.Logf("fetch resource %v", resourceName)

	var resourceSender = &testResourceSender{
		Errors: []string{},
	}

	if err := resource.Provider.FetchResources(context.Background(),
		&cqproto.FetchResourcesRequest{
			Resources:             []string{resourceName},
			ParallelFetchingLimit: resource.ParallelFetchingLimit,
		},
		resourceSender,
	); err != nil {
		return err
	}

	if len(resourceSender.Errors) > 0 {
		return fmt.Errorf("error/s occurred during test, %s", strings.Join(resourceSender.Errors, ", "))
	}

	return nil
}

func verifyNoEmptyColumns(t *testing.T, table *schema.Table, conn pgxscan.Querier, shouldSkipIgnoreInTest bool) {
	t.Helper()
	t.Run(table.Name, func(t *testing.T) {
		t.Helper()

		if !shouldSkipIgnoreInTest && table.IgnoreInTests {
			t.Skipf("table %s marked as IgnoreInTest. Skipping...", table.Name)
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
			if shouldSkipIgnoreInTest || !c.IgnoreInTests {
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
			t.Errorf("found nil column in table %s. columns=%s", table.Name, strings.Join(nilColumnsArr, ","))
		}
		for _, childTable := range table.Relations {
			verifyNoEmptyColumns(t, childTable, conn, shouldSkipIgnoreInTest)
		}
	})
}

func dropAndCreateTable(ctx context.Context, conn execution.QueryExecer, table *schema.Table) error {
	ups, err := migration.CreateTableDefinitions(ctx, schema.PostgresDialect{}, table, nil)
	if err != nil {
		return err
	}

	if err := dropTables(ctx, conn, table); err != nil {
		return err
	}

	for _, sql := range ups {
		if err := conn.Exec(ctx, sql); err != nil {
			return err
		}
	}

	return nil
}

func dropTables(ctx context.Context, db execution.QueryExecer, table *schema.Table) error {
	if err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", strconv.Quote(table.Name))); err != nil {
		return err
	}
	for _, rel := range table.Relations {
		if err := dropTables(ctx, db, rel); err != nil {
			return err
		}
	}
	return nil
}

func (f *testResourceSender) Send(r *cqproto.FetchResourcesResponse) error {
	if r.Error != "" {
		fmt.Printf(r.Error)
		f.Errors = append(f.Errors, r.Error)
	}
	for _, d := range r.Summary.Diagnostics {
		if d.Severity() != diag.IGNORE {
			f.Errors = append(f.Errors, fmt.Sprintf("resource: %s. summary: %s, details %s", d.Description().Resource, d.Description().Summary, d.Description().Detail))
		}
	}
	return nil
}

func setupDatabase() (execution.QueryExecer, error) {
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
