package codegen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/thoas/go-funk"
	"golang.org/x/exp/slices"
)

type (
	ColumnDefinition struct {
		Name          string
		Type          schema.ValueType
		Resolver      string
		Description   string
		IgnoreInTests bool
		Options       schema.ColumnCreationOptions
	}
	ColumnDefinitions []ColumnDefinition
)

// sorted returns a copy of sorted columns
func (c ColumnDefinitions) sorted() ColumnDefinitions {
	columns := make(ColumnDefinitions, len(c))
	copy(columns, c)

	const cqPrefix = "_cq"
	slices.SortStableFunc(columns, func(a, b ColumnDefinition) bool {
		switch {
		case strings.HasPrefix(a.Name, cqPrefix):
			return !strings.HasPrefix(b.Name, cqPrefix) || a.Name < b.Name
		case strings.HasPrefix(b.Name, cqPrefix):
			return false
		case a.Options.PrimaryKey:
			return !b.Options.PrimaryKey || a.Name < b.Name
		case b.Options.PrimaryKey:
			return false
		default:
			return a.Name < b.Name
		}
	})

	return columns
}

func (c ColumnDefinitions) GetByName(name string) *ColumnDefinition {
	for _, col := range c {
		if col.Name == name {
			return &col
		}
	}
	return nil
}

func (t *TableDefinition) addColumnFromField(field reflect.StructField, parent *reflect.StructField) error {
	if t.ignoreField(field) {
		return nil
	}

	columnType, err := t.typeTransformer(field)
	if err != nil {
		return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
	}

	if columnType == schema.TypeInvalid {
		columnType, err = DefaultTypeTransformer(field)
		if err != nil {
			return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
	}

	path := field.Name
	name, err := t.nameTransformer(field)
	if err != nil {
		return fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
	}
	// skip field if there is no name
	if name == "" {
		return nil
	}
	if parent != nil {
		parentName, err := t.nameTransformer(*parent)
		if err != nil {
			return fmt.Errorf("failed to transform field name for parent field %s: %w", parent.Name, err)
		}
		name = parentName + "_" + name
		path = parent.Name + `.` + path
	}
	resolver, err := t.resolverTransformer(field, path)
	if err != nil {
		return err
	}
	if resolver == "" {
		resolver = defaultResolver(path)
	}

	t.Columns = append(t.Columns,
		ColumnDefinition{
			Name:     name,
			Type:     columnType,
			Resolver: resolver,
		},
	)
	return nil
}

func (t *TableDefinition) postProcessColumns() error {
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
		return fmt.Errorf("%s table definition has %d extra PK keys: %v", t.Name, len(t.extraPKColumns), funk.Keys(t.extraPKColumns))
	}

	t.Columns = columns.sorted()

	return nil
}
