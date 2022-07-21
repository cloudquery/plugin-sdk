package helpers

import "reflect"

// ToPointer takes an interface{} object and will return a pointer to this object
// if the object is not already a pointer. Otherwise, it will return the original value.
// It is safe to typecast the return-value of GetPointer into a pointer of the right type,
// except in very special cases (such as passing in nil without an explicit type)
func ToPointer(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		return v
	}
	if !val.IsValid() {
		return v
	}
	p := reflect.New(val.Type())
	p.Elem().Set(val)
	return p.Interface()
}
