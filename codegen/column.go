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

// The returned slice of columns is always of length 0 or 1.
func (t *TableDefinition) makeColumnFromField(field reflect.StructField, parent *reflect.StructField) (ColumnDefinitions, error) {
	if t.shouldIgnoreField(field) {
		return nil, nil
	}

	typeTransformer := t.TypeTransformer
	if typeTransformer == nil {
		typeTransformer = DefaultTypeTransformer
	}

	columnType, err := typeTransformer(field)
	if err != nil {
		return nil, fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
	}

	if columnType == schema.TypeInvalid {
		columnType, err = DefaultTypeTransformer(field)
		if err != nil {
			return nil, fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
	}

	nameTransformer := t.NameTransformer
	if nameTransformer == nil {
		nameTransformer = DefaultNameTransformer
	}

	path := field.Name
	name, err := nameTransformer(field)
	if err != nil {
		return nil, fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
	}
	// skip field if there is no name
	if name == "" {
		return nil, nil
	}
	if parent != nil {
		parentName, err := nameTransformer(*parent)
		if err != nil {
			return nil, fmt.Errorf("failed to transform field name for parent field %s: %w", parent.Name, err)
		}
		name = parentName + "_" + name
		path = parent.Name + `.` + path
	}

	resolverTransformer := t.ResolverTransformer
	if resolverTransformer == nil {
		resolverTransformer = DefaultResolverTransformer
	}

	resolver, err := resolverTransformer(field, path)
	if err != nil {
		return nil, err
	}
	if resolver == "" {
		resolver = defaultResolver(path)
	}

	return ColumnDefinitions{
		ColumnDefinition{
			Name:     name,
			Type:     columnType,
			Resolver: resolver,
		},
	}, nil

}
