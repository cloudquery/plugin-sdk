package plugins

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v1/helpers"
	"github.com/cloudquery/plugin-sdk/v1/schema"
)

type DefaultReverseTransformer struct {
}

// DefaultReverseTransformer tries best effort to convert a slice of values to CQTypes
// based on the provided table columns.
func (*DefaultReverseTransformer) ReverseTransformValues(table *schema.Table, values []interface{}) (schema.CQTypes, error) {
	valuesSlice := helpers.InterfaceSlice(values)
	res := make(schema.CQTypes, len(valuesSlice))

	for i, v := range valuesSlice {
		t := schema.NewCqTypeFromValueType(table.Columns[i].Type)
		if err := t.Set(v); err != nil {
			return nil, fmt.Errorf("failed to convert value %v to type %s: %w", v, table.Columns[i].Type, err)
		}
		res[i] = t
	}
	return res, nil
}
