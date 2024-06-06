package transformers

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/thoas/go-funk"
)

type structTransformer struct {
	table                         *schema.Table
	skipFields                    []string
	nameTransformer               NameTransformer
	typeTransformer               TypeTransformer
	resolverTransformer           ResolverTransformer
	ignoreInTestsTransformer      IgnoreInTestsTransformer
	unwrapAllEmbeddedStructFields bool
	structFieldsToUnwrap          []string
	pkFields                      []string
	pkFieldsFound                 []string
	pkComponentFields             []string
	pkComponentFieldsFound        []string
}

func isFieldStruct(reflectType reflect.Type) bool {
	switch reflectType.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		return reflectType.Elem().Kind() == reflect.Struct
	default:
		return false
	}
}

func isTypeIgnored(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Func,
		reflect.Chan,
		reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

func (t *structTransformer) getUnwrappedFields(field reflect.StructField) []reflect.StructField {
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

func (t *structTransformer) unwrapField(field reflect.StructField) error {
	unwrappedFields := t.getUnwrappedFields(field)
	var parent *reflect.StructField
	// For non-embedded structs we need to add the parent field name to the path
	if !field.Anonymous {
		parent = &field
	}
	for _, f := range unwrappedFields {
		switch {
		case t.shouldUnwrapField(f):
			if err := t.unwrapField(f); err != nil {
				return err
			}
		default:
			if err := t.addColumnFromField(f, parent); err != nil {
				return fmt.Errorf("failed to add column from field %s: %w", f.Name, err)
			}
		}
	}
	return nil
}

func (t *structTransformer) shouldUnwrapField(field reflect.StructField) bool {
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

func (t *structTransformer) ignoreField(field reflect.StructField) bool {
	switch {
	case len(field.Name) == 0,
		slices.Contains(t.skipFields, field.Name),
		isTypeIgnored(field.Type):
		return true
	case !field.IsExported():
		return !t.shouldUnwrapField(field)
	default:
		return false
	}
}

func (t *structTransformer) addColumnFromField(field reflect.StructField, parent *reflect.StructField) error {
	if t.ignoreField(field) {
		return nil
	}

	columnType, err := t.typeTransformer(field)
	if err != nil {
		return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
	}

	if columnType == nil {
		columnType, err = DefaultTypeTransformer(field)
		if err != nil {
			return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
	}
	if columnType == nil {
		return nil // ignored
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
	if t.table.Columns.Get(name) != nil {
		return nil
	}

	resolver := t.resolverTransformer(field, path)
	if resolver == nil {
		resolver = DefaultResolverTransformer(field, path)
	}

	column := schema.Column{
		Name:          name,
		Type:          columnType,
		Resolver:      resolver,
		IgnoreInTests: t.ignoreInTestsTransformer(field),
	}

	for _, pk := range t.pkFields {
		if pk == path {
			// use path to allow the following
			// 1. Don't duplicate the PK fields if the unwrapped struct contains a fields with the same name
			// 2. Allow specifying the nested unwrapped field as part of the PK.
			column.PrimaryKey = true
			t.pkFieldsFound = append(t.pkFieldsFound, pk)
		}
	}

	for _, pk := range t.pkComponentFields {
		if pk == path {
			// use path to allow the following
			// 1. Don't duplicate the PK fields if the unwrapped struct contains a fields with the same name
			// 2. Allow specifying the nested unwrapped field as part of the PK.
			column.PrimaryKeyComponent = true
			t.pkComponentFieldsFound = append(t.pkComponentFieldsFound, pk)
		}
	}

	t.table.Columns = append(t.table.Columns, column)

	return nil
}

func TransformWithStruct(st any, opts ...StructTransformerOption) schema.Transform {
	t := &structTransformer{
		nameTransformer:          DefaultNameTransformer,
		typeTransformer:          DefaultTypeTransformer,
		resolverTransformer:      DefaultResolverTransformer,
		ignoreInTestsTransformer: DefaultIgnoreInTestsTransformer,
	}
	for _, opt := range opts {
		opt(t)
	}

	return func(table *schema.Table) error {
		t.table = table
		e := reflect.ValueOf(st)
		if e.Kind() == reflect.Pointer {
			e = e.Elem()
		}
		if e.Kind() == reflect.Slice {
			e = reflect.MakeSlice(e.Type(), 1, 1).Index(0)
		}
		if e.Kind() != reflect.Struct {
			return fmt.Errorf("expected struct, got %s", e.Kind())
		}
		if err := t.unwrapField(reflect.StructField{Anonymous: true, Type: e.Type()}); err != nil {
			return err
		}
		// Validate that all expected PK fields were found
		if diff := funk.SubtractString(t.pkFields, t.pkFieldsFound); len(diff) > 0 {
			return fmt.Errorf("failed to create all of the desired primary keys: %v", diff)
		}

		if diff := funk.SubtractString(t.pkComponentFields, t.pkComponentFieldsFound); len(diff) > 0 {
			return fmt.Errorf("failed to find all of the desired primary key components: %v", diff)
		}
		return nil
	}
}
