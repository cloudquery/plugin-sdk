package migrator

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/database/dsn"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

const (
	defaultQuery = "select 1;"
	emptyQuery   = ""
)

var (
	simpleMigrations = map[string]map[string][]byte{
		"postgres": {
			"1_v0.0.1.up.sql":        []byte(defaultQuery),
			"1_v0.0.1.down.sql":      []byte(defaultQuery),
			"3_v0.0.2.up.sql":        []byte(defaultQuery),
			"3_v0.0.2.down.sql":      []byte(defaultQuery),
			"2_v0.0.2-beta.up.sql":   []byte(defaultQuery),
			"2_v0.0.2-beta.down.sql": []byte(defaultQuery),
			"4_v0.0.3.up.sql":        []byte(defaultQuery),
			"4_v0.0.3.down.sql":      []byte(defaultQuery),
			"5_v0.0.4.up.sql":        []byte(emptyQuery),
			"5_v0.0.4.down.sql":      []byte(defaultQuery),
		},
	}

	complexMigrations = map[string]map[string][]byte{
		"postgres": {
			"1_v0.0.2.up.sql":        []byte(defaultQuery),
			"1_v0.0.2.down.sql":      []byte(defaultQuery),
			"2_v0.0.3-beta.up.sql":   []byte(defaultQuery),
			"2_v0.0.3-beta.down.sql": []byte(defaultQuery),
			"3_v0.0.3.up.sql":        []byte(defaultQuery),
			"3_v0.0.3.down.sql":      []byte(defaultQuery),
			"4_v0.0.6.up.sql":        []byte(defaultQuery),
			"4_v0.0.6.down.sql":      []byte(defaultQuery),
			"5_v0.1.4.up.sql":        []byte(emptyQuery),
			"5_v0.1.4.down.sql":      []byte(defaultQuery),
		},
	}
)

func getDBUrl() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable"
}

func TestMigrations(t *testing.T) {
	m, err := New(hclog.Default(), schema.Postgres, simpleMigrations, getDBUrl(), "test", nil)
	assert.Nil(t, err)

	err = m.DropProvider(context.Background(), nil)
	assert.Nil(t, err)

	err = m.UpgradeProvider(Latest)
	assert.Nil(t, err)

	err = m.UpgradeProvider(Latest)
	assert.Equal(t, err, migrate.ErrNoChange)

	err = m.DowngradeProvider("v0.0.2-beta")
	assert.Nil(t, err)

	err = m.UpgradeProvider("v0.0.3")
	assert.Nil(t, err)

	version, dirty, err := m.Version()
	assert.Equal(t, []interface{}{"v0.0.3", false, nil}, []interface{}{version, dirty, err})

	err = m.UpgradeProvider(Latest)
	assert.Nil(t, err)

	version, dirty, err = m.Version()
	assert.Equal(t, []interface{}{"v0.0.4", false, nil}, []interface{}{version, dirty, err})

	err = m.UpgradeProvider("v0.0.4")
	assert.Equal(t, err, migrate.ErrNoChange)

	version, dirty, err = m.Version()
	assert.Equal(t, []interface{}{"v0.0.4", false, nil}, []interface{}{version, dirty, err})
}

// TestMigrationJumps tests an edge case we request a higher version but latest migration is a previous version
func TestMigrationJumps(t *testing.T) {
	m, err := New(hclog.Default(), schema.Postgres, complexMigrations, getDBUrl(), "test", nil)
	assert.Nil(t, err)

	err = m.DropProvider(context.Background(), nil)
	assert.Nil(t, err)

	err = m.UpgradeProvider("v0.2.0")
	assert.Nil(t, err)

	version, dirty, err := m.Version()
	assert.Equal(t, []interface{}{"v0.1.4", false, nil}, []interface{}{version, dirty, err})
}

func TestMultiProviderMigrations(t *testing.T) {
	mtest, err := New(hclog.Default(), schema.Postgres, simpleMigrations, getDBUrl(), "test", nil)
	assert.Nil(t, err)

	mtest2, err := New(hclog.Default(), schema.Postgres, simpleMigrations, getDBUrl(), "test2", nil)
	assert.Nil(t, err)

	err = mtest.DropProvider(context.Background(), nil)
	assert.Nil(t, err)
	err = mtest2.DropProvider(context.Background(), nil)
	assert.Nil(t, err)

	err = mtest.UpgradeProvider(Latest)
	assert.Nil(t, err)
	err = mtest.UpgradeProvider(Latest)
	assert.Equal(t, err, migrate.ErrNoChange)
	version, dirty, err := mtest.Version()
	assert.Equal(t, []interface{}{"v0.0.4", false, nil}, []interface{}{version, dirty, err})

	version, dirty, err = mtest2.Version()
	assert.Equal(t, []interface{}{"v0.0.0", false, migrate.ErrNilVersion}, []interface{}{version, dirty, err})
	err = mtest2.UpgradeProvider("v0.0.3")
	assert.Nil(t, err)
	version, dirty, err = mtest2.Version()
	assert.Equal(t, []interface{}{"v0.0.3", false, nil}, []interface{}{version, dirty, err})

	err = mtest.DropProvider(context.Background(), nil)
	assert.Nil(t, err)

	version, dirty, err = mtest2.Version()
	assert.Equal(t, []interface{}{"v0.0.3", false, nil}, []interface{}{version, dirty, err})
	version, dirty, err = mtest.Version()
	assert.Equal(t, []interface{}{"v0.0.0", false, migrate.ErrNilVersion}, []interface{}{version, dirty, err})
}

