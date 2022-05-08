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
	intTestTable = Table{
		Name: "test_table_validator",
		Columns: []Column{
			{
				Name: "int32",
				Type: TypeInt,
			},
			{
				Name: "int64",
				Type: TypeInt,
			},
		},
	}
	resources = []Resource{
		{
			data: map[string]interface{}{
				"test":    stringJson,
				"cq_meta": make(map[string]string),
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

	intResources = []Resource{
		{
			data: map[string]interface{}{
				"int32": 123,
				"int64": int64(123),
			},
			table: &intTestTable,
		},
		{
			data: map[string]interface{}{
				"int32": 123,
				"int64": int64(9223372036854775807),
			},
			table: &intTestTable,
		},
	}
)

func TestJsonColumn(t *testing.T) {
	for _, r := range resources {
		_, err := PostgresDialect{}.GetResourceValues(&r)
		assert.Nil(t, err)
	}

	for _, r := range failResources {
		_, err := PostgresDialect{}.GetResourceValues(&r)
		assert.Error(t, err)
	}
}

func TestIntColumn(t *testing.T) {
	for _, r := range intResources {
		_, err := PostgresDialect{}.GetResourceValues(&r)
		assert.Nil(t, err)
	}
}
