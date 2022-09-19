package schema

import (
	"context"

	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PathResolver resolves a field in the Resource.Item
//
// Examples:
// PathResolver("Field")
// PathResolver("InnerStruct.Field")
// PathResolver("InnerStruct.InnerInnerStruct.Field")
func PathResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		data := funk.Get(r.Item, path, funk.WithAllowZero())
		// special case for timestamppb.Timestamp
		if ts, ok := data.(*timestamppb.Timestamp); ok {
			data = ts.AsTime()
		}
		return r.Set(c.Name, data)
	}
}

func CQUUIDResolver() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		uuidGen := uuid.New()
		return r.Set(c.Name, uuidGen)
	}
}

// ParentIDResolver resolves the cq_id from the parent
// if you want to reference the parent's primary keys use ParentResourceFieldResolver as required.
func ParentIDResolver(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
	return r.Set(c.Name, r.Parent.ID())
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
