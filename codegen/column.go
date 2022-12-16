package codegen

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/schema"
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

func (t *TableDefinition) pkOrder() error {
	if len(t.pkColumns) == 0 {
		// no need for extra work
		return nil
	}

	// check for dups
	pkColumns := slices.Clone(t.pkColumns)
	slices.Sort(pkColumns)
	pkColumns = slices.Compact(pkColumns)
	if len(pkColumns) != len(t.pkColumns) {
		return fmt.Errorf("%s table definition has %d duplicate PK columns", t.Name, len(t.pkColumns)-len(pkColumns))
	}

	// Add PK columns first in the order they were specified
	columns := slices.Clone(t.Columns)
	t.Columns = make(ColumnDefinitions, 0, len(t.Columns)) // reset
	for _, name := range t.pkColumns {
		idx := slices.IndexFunc(columns,
			func(def ColumnDefinition) bool {
				return def.Name == name
			},
		)
		if idx < 0 {
			return fmt.Errorf("%s table definition missing %q column required for PK", t.Name, name)
		}

		col := columns[idx]
		col.Options.PrimaryKey = true
		t.Columns = append(t.Columns, col)
		columns = slices.Delete(columns, idx, idx+1)
	}

	// If we got here, all t.pkColumns have been processed and there are no duplicates
	// Now we need to take all remaining PKs
	var rest ColumnDefinitions
	for _, col := range columns {
		if col.Options.PrimaryKey {
			t.Columns = append(t.Columns, col)
		} else {
			rest = append(rest, col)
		}
	}

	t.Columns = append(t.Columns, rest...)
	return nil
}
