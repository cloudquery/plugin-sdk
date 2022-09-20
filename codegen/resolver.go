package codegen

import (
	"github.com/cloudquery/plugin-sdk/schema"
)

type ResolverProvider func(path string) schema.ColumnResolver
