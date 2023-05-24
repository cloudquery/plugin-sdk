package transformers

import (
	"reflect"

	"github.com/cloudquery/plugin-sdk/v3/schema"
)

type ResolverTransformer func(field reflect.StructField, path string) schema.ColumnResolver

func DefaultResolverTransformer(_ reflect.StructField, path string) schema.ColumnResolver {
	return schema.PathResolver(path)
}
