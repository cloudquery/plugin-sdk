package codegen

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
)

type testStruct struct {
	IntCol    int     `json:"int_col,omitempty"`
	StringCol string  `json:"string_col,omitempty"`
	FloatCol  float64 `json:"float_col,omitempty"`
	BoolCol   bool    `json:"bool_col,omitempty"`
	JsonCol   struct {
		IntCol    int    `json:"int_col,omitempty"`
		StringCol string `json:"string_col,omitempty"`
	}
	IntArrayCol    []int    `json:"int_array_col,omitempty"`
	StringArrayCol []string `json:"string_array_col,omitempty"`
}

var expectedTestTable = TableDefinition{
	Name: "test_struct",
	Columns: []ColumnDefinition{
		{
			Name: "int_col",
			Type: schema.TypeInt,
		},
		{
			Name: "string_col",
			Type: schema.TypeString,
		},
		{
			Name: "float_col",
			Type: schema.TypeFloat,
		},
		{
			Name: "bool_col",
			Type: schema.TypeBool,
		},
		{
			Name: "json_col",
			Type: schema.TypeJSON,
		},
		{
			Name: "int_array_col",
			Type: schema.TypeIntArray,
		},
		{
			Name: "string_array_col",
			Type: schema.TypeStringArray,
		},
	},
	nameTransformer: defaultTransformer,
}

func TestTableFromGoStruct(t *testing.T) {
	table, err := NewTableFromStruct("test_struct", testStruct{})
	if err != nil {
		t.Fatal(err)
	}
	// if !reflect.DeepEqual(table, &expectedTestTable) {
	// 	t.Fatalf("expected:\n%+v, got:\n%+v", expectedTestTable, table)
	// }
	buf := bytes.NewBufferString("")
	if err := table.GenerateTemplate(buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
}
