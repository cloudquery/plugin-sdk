package codegen

import (
	"github.com/rs/zerolog"
)

// WithNameTransformer overrides how column name will be determined.
func WithNameTransformer(transformer NameTransformer) TableOptions {
	return func(t *TableDefinition) {
		t.nameTransformer = transformer
	}
}

// WithSkipFields allows to specify what struct fields should be skipped.
func WithSkipFields(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.skipFields = fields
	}
}

// WithExtraColumns allows passing additional ColumnDefinitions
func WithExtraColumns(columns ColumnDefinitions) TableOptions {
	return func(t *TableDefinition) {
		t.extraColumns = columns
	}
}

// WithUnwrapFieldsStructs allows to unwrap specific struct fields (1 level deep only)
func WithUnwrapFieldsStructs(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.structFieldsToUnwrap = fields
	}
}

// WithUnwrapAllEmbeddedStructs instructs codegen to unwrap all embedded fields (1 level deep only)
func WithUnwrapAllEmbeddedStructs() TableOptions {
	return func(t *TableDefinition) {
		t.unwrapAllEmbeddedStructFields = true
	}
}

// WithLogger allows passing custom logger
func WithLogger(logger zerolog.Logger) TableOptions {
	return func(t *TableDefinition) {
		t.logger = logger
	}
}

// WithTypeTransformer sets a function that can override the schema type for specific fields. Return `schema.TypeInvalid` to fall back to default behavior.
func WithTypeTransformer(transformer TypeTransformer) TableOptions {
	return func(t *TableDefinition) {
		t.typeTransformer = transformer
	}
}
