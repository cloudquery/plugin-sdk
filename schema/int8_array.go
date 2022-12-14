//nolint:revive,gocritic
package schema

import (
	"fmt"
	"reflect"
	"strings"
)

type Int8ArrayTransformer interface {
	TransformInt8Array(*Int8Array) any
}

type Int8Array struct {
	Elements   []Int8
	Dimensions []ArrayDimension
	Status     Status
}

func (dst *Int8Array) GetStatus() Status {
	return dst.Status
}

func (*Int8Array) Type() ValueType {
	return TypeIntArray
}

func (dst *Int8Array) Size() int {
	return 8 * len(dst.Elements)
}

func (dst *Int8Array) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Int8Array)
	if !ok {
		return false
	}
	if dst.Status != s.Status {
		return false
	}
	if len(dst.Elements) != len(s.Elements) {
		return false
	}

	for i := range dst.Elements {
		if !(dst.Elements[i]).Equal(&s.Elements[i]) {
			return false
		}
	}

	return true
}

func (dst *Int8Array) fromString(value string) error {
	// this is basically back from string encoding
	if !strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		return &ValidationError{Type: TypeIntArray, Msg: cannotDecodeString, Value: value}
	}

	value = value[1 : len(value)-1]
	strs := strings.Split(value, ",")
	elements := make([]Int8, len(strs))
	for i := range strs {
		if err := elements[i].Set(strs[i]); err != nil {
			return err
		}
	}

	*dst = Int8Array{
		Elements:   elements,
		Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
		Status:     Present,
	}
	return nil
}

func (dst *Int8Array) String() string {
	var sb strings.Builder
	if dst.Status == Present {
		sb.WriteString("{")
		for i, element := range dst.Elements {
			if i != 0 {
				sb.WriteString(",")
			}
			sb.WriteString(element.String())
		}
		sb.WriteString("}")
	} else {
		return ""
	}
	return sb.String()
}

func (dst *Int8Array) Set(src any) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		*dst = Int8Array{Status: Null}
		return nil
	}

	if value, ok := src.(CQType); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	// Attempt to match to select common types:
	switch value := src.(type) {
	case []int16:
		//nolint:gocritic
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*int16:
		//nolint:gocritic
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint16:
		//nolint:gocritic
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*uint16:
		//nolint:gocritic
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []int32:
		//nolint:gocritic
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*int32:
		//nolint:gocritic
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint32:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*uint32:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []int64:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*int64:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint64:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*uint64:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []int:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*int:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*uint:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			elements := make([]Int8, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int8Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []Int8:
		if value == nil {
			*dst = Int8Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int8Array{Status: Present}
		} else {
			*dst = Int8Array{
				Elements:   value,
				Dimensions: []ArrayDimension{{Length: int32(len(value)), LowerBound: 1}},
				Status:     Present,
			}
		}
	case string:
		if err := dst.fromString(value); err != nil {
			return err
		}
	case *string:
		if err := dst.fromString(*value); err != nil {
			return err
		}
	default:
		// Fallback to reflection if an optimised match was not found.
		// The reflection is necessary for arrays and multidimensional slices,
		// but it comes with a 20-50% performance penalty for large arrays/slices
		reflectedValue := reflect.ValueOf(src)
		if !reflectedValue.IsValid() || reflectedValue.IsZero() {
			*dst = Int8Array{Status: Null}
			return nil
		}

		dimensions, elementsLength, ok := findDimensionsFromValue(reflectedValue, nil, 0)
		if !ok {
			return &ValidationError{Type: TypeIntArray, Msg: cannotFindDimensions, Value: src}
		}
		if elementsLength == 0 {
			*dst = Int8Array{Status: Present}
			return nil
		}
		if len(dimensions) == 0 {
			if originalSrc, ok := underlyingSliceType(src); ok {
				return dst.Set(originalSrc)
			}
			return &ValidationError{Type: TypeIntArray, Msg: noConversion, Value: src}
		}

		*dst = Int8Array{
			Elements:   make([]Int8, elementsLength),
			Dimensions: dimensions,
			Status:     Present,
		}
		elementCount, err := dst.setRecursive(reflectedValue, 0, 0)
		if err != nil {
			// Maybe the target was one dimension too far, try again:
			if len(dst.Dimensions) > 1 {
				dst.Dimensions = dst.Dimensions[:len(dst.Dimensions)-1]
				elementsLength = 0
				for _, dim := range dst.Dimensions {
					if elementsLength == 0 {
						elementsLength = int(dim.Length)
					} else {
						elementsLength *= int(dim.Length)
					}
				}
				dst.Elements = make([]Int8, elementsLength)
				elementCount, err = dst.setRecursive(reflectedValue, 0, 0)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if elementCount != len(dst.Elements) {
			return &ValidationError{Type: TypeIntArray, Msg: fmt.Sprintf(expectedElements, len(dst.Elements), elementCount), Value: src}
		}
	}

	return nil
}

func (dst *Int8Array) setRecursive(value reflect.Value, index, dimension int) (int, error) {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if len(dst.Dimensions) == dimension {
			break
		}

		valueLen := value.Len()
		if int32(valueLen) != dst.Dimensions[dimension].Length {
			return 0, &ValidationError{Type: TypeIntArray, Msg: fmt.Sprintf(expectedElementsInDimension, dst.Dimensions[dimension].Length, dimension, valueLen), Value: value}
		}
		for i := 0; i < valueLen; i++ {
			var err error
			index, err = dst.setRecursive(value.Index(i), index, dimension+1)
			if err != nil {
				return 0, err
			}
		}

		return index, nil
	}
	if !value.CanInterface() {
		return 0, &ValidationError{Type: TypeIntArray, Msg: notInterface, Value: value}
	}
	if err := dst.Elements[index].Set(value.Interface()); err != nil {
		return 0, &ValidationError{Type: TypeIntArray, Msg: fmt.Sprintf(cannotSetIndex, index), Value: value.Interface(), Err: err}
	}
	index++

	return index, nil
}

func (dst Int8Array) Get() any {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}
