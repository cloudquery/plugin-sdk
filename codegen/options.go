package codegen

import (
	"github.com/rs/zerolog"
)

type TableOption func(*TableDefinition)

// WithNameTransformer overrides how column name will be determined.
func WithNameTransformer(transformer NameTransformer) TableOption {
	return func(t *TableDefinition) {
		t.nameTransformer = transformer
	}
}

// WithSkipFields allows to specify what struct fields should be skipped.
func WithSkipFields(fields []string) TableOption {
	return func(t *TableDefinition) {
		t.skipFields = fields
	}
}

// WithExtraColumns allows passing additional ColumnDefinitions
func WithExtraColumns(columns ColumnDefinitions) TableOption {
	return func(t *TableDefinition) {
		t.extraColumns = columns
	}
}

// WithUnwrapStructFields allows to unwrap specific struct fields (1 level deep only)
func WithUnwrapStructFields(fields []string) TableOption {
	return func(t *TableDefinition) {
		t.structFieldsToUnwrap = fields
	}
}

// WithUnwrapAllEmbeddedStructs instructs codegen to unwrap all embedded fields (1 level deep only)
func WithUnwrapAllEmbeddedStructs() TableOption {
	return func(t *TableDefinition) {
		t.unwrapAllEmbeddedStructFields = true
	}
}

// WithLogger allows passing custom logger
func WithLogger(logger zerolog.Logger) TableOption {
	return func(t *TableDefinition) {
		t.logger = logger
	}
}

// WithTypeTransformer sets a function that can override the schema type for specific fields. Return `schema.TypeInvalid` to fall back to default behavior.
func WithTypeTransformer(transformer TypeTransformer) TableOption {
	return func(t *TableDefinition) {
		t.typeTransformer = transformer
	}
}
