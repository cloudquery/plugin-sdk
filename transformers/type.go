package transformers

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type TypeTransformer func(reflect.StructField) (arrow.DataType, error)

func defaultGoTypeToSchemaType(v reflect.Type) (arrow.DataType, error) {
	// Non-primitive types
	if v == reflect.TypeOf(net.IP{}) {
		return types.ExtensionTypes.Inet, nil
	}

	k := v.Kind()
	switch k {
	case reflect.Pointer:
		return defaultGoTypeToSchemaType(v.Elem())
	case reflect.String:
		return arrow.BinaryTypes.String, nil
	case reflect.Bool:
		return arrow.FixedWidthTypes.Boolean, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return arrow.PrimitiveTypes.Int64, nil
	case reflect.Float32, reflect.Float64:
		return arrow.PrimitiveTypes.Float64, nil
	case reflect.Map:
		return types.ExtensionTypes.JSON, nil
	case reflect.Struct:
		if v == reflect.TypeOf(time.Time{}) {
			return arrow.FixedWidthTypes.Timestamp_us, nil
		}
		return types.ExtensionTypes.JSON, nil
	case reflect.Slice:
		switch v.Elem().Kind() {
		case reflect.Uint8:
			return arrow.BinaryTypes.Binary, nil
		case reflect.Interface:
			return types.ExtensionTypes.JSON, nil
		default:
			elemValueType, err := defaultGoTypeToSchemaType(v.Elem())
			if err != nil {
				return nil, err
			}
			// if it's already JSON then we don't want to create list of JSON
			if arrow.TypeEqual(elemValueType, types.ExtensionTypes.JSON) {
				return elemValueType, nil
			}
			return arrow.ListOf(elemValueType), nil
		}

	case reflect.Interface:
		return nil, nil // silently ignore

	default:
		return nil, fmt.Errorf("unsupported type: %s", k)
	}
}

func DefaultTypeTransformer(v reflect.StructField) (arrow.DataType, error) {
	return defaultGoTypeToSchemaType(v.Type)
}

var _ TypeTransformer = DefaultTypeTransformer
