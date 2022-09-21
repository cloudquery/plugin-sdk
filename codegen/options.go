package codegen

import (
	"reflect"
)

type TableOptions func(*TableDefinition)

func WithNameTransformer(transformer func(field reflect.StructField) string) TableOptions {
	return func(t *TableDefinition) {
		t.nameTransformer = transformer
	}
}

func WithSkipFields(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.skipFields = fields
	}
}

func WithExtraColumns(columns []ColumnDefinition) TableOptions {
	return func(t *TableDefinition) {
		t.extraColumns = columns
	}
}

// WithUnwrapFieldsStructs unwraps specific struct fields (1 level deep only)
func WithUnwrapFieldsStructs(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.structFieldsToUnwrap = fields
	}
}

// WithUnwrapAllEmbeddedStructs unwraps all fields that are embedded structs (1 level deep only)
func WithUnwrapAllEmbeddedStructs() TableOptions {
	return func(t *TableDefinition) {
		t.unwrapAllEmbeddedStructFields = true
	}
}
