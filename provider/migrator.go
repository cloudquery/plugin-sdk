package provider

import (
	"context"
	"reflect"
	"runtime"
	"strings"

	"github.com/georgysavva/scany/pgxscan"

	"github.com/cloudquery/go-funk"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/huandu/go-sqlbuilder"
)

const queryTableColumns = `SELECT array_agg(column_name::text) as columns FROM information_schema.columns WHERE table_name = $1`

const addColumnToTable = `ALTER TABLE %s ADD COLUMN %v %v;`

// Migrator handles creation of schema.Table in database if they don't exist
type Migrator struct {
	db  schema.Database
	log hclog.Logger
}

func NewMigrator(db schema.Database, log hclog.Logger) Migrator {
	return Migrator{db, log}
}

func (m Migrator) upgradeTable(ctx context.Context, t *schema.Table) error {

	rows, err := m.db.Query(ctx, queryTableColumns, t.Name)
	if err != nil {
		return err
	}

	var existingColumns struct {
		Columns []string
	}

	if err := pgxscan.ScanOne(&existingColumns, rows); err != nil {
		return err
	}

	columnsToAdd, _ := funk.DifferenceString(t.ColumnNames(), existingColumns.Columns)
	for _, d := range columnsToAdd {
		m.log.Debug("adding column", "column", d)

		col := t.Column(d)
		sql, _ := sqlbuilder.Buildf(addColumnToTable, sqlbuilder.Raw(t.Name), sqlbuilder.Raw(d), sqlbuilder.Raw(schema.GetPgTypeFromType(col.Type))).BuildWithFlavor(sqlbuilder.PostgreSQL)
		if err := m.db.Exec(ctx, sql); err != nil {
			return err
		}
	}
	return nil

}

func (m Migrator) CreateTable(ctx context.Context, t *schema.Table, parent *schema.Table) error {
	// Build a SQL to create a table.
	ctb := sqlbuilder.CreateTable(t.Name).IfNotExists()
	ctb.Define("id", "uuid", "NOT NULL", "PRIMARY KEY")
	m.buildColumns(ctb, t.Columns, parent)
	sql, _ := ctb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	m.log.Debug("creating table if not exists", "table", t.Name)
	if err := m.db.Exec(ctx, sql); err != nil {
		return err
	}

	m.log.Debug("migrating table columns if required", "table", t.Name)
	if err := m.upgradeTable(ctx, t); err != nil {
		return err
	}

	if t.Relations == nil {
		return nil
	}

	m.log.Debug("creating table relations", "table", t.Name)
	// Create relation tables
	for _, r := range t.Relations {
		m.log.Debug("creating table relation", "table", r.Name)
		if err := m.CreateTable(ctx, r, t); err != nil {
			return err
		}
	}
	return nil
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (m Migrator) buildColumns(ctb *sqlbuilder.CreateTableBuilder, cc []schema.Column, parent *schema.Table) {
	for _, c := range cc {
		defs := []string{c.Name, schema.GetPgTypeFromType(c.Type)}
		// TODO: This is a bit ugly. Think of a better way
		if strings.HasSuffix(GetFunctionName(c.Resolver), "ParentIdResolver") {
			defs = append(defs, "REFERENCES", parent.Name, "ON DELETE CASCADE")
		}

		ctb.Define(defs...)
	}
}
