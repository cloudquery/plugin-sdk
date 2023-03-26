package schema

import "reflect"

// Information on the internals of PostgreSQL arrays can be found in
// src/include/utils/array.h and src/backend/utils/adt/arrayfuncs.c. Of
// particular interest is the array_send function.

// type ArrayHeader struct {
// 	ContainsNull bool
// 	ElementOID   int32
// 	Dimensions   []ArrayDimension
// }

type ArrayDimension struct {
	Length     int32
	LowerBound int32
}

// nolint:unparam
func findDimensionsFromValue(value reflect.Value, dimensions []ArrayDimension, elementsLength int) ([]ArrayDimension, int, bool) {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		length := value.Len()
		if elementsLength == 0 {
			elementsLength = length
		} else {
			elementsLength *= length
		}
		dimensions = append(dimensions, ArrayDimension{Length: int32(length), LowerBound: 1})
		for i := 0; i < length; i++ {
			if d, l, ok := findDimensionsFromValue(value.Index(i), dimensions, elementsLength); ok {
				return d, l, true
			}
		}
	}
	return dimensions, elementsLength, true
}
