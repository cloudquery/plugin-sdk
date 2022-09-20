package codegen

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type (
	embeddedStruct struct {
		EmbeddedString string
	}
	testStruct struct {
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
	testStructWithEmbeddedStruct struct {
		*testStruct
		*embeddedStruct
	}
	testStructWithNonEmbeddedStruct struct {
		TestStruct  *testStruct
		NonEmbedded *embeddedStruct
	}
	testStructForCustomResolvers struct {
		IntCol         int        `json:"int_col,omitempty"`
		StringCol      string     `json:"string_col,omitempty"`
		TimeCol        time.Time  `json:"time_col,omitempty"`
		TimePointerCol *time.Time `json:"time_pointer_col,omitempty"`
	}
)

var (
	expectedColumns = []ColumnDefinition{
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
	expectedTestTable = TableDefinition{
		Name:            "test_struct",
		Columns:         expectedColumns,
		nameTransformer: defaultTransformer,
	}
	expectedTestTableEmbeddedStruct = TableDefinition{
		Name:            "test_struct",
		Columns:         append(expectedColumns, ColumnDefinition{Name: "embedded_string", Type: schema.TypeString, Resolver: `schema.PathResolver("EmbeddedString")`}),
		nameTransformer: defaultTransformer,
	}
	expectedTestTableNonEmbeddedStruct = TableDefinition{
		Name: "test_struct",
		Columns: ColumnDefinitions{
			// Should not be unwrapped
			ColumnDefinition{Name: "test_struct", Type: schema.TypeJSON, Resolver: `schema.PathResolver("TestStruct")`},
			// Should be unwrapped
			ColumnDefinition{Name: "non_embedded_embedded_string", Type: schema.TypeString, Resolver: `schema.PathResolver("NonEmbedded.EmbeddedString")`},
		},
		nameTransformer: defaultTransformer,
	}
	expectedTestTableStructForCustomResolvers = TableDefinition{
		Name: "test_struct",
		Columns: ColumnDefinitions{
			{
				Name:     "int_col",
				Type:     schema.TypeInt,
				Resolver: `schema.PathResolver("IntCol")`,
			},
			{
				Name:     "string_col",
				Type:     schema.TypeString,
				Resolver: `stringResolver("StringCol")`,
			},
			{
				Name:     "time_col",
				Type:     schema.TypeTimestamp,
				Resolver: `timestampResolver("TimeCol")`,
			},
			{
				Name:     "time_pointer_col",
				Type:     schema.TypeTimestamp,
				Resolver: `timestampResolver("TimePointerCol")`,
			},
		},
		nameTransformer: defaultTransformer,
		customTypeResolvers: map[schema.ValueType]string{
			schema.TypeString:    `stringResolver`,
			schema.TypeTimestamp: `timestampResolver`,
		},
	}
)

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
			name: "should_unwrap_specific_structs_when_option_is_set",
			args: args{
				testStruct: testStructWithNonEmbeddedStruct{},
				options:    []TableOptions{WithUnwrapFieldsStructs([]string{"NonEmbedded"})},
			},
			want: expectedTestTableNonEmbeddedStruct,
		},
		{
			name: "should use custom resolvers",
			args: args{
				testStruct: testStructForCustomResolvers{},
				options: []TableOptions{
					WithCustomValueTypeResolver(schema.TypeString, `stringResolver`),
					WithCustomValueTypeResolver(schema.TypeTimestamp, `timestampResolver`),
				},
			},
			want: expectedTestTableStructForCustomResolvers,
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
