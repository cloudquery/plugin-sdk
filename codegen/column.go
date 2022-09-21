package codegen

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/schema"
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

	// generate a PathResolver to use by default
	pathResolver := fmt.Sprintf(`schema.PathResolver("%s")`, field.Name)
	name, err := t.nameTransformer(field)
	if err != nil {
		return fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
	}
	// skip field if there is no name
	if name == "" {
		return nil
	}
	if parent != nil {
		pathResolver = fmt.Sprintf(`schema.PathResolver("%s.%s")`, parent.Name, field.Name)
		parentName, err := t.nameTransformer(*parent)
		if err != nil {
			return fmt.Errorf("failed to transform field name for parent field %s: %w", parent.Name, err)
		}
		name = fmt.Sprintf("%s_%s", parentName, name)
	}

	column := ColumnDefinition{
		Name:     name,
		Type:     columnType,
		Resolver: pathResolver,
	}
	t.Columns = append(t.Columns, column)
	return nil
}
