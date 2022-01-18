package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/modern-go/reflect2"
)

type DialectType string

const (
	Postgres = DialectType("postgres")
	TSDB     = DialectType("timescale")
)

func (t DialectType) MigrationDirectory() string {
	return string(t)
}

type Dialect interface {
	// PrimaryKeys returns the primary keys of table according to dialect
	PrimaryKeys(t *Table) []string

	// Columns returns the columns of table according to dialect
	Columns(t *Table) ColumnList

	// Constraints returns constraint definitions for table, according to dialect
	Constraints(t, parent *Table) []string

	// Extra returns additional definitions for table outside the CREATE TABLE statement, according to dialect
	Extra(t, parent *Table) []string

	// DBTypeFromType returns the database type from the given ValueType. Always lowercase.
	DBTypeFromType(v ValueType) string

	// GetResourceValues will return column values from the resource, ready to go in pgx.CopyFromSlice
	GetResourceValues(r *Resource) ([]interface{}, error)
}

var (
	_ Dialect = (*PostgresDialect)(nil)
	_ Dialect = (*TSDBDialect)(nil)
)

// GetDialect creates and returns a dialect specified by the DialectType
func GetDialect(t DialectType) (Dialect, error) {
	switch t {
	case Postgres:
		return PostgresDialect{}, nil
	case TSDB:
		return TSDBDialect{}, nil
	default:
		return nil, fmt.Errorf("unknown dialect %q", t)
	}
}

type PostgresDialect struct{}

func (d PostgresDialect) PrimaryKeys(t *Table) []string {
	if len(t.Options.PrimaryKeys) > 0 {
		return t.Options.PrimaryKeys
	}
	return []string{cqIdColumn.Name}
}

func (d PostgresDialect) Columns(t *Table) ColumnList {
	return append([]Column{cqIdColumn, cqMeta}, t.Columns...)
}

func (d PostgresDialect) Constraints(t, parent *Table) []string {
	ret := make([]string, 0, len(t.Columns))

	ret = append(ret, fmt.Sprintf("CONSTRAINT %s_pk PRIMARY KEY(%s)", truncatePKConstraint(t.Name), strings.Join(d.PrimaryKeys(t), ",")))

	for _, c := range d.Columns(t) {
		if !c.CreationOptions.Unique {
			continue
		}

		ret = append(ret, fmt.Sprintf("UNIQUE(%s)", c.Name))
	}

	if parent != nil {
		pc := findParentIdColumn(t)
		if pc != nil {
			ret = append(ret, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE", pc.Name, parent.Name, cqIdColumn.Name))
		}
	}

	return ret
}

func (d PostgresDialect) Extra(_, _ *Table) []string {
	return nil
}

func (d PostgresDialect) DBTypeFromType(v ValueType) string {
	switch v {
	case TypeBool:
		return "boolean"
	case TypeInt:
		return "integer"
	case TypeBigInt:
		return "bigint"
	case TypeSmallInt:
		return "smallint"
	case TypeFloat:
		return "float"
	case TypeUUID:
		return "uuid"
	case TypeString:
		return "text"
	case TypeJSON:
		return "jsonb"
	case TypeIntArray:
		return "integer[]"
	case TypeStringArray:
		return "text[]"
	case TypeTimestamp:
		return "timestamp without time zone"
	case TypeByteArray:
		return "bytea"
	case TypeInvalid:
		fallthrough
	case TypeInet:
		return "inet"
	case TypeMacAddr:
		return "mac"
	case TypeInetArray:
		return "inet[]"
	case TypeMacAddrArray:
		return "mac[]"
	case TypeCIDR:
		return "cidr"
	case TypeCIDRArray:
		return "cidr[]"
	default:
		panic("invalid type")
	}
}

func (d PostgresDialect) GetResourceValues(r *Resource) ([]interface{}, error) {
	return doResourceValues(d, r)
}

type TSDBDialect struct {
	pg PostgresDialect
}

func (d TSDBDialect) PrimaryKeys(t *Table) []string {
	return append([]string{cqFetchDateColumn.Name}, d.pg.PrimaryKeys(t)...)
}

func (d TSDBDialect) Columns(t *Table) ColumnList {
	return append([]Column{cqIdColumn, cqMeta, cqFetchDateColumn}, t.Columns...)
}

func (d TSDBDialect) Constraints(t, _ *Table) []string {
	ret := make([]string, 0, len(t.Columns))

	ret = append(ret, fmt.Sprintf("CONSTRAINT %s_pk PRIMARY KEY(%s)", truncatePKConstraint(t.Name), strings.Join(d.PrimaryKeys(t), ",")))

	for _, c := range d.Columns(t) {
		if !c.CreationOptions.Unique {
			continue
		}

		ret = append(ret, fmt.Sprintf("UNIQUE(%s,%s)", cqFetchDateColumn.Name, c.Name))
	}

	return ret
}

func (d TSDBDialect) Extra(t, parent *Table) []string {
	pc := findParentIdColumn(t)

	if parent == nil || pc == nil {
		return []string{
			fmt.Sprintf("SELECT setup_tsdb_parent('%s');", t.Name),
		}
	}

	return []string{
		fmt.Sprintf("CREATE INDEX ON %s (%s, %s);", t.Name, cqFetchDateColumn.Name, pc.Name),
		fmt.Sprintf("SELECT setup_tsdb_child('%s', '%s', '%s', '%s');", t.Name, pc.Name, parent.Name, cqIdColumn.Name),
	}
}

func (d TSDBDialect) DBTypeFromType(v ValueType) string {
	return d.pg.DBTypeFromType(v)
}

func (d TSDBDialect) GetResourceValues(r *Resource) ([]interface{}, error) {
	return doResourceValues(d, r)
}

func doResourceValues(dialect Dialect, r *Resource) ([]interface{}, error) {
	values := make([]interface{}, 0)
	for _, c := range dialect.Columns(r.table) {
		v := r.Get(c.Name)
		if err := c.ValidateType(v); err != nil {
			return nil, err
		}
		if c.Type == TypeJSON {
			if v == nil {
				values = append(values, v)
				continue
			}
			if reflect2.TypeOf(v).Kind() == reflect.Map {
				values = append(values, v)
				continue
			}
			switch data := v.(type) {
			case map[string]interface{}:
				values = append(values, data)
			case string:
				newV := make(map[string]interface{})
				err := json.Unmarshal([]byte(data), &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			case *string:
				var newV interface{}
				err := json.Unmarshal([]byte(*data), &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			case []byte:
				var newV interface{}
				err := json.Unmarshal(data, &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			default:
				d, err := json.Marshal(data)
				if err != nil {
					return nil, err
				}
				var newV interface{}
				err = json.Unmarshal(d, &newV)
				if err != nil {
					return nil, err
				}
				values = append(values, newV)
			}
		} else {
			values = append(values, v)
		}
	}
	return values, nil
}

func findParentIdColumn(t *Table) (ret *Column) {
	for _, c := range t.Columns {
		if c.Meta().Resolver != nil && c.Meta().Resolver.Name == "schema.ParentIdResolver" {
			return &c
		}
	}

	return nil
}

func truncatePKConstraint(name string) string {
	const (
		// MaxTableLength in postgres is 63 when building _fk or _pk we want to truncate the name to 60 chars max
		maxTableNamePKConstraint = 60
	)

	if len(name) > maxTableNamePKConstraint {
		return name[:maxTableNamePKConstraint]
	}
	return name
}
