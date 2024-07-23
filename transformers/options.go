package transformers

type StructTransformerOption func(*structTransformer)

// WithUnwrapAllEmbeddedStructs instructs codegen to unwrap all embedded fields (1 level deep only)
func WithUnwrapAllEmbeddedStructs() StructTransformerOption {
	return func(t *structTransformer) {
		t.unwrapAllEmbeddedStructFields = true
	}
}

// WithUnwrapStructFields allows to unwrap specific struct fields (1 level deep only)
func WithUnwrapStructFields(fields ...string) StructTransformerOption {
	return func(t *structTransformer) {
		t.structFieldsToUnwrap = fields
	}
}

// WithSkipFields allows to specify what struct fields should be skipped.
func WithSkipFields(fields ...string) StructTransformerOption {
	return func(t *structTransformer) {
		t.skipFields = fields
	}
}

// WithNameTransformer overrides how column name will be determined.
// DefaultNameTransformer is used as the default.
func WithNameTransformer(transformer NameTransformer) StructTransformerOption {
	return func(t *structTransformer) {
		t.nameTransformer = transformer
	}
}

// WithTypeTransformer overrides how column type will be determined.
// DefaultTypeTransformer is used as the default.
func WithTypeTransformer(transformer TypeTransformer) StructTransformerOption {
	return func(t *structTransformer) {
		t.typeTransformer = transformer
	}
}

// WithResolverTransformer overrides how column resolver will be determined.
// DefaultResolverTransformer is used as the default.
func WithResolverTransformer(transformer ResolverTransformer) StructTransformerOption {
	return func(t *structTransformer) {
		t.resolverTransformer = transformer
	}
}

// WithIgnoreInTestsTransformer overrides how column ignoreInTests will be determined.
// DefaultIgnoreInTestsTransformer is used as the default.
func WithIgnoreInTestsTransformer(transformer IgnoreInTestsTransformer) StructTransformerOption {
	return func(t *structTransformer) {
		t.ignoreInTestsTransformer = transformer
	}
}

// WithPrimaryKeys allows to specify what struct fields should be used as primary keys
func WithPrimaryKeys(fields ...string) StructTransformerOption {
	return func(t *structTransformer) {
		t.pkFields = fields
	}
}

// WithPrimaryKeyComponents allows to specify what struct fields should be used as primary key components
func WithPrimaryKeyComponents(fields ...string) StructTransformerOption {
	return func(t *structTransformer) {
		t.pkComponentFields = fields
	}
}
