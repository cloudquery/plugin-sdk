package transformers

import "reflect"

type IgnoreInTestsTransformer func(field reflect.StructField) bool

func DefaultIgnoreInTestsTransformer(_ reflect.StructField) bool {
	return false
}

var _ IgnoreInTestsTransformer = DefaultIgnoreInTestsTransformer
