package codegen

import (
	"github.com/rs/zerolog"
)

type TableOption func(*TableDefinition)

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

// WithPKColumns allows to specify what columns should be considered PKs without need for WithExtraColumns + WithSkipFields
func WithPKColumns(columnNames ...string) TableOption {
	return func(t *TableDefinition) {
		if t.extraPKColumns == nil {
			t.extraPKColumns = make(map[string]struct{}, len(columnNames))
		}
		for _, name := range columnNames {
			t.extraPKColumns[name] = struct{}{}
		}
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

// WithNameTransformer overrides how column name will be determined.
// DefaultNameTransformer is used as the default.
func WithNameTransformer(transformer NameTransformer) TableOption {
	return func(t *TableDefinition) {
		t.nameTransformer = transformer
	}
}

// WithTypeTransformer sets a function that can override the schema type for specific fields. Return `schema.TypeInvalid` to fall back to default behavior.
// DefaultTypeTransformer is used as the default.
func WithTypeTransformer(transformer TypeTransformer) TableOption {
	return func(t *TableDefinition) {
		t.typeTransformer = transformer
	}
}

// WithResolverTransformer sets a function that can override the resolver for a field.
// DefaultResolverTransformer is used as the default.
// If the transformer provided returns error, gen will fail.
// To fallback onto the default resolver return "", nil.
func WithResolverTransformer(transformer ResolverTransformer) TableOption {
	return func(t *TableDefinition) {
		t.resolverTransformer = transformer
	}
}
