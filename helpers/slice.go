package helpers

import "reflect"

// InterfaceSlice converts any interface{} into a []interface{} slice
func InterfaceSlice(slice interface{}) []interface{} {
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
