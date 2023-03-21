package schema

import (
	"context"
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
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

// ParentColumnResolver resolves a column from the parent's table data, if name isn't set the column will be set to null
func ParentColumnResolver(name string) ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, r.Parent.Get(name))
	}
}

func ObjectHashResolve() ColumnResolver {
	return func(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
		hash, err := hashstructure.Hash(r.Item, hashstructure.FormatV2, nil)
		if err != nil {
			return err
		}
		return r.Set(c.Name, fmt.Sprint(hash))
	}
}
