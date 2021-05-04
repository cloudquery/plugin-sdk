package schema

import (
	"context"
	"reflect"

	"github.com/cloudquery/go-funk"
)

func PathResolver(path string) ColumnResolver {
	return func(_ context.Context, meta ClientMeta, r *Resource, c Column) error {
		return r.Set(c.Name, funk.GetAllowZero(r.Item, path))
	}
}

func ParentIdResolver(_ context.Context, _ ClientMeta, r *Resource, c Column) error {
	return r.Set(c.Name, r.Parent.Id())
}

func interfaceSlice(slice interface{}) []interface{} {
	// if value is nil return nil
	if slice == nil {
		return nil
	}
	s := reflect.ValueOf(slice)
	// Keep the distinction between nil and empty slice input
	if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Slice && s.Elem().IsNil() {
		return nil
	}
	if s.Kind() != reflect.Slice {
		return []interface{}{slice}
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
