package migration

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/migration/longestcommon"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thoas/go-funk"
)

const (
	queryTableColumns   = `SELECT ARRAY_AGG(column_name::text) AS columns, ARRAY_AGG(udt_name::regtype::text) AS types FROM information_schema.columns WHERE table_schema=$1 AND table_name=$2`
	addColumnToTable    = `ALTER TABLE IF EXISTS %s ADD COLUMN IF NOT EXISTS %v %v;`
	dropColumnFromTable = `ALTER TABLE IF EXISTS %s DROP COLUMN IF EXISTS %v;`
	renameColumnInTable = `-- ALTER TABLE %s RENAME COLUMN %v TO %v; -- uncomment to activate, remove ADD/DROP COLUMN above and below` // Can't have IF EXISTS here

	queryTablePKs           = `SELECT c.constraint_name, ARRAY_AGG(k.column_name::text ORDER BY k.ordinal_position) AS columns FROM information_schema.table_constraints c JOIN information_schema.key_column_usage k ON c.table_catalog=k.table_catalog AND c.table_schema=k.table_schema AND c.table_name=k.table_name AND c.constraint_name=k.constraint_name WHERE c.table_schema=$1 AND c.table_name=$2 AND c.constraint_type='PRIMARY KEY' GROUP BY 1`
	addPKToTable            = `ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s PRIMARY KEY (%s);`
	dropConstraintFromTable = `ALTER TABLE IF EXISTS %s DROP CONSTRAINT %s;`

	dropTable = `DROP TABLE IF EXISTS %s;`

	fakeTSDBAssumeColumn = `cq_fetch_date`
)

// TableCreator handles creation of schema.Table in database as SQL strings
type TableCreator struct {
	log     hclog.Logger
	dialect schema.Dialect
}

func NewTableCreator(log hclog.Logger, dialect schema.Dialect) *TableCreator {
	return &TableCreator{
		log:     log,
		dialect: dialect,
	}
}

// CreateTableDefinitions reads schema.Table and builds the CREATE TABLE and DROP TABLE statements for it, also processing and returning subrelation tables
func (m TableCreator) CreateTableDefinitions(ctx context.Context, t *schema.Table, parent *schema.Table) (up, down []string, err error) {
	b := &strings.Builder{}

	// Build a SQL to create a table
	b.WriteString("CREATE TABLE IF NOT EXISTS " + strconv.Quote(t.Name) + " (\n")

	for _, c := range m.dialect.Columns(t) {
		b.WriteByte('\t')
		b.WriteString(strconv.Quote(c.Name) + " " + m.dialect.DBTypeFromType(c.Type))
		if c.CreationOptions.NotNull {
			b.WriteString(" NOT NULL")
		}
		// c.CreationOptions.Unique is handled in the Constraints() call below
		b.WriteString(",\n")
	}

	cons := m.dialect.Constraints(t, parent)
	for i, cn := range cons {
		b.WriteByte('\t')
		b.WriteString(cn)

		if i < len(cons)-1 {
			b.WriteByte(',')
		}

		b.WriteByte('\n')
	}

	b.WriteString(");")

	up, down = make([]string, 0, 1+len(t.Relations)), make([]string, 0, 1+len(t.Relations))
	up = append(up, b.String())
	up = append(up, m.dialect.Extra(t, parent)...)

	// Create relation tables
	for _, r := range t.Relations {
		if cr, dr, err := m.CreateTableDefinitions(ctx, r, t); err != nil {
			return nil, nil, err
		} else {
			up = append(up, cr...)
			down = append(down, dr...)
		}
	}

	down = append(down, fmt.Sprintf(dropTable, t.Name))

	return up, down, nil
}

