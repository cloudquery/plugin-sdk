package transformers

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/codegen"
	"github.com/cloudquery/plugin-sdk/schema"
	"golang.org/x/exp/slices"
)

type structTransformer struct {
	table *schema.Table
	skipFields []string
	nameTransformer     NameTransformer
	typeTransformer     TypeTransformer
}

type NameTransformer func(reflect.StructField) (string, error)

type TypeTransformer func(reflect.StructField) (schema.ValueType, error)

type StructTransformerOption func(*structTransformer)

// WithSkipFields allows to specify what struct fields should be skipped.
func WithSkipFields(fields []string) StructTransformerOption {
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

// WithTypeTransformer sets a function that can override the schema type for specific fields. Return `schema.TypeInvalid` to fall back to default behavior.
// DefaultTypeTransformer is used as the default.
func WithTypeTransformer(transformer TypeTransformer) StructTransformerOption {
	return func(t *structTransformer) {
		t.typeTransformer = transformer
	}
}

func TransformWithStruct(st any, opts ...StructTransformerOption) schema.Transform {
	t := &structTransformer{
		nameTransformer:     codegen.DefaultNameTransformer,
		typeTransformer:     codegen.DefaultTypeTransformer,
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
	
			if err := t.addColumnFromField(field, nil); err != nil {
				return fmt.Errorf("failed to add column for field %s: %w", field.Name, err)
			}
		}
		return nil
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
	t.table.Columns = append(t.table.Columns,
		schema.Column{
			Name:     name,
			Type:     columnType,
			Resolver: schema.PathResolver(path),
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