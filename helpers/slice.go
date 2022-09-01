package helpers

import "reflect"

// InterfaceSlice converts any interface{} into a []interface{} slice
func InterfaceSlice(slice interface{}) []interface{} {
	// if value is nil return nil
	if slice == nil {
		return nil
	}
	s := reflect.ValueOf(slice)
	//handle slice behind pointer
	if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Slice {
		// Keep the distinction between nil and empty slice input
		if s.Elem().IsNil() {
			return nil
		}

		ret := make([]interface{}, s.Elem().Len())
		for i := 0; i < s.Elem().Len(); i++ {
			ret[i] = s.Elem().Index(i).Interface()
		}
		return ret
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
