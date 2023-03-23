package helpers

import "reflect"

// InterfaceSlice converts any any into a []any slice
func InterfaceSlice(slice any) []any {
	// if value is nil return nil
	if slice == nil {
		return nil
	}
	s := reflect.ValueOf(slice)
	// handle slice behind pointer
	if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Slice {
		// Keep the distinction between nil and empty slice input
		if s.Elem().IsNil() {
			return nil
		}

		ret := make([]any, s.Elem().Len())
		for i := 0; i < s.Elem().Len(); i++ {
			ret[i] = s.Elem().Index(i).Interface()
		}
		return ret
	}
	if s.Kind() != reflect.Slice {
		return []any{slice}
	}

	ret := make([]any, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
