package codegen

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

type embeddedStruct struct {
	EmbeddedString string
}

type testStruct struct {
	// IntCol this is an example documentation comment
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
	JSONTAG        *string    `json:"json_tag"`
	SKIPJSONTAG    *string    `json:"-"`
	NOJSONTAG      *string

	*embeddedStruct
}

type testStructWithEmbeddedStruct struct {
	*testStruct
	*embeddedStruct
}

type testStructWithNonEmbeddedStruct struct {
	TestStruct  *testStruct
	NonEmbedded *embeddedStruct
}

type testStructWithCustomType struct {
	TimeCol time.Time `json:"time_col,omitempty"`
}

var expectedColumns = []ColumnDefinition{
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
	{
		Name:     "json_tag",
		Type:     schema.TypeString,
		Resolver: `schema.PathResolver("JSONTAG")`,
	},
	{
		Name:     "nojsontag",
		Type:     schema.TypeString,
		Resolver: `schema.PathResolver("NOJSONTAG")`,
	},
}

var expectedTestTable = TableDefinition{
	Name:            "test_struct",
	Columns:         expectedColumns,
	nameTransformer: DefaultTransformer,
}

var expectedTestTableEmbeddedStruct = TableDefinition{
	Name:            "test_struct",
	Columns:         append(expectedColumns, ColumnDefinition{Name: "embedded_string", Type: schema.TypeString, Resolver: `schema.PathResolver("EmbeddedString")`}),
	nameTransformer: DefaultTransformer,
}

var expectedTestTableNonEmbeddedStruct = TableDefinition{
	Name: "test_struct",
	Columns: ColumnDefinitions{
		// Should not be unwrapped
		ColumnDefinition{Name: "test_struct", Type: schema.TypeJSON, Resolver: `schema.PathResolver("TestStruct")`},
		// Should be unwrapped
		ColumnDefinition{Name: "non_embedded_embedded_string", Type: schema.TypeString, Resolver: `schema.PathResolver("NonEmbedded.EmbeddedString")`},
	},
	nameTransformer: DefaultTransformer,
}

func TestTableFromGoStruct(t *testing.T) {
	type args struct {
		testStruct interface{}
		options    []TableOptions
	}

	tests := []struct {
		name string
		args args
		want TableDefinition
	}{
		{
			name: "should generate table from struct with default options",
			args: args{
				testStruct: testStruct{},
			},
			want: expectedTestTable,
		},
		{
			name: "should unwrap all embedded structs when option is set",
			args: args{
				testStruct: testStructWithEmbeddedStruct{},
				options:    []TableOptions{WithUnwrapAllEmbeddedStructs()},
			},
			want: expectedTestTableEmbeddedStruct,
		},
		{
			name: "should unwrap specific structs when option is set",
			args: args{
				testStruct: testStructWithNonEmbeddedStruct{},
				options:    []TableOptions{WithUnwrapFieldsStructs([]string{"NonEmbedded"})},
			},
			want: expectedTestTableNonEmbeddedStruct,
		},
		{
			name: "should override schema type when option is set",
			args: args{
				testStruct: testStructWithCustomType{},
				options: []TableOptions{WithTypeTransformer(func(t reflect.StructField) (schema.ValueType, error) {
					switch t.Type {
					case reflect.TypeOf(time.Time{}), reflect.TypeOf(&time.Time{}):
						return schema.TypeJSON, nil
					default:
						return schema.TypeInvalid, nil
					}
				})},
			},
			want: TableDefinition{Name: "test_struct",
				// We expect the time column to be of type JSON, since we override the type of `time.Time` to be JSON
				Columns:         ColumnDefinitions{{Name: "time_col", Type: schema.TypeJSON, Resolver: `schema.PathResolver("TimeCol")`}},
				nameTransformer: DefaultTransformer},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := NewTableFromStruct("test_struct", tt.args.testStruct, tt.args.options...)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(table, &tt.want,
				cmpopts.IgnoreUnexported(TableDefinition{})); diff != "" {
				t.Fatalf("table does not match expected. diff (-got, +want): %v", diff)
			}
			buf := bytes.NewBufferString("")
			if err := table.GenerateTemplate(buf); err != nil {
				t.Fatal(err)
			}
			fmt.Println(buf.String())
		})
	}
}

func TestGenerateTemplate(t *testing.T) {
	type args struct {
		table *TableDefinition
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "should add comma between relations",
			args: args{
				table: &TableDefinition{
					Name:      "with relations",
					Relations: []string{"relation1", "relation2"},
				},
			},
		},
		{
			name: "should add ignore_in_tests to columns",
			args: args{
				table: &TableDefinition{
					Name: "with relations",
					Columns: []ColumnDefinition{
						{Name: "ignore_in_tests", Type: schema.TypeString, IgnoreInTests: true},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			err := tt.args.table.GenerateTemplate(buf)
			require.NoError(t, err)
			cupaloy.SnapshotT(t, buf.Bytes())
		})
	}
}
