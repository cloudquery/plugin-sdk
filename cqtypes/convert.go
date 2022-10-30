//nolint:revive,gocritic
package cqtypes

import (
	"reflect"
	"time"
)

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

// var kindTypes map[reflect.Kind]reflect.Type

// underlyingNumberType gets the underlying type that can be converted to Int2, Int4, Int8, Float4, or Float8
func underlyingNumberType(val interface{}) (interface{}, bool) {
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
func underlyingBoolType(val interface{}) (interface{}, bool) {
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
func underlyingBytesType(val interface{}) (interface{}, bool) {
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
func underlyingTimeType(val interface{}) (interface{}, bool) {
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
func underlyingStringType(val interface{}) (interface{}, bool) {
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
func underlyingUUIDType(val interface{}) (interface{}, bool) {
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
func underlyingSliceType(val interface{}) (interface{}, bool) {
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
func underlyingPtrType(val interface{}) (interface{}, bool) {
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

// func init() {
// 	kindTypes = map[reflect.Kind]reflect.Type{
// 		reflect.Bool:    reflect.TypeOf(false),
// 		reflect.Float32: reflect.TypeOf(float32(0)),
// 		reflect.Float64: reflect.TypeOf(float64(0)),
// 		reflect.Int:     reflect.TypeOf(int(0)),
// 		reflect.Int8:    reflect.TypeOf(int8(0)),
// 		reflect.Int16:   reflect.TypeOf(int16(0)),
// 		reflect.Int32:   reflect.TypeOf(int32(0)),
// 		reflect.Int64:   reflect.TypeOf(int64(0)),
// 		reflect.Uint:    reflect.TypeOf(uint(0)),
// 		reflect.Uint8:   reflect.TypeOf(uint8(0)),
// 		reflect.Uint16:  reflect.TypeOf(uint16(0)),
// 		reflect.Uint32:  reflect.TypeOf(uint32(0)),
// 		reflect.Uint64:  reflect.TypeOf(uint64(0)),
// 		reflect.String:  reflect.TypeOf(""),
// 	}
// }
