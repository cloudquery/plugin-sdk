package migration

import (
	"context"
	"os"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/database"
	"github.com/cloudquery/cq-provider-sdk/migration/migrator"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// RunMigrationsTest helper tests the migration files of the provider using the database (and dialect) specified in CQ_MIGRATION_TEST_DSN
func RunMigrationsTest(t *testing.T, prov *provider.Provider, additionalVersionsToTest []string) {
	dsn := os.Getenv("CQ_MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CQ_MIGRATION_TEST_DSN not set")
		return
	}

	doMigrationsTest(t, context.Background(), dsn, prov, additionalVersionsToTest)
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

	mig, err := migrator.New(hclog.L(), dialect, migFiles, dsn, prov.Name, nil)
	assert.NoError(t, err)

	// clean up first... just as a precaution
	assert.NoError(t, mig.DropProvider(ctx, prov.ResourceMap))

	assert.NoError(t, mig.UpgradeProvider(migrator.Latest))
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
