package plugins

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
)

type DefaultReverseTransformer struct {
}

// DefaultReverseTransformer tries best effort to convert a slice of values to CQTypes
// based on the provided table columns.
func (*DefaultReverseTransformer) ReverseTransformValues(table *schema.Table, values []interface{}) (schema.CQTypes, error) {
	res := make(schema.CQTypes, len(values))

	for i, v := range values {
		t := schema.NewCqTypeFromValueType(table.Columns[i].Type)
		if err := t.Set(v); err != nil {
			return nil, fmt.Errorf("failed to convert value %v to type %s: %w", v, table.Columns[i].Type, err)
		}
		res[i] = t
	}
	return res, nil
}
