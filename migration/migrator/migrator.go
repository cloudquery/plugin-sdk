package migrator

import (
	"context"
	"embed"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/database/dsn"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"

	"github.com/golang-migrate/migrate/v4"
	mpg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/xo/dburl"
)

const (
	Latest  = "latest"
	Initial = "initial"
	Down    = "down_testing" // used in testing

	migrationsEmbeddedDirectoryPath = "migrations"
	dropTableSQL                    = "DROP TABLE IF EXISTS %s CASCADE"
)

// ReadMigrationFiles reads the given embed.FS for the migration files and returns a map of dialect directories vs. filenames vs. data
func ReadMigrationFiles(log hclog.Logger, migrationFiles embed.FS) (map[string]map[string][]byte, error) {
	dirs, err := migrationFiles.ReadDir(migrationsEmbeddedDirectoryPath)
	if err != nil {
		log.Info("Provider doesn't define any migration files")
		return nil, nil
	}

	migrations := make(map[string]map[string][]byte)

	for _, d := range dirs {
		if !d.IsDir() {
			return nil, fmt.Errorf("bad migrations structure: missing dialect directories")
		}

		dialectMigrations := make(map[string][]byte)

		basePath := path.Join(migrationsEmbeddedDirectoryPath, d.Name())
		files, err := migrationFiles.ReadDir(basePath)
		if err != nil {
			log.Info("Provider doesn't define any migration files for dialect")
			continue
		}
		for _, m := range files {
			f, err := migrationFiles.Open(path.Join(basePath, m.Name()))
			if err != nil {
				return nil, err
			}

			info, _ := m.Info()
			if info.Size() == 0 {
				dialectMigrations[m.Name()] = []byte("")
				continue
			}
			data := make([]byte, info.Size())
			if _, err := f.Read(data); err != nil {
				return nil, err
			}
			dialectMigrations[m.Name()] = data
		}

		migrations[d.Name()] = dialectMigrations
	}
	return migrations, nil
}

type Migrator struct {
	provider    string
	dsn         string
	migratorUrl *dburl.URL
	log         hclog.Logger
	m           *migrate.Migrate
	driver      source.Driver
	// maps between semantic version to the timestamp it was created at
	versionMapper map[string]uint
	versions      version.Collection

	postHook func(context.Context) error
}

func New(log hclog.Logger, dt schema.DialectType, migrationFiles map[string]map[string][]byte, dsnURI, providerName string, postHook func(context.Context) error) (*Migrator, error) {
	versionMapper := make(map[string]uint)
	versions := make(version.Collection, 0)
	mm := afero.NewMemMapFs()
	_ = mm.Mkdir("migrations", 0755)

	for k, data := range migrationFiles[dt.MigrationDirectory()] {
		log.Debug("adding migration file", "file", k)
		if err := afero.WriteFile(mm, path.Join(migrationsEmbeddedDirectoryPath, k), data, 0644); err != nil {
			return nil, err
		}
		raw := strings.Split(strings.TrimSuffix(strings.TrimSuffix(k, ".up.sql"), ".down.sql"), "_")
		if len(raw) == 1 {
			return nil, fmt.Errorf("invalid migration filename %q: should be in format <int>_v<version>.up|down.sql", k)
		}

		// add version once to mapper, up/down should have same migration number anyway
		if _, ok := versionMapper[raw[1]]; !ok {
			versionMapper[raw[1]] = cast.ToUint(raw[0])
			v, err := version.NewVersion(raw[1])
			if err != nil {
				return nil, err
			}
			versions = append(versions, v)
		}
	}
	sort.Sort(versions)
	driver, err := iofs.New(afero.NewIOFS(mm), migrationsEmbeddedDirectoryPath)
	if err != nil {
		return nil, err
	}
	u, err := dsn.ParseConnectionString(dsnURI)
	if err != nil {
		return nil, err
	}
	if u.RawQuery != "" {
		u.RawQuery += fmt.Sprintf("&x-migrations-table=%s_schema_migrations", providerName)
	} else {
		u.RawQuery += fmt.Sprintf("x-migrations-table=%s_schema_migrations", providerName)
	}
	m, err := migrate.NewWithSourceInstance(providerName, driver, u.String())
	if err != nil {
		return nil, convertMigrateError(u.String(), err)
	}
	return &Migrator{
		log:           log,
		provider:      providerName,
		dsn:           dsnURI,
		migratorUrl:   u,
		m:             m,
		driver:        driver,
		versionMapper: versionMapper,
		versions:      versions,
		postHook:      postHook,
	}, nil
}

func (m *Migrator) callPostHook(ctx context.Context) error {
	if m.postHook == nil {
		return nil
	}
	return m.postHook(ctx)
}

func (m *Migrator) Close() error {
	if m.m == nil {
		return nil
	}

	_, dbErr := m.m.Close()
	return dbErr
}

func (m *Migrator) UpgradeProvider(version string) (retErr error) {
	defer func() {
		if retErr != nil {
			return
		}
		retErr = m.callPostHook(context.Background())
	}()

	if version == Latest {
		return m.m.Up()
	}

	mv, err := m.FindLatestMigration(version)
	if err != nil {
		return fmt.Errorf("version %s upgrade doesn't exist", version)
	}
	m.log.Debug("upgrading provider version", "version", version, "migrator_version", mv)
	return m.m.Migrate(mv)
}

