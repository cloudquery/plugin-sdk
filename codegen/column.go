package codegen

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/v2/schema"
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
