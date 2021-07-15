package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	stringJson = "{\"test\":true}"
	resources  = []Resource{
		{
			data: map[string]interface{}{
				"test": stringJson,
				"meta": make(map[string]string),
			},
			table: &Table{
				Name: "test_table_validator",
				Columns: []Column{
					{
						Name: "test",
						Type: TypeJSON,
					},
				},
			}},
		{
			data: map[string]interface{}{
				"test": &stringJson,
				"meta": make(map[string]string),
			},
			table: &Table{
				Name: "test_table_validator",
				Columns: []Column{
					{
						Name: "test",
						Type: TypeJSON,
					},
				},
			}},
		{
			data: map[string]interface{}{
				"test": map[string]interface{}{
					"test": 1,
					"test1": map[string]interface{}{
						"test": 1,
					},
				},
				"meta": make(map[string]string),
			},
			table: &Table{
				Name: "test_table_validator",
				Columns: []Column{
					{
						Name: "test",
						Type: TypeJSON,
					},
				},
			}},
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
				"meta": make(map[string]string),
			},
			table: &Table{
				Name: "test_table_validator",
				Columns: []Column{
					{
						Name: "test",
						Type: TypeJSON,
					},
				},
			}},
	}
)

func TestJsonColumn(t *testing.T) {
	for _, r := range resources {
		_, err := getResourceValues(&r)
		assert.Nil(t, err)
	}
}