// DiffTable reads current table info from the given conn for the given table, and returns ALTER TABLE ADD COLUMN statements for the missing columns.
// Newly appearing tables will return a CREATE TABLE statement.
// Column renames are detected (best effort) and ALTER TABLE RENAME COLUMN statements are generated as comments.
// Table renames or removals are not detected.
// FK changes are not detected.
// if fakeTSDB is set, PKs for existing resources are assumed to have fakeTSDBAssumeColumn (`cq_fetch_date`) as the first part of composite PK.
func (m TableCreator) DiffTable(ctx context.Context, conn *pgxpool.Conn, schemaName string, t, parent *schema.Table, fakeTSDB bool) (up, down []string, err error) {
	rows, err := conn.Query(ctx, queryTableColumns, schemaName, t.Name)
	if err != nil {
		return nil, nil, err
	}

	var existingColumns struct {
		Columns []string
		Types   []string
	}

	if err := pgxscan.ScanOne(&existingColumns, rows); err != nil {
		return nil, nil, err
	}

	if len(existingColumns.Columns) == 0 {
		// Table does not exist, CREATE TABLE instead
		u, d, err := m.CreateTableDefinitions(ctx, t, parent)
		if err != nil {
			return nil, nil, fmt.Errorf("CreateTable: %w", err)
		}
		return u, d, nil
	}

	dbColTypes := make(map[string]string, len(existingColumns.Columns))
	for i := range existingColumns.Columns {
		dbColTypes[existingColumns.Columns[i]] = strings.ToLower(existingColumns.Types[i])
	}

	tableColsWithDialect := m.dialect.Columns(t)

	columnsToAdd, columnsToRemove := funk.DifferenceString(tableColsWithDialect.Names(), existingColumns.Columns)
	similars := getSimilars(m.dialect, t, columnsToAdd, columnsToRemove, dbColTypes)

	capSize := len(columnsToAdd) + len(columnsToRemove) // relations not included...
	up, down = make([]string, 0, capSize), make([]string, 0, capSize)
	downLast := make([]string, 0, capSize)

	cUp, cDown, err := m.diffConstraints(ctx, conn, schemaName, t, fakeTSDB)
	if err != nil {
		return nil, nil, fmt.Errorf("diffConstraints failed: %w", err)
	}
	up = append(up, cUp.Squash().removals...)
	down = append(down, cDown.Squash().removals...)

	for _, d := range columnsToAdd {
		m.log.Debug("adding column", "column", d)
		col := tableColsWithDialect.Get(d)
		if col == nil {
			m.log.Warn("column missing from table, not adding it", "table", t.Name, "column", d)
			continue
		}
		if fakeTSDB && col.Name == fakeTSDBAssumeColumn {
			continue
		}
		if col.Internal() {
			m.log.Warn("table missing internal column, not adding it", "table", t.Name, "column", d)
			continue
		}

		var notice string
		if v, ok := similars[d]; ok {
			notice = " -- could this be " + strconv.Quote(v) + " ?"
		}

		up = append(up, fmt.Sprintf(addColumnToTable, strconv.Quote(t.Name), strconv.Quote(d), m.dialect.DBTypeFromType(col.Type))+notice)
		downLast = append(downLast, fmt.Sprintf(dropColumnFromTable, strconv.Quote(t.Name), strconv.Quote(d))+notice)

		if v, ok := similars[d]; ok {
			up = append(up, fmt.Sprintf(renameColumnInTable, strconv.Quote(t.Name), strconv.Quote(v), strconv.Quote(d)))
			downLast = append(downLast, fmt.Sprintf(renameColumnInTable, strconv.Quote(t.Name), strconv.Quote(d), strconv.Quote(v)))
		}
	}

	for _, d := range columnsToRemove {
		m.log.Debug("removing column", "column", d)
		if col := t.Column(d); col != nil {
			if !col.Internal() {
				m.log.Warn("column still in table, not removing it", "table", t.Name, "column", d)
			}
			continue
		}

		var notice string
		if v, ok := similars[d]; ok {
			notice = " -- could this be " + strconv.Quote(v) + " ? Check the RENAME COLUMN statement above"
		}

		up = append(up, fmt.Sprintf(dropColumnFromTable, strconv.Quote(t.Name), strconv.Quote(d))+notice)
		downLast = append(downLast, fmt.Sprintf(addColumnToTable, strconv.Quote(t.Name), strconv.Quote(d), dbColTypes[d])+notice)
	}

	up = append(up, cUp.Squash().additions...)
	down = append(down, cDown.Squash().additions...)

	// Do relation tables
	for _, r := range t.Relations {
		if cr, dr, err := m.DiffTable(ctx, conn, schemaName, r, t, fakeTSDB); err != nil {
			return nil, nil, err
		} else {
			up = append(up, cr...)
			down = append(down, dr...)
		}
	}

	down = append(down, downLast...)

	return up, down, nil
}