func (m *Migrator) DowngradeProvider(version string) (retErr error) {
	defer func() {
		if retErr != nil {
			return
		}
		retErr = m.callPostHook(context.Background())
	}()

	if version == Down { // Used in testing
		return m.m.Down()
	}

	mv, err := m.FindLatestMigration(version)
	if err != nil {
		return fmt.Errorf("version %s upgrade doesn't exist", version)
	}
	m.log.Debug("downgrading provider version", "version", version, "migrator_version", mv)

	return m.m.Migrate(mv)
}

func (m *Migrator) DropProvider(ctx context.Context, schema map[string]*schema.Table) (retErr error) {
	defer func() {
		if retErr != nil {
			return
		}
		retErr = m.callPostHook(context.Background())
	}()

	// we don't use go-migrate's drop since its too violent and it will remove all tables of other providers,
	// instead we will only drop the migration table and all schema's tables
	// we additionally don't use a transaction since this results quite often in out of shared memory errors
	conn, err := pgx.Connect(ctx, m.dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	q := fmt.Sprintf(dropTableSQL, strconv.Quote(fmt.Sprintf("%s_schema_migrations", m.provider)))
	if _, err := conn.Exec(ctx, q); err != nil {
		return err
	}
	for name, table := range schema {
		m.log.Debug("deleting table and all relations", "table", name, "provider", m.provider)
		if err := dropTables(ctx, conn, table); err != nil {
			return err
		}
	}

	if _, dbErr := m.m.Close(); dbErr != nil {
		m.log.Warn("error closing migrator", "error", dbErr)
	}

	newM, err := migrate.NewWithSourceInstance(m.provider, m.driver, m.migratorUrl.String())
	if err != nil {
		return convertMigrateError(m.migratorUrl.String(), err)
	}
	// reset migrator
	m.m = newM
	return nil
}

func (m *Migrator) Version() (string, bool, error) {
	ver, dirty, err := m.m.Version()
	for k, v := range m.versionMapper {
		if ver == v {
			return k, dirty, err
		}
	}
	return "v0.0.0", dirty, err
}

func (m *Migrator) SetVersion(requestedVersion string) (retErr error) {
	defer func() {
		if retErr != nil {
			return
		}
		retErr = m.callPostHook(context.Background())
	}()

	mv, err := m.FindLatestMigration(requestedVersion)
	if err != nil {
		return err
	}
	return m.m.Force(int(mv))
}

// FindLatestMigration finds closet migration to the requested version
//  For example we have the following migrations:
//  1_001, 2_005, 3_009
// if we ask for 007 we get 005
// if we ask for 004 we get 001
// if we ask for 005 we get 005
func (m *Migrator) FindLatestMigration(requestedVersion string) (uint, error) {
	if requestedVersion == Latest {
		mv := m.versionMapper[m.versions[len(m.versions)-1].Original()]
		return mv, nil
	} else if requestedVersion == Initial {
		mv := m.versionMapper[m.versions[0].Original()]
		return mv, nil
	}
	// if we have a migration for specific version return that mv number
	mv, ok := m.versionMapper[requestedVersion]
	if ok {
		return mv, nil
	}
	ov, err := version.NewVersion(requestedVersion)
	if err != nil {
		return 0, fmt.Errorf("version %s doesn't exist", requestedVersion)
	}
	// find closest migration level
	for i, v := range m.versions {
		if v.GreaterThan(ov) {
			// edge case that requested version is smaller than all migrations
			if i == 0 {
				return m.versionMapper[m.versions[0].Original()], nil
			}
			mv = m.versionMapper[m.versions[i-1].Original()]
			return mv, nil
		}
	}
	mv = m.versionMapper[m.versions[len(m.versions)-1].Original()]
	return mv, nil
}

func dropTables(ctx context.Context, conn *pgx.Conn, table *schema.Table) error {
	if _, err := conn.Exec(ctx, fmt.Sprintf(dropTableSQL, strconv.Quote(table.Name))); err != nil {
		return err
	}
	for _, rel := range table.Relations {
		if err := dropTables(ctx, conn, rel); err != nil {
			return err
		}
	}
	return nil
}

func convertMigrateError(dsnURI string, err error) error {
	if err == nil {
		return err
	}

	// https://github.com/golang-migrate/migrate/issues/696
	if err != mpg.ErrNoSchema && !strings.Contains(err.Error(), `"current_schema": converting NULL to string`) {
		return err
	}

	const errFmt = "CURRENT_SCHEMA seems empty, possibly due to empty search_path. Try `GRANT ALL PRIVILEGES ON %s TO <user>`"

	u, err2 := dsn.ParseConnectionString(dsnURI)
	if err2 != nil {
		return fmt.Errorf(errFmt, `<schema>`)
	}
	p := u.Query().Get("search_path")
	if p == "" {
		p = "public"
	}

	return fmt.Errorf(errFmt, p)
}
