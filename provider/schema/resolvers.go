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

// ParentIdResolver resolves the cq_id from the parent
// if you want to reference the parent's primary keys use ParentResourceFieldResolver as required.
func ParentIdResolver(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
	return r.Set(c.Name, r.Parent.Id())
}

// ParentResourceFieldResolver resolves a field from the parent's resource, the value is expected to be set
// if name isn't set the field will be set to null
func ParentResourceFieldResolver(name string) ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, r.Parent.Get(name))
	}
}

// ParentPathResolver resolves a field from the parent
func ParentPathResolver(path string) ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, funk.Get(r.Parent.Item, path, funk.WithAllowZero()))
	}
}
