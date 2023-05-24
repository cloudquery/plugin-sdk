package transformers

import (
	"reflect"
)

func Nullable(typ reflect.Type) bool {
	if typ == nil {
		// can happen in some occasions
		return true
	}
	switch typ.Kind() {
	case reflect.Invalid,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
