package transformers

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v3/types"
)

func (t *structTransformer) getColumnType(field reflect.StructField) (arrow.DataType, error) {
	columnType, err := t.typeTransformer(field)
	if err != nil {
		return nil, err
	}

	if columnType != nil {
		return columnType, nil
	}
	return DefaultTypeTransformer(field)
}

type TypeTransformer func(reflect.StructField) (arrow.DataType, error)

func DefaultTypeTransformer(v reflect.StructField) (arrow.DataType, error) {
	return defaultGoTypeToSchemaType(v.Type)
}

func defaultGoTypeToSchemaType(v reflect.Type) (arrow.DataType, error) {
	// Non primitive types
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
		if v.Elem().Kind() == reflect.Uint8 {
			return arrow.BinaryTypes.Binary, nil
		}
		elemValueType, err := defaultGoTypeToSchemaType(v.Elem())
		if err != nil {
			return nil, err
		}
		// if it's already JSON then we don't want to create list of JSON
		if arrow.TypeEqual(elemValueType, types.ExtensionTypes.JSON) {
			return elemValueType, nil
		}
		return arrow.ListOf(elemValueType), nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", k)
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
