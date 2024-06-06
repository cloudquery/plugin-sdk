package transformers

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type (
	TypeTransformer func(reflect.StructField) (arrow.DataType, error)
	transformed     = map[reflect.Type]struct{}
)

var (
	jsonObjectType = reflect.TypeOf((map[string]any)(nil))
	netIPType      = reflect.TypeOf(net.IP{})
	timeType       = reflect.TypeOf(time.Time{})
	durationType   = reflect.TypeOf(time.Duration(0))
)

func (t *structTransformer) transformType(typ reflect.Type) (arrow.DataType, error) {
	return t.typeTransformer(typ, make(transformed))
}

func (t *structTransformer) typeTransformer(typ reflect.Type, visited transformed) (arrow.DataType, error) {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if t.customTypeTransformer != nil {
		dt, err := t.customTypeTransformer(reflect.StructField{Type: typ})
		if dt != nil || err != nil {
			return dt, err
		}
	}

	// Non-primitive types
	switch typ {
	case netIPType:
		return types.ExtensionTypes.Inet, nil
	case timeType:
		return arrow.FixedWidthTypes.Timestamp_us, nil
	case durationType:
		return arrow.FixedWidthTypes.Duration_us, nil
	case jsonObjectType:
		return types.ExtensionTypes.JSON, nil
	}

	switch k := typ.Kind(); k {
	case reflect.Map:
		return t.mapTypeTransformer(typ, visited)
	case reflect.Struct:
		return t.structTypeTransformer(typ, visited)
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return arrow.BinaryTypes.Binary, nil // []byte
		}

		elem, err := t.typeTransformer(typ.Elem(), visited)
		if err != nil {
			return nil, err
		}
		if elem == nil {
			return nil, nil // []ignored -> ignored
		}

		// if it's a JSON then we don't want to create list of JSON
		if arrow.TypeEqual(elem, types.ExtensionTypes.JSON) {
			return types.ExtensionTypes.JSON, nil
		}
		return arrow.ListOf(elem), nil
	case reflect.Interface:
		return types.ExtensionTypes.JSON, nil
	case reflect.String:
		return arrow.BinaryTypes.String, nil
	case reflect.Bool:
		return arrow.FixedWidthTypes.Boolean, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return arrow.PrimitiveTypes.Int64, nil
	case reflect.Float32, reflect.Float64:
		return arrow.PrimitiveTypes.Float64, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", k)
	}
}

func (t *structTransformer) mapTypeTransformer(typ reflect.Type, visited transformed) (arrow.DataType, error) {
	key, err := t.typeTransformer(typ.Key(), visited)
	if err != nil {
		return nil, err
	}
	item, err := t.typeTransformer(typ.Elem(), visited)
	if err != nil {
		return nil, err
	}
	return arrow.MapOf(key, item), nil
}

func (t *structTransformer) structTypeTransformer(typ reflect.Type, visited transformed) (arrow.DataType, error) {
	if _, ok := visited[typ]; ok {
		return types.ExtensionTypes.JSON, nil
	}
	visited[typ] = struct{}{}

	fields := make([]arrow.Field, 0, typ.NumField()) // at most, we have NumField fields (but need to check for unexported)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if t.ignoreField(field) {
			continue
		}

		name, err := t.nameTransformer(field)
		if err != nil {
			return nil, fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
		}
		// skip field if there is no name
		if name == "" {
			continue
		}

		dt, err := t.typeTransformer(field.Type, visited)
		if err != nil {
			return nil, fmt.Errorf("failed to transform %q struct field: %w", field.Name, err)
		}
		if dt == nil {
			// ignored type
			continue
		}

		fields = append(fields, arrow.Field{Name: name, Type: dt})
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields to propagate to struct")
	}

	return arrow.StructOf(fields...), nil
}

// DefaultTypeTransformer is a nop implementation
func DefaultTypeTransformer(reflect.StructField) (arrow.DataType, error) { return nil, nil }

var _ TypeTransformer = DefaultTypeTransformer
