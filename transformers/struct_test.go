package transformers

import (
	"net"
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
		Int64Col  int64   `json:"int64_col,omitempty"`
		StringCol string  `json:"string_col,omitempty"`
		FloatCol  float64 `json:"float_col,omitempty"`
		BoolCol   bool    `json:"bool_col,omitempty"`
		JSONCol   struct {
			IntCol    int    `json:"int_col,omitempty"`
			StringCol string `json:"string_col,omitempty"`
		}
		IntArrayCol        []int  `json:"int_array_col,omitempty"`
		IntPointerArrayCol []*int `json:"int_pointer_array_col,omitempty"`

		StringArrayCol        []string  `json:"string_array_col,omitempty"`
		StringPointerArrayCol []*string `json:"string_pointer_array_col,omitempty"`

		InetCol        net.IP  `json:"inet_col,omitempty"`
		InetPointerCol *net.IP `json:"inet_pointer_col,omitempty"`

		ByteArrayCol []byte `json:"byte_array_col,omitempty"`

		TimeCol        time.Time  `json:"time_col,omitempty"`
		TimePointerCol *time.Time `json:"time_pointer_col,omitempty"`
		JSONTag        *string    `json:"json_tag"`
		SkipJSONTag    *string    `json:"-"`
		NoJSONTag      *string
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

	testSliceStruct []struct {
		IntCol int
	}
)

var (
	expectedColumns = []schema.Column{
		{
			Name: "int_col",
			Type: schema.TypeInt,
		},
		{
			Name: "int64_col",
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
			Name: "int_pointer_array_col",
			Type: schema.TypeIntArray,
		},
		{
			Name: "string_array_col",
			Type: schema.TypeStringArray,
		},
		{
			Name: "string_pointer_array_col",
			Type: schema.TypeStringArray,
		},
		{
			Name: "inet_col",
			Type: schema.TypeInet,
		},
		{
			Name: "inet_pointer_col",
			Type: schema.TypeInet,
		},
		{
			Name: "byte_array_col",
			Type: schema.TypeByteArray,
		},
		{
			Name: "time_col",
			Type: schema.TypeTimestamp,
		},
		{
			Name: "time_pointer_col",
			Type: schema.TypeTimestamp,
		},
		{
			Name: "json_tag",
			Type: schema.TypeString,
		},
		{
			Name: "no_json_tag",
			Type: schema.TypeString,
		},
	}
	expectedTestTable = schema.Table{
		Name:    "test_struct",
		Columns: expectedColumns,
	}
	expectedTestTableEmbeddedStruct = schema.Table{
		Name: "test_struct",
		Columns: append(
			expectedColumns, schema.Column{
				Name: "embedded_string",
				Type: schema.TypeString,
			}),
	}
	expectedTestTableNonEmbeddedStruct = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			// Should not be unwrapped
			schema.Column{Name: "test_struct", Type: schema.TypeJSON},
			// Should be unwrapped
			schema.Column{
				Name: "non_embedded_embedded_string",
				Type: schema.TypeString,
			},
		},
	}
	expectedTestTableStructForCustomResolvers = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			{
				Name: "time_col",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "custom",
				Type: schema.TypeTimestamp,
			},
		},
	}
	expectedTestSliceStruct = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			{
				Name: "int_col",
				Type: schema.TypeInt,
			},
		},
	}
)

func TestTableFromGoStruct(t *testing.T) {
	type args struct {
		testStruct any
		options    []StructTransformerOption
	}

	tests := []struct {
		name    string
		args    args
		want    schema.Table
		wantErr bool
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
				options: []StructTransformerOption{
					WithUnwrapAllEmbeddedStructs(),
				},
			},
			want: expectedTestTableEmbeddedStruct,
		},
		{
			name: "should unwrap specific structs when option is set",
			args: args{
				testStruct: testStructWithNonEmbeddedStruct{},
				options: []StructTransformerOption{
					WithUnwrapStructFields("NonEmbedded"),
				},
			},
			want: expectedTestTableNonEmbeddedStruct,
		},
		{
			name: "should generate table from slice struct",
			args: args{
				testStruct: testSliceStruct{},
			},
			want: expectedTestSliceStruct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := schema.Table{
				Name:    "test",
				Columns: schema.ColumnList{},
			}
			transformer := TransformWithStruct(tt.args.testStruct, tt.args.options...)
			if err := transformer(&table); err != nil {
				t.Fatal(err)
			}
			// table, err := NewTableFromStruct("test_struct", tt.args.testStruct, tt.args.options...)
			// if (err != nil) != tt.wantErr {
			// 	t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			// }
			// if tt.wantErr {
			// 	return
			// }
			if diff := cmp.Diff(table.Columns, tt.want.Columns,
				cmpopts.IgnoreFields(schema.Column{}, "Resolver")); diff != "" {
				t.Fatalf("table does not match expected. diff (-got, +want): %v", diff)
			}
		})
	}
}
