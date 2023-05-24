package transformers

import (
	"fmt"
	"reflect"
)

type NameTransformer func(reflect.StructField) (string, error)

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
