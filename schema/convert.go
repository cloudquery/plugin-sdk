//nolint:revive,gocritic,unused
package schema

import (
	"reflect"
	"time"
)

// underlyingNumberType gets the underlying type that can be converted to Int2, Int4, Int8, Float4, or Float8
func underlyingNumberType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	case reflect.Int:
		convVal := int(refVal.Int())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Int8:
		convVal := int8(refVal.Int())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Int16:
		convVal := int16(refVal.Int())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Int32:
		convVal := int32(refVal.Int())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Int64:
		convVal := refVal.Int()
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Uint:
		convVal := uint(refVal.Uint())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Uint8:
		convVal := uint8(refVal.Uint())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Uint16:
		convVal := uint16(refVal.Uint())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Uint32:
		convVal := uint32(refVal.Uint())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Uint64:
		convVal := refVal.Uint()
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Float32:
		convVal := float32(refVal.Float())
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.Float64:
		convVal := refVal.Float()
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	case reflect.String:
		convVal := refVal.String()
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	}

	return nil, false
}

// underlyingBoolType gets the underlying type that can be converted to Bool
func underlyingBoolType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	case reflect.Bool:
		convVal := refVal.Bool()
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	}

	return nil, false
}

// underlyingBytesType gets the underlying type that can be converted to []byte
func underlyingBytesType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	case reflect.Slice:
		if refVal.Type().Elem().Kind() == reflect.Uint8 {
			convVal := refVal.Bytes()
			return convVal, reflect.TypeOf(convVal) != refVal.Type()
		}
	}

	return nil, false
}

// underlyingTimeType gets the underlying type that can be converted to time.Time
func underlyingTimeType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	}

	timeType := reflect.TypeOf(time.Time{})
	if refVal.Type().ConvertibleTo(timeType) {
		return refVal.Convert(timeType).Interface(), true
	}

	return nil, false
}

// underlyingStringType gets the underlying type that can be converted to String
func underlyingStringType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	case reflect.String:
		convVal := refVal.String()
		return convVal, reflect.TypeOf(convVal) != refVal.Type()
	}

	return nil, false
}

// underlyingUUIDType gets the underlying type that can be converted to [16]byte
func underlyingUUIDType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	}

	uuidType := reflect.TypeOf([16]byte{})
	if refVal.Type().ConvertibleTo(uuidType) {
		return refVal.Convert(uuidType).Interface(), true
	}

	return nil, false
}

// underlyingSliceType gets the underlying slice type
func underlyingSliceType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	case reflect.Slice:
		baseSliceType := reflect.SliceOf(refVal.Type().Elem())
		if refVal.Type().ConvertibleTo(baseSliceType) {
			convVal := refVal.Convert(baseSliceType)
			return convVal.Interface(), reflect.TypeOf(convVal.Interface()) != refVal.Type()
		}
	}

	return nil, false
}

// underlyingPtrType dereferences a pointer
func underlyingPtrType(val any) (any, bool) {
	refVal := reflect.ValueOf(val)

	//nolint:gocritic
	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	}

	return nil, false
}
