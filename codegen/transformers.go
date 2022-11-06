package codegen

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cloudquery/plugin-sdk/v2/caser"
	"github.com/cloudquery/plugin-sdk/v2/schema"
)

type NameTransformer func(reflect.StructField) (string, error)

var defaultCaser = caser.New()

func DefaultNameTransformer(field reflect.StructField) (string, error) {
	name := field.Name
	if jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]; len(jsonTag) > 0 {
		// return empty string if the field is not related api response
		if jsonTag == "-" {
			return "", nil
		}
		name = jsonTag
	}
	return defaultCaser.ToSnake(name), nil
}

type TypeTransformer func(reflect.StructField) (schema.ValueType, error)

func DefaultTypeTransformer(v reflect.StructField) (schema.ValueType, error) {
	return defaultGoTypeToSchemaType(v.Type)
}

func defaultGoTypeToSchemaType(v reflect.Type) (schema.ValueType, error) {
	k := v.Kind()
	switch k {
	case reflect.Pointer:
		return defaultGoTypeToSchemaType(v.Elem())
	case reflect.String:
		return schema.TypeString, nil
	case reflect.Bool:
		return schema.TypeBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return schema.TypeInt, nil
	case reflect.Float32, reflect.Float64:
		return schema.TypeFloat, nil
	case reflect.Map:
		return schema.TypeJSON, nil
	case reflect.Struct:
		if v == reflect.TypeOf(time.Time{}) {
			return schema.TypeTimestamp, nil
		}
		return schema.TypeJSON, nil
	case reflect.Slice:
		switch v.Elem().Kind() {
		case reflect.String:
			return schema.TypeStringArray, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return schema.TypeIntArray, nil
		default:
			return schema.TypeJSON, nil
		}
	default:
		return schema.TypeInvalid, fmt.Errorf("unsupported type: %s", k)
	}
}

type ResolverTransformer func(field reflect.StructField, path string) (string, error)

func DefaultResolverTransformer(_ reflect.StructField, path string) (string, error) {
	return defaultResolver(path), nil
}

func defaultResolver(path string) string {
	return `schema.PathResolver("` + path + `")`
}
