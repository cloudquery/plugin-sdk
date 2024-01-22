package scalar

import "reflect"

func IsNil(value any) bool {
	if value == nil {
		return true
	}

	// typed nil, here we go again
	return isReflectValueNil(reflect.ValueOf(value))
}

func isReflectValueNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Map,
		reflect.Pointer,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice,
		reflect.Chan,
		reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
