package transformers

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/codegen"
	"github.com/cloudquery/plugin-sdk/schema"
	"golang.org/x/exp/slices"
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
}

type NameTransformer func(reflect.StructField) (string, error)

type TypeTransformer func(reflect.StructField) (schema.ValueType, error)

type ResolverTransformer func(field reflect.StructField, path string) schema.ColumnResolver

func DefaultResolverTransformer(_ reflect.StructField, path string) schema.ColumnResolver {
	return schema.PathResolver(path)
}

type IgnoreInTestsTransformer func(field reflect.StructField) bool

func DefaultIgnoreInTestsTransformer(_ reflect.StructField) bool {
	return false
}

type StructTransformerOption func(*structTransformer)

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

func TransformWithStruct(st any, opts ...StructTransformerOption) schema.Transform {
	t := &structTransformer{
		nameTransformer:          codegen.DefaultNameTransformer,
		typeTransformer:          codegen.DefaultTypeTransformer,
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
		eType := e.Type()
		for i := 0; i < e.NumField(); i++ {
			field := eType.Field(i)

			switch {
			case t.shouldUnwrapField(field):
				if err := t.unwrapField(field); err != nil {
					return err
				}
			default:
				if err := t.addColumnFromField(field, nil); err != nil {
					return fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
				}
			}
		}
		return nil
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
	// For non embedded structs we need to add the parent field name to the path
	if !field.Anonymous {
		parent = &field
	}
	for _, f := range unwrappedFields {
		if err := t.addColumnFromField(f, parent); err != nil {
			return fmt.Errorf("failed to add column from field %s: %w", f.Name, err)
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
		!field.IsExported(),
		isTypeIgnored(field.Type):
		return true
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

	if columnType == schema.TypeInvalid {
		columnType, err = codegen.DefaultTypeTransformer(field)
		if err != nil {
			return fmt.Errorf("failed to transform type for field %s: %w", field.Name, err)
		}
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

	t.table.Columns = append(t.table.Columns,
		schema.Column{
			Name:          name,
			Type:          columnType,
			Resolver:      resolver,
			IgnoreInTests: t.ignoreInTestsTransformer(field),
		},
	)

	return nil
}

func isTypeIgnored(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface,
		reflect.Func,
		reflect.Chan,
		reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
