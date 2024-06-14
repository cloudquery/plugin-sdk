package transformers

import (
	"reflect"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type ResolverTransformer func(field reflect.StructField, path string) schema.ColumnResolver

func DefaultResolverTransformer(_ reflect.StructField, path string) schema.ColumnResolver {
	return schema.PathResolver(path)
}

var _ ResolverTransformer = DefaultResolverTransformer
