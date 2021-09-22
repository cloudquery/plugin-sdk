package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type jsonTestType struct {
	Name        string `json:"name"`
	Description string `json:"decription"`
	Version     int    `json:"version"`
}

type jsonNoTags struct {
	Name        string
	Description string
	Version     int
}

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
		{
			data: map[string]interface{}{
				"test": "{\"hello\":123}",
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": jsonTestType{
					Name:        "test",
					Description: "test1",
					Version:     10,
				},
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": jsonNoTags{
					Name:        "test",
					Description: "test1",
					Version:     10,
				},
			},
			table: &jsonTestTable,
		},
	}

	failResources = []Resource{
		{
			data: map[string]interface{}{
				"test": true,
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": 10.1,
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": "true_test",
			},
			table: &jsonTestTable,
		},
		{
			data: map[string]interface{}{
				"test": "{\"hello\":123}1",
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

	for _, r := range failResources {
		_, err := getResourceValues(&r)
		assert.Error(t, err)
	}
}
