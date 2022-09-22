package codegen

import (
	"reflect"
)

func isFieldStruct(reflectType reflect.Type) bool {
	switch reflectType.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		return reflectType.Elem().Kind() == reflect.Struct
	default:
		return false
	}
}

func isTypeIgnored(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface,
		reflect.Func,
		reflect.Chan,
		reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
