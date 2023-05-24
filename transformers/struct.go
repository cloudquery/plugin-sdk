package transformers

import (
	"fmt"
	"reflect"

	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/thoas/go-funk"
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
	pkFields                      []string
	pkFieldsFound                 []string
}

func isFieldStruct(reflectType reflect.Type) bool {
	switch reflectType.Kind() {
	case reflect.Pointer:
		return isFieldStruct(reflectType.Elem())
	case reflect.Struct:
		return true
	default:
		return false
	}
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
		// Validate that all expected PK fields were found
		if diff := funk.SubtractString(t.pkFields, t.pkFieldsFound); len(diff) > 0 {
			return fmt.Errorf("failed to create all of the desired primary keys: %v", diff)
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
