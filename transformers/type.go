package transformers

import (
	"reflect"

	"github.com/apache/arrow/go/v13/arrow"
)

type TypeTransformer func(reflect.StructField) (arrow.DataType, error)

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
