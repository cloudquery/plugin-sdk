package transformers

import "reflect"

type NullableFieldTransformer func(field reflect.StructField) bool

func DefaultNullableFieldTransformer(_ reflect.StructField) bool {
	return true
}

var _ NullableFieldTransformer = DefaultNullableFieldTransformer
