package schema

import (
	"context"
	"reflect"

	"github.com/thoas/go-funk"
)

func PathResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		r.Set(c.Name, funk.GetAllowZero(r.Item, path))
		return nil
	}
}

func ParentIdResolver(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
	r.Set(c.Name, r.Parent.Id())
	return nil
}

func interfaceSlice(slice interface{}) []interface{} {
	// if value is nil return nil
	if slice == nil {
		return nil
	}

	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return []interface{}{slice}
	}
	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
