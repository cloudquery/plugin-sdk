package codegen

import (
	"reflect"
	"strings"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/iancoleman/strcase"
)

type TableOptions func(*TableDefinition)

func WithNameTransformer(transformer func(field reflect.StructField) string) TableOptions {
	return func(t *TableDefinition) {
		t.nameTransformer = transformer
	}
}

func defaultTransformer(field reflect.StructField) string {
	name := field.Name
	if jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]; len(jsonTag) > 0 {
		// return empty string if the field is not related api response
		if jsonTag == "-" {
			return ""
		}
		name = jsonTag
	}
	return strcase.ToSnake(name)
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

func WithCustomValueTypeResolver(t schema.ValueType, fn string) TableOptions {
	return func(definition *TableDefinition) {
		if definition.customTypeResolvers == nil {
			definition.customTypeResolvers = map[schema.ValueType]string{
				t: fn,
			}
			return
		}
		definition.customTypeResolvers[t] = fn
	}
}
