//nolint:dupl,revive,gocritic
package schema

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

type CIDRArrayTransformer interface {
	TransformCIDRArray(*CIDRArray) interface{}
}

type CIDRArray struct {
	Elements   []CIDR
	Dimensions []ArrayDimension
	Status     Status
}

func (*CIDRArray) Type() ValueType {
	return TypeCIDRArray
}

func (dst *CIDRArray) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*CIDRArray)
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

func (dst *CIDRArray) String() string {
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

func (dst *CIDRArray) Set(src interface{}) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		*dst = CIDRArray{Status: Null}
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
	case []*net.IPNet:
		if value == nil {
			*dst = CIDRArray{Status: Null}
		} else if len(value) == 0 {
			*dst = CIDRArray{Status: Present}
		} else {
			elements := make([]CIDR, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = CIDRArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []net.IP:
		if value == nil {
			*dst = CIDRArray{Status: Null}
		} else if len(value) == 0 {
			*dst = CIDRArray{Status: Present}
		} else {
			elements := make([]CIDR, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = CIDRArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []*net.IP:
		if value == nil {
			*dst = CIDRArray{Status: Null}
		} else if len(value) == 0 {
			*dst = CIDRArray{Status: Present}
		} else {
			elements := make([]CIDR, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = CIDRArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []CIDR:
		if value == nil {
			*dst = CIDRArray{Status: Null}
		} else if len(value) == 0 {
			*dst = CIDRArray{Status: Present}
		} else {
			*dst = CIDRArray{
				Elements:   value,
				Dimensions: []ArrayDimension{{Length: int32(len(value)), LowerBound: 1}},
				Status:     Present,
			}
		}
	default:
		// Fallback to reflection if an optimised match was not found.
		// The reflection is necessary for arrays and multidimensional slices,
		// but it comes with a 20-50% performance penalty for large arrays/slices
		reflectedValue := reflect.ValueOf(src)
		if !reflectedValue.IsValid() || reflectedValue.IsZero() {
			*dst = CIDRArray{Status: Null}
			return nil
		}

		dimensions, elementsLength, ok := findDimensionsFromValue(reflectedValue, nil, 0)
		if !ok {
			return fmt.Errorf("cannot find dimensions of %v for CIDRArray", src)
		}
		if elementsLength == 0 {
			*dst = CIDRArray{Status: Present}
			return nil
		}
		if len(dimensions) == 0 {
			if originalSrc, ok := underlyingSliceType(src); ok {
				return dst.Set(originalSrc)
			}
			return fmt.Errorf("cannot convert %v to CIDRArray", src)
		}

		*dst = CIDRArray{
			Elements:   make([]CIDR, elementsLength),
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
				dst.Elements = make([]CIDR, elementsLength)
				elementCount, err = dst.setRecursive(reflectedValue, 0, 0)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if elementCount != len(dst.Elements) {
			return fmt.Errorf("cannot convert %v to CIDRArray, expected %d dst.Elements, but got %d instead", src, len(dst.Elements), elementCount)
		}
	}

	return nil
}

func (dst *CIDRArray) setRecursive(value reflect.Value, index, dimension int) (int, error) {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if len(dst.Dimensions) == dimension {
			break
		}

		valueLen := value.Len()
		if int32(valueLen) != dst.Dimensions[dimension].Length {
			return 0, fmt.Errorf("multidimensional arrays must have array expressions with matching dimensions")
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
		return 0, fmt.Errorf("cannot convert all values to CIDRArray")
	}
	if err := dst.Elements[index].Set(value.Interface()); err != nil {
		return 0, fmt.Errorf("%v in CIDRArray", err)
	}
	index++

	return index, nil
}

func (dst CIDRArray) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}
