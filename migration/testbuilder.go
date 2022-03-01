package migration

import (
	"context"
	"os"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/database"
	"github.com/cloudquery/cq-provider-sdk/database/dsn"
	"github.com/cloudquery/cq-provider-sdk/migration/migrator"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// queryPrimaryKeys lists all tables in the schema $1 which don't have the column $2 in their primary keys. Ignores *_schema_migrations tables.
const queryPrimaryKeys = `SELECT kcu.table_name, c.constraint_name, ARRAY_AGG(kcu.column_name::text ORDER BY kcu.ordinal_position) AS pk_cols
FROM information_schema.table_constraints c
JOIN information_schema.key_column_usage kcu ON kcu.constraint_name = c.constraint_name AND kcu.constraint_schema = c.constraint_schema AND kcu.constraint_name = c.constraint_name
WHERE kcu.table_schema=$1 AND c.constraint_type = 'PRIMARY KEY' AND kcu.table_name NOT LIKE '%_schema_migrations'
GROUP BY 1,2
HAVING NOT ($2 = ANY(ARRAY_AGG(kcu.column_name::text)))
ORDER BY 1,2;`

// RunMigrationsTest helper tests the migration files of the provider using the database (and dialect) specified in CQ_MIGRATION_TEST_DSN
func RunMigrationsTest(t *testing.T, prov *provider.Provider, additionalVersionsToTest []string) {
	dbDSN := os.Getenv("CQ_MIGRATION_TEST_DSN")
	if dbDSN == "" {
		t.Skip("CQ_MIGRATION_TEST_DSN not set")
		return
	}

	doMigrationsTest(t, context.Background(), dbDSN, prov, additionalVersionsToTest)
}

func RunMigrationsTestWithNewDB(t *testing.T, dbDSN string, newDBName string, prov *provider.Provider, additionalVersionsToTest []string) {
	ctx := context.Background()
	pool, _, err := connect(ctx, dbDSN)
	assert.NoError(t, err)

	_, err = pool.Exec(ctx, "CREATE DATABASE "+newDBName)
	assert.NoError(t, err)
	if t.Failed() {
		t.FailNow()
	}

	defer func() {
		if _, err := pool.Exec(ctx, "DROP DATABASE "+newDBName); err != nil {
			t.Logf("DROP DATABASE failed: %v", err)
		}
	}()

	u, err := dsn.ParseConnectionString(dbDSN)
	assert.NoError(t, err)
	u.Path = "/" + newDBName
	newDSN := u.String()

	doMigrationsTest(t, ctx, newDSN, prov, additionalVersionsToTest)
}

func doMigrationsTest(t *testing.T, ctx context.Context, dsn string, prov *provider.Provider, additionalVersionsToTest []string) {
	var dialect schema.DialectType

	const (
		setupTSDBChildFnMock = `CREATE OR REPLACE FUNCTION setup_tsdb_child(_table_name text, _column_name text, _parent_table_name text, _parent_column_name text)
					RETURNS integer
					LANGUAGE 'plpgsql'
					COST 100
					VOLATILE PARALLEL UNSAFE
				AS $BODY$
				BEGIN
					return 0;
				END;
				$BODY$;`
		setupTSDBParentFnMock = `CREATE OR REPLACE FUNCTION setup_tsdb_parent(_table_name text)
					RETURNS integer
					LANGUAGE 'plpgsql'
					COST 100
					VOLATILE PARALLEL UNSAFE
				AS $BODY$
				DECLARE
					result integer;
				BEGIN
					return 0;
				END;
				$BODY$;`
	)

	pool, _, err := connect(ctx, dsn)
	assert.NoError(t, err)
	defer pool.Close()

	dialect, dsn, err = database.ParseDialectDSN(dsn)
	assert.Nil(t, err)

	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	if dialect == schema.TSDB {
		// mock history functions... in the default schema
		for _, sql := range []string{
			setupTSDBChildFnMock,
			setupTSDBParentFnMock,
		} {
			_, err := conn.Exec(ctx, sql)
			assert.NoError(t, err)
		}
	}
	assert.NoError(t, err)

	migFiles, err := migrator.ReadMigrationFiles(hclog.L(), prov.Migrations)
	assert.NoError(t, err)

	mig, err := migrator.New(hclog.L(), dialect, migFiles, dsn, prov.Name)
	assert.NoError(t, err)
	if t.Failed() {
		t.FailNow()
	}

	defer mig.Close()

	// clean up first... just as a precaution
	assert.NoError(t, mig.DropProvider(ctx, prov.ResourceMap))

	assert.NoError(t, mig.UpgradeProvider(migrator.Latest))

	if dialect == schema.TSDB {
		// while we're at latest, check PK validity: all PKs should contain cq_fetch_date
		t.Run("RequireCQFetchDate", func(t *testing.T) {
			requireAllPKsToHaveColumn(t, ctx, conn, "public", "cq_fetch_date")
		})
	}

	err = mig.DowngradeProvider(migrator.Initial)
	if err == migrate.ErrNoChange {
		err = nil
	}
	assert.NoError(t, err)
	assert.NoError(t, mig.DowngradeProvider(migrator.Down))

	// Run user supplied versions
	for _, v := range additionalVersionsToTest {
		assert.NoError(t, mig.UpgradeProvider(v))
	}

	// Go to latest again and check if we have missing migrations
	{
		if err := mig.UpgradeProvider(migrator.Latest); err != migrate.ErrNoChange {
			assert.NoError(t, err)
		}

		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		dialectType, err := schema.GetDialect(dialect)
		assert.NoError(t, err)

		if err := generateDiffForDialect(ctx, hclog.NewNullLogger(), fs, conn, "public", dialectType, prov, "/", ""); err != errNoChange {
			assert.NoError(t, err)

			mig, err := fs.ReadFile("/up.sql")
			assert.NoError(t, err)
			assert.Empty(t, string(mig), "Found missing migrations")
		}
	}

	assert.NoError(t, mig.DropProvider(ctx, prov.ResourceMap))
}

func requireAllPKsToHaveColumn(t *testing.T, ctx context.Context, conn *pgxpool.Conn, schema, column string) {
	var res []struct {
		TableName string   `db:"table_name"`
		ConstName string   `db:"constraint_name"`
		PKCols    []string `db:"pk_cols"`
	}
	err := pgxscan.Select(ctx, conn, &res, queryPrimaryKeys, schema, column)
	assert.NoError(t, err)
	assert.Empty(t, res)
}
