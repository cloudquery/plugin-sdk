package transformers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudquery/plugin-sdk/v3/caser"
	"github.com/cloudquery/plugin-sdk/v3/schema"
)

func (t *structTransformer) getFieldNamePath(field reflect.StructField, parent *reflect.StructField) (name, path string, err error) {
	path = field.Name
	name, err = t.nameTransformer(field)
	if err != nil {
		return "", "", fmt.Errorf("failed to transform field name for field %s: %w", field.Name, err)
	}
	// skip field if there is no name
	if name == "" {
		return "", "", nil
	}

	if parent == nil {
		return name, path, nil
	}

	parentName, err := t.nameTransformer(*parent)
	if err != nil {
		return "", "", fmt.Errorf("failed to transform field name for parent field %s: %w", parent.Name, err)
	}

	return parentName + "_" + name, parent.Name + `.` + path, nil
}

type NameTransformer func(reflect.StructField) (string, error)

var defaultCaser = caser.New()

func DefaultNameTransformer(field reflect.StructField) (string, error) {
	name := field.Name
	if jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]; len(jsonTag) > 0 {
		// return empty string if the field is not related api response
		if jsonTag == "-" {
			return "", nil
		}
		if nameFromJSONTag := defaultCaser.ToSnake(jsonTag); schema.ValidColumnName(nameFromJSONTag) {
			return nameFromJSONTag, nil
		}
	}
	return defaultCaser.ToSnake(name), nil
}