func TestFindLatestMigration(t *testing.T) {
	mtest, err := New(hclog.Default(), schema.Postgres, complexMigrations, getDBUrl(), "test", nil)
	assert.Nil(t, err)
	mv, err := mtest.FindLatestMigration("v0.0.3")
	assert.Nil(t, err)
	assert.Equal(t, uint(3), mv)

	mv, err = mtest.FindLatestMigration("v0.0.3-alpha")
	assert.Nil(t, err)
	assert.Equal(t, uint(1), mv)

	mv, err = mtest.FindLatestMigration("v0.1.3-alpha")
	assert.Nil(t, err)
	assert.Equal(t, uint(4), mv)

	mv, err = mtest.FindLatestMigration("v0.1.5")
	assert.Nil(t, err)
	assert.Equal(t, uint(5), mv)

	mv, err = mtest.FindLatestMigration("v0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, uint(1), mv)

	mv, err = mtest.FindLatestMigration(Latest)
	assert.Nil(t, err)
	assert.Equal(t, uint(5), mv)
}

func TestNoSchemaError(t *testing.T) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, getDBUrl())
	assert.NoError(t, err)
	defer conn.Close(ctx)

	const newDBName = "testschemadb"

	if _, err := conn.Exec(ctx, "DROP DATABASE IF EXISTS "+newDBName); err != nil {
		t.Logf("DROP DATABASE failed: %v", err)
	}

	_, err = conn.Exec(ctx, "CREATE DATABASE "+newDBName)
	assert.NoError(t, err)
	if t.Failed() {
		t.FailNow()
	}

	defer func() {
		if _, err := conn.Exec(ctx, "DROP DATABASE "+newDBName+" WITH(FORCE)"); err != nil {
			t.Logf("DROP DATABASE failed: %v", err)
		}
	}()

	u, err := dsn.ParseConnectionString(getDBUrl())
	assert.NoError(t, err)
	u.Path = "/" + newDBName
	newDSN := u.String()

	newConn, err := pgx.Connect(ctx, newDSN)
	assert.NoError(t, err)
	defer newConn.Close(ctx)

	for _, q := range []string{
		"CREATE USER weakuser WITH PASSWORD 'weak'",
		"REVOKE ALL ON SCHEMA public FROM PUBLIC",
	} {
		_, err = newConn.Exec(ctx, q)
		assert.NoError(t, err)
		if t.Failed() {
			t.FailNow()
		}
	}
	defer func() {
		if _, err := newConn.Exec(ctx, "DROP USER weakuser"); err != nil {
			t.Logf("DROP USER failed: %v", err)
		}
	}()

	u.User = url.UserPassword("weakuser", "weak")
	weakDSN := u.String()
	weakConn, err := pgx.Connect(ctx, weakDSN)
	assert.NoError(t, err)
	defer weakConn.Close(ctx)

	var results []struct {
		Name    string `db:"name"`
		Create  bool   `db:"create"`
		Usage   bool   `db:"usage"`
		Current *bool  `db:"current"`
	}
	err = pgxscan.Select(ctx, weakConn, &results, `WITH "names"("name") AS (
  SELECT n.nspname AS "name"
    FROM pg_catalog.pg_namespace n
      WHERE n.nspname !~ '^pg_'
        AND n.nspname <> 'information_schema'
) SELECT "name",
  pg_catalog.has_schema_privilege(current_user, "name", 'CREATE') AS "create",
  pg_catalog.has_schema_privilege(current_user, "name", 'USAGE')  AS "usage",
  "name" = pg_catalog.current_schema() AS "current"
    FROM "names"`)
	assert.NoError(t, err)
	for _, row := range results {
		//t.Logf("%s\t%v\t%v\t%v\n", row.Name, row.Create, row.Usage, row.Current)
		if row.Name == "public" {
			assert.Nil(t, row.Current)
			if t.Failed() {
				t.FailNow()
			}
		}
	}

	m, err := New(hclog.Default(), schema.Postgres, simpleMigrations, weakDSN, "test", nil)
	assert.Nil(t, m)
	if t.Failed() {
		m.Close()
	}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `CURRENT_SCHEMA seems empty, possibly due to empty search_path`)
}
