//nolint:gocritic,revive
package schema

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

type MacaddrArrayTransformer interface {
	TransformMacaddrArray(*MacaddrArray) any
}

type MacaddrArray struct {
	Elements   []Macaddr
	Dimensions []ArrayDimension
	Status     Status
}

func (dst *MacaddrArray) GetStatus() Status {
	return dst.Status
}

func (*MacaddrArray) Type() ValueType {
	return TypeMacAddrArray
}

func (dst *MacaddrArray) Size() int {
	totalSize := 0
	for _, element := range dst.Elements {
		totalSize += element.Size()
	}
	return totalSize
}

func (dst *MacaddrArray) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*MacaddrArray)
	if !ok {
		return false
	}
	if dst.Status != s.Status {
		return false
	}
	if len(dst.Elements) != len(s.Elements) || len(dst.Dimensions) != len(s.Dimensions) {
		return false
	}

	for i := range dst.Elements {
		if !(dst.Elements[i]).Equal(&s.Elements[i]) {
			return false
		}
	}

	return true
}

func (dst *MacaddrArray) fromString(value string) error {
	// this is basically back from string encoding
	if !strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		return &ValidationError{Type: TypeMacAddrArray, Msg: cannotDecodeString, Value: value}
	}

	value = value[1 : len(value)-1]
	strs := strings.Split(value, ",")
	elements := make([]Macaddr, len(strs))
	for i := range strs {
		if err := elements[i].Set(strs[i]); err != nil {
			return err
		}
	}

	*dst = MacaddrArray{
		Elements:   elements,
		Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
		Status:     Present,
	}
	return nil
}

func (dst *MacaddrArray) String() string {
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

func (dst *MacaddrArray) Set(src any) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		*dst = MacaddrArray{Status: Null}
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
	case []net.HardwareAddr:
		if value == nil {
			*dst = MacaddrArray{Status: Null}
		} else if len(value) == 0 {
			*dst = MacaddrArray{Status: Present}
		} else {
			elements := make([]Macaddr, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = MacaddrArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*net.HardwareAddr:
		if value == nil {
			*dst = MacaddrArray{Status: Null}
		} else if len(value) == 0 {
			*dst = MacaddrArray{Status: Present}
		} else {
			elements := make([]Macaddr, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = MacaddrArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []Macaddr:
		if value == nil {
			*dst = MacaddrArray{Status: Null}
		} else if len(value) == 0 {
			*dst = MacaddrArray{Status: Present}
		} else {
			*dst = MacaddrArray{
				Elements:   value,
				Dimensions: []ArrayDimension{{Length: int32(len(value)), LowerBound: 1}},
				Status:     Present,
			}
		}
	case string:
		return dst.fromString(value)
	case *string:
		return dst.fromString(*value)
	default:
		// Fallback to reflection if an optimised match was not found.
		// The reflection is necessary for arrays and multidimensional slices,
		// but it comes with a 20-50% performance penalty for large arrays/slices
		reflectedValue := reflect.ValueOf(src)
		if !reflectedValue.IsValid() || reflectedValue.IsZero() {
			*dst = MacaddrArray{Status: Null}
			return nil
		}

		dimensions, elementsLength, ok := findDimensionsFromValue(reflectedValue, nil, 0)
		if !ok {
			return &ValidationError{Type: TypeMacAddrArray, Msg: cannotFindDimensions, Value: src}
		}
		if elementsLength == 0 {
			*dst = MacaddrArray{Status: Present}
			return nil
		}
		if len(dimensions) == 0 {
			if originalSrc, ok := underlyingSliceType(src); ok {
				return dst.Set(originalSrc)
			}
			return &ValidationError{Type: TypeMacAddrArray, Msg: noConversion, Value: src}
		}

		*dst = MacaddrArray{
			Elements:   make([]Macaddr, elementsLength),
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
				dst.Elements = make([]Macaddr, elementsLength)
				elementCount, err = dst.setRecursive(reflectedValue, 0, 0)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if elementCount != len(dst.Elements) {
			return &ValidationError{Type: TypeMacAddrArray, Msg: fmt.Sprintf(expectedElements, len(dst.Elements), elementCount), Value: src}
		}
	}

	return nil
}

func (dst *MacaddrArray) setRecursive(value reflect.Value, index, dimension int) (int, error) {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if len(dst.Dimensions) == dimension {
			break
		}

		valueLen := value.Len()
		if int32(valueLen) != dst.Dimensions[dimension].Length {
			return 0, &ValidationError{Type: TypeMacAddrArray, Msg: fmt.Sprintf(expectedElementsInDimension, dst.Dimensions[dimension].Length, dimension, valueLen), Value: value.Interface()}
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
		return 0, &ValidationError{Type: TypeMacAddrArray, Msg: notInterface, Value: value}
	}
	if err := dst.Elements[index].Set(value.Interface()); err != nil {
		return 0, &ValidationError{Type: TypeMacAddrArray, Msg: fmt.Sprintf(cannotSetIndex, index), Value: value.Interface(), Err: err}
	}
	index++

	return index, nil
}

func (dst MacaddrArray) Get() any {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}
