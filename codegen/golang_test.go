package codegen

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
)

type testStruct struct {
	IntCol    int     `json:"int_col,omitempty"`
	StringCol string  `json:"string_col,omitempty"`
	FloatCol  float64 `json:"float_col,omitempty"`
	BoolCol   bool    `json:"bool_col,omitempty"`
	JSONCol   struct {
		IntCol    int    `json:"int_col,omitempty"`
		StringCol string `json:"string_col,omitempty"`
	}
	IntArrayCol    []int      `json:"int_array_col,omitempty"`
	StringArrayCol []string   `json:"string_array_col,omitempty"`
	TimeCol        time.Time  `json:"time_col,omitempty"`
	TimePointerCol *time.Time `json:"time_pointer_col,omitempty"`
}

var expectedTestTable = TableDefinition{
	Name: "test_struct",
	Columns: []ColumnDefinition{
		{
			Name:     "int_col",
			Type:     schema.TypeInt,
			Resolver: `schema.PathResolver("IntCol")`,
		},
		{
			Name:     "string_col",
			Type:     schema.TypeString,
			Resolver: `schema.PathResolver("StringCol")`,
		},
		{
			Name:     "float_col",
			Type:     schema.TypeFloat,
			Resolver: `schema.PathResolver("FloatCol")`,
		},
		{
			Name:     "bool_col",
			Type:     schema.TypeBool,
			Resolver: `schema.PathResolver("BoolCol")`,
		},
		{
			Name:     "json_col",
			Type:     schema.TypeJSON,
			Resolver: `schema.PathResolver("JSONCol")`,
		},
		{
			Name:     "int_array_col",
			Type:     schema.TypeIntArray,
			Resolver: `schema.PathResolver("IntArrayCol")`,
		},
		{
			Name:     "string_array_col",
			Type:     schema.TypeStringArray,
			Resolver: `schema.PathResolver("StringArrayCol")`,
		},
		{
			Name:     "time_col",
			Type:     schema.TypeTimestamp,
			Resolver: `schema.PathResolver("TimeCol")`,
		},
		{
			Name:     "time_pointer_col",
			Type:     schema.TypeTimestamp,
			Resolver: `schema.PathResolver("TimePointerCol")`,
		},
	},
	nameTransformer: defaultTransformer,
}

func TestTableFromGoStruct(t *testing.T) {
	table, err := NewTableFromStruct("test_struct", testStruct{})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(table, &expectedTestTable,
		cmpopts.IgnoreUnexported(TableDefinition{})); diff != "" {
		t.Fatalf("table does not match expected. diff (-got, +want): %v", diff)
	}
	buf := bytes.NewBufferString("")
	if err := table.GenerateTemplate(buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
}
