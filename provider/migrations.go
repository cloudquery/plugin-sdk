package provider

import (
	"context"
	"embed"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/xo/dburl"
)

const (
	migrationsEmbeddedDirectoryPath = "migrations"
	dropTableSQL                    = "DROP TABLE IF EXISTS %s CASCADE"
)

func readProviderMigrationFiles(log hclog.Logger, migrationFiles embed.FS) (map[string][]byte, error) {
	var (
		err        error
		migrations = make(map[string][]byte)
	)
	files, err := migrationFiles.ReadDir(migrationsEmbeddedDirectoryPath)
	if err != nil {
		log.Info("Provider doesn't define any migration files")
		return migrations, nil
	}
	for _, m := range files {
		f, err := migrationFiles.Open(path.Join(migrationsEmbeddedDirectoryPath, m.Name()))
		if err != nil {
			return nil, err
		}
		info, _ := m.Info()
		if info.Size() == 0 {
			migrations[m.Name()] = []byte("")
			continue
		}
		data := make([]byte, info.Size())
		if _, err := f.Read(data); err != nil {
			return nil, err
		}
		migrations[m.Name()] = data
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
}

func NewMigrator(log hclog.Logger, migrationFiles map[string][]byte, dsn string, providerName string) (*Migrator, error) {
	versionMapper := make(map[string]uint)
	versions := make(version.Collection, 0)
	mm := afero.NewMemMapFs()
	_ = mm.Mkdir("migrations", 0755)
	for k, data := range migrationFiles {
		log.Debug("adding migration file", "file", k)
		if err := afero.WriteFile(mm, path.Join(migrationsEmbeddedDirectoryPath, k), data, 0644); err != nil {
			return nil, err
		}
		raw := strings.Split(strings.TrimSuffix(strings.TrimSuffix(k, ".up.sql"), ".down.sql"), "_")
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
	u, err := dburl.Parse(dsn)
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
		return nil, err
	}
	return &Migrator{
		log:           log,
		provider:      providerName,
		dsn:           dsn,
		migratorUrl:   u,
		m:             m,
		driver:        driver,
		versionMapper: versionMapper,
		versions:      versions,
	}, nil
}

func (m *Migrator) UpgradeProvider(version string) error {
	if version == "latest" {
		return m.m.Up()
	}
	mv, ok := m.versionMapper[version]
	if !ok {
		return fmt.Errorf("version %s upgrade doesn't exist", version)
	}
	m.log.Debug("upgrading provider version", "version", version, "migrator_version", mv)
	return m.m.Migrate(mv)
}

func (m *Migrator) DowngradeProvider(version string) error {
	mv, ok := m.versionMapper[version]
	if !ok {
		return fmt.Errorf("version %s downgrade doesn't exist", version)
	}
	m.log.Debug("downgrading provider version", "version", version, "migrator_version", mv)
	return m.m.Migrate(mv)
}

func (m *Migrator) DropProvider(ctx context.Context, schema map[string]*schema.Table) error {
	// we don't use go-migrate's drop since its too violent and it will remove all tables of other providers,
	// instead we will only drop the migration table and all schema's tables
	conn, err := pgx.Connect(ctx, m.dsn)
	if err != nil {
		return err
	}
	err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		q := fmt.Sprintf(dropTableSQL, strconv.Quote(fmt.Sprintf("%s_schema_migrations", m.provider)))
		if _, err := tx.Exec(ctx, q); err != nil {
			return err
		}
		for name, table := range schema {
			m.log.Debug("deleting table and all relations", "table", name, "provider", m.provider)
			if err := dropTables(ctx, tx, table); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	newM, err := migrate.NewWithSourceInstance(m.provider, m.driver, m.migratorUrl.String())
	if err != nil {
		return err
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

func (m *Migrator) SetVersion(requestedVersion string) error {
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
	if requestedVersion == "latest" {
		mv := m.versionMapper[m.versions[len(m.versions)-1].Original()]
		return mv, nil
	}
	// if we have a migration for specific version return that mv number
	mv, ok := m.versionMapper[requestedVersion]
	if ok {
		return mv, nil
	}
	ov, err := version.NewVersion(requestedVersion)
	if err != nil {
		return 0, err
	}
	// find closest migration level
	for i, v := range m.versions {
		if v.GreaterThan(ov) {
			// edge case that requested version is smaller than ll migrations
			if i == 0 {
				return 0, nil
			}
			mv = m.versionMapper[m.versions[i-1].Original()]
			return mv, nil
		}
	}
	mv = m.versionMapper[m.versions[len(m.versions)-1].Original()]
	return mv, nil
}

func dropTables(ctx context.Context, tx pgx.Tx, table *schema.Table) error {
	if _, err := tx.Exec(ctx, fmt.Sprintf(dropTableSQL, strconv.Quote(table.Name))); err != nil {
		return err
	}
	for _, rel := range table.Relations {
		if err := dropTables(ctx, tx, rel); err != nil {
			return err
		}
	}
	return nil
}