func getSimilars(dialect schema.Dialect, t *schema.Table, columnsToAdd, columnsToRemove []string, existingColsTypes map[string]string) map[string]string {
	upColsByType, downColsByType := make(map[string][]string), make(map[string][]string)

	for _, d := range columnsToAdd {
		col := t.Column(d)
		if col == nil {
			continue
		}
		upColsByType[dialect.DBTypeFromType(col.Type)] = append(upColsByType[dialect.DBTypeFromType(col.Type)], d)
	}
	for _, d := range columnsToRemove {
		if col := t.Column(d); col != nil {
			continue
		}
		downColsByType[existingColsTypes[d]] = append(downColsByType[existingColsTypes[d]], d)
	}

	return findSimilarColumnsWithSameType(upColsByType, downColsByType)
}

func findSimilarColumnsWithSameType(setA, setB map[string][]string) map[string]string {
	const threshold = 4 // minimum common prefix/suffix length

	ret := make(map[string]string)

	for typeKey, alist := range setA {
		blist, ok := setB[typeKey]
		if !ok {
			continue
		}

		for _, A := range alist {
			for _, B := range blist {
				if A == B {
					panic("passed equal sets") // should not happen
				}

				pref := longestcommon.Prefix([]string{A, B})
				suf := longestcommon.Suffix([]string{A, B})
				if len(suf) > len(pref) {
					pref = suf
				}
				if len(pref) < threshold {
					continue
				}

				ret[A] = B
				ret[B] = A
			}
		}
	}

	return ret
}

type constraintMigration struct {
	removals  []string
	additions []string
}

type constraintMigrations []constraintMigration

func (c constraintMigrations) Squash() constraintMigration {
	var ret constraintMigration
	for _, m := range c {
		ret.additions = append(ret.additions, m.additions...)
		ret.removals = append(ret.removals, m.removals...)
	}
	return ret
}

func (m TableCreator) diffConstraints(ctx context.Context, conn *pgxpool.Conn, schemaName string, t *schema.Table, fakeTSDB bool) (up, down constraintMigrations, err error) {
	rows, err := conn.Query(ctx, queryTablePKs, schemaName, t.Name)
	if err != nil {
		return nil, nil, err
	}

	var existingPKs []struct {
		ConstraintName string
		Columns        []string
	}

	if err := pgxscan.ScanAll(&existingPKs, rows); err != nil {
		return nil, nil, err
	}

	pks := m.dialect.PrimaryKeys(t)
	if len(pks) == 0 {
		return nil, nil, fmt.Errorf("dialect returned no primary keys for table")
	}

	if l := len(existingPKs); l > 1 {
		return nil, nil, fmt.Errorf("query found more than one PK constraint")
	} else if l == 0 {
		return []constraintMigration{
				{
					additions: []string{
						fmt.Sprintf(addPKToTable, t.Name, t.Name+"_pk", strings.Join(pks, ",")),
					},
				},
			}, []constraintMigration{
				{
					removals: []string{
						fmt.Sprintf(dropConstraintFromTable, t.Name, t.Name+"_pk"),
					},
				},
			}, nil
	}

	if fakeTSDB && existingPKs[0].Columns[0] != fakeTSDBAssumeColumn {
		existingPKs[0].Columns = append([]string{fakeTSDBAssumeColumn}, existingPKs[0].Columns...)
	}

	if reflect.DeepEqual(existingPKs[0].Columns, pks) {
		return nil, nil, nil
	}

	return []constraintMigration{
			{
				removals: []string{
					fmt.Sprintf(dropConstraintFromTable, t.Name, existingPKs[0].ConstraintName),
				},
				additions: []string{
					fmt.Sprintf(addPKToTable, t.Name, t.Name+"_pk", strings.Join(pks, ",")),
				},
			},
		}, []constraintMigration{
			{
				removals: []string{
					fmt.Sprintf(dropConstraintFromTable, t.Name, t.Name+"_pk"),
				},
				additions: []string{
					fmt.Sprintf(addPKToTable, t.Name, existingPKs[0].ConstraintName, strings.Join(existingPKs[0].Columns, ",")),
				},
			},
		}, nil
}
