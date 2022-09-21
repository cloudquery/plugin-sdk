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

		nameTransformer NameTransformer
		typeTransformer TypeTransformer

		extraColumns ColumnDefinitions

		skipFields           []string
		structFieldsToUnwrap []string

		unwrapAllEmbeddedStructFields bool

		logger zerolog.Logger
	}
)

func (t *TableDefinition) shouldUnwrapField(field reflect.StructField) bool {
	switch {
	case !isFieldStruct(field.Type):
		return false
	case slices.Contains(t.structFieldsToUnwrap, field.Name):
		return true
	case !field.Anonymous:
		return false
	case t.unwrapAllEmbeddedStructFields:
		return true
	default:
		return false
	}
}

func (t *TableDefinition) getUnwrappedFields(field reflect.StructField) []reflect.StructField {
	reflectType := field.Type
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	fields := make([]reflect.StructField, 0)
	for i := 0; i < reflectType.NumField(); i++ {
		sf := reflectType.Field(i)
		if t.ignoreField(sf) {
			continue
		}

		fields = append(fields, sf)
	}
	return fields
}

func (t *TableDefinition) ignoreField(field reflect.StructField) bool {
	switch {
	case len(field.Name) == 0,
		slices.Contains(t.skipFields, field.Name),
		!field.IsExported():
		return true
	default:
		return false
	}
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

// NewTableFromStruct creates a new TableDefinition from a struct by inspecting its fields
func NewTableFromStruct(name string, obj interface{}, opts ...TableOptions) (*TableDefinition, error) {
	t := &TableDefinition{
		Name:            name,
		nameTransformer: DefaultNameTransformer,
		typeTransformer: DefaultTypeTransformer,
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

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)

		if t.shouldUnwrapField(field) {
			unwrappedFields := t.getUnwrappedFields(field)
			var parent *reflect.StructField
			// For non embedded structs we need to add the parent field name to the path
			if !field.Anonymous {
				parent = &field
			}
			for _, f := range unwrappedFields {
				if err := t.addColumnFromField(f, parent); err != nil {
					return nil, fmt.Errorf("failed to add column from field %s: %w", f.Name, err)
				}
			}
		} else {
			if err := t.addColumnFromField(field, nil); err != nil {
				return nil, fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
			}
		}
	}

	return t, nil
}
