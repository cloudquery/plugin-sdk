package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/huandu/go-sqlbuilder"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thoas/go-funk"
)

const (
	queryTableColumns = `SELECT array_agg(column_name::text) as columns FROM information_schema.columns WHERE table_name = $1`
	addColumnToTable  = `ALTER TABLE %s ADD COLUMN IF NOT EXISTS %v %v;`
)

// TableCreator handles creation of schema.Table in database if they don't exist, and migration of tables if provider was upgraded.
type TableCreator struct {
	log hclog.Logger
}

func NewTableCreator(log hclog.Logger) *TableCreator {
	return &TableCreator{
		log,
	}
}

func (m TableCreator) CreateTable(ctx context.Context, conn *pgxpool.Conn, t *schema.Table, parent *schema.Table) error {
	// Build a SQL to create a table.
	ctb := sqlbuilder.CreateTable(t.Name).IfNotExists()
	for _, c := range schema.GetDefaultSDKColumns() {
		if c.CreationOptions.Unique {
			ctb.Define(c.Name, schema.GetPgTypeFromType(c.Type), "unique")
		} else {
			ctb.Define(c.Name, schema.GetPgTypeFromType(c.Type))
		}

	}

	m.buildColumns(ctb, t.Columns, parent)
	ctb.Define(fmt.Sprintf("constraint %s_pk primary key(%s)", schema.TruncateTableConstraint(t.Name), strings.Join(t.PrimaryKeys(), ",")))
	sql, _ := ctb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	m.log.Debug("creating table if not exists", "table", t.Name)
	if _, err := conn.Exec(ctx, sql); err != nil {
		return err
	}

	m.log.Debug("migrating table columns if required", "table", t.Name)
	if err := m.upgradeTable(ctx, conn, t); err != nil {
		return err
	}

	if t.Relations == nil {
		return nil
	}

	m.log.Debug("creating table relations", "table", t.Name)
	// Create relation tables
	for _, r := range t.Relations {
		m.log.Debug("creating table relation", "table", r.Name)
		if err := m.CreateTable(ctx, conn, r, t); err != nil {
			return err
		}
	}
	return nil
}

func (m TableCreator) upgradeTable(ctx context.Context, conn *pgxpool.Conn, t *schema.Table) error {
	rows, err := conn.Query(ctx, queryTableColumns, t.Name)
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
		if col == nil {
			m.log.Warn("column missing from table, not adding it", "table", t.Name, "column", d)
			continue
		}
		sql, _ := sqlbuilder.Buildf(addColumnToTable, sqlbuilder.Raw(t.Name), sqlbuilder.Raw(d), sqlbuilder.Raw(schema.GetPgTypeFromType(col.Type))).BuildWithFlavor(sqlbuilder.PostgreSQL)
		if _, err := conn.Exec(ctx, sql); err != nil {
			return err
		}
	}
	return nil

}

func (m TableCreator) buildColumns(ctb *sqlbuilder.CreateTableBuilder, cc []schema.Column, parent *schema.Table) {
	for _, c := range cc {
		defs := []string{strconv.Quote(c.Name), schema.GetPgTypeFromType(c.Type)}
		if c.CreationOptions.Unique {
			defs = []string{strconv.Quote(c.Name), schema.GetPgTypeFromType(c.Type), "unique"}
		}
		if strings.HasSuffix(c.Name, "cq_id") && c.Name != "cq_id" {
			defs = append(defs, "REFERENCES", fmt.Sprintf("%s(cq_id)", parent.Name), "ON DELETE CASCADE")
		}
		ctb.Define(defs...)
	}
}
