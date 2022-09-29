package schema

import (
	"context"

	"github.com/thoas/go-funk"
)

// PathResolver resolves a field in the Resource.Item
//
// Examples:
// PathResolver("Field")
// PathResolver("InnerStruct.Field")
// PathResolver("InnerStruct.InnerInnerStruct.Field")
func PathResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, funk.Get(r.Item, path, funk.WithAllowZero()))
	}
}

// ParentIDResolver resolves the cq_id from the parent
// if you want to reference the parent's primary keys use ParentResourceFieldResolver as required.
func ParentIDResolver(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
	return r.Set(c.Name, r.Parent.ID())
}

// ParentColumnResolver resolves a column from the parent's table data, if name isn't set the column will be set to null
func ParentColumnResolver(name string) ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, r.Parent.Get(name))
	}
}
