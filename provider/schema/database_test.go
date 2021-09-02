package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	stringJson    = "{\"test\":true}"
	jsonTestTable = Table{
		Name: "test_table_validator",
		Columns: []Column{
			{
				Name: "test",
				Type: TypeJSON,
			},
		},
	}
	resources = []Resource{
		{
			data: map[string]interface{}{
				"test": stringJson,
				"meta": make(map[string]string),
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": &stringJson,
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": map[string]interface{}{
					"test": 1,
					"test1": map[string]interface{}{
						"test": 1,
					},
				},
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": []interface{}{
					map[string]interface{}{
						"test":  1,
						"test1": true,
					},
					map[string]interface{}{
						"test":  1,
						"test1": true,
					},
				},
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": nil,
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": []interface{}{
					nil,
				},
			},
			table: &jsonTestTable,
		},
	}
)

func TestJsonColumn(t *testing.T) {
	for _, r := range resources {
		_, err := getResourceValues(&r)
		assert.Nil(t, err)
	}
}
