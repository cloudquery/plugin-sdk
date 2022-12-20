package codegen

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

func (t *TableDefinition) shouldUnwrapField(field reflect.StructField) bool {
	switch {
	case !isFieldStruct(field.Type):
		return false
	case slices.Contains(t.StructFieldsToUnwrap, field.Name):
		return true
	case !field.Anonymous:
		return false
	case t.UnwrapAllEmbeddedStructFields:
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
		if t.shouldIgnoreField(sf) {
			continue
		}

		fields = append(fields, sf)
	}
	return fields
}

func (t *TableDefinition) unwrapField(field reflect.StructField) (ColumnDefinitions, error) {
	var columns ColumnDefinitions

	unwrappedFields := t.getUnwrappedFields(field)
	var parent *reflect.StructField
	// For non embedded structs we need to add the parent field name to the path
	if !field.Anonymous {
		parent = &field
	}
	for _, f := range unwrappedFields {
		fieldColumns, err := t.makeColumnFromField(f, parent)
		if err != nil {
			return nil, fmt.Errorf("failed to add column from field %s: %w", f.Name, err)
		}
		columns = append(columns, fieldColumns...)
	}

	return columns, nil
}
