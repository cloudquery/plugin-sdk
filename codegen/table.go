package codegen

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type (
	ResourceDefinition struct {
		Name  string
		Table *TableDefinition
	}
	TableDefinition struct {
		Name        string
		Columns     ColumnDefinitions
		Description string
		Relations   []string

		Resolver             string
		IgnoreError          string
		Multiplex            string
		PostResourceResolver string
		PreResourceResolver  string

		nameTransformer     NameTransformer
		typeTransformer     TypeTransformer
		resolverTransformer ResolverTransformer

		extraColumns   ColumnDefinitions
		extraPKColumns map[string]struct{}

		skipFields           []string
		structFieldsToUnwrap []string

		unwrapAllEmbeddedStructFields bool

		logger zerolog.Logger
	}
)

func (t *TableDefinition) ignoreField(field reflect.StructField) bool {
	switch {
	case len(field.Name) == 0,
		slices.Contains(t.skipFields, field.Name),
		!field.IsExported(),
		isTypeIgnored(field.Type):
		return true
	default:
		return false
	}
}

// NewTableFromStruct creates a new TableDefinition from a struct by inspecting its fields
func NewTableFromStruct(name string, obj interface{}, opts ...TableOption) (*TableDefinition, error) {
	t := &TableDefinition{
		Name:                name,
		nameTransformer:     DefaultNameTransformer,
		typeTransformer:     DefaultTypeTransformer,
		resolverTransformer: DefaultResolverTransformer,
	}
	for _, opt := range opts {
		opt(t)
	}

	e := reflect.ValueOf(obj)
	if e.Kind() == reflect.Pointer {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", e.Kind())
	}

	t.Columns = append(t.Columns, t.extraColumns...)

	eType := e.Type()
	for i := 0; i < e.NumField(); i++ {
		field := eType.Field(i)

		switch {
		case t.shouldUnwrapField(field):
			if err := t.unwrapField(field); err != nil {
				return nil, err
			}
		default:
			if err := t.addColumnFromField(field, nil); err != nil {
				return nil, fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
			}
		}
	}

	// add PK options
	columns := make(ColumnDefinitions, 0, len(t.Columns))
	for _, column := range t.Columns {
		if _, ok := t.extraPKColumns[column.Name]; ok {
			column.Options.PrimaryKey = true
			delete(t.extraPKColumns, column.Name)
		}
		columns = append(columns, column)
	}
	if len(t.extraPKColumns) > 0 {
		return nil, fmt.Errorf("%s table definition has %d extra PK keys", t.Name, len(t.extraPKColumns))
	}
	t.Columns = columns

	return t, t.Check()
}

// Check that the resulting TableDefinition is correct (e.g., no overlapping column names).
func (t *TableDefinition) Check() error {
	if t == nil {
		return fmt.Errorf("nil table definition")
	}

	if len(t.Name) == 0 {
		return fmt.Errorf("empty table name")
	}

	if len(t.Columns) == 0 {
		return fmt.Errorf("no columns for table %s", t.Name)
	}

	columns := make(map[string]bool, len(t.Columns))
	for _, column := range t.Columns {
		switch {
		case column.Type == schema.TypeInvalid:
			return fmt.Errorf("%s->%s: invalid column type", t.Name, column.Name)
		case columns[column.Name]:
			return fmt.Errorf("%s->%s: duplicate column name", t.Name, column.Name)
		default:
			columns[column.Name] = true
		}
	}

	return nil
}
