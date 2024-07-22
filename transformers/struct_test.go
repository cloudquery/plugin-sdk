package transformers

import (
	"net"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/google/go-cmp/cmp"
)

type (
	embeddedStruct struct {
		EmbeddedString string
		IntCol         int `json:"int_col,omitempty"`
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

		AnyArrayCol []any `json:"any_array_col,omitempty"`

		TimeCol        time.Time  `json:"time_col,omitempty"`
		TimePointerCol *time.Time `json:"time_pointer_col,omitempty"`
		JSONTag        *string    `json:"json_tag"`
		SkipJSONTag    *string    `json:"-"`
		NoJSONTag      *string
		*embeddedStruct
	}
	testStructWithEmbeddedStruct struct {
		IntCol int `json:"int_col,omitempty"`
		*testStruct
		*embeddedStruct
	}
	testStructWithNonEmbeddedStruct struct {
		IntCol      int `json:"int_col,omitempty"`
		TestStruct  *testStruct
		NonEmbedded *embeddedStruct
	}

	testSliceStruct []struct {
		IntCol int
	}

	testPKStruct struct {
		Parent  string `json:"parent"`
		Name    string `json:"name"`
		Version int    `json:"version"`
	}

	testFunnyStruct struct {
		AFunnyLookingField      string `json:"OS-EXT:a-funny-looking-field"`
		AFieldWithCamelCaseName string `json:"camelCaseName"`
	}

	testStructWithAny struct {
		IntCol     int `json:"int_col"`
		Properties any
	}
)

var (
	expectedColumns = []schema.Column{
		{
			Name: "int_col",
			Type: arrow.PrimitiveTypes.Int64,
		},
		{
			Name: "int64_col",
			Type: arrow.PrimitiveTypes.Int64,
		},
		{
			Name: "string_col",
			Type: arrow.BinaryTypes.String,
		},
		{
			Name: "float_col",
			Type: arrow.PrimitiveTypes.Float64,
		},
		{
			Name: "bool_col",
			Type: arrow.FixedWidthTypes.Boolean,
		},
		{
			Name: "json_col",
			Type: types.ExtensionTypes.JSON,
		},
		{
			Name: "int_array_col",
			Type: arrow.ListOf(arrow.PrimitiveTypes.Int64),
		},
		{
			Name: "int_pointer_array_col",
			Type: arrow.ListOf(arrow.PrimitiveTypes.Int64),
		},
		{
			Name: "string_array_col",
			Type: arrow.ListOf(arrow.BinaryTypes.String),
		},
		{
			Name: "string_pointer_array_col",
			Type: arrow.ListOf(arrow.BinaryTypes.String),
		},
		{
			Name: "inet_col",
			Type: types.ExtensionTypes.Inet,
		},
		{
			Name: "inet_pointer_col",
			Type: types.ExtensionTypes.Inet,
		},
		{
			Name: "byte_array_col",
			Type: arrow.BinaryTypes.Binary,
		},
		{
			Name: "any_array_col",
			Type: types.ExtensionTypes.JSON,
		},
		{
			Name: "time_col",
			Type: arrow.FixedWidthTypes.Timestamp_us,
		},
		{
			Name: "time_pointer_col",
			Type: arrow.FixedWidthTypes.Timestamp_us,
		},
		{
			Name: "json_tag",
			Type: arrow.BinaryTypes.String,
		},
		{
			Name: "no_json_tag",
			Type: arrow.BinaryTypes.String,
		},
	}
	expectedTestTable = schema.Table{
		Name:    "test_struct",
		Columns: expectedColumns,
	}
	expectedTestTableEmbeddedStruct = schema.Table{
		Name:    "test_struct",
		Columns: append(expectedColumns, schema.Column{Name: "embedded_string", Type: arrow.BinaryTypes.String}),
	}
	expectedTestTableEmbeddedStructWithTopLevelPK = schema.Table{
		Name: "test_struct",
		Columns: func(base schema.ColumnList) schema.ColumnList {
			cols := slices.Clone(base)
			cols = append(cols, schema.Column{Name: "embedded_string", Type: arrow.BinaryTypes.String})
			cols[cols.Index("int_col")].PrimaryKey = true
			return cols
		}(expectedColumns),
	}
	expectedTestTableEmbeddedStructWithUnwrappedPK = schema.Table{
		Name: "test_struct",
		Columns: append(
			expectedColumns, schema.Column{
				Name:       "embedded_string",
				Type:       arrow.BinaryTypes.String,
				PrimaryKey: true,
			}),
	}
	expectedTestTableNonEmbeddedStruct = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			schema.Column{Name: "int_col", Type: arrow.PrimitiveTypes.Int64},
			// Should not be unwrapped
			schema.Column{Name: "test_struct", Type: types.ExtensionTypes.JSON},
			// Should be unwrapped
			schema.Column{Name: "non_embedded_embedded_string", Type: arrow.BinaryTypes.String},
			schema.Column{Name: "non_embedded_int_col", Type: arrow.PrimitiveTypes.Int64},
		},
	}
	expectedTestTableNonEmbeddedStructWithTopLevelPK = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			schema.Column{
				Name:       "int_col",
				Type:       arrow.PrimitiveTypes.Int64,
				PrimaryKey: true,
			},
			// Should not be unwrapped
			schema.Column{Name: "test_struct", Type: types.ExtensionTypes.JSON},
			// Should be unwrapped
			schema.Column{
				Name: "non_embedded_embedded_string",
				Type: arrow.BinaryTypes.String,
			},
			schema.Column{Name: "non_embedded_int_col", Type: arrow.PrimitiveTypes.Int64},
		},
	}
	expectedTestTableNonEmbeddedStructWithUnwrappedPK = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			// shouldn't be PK
			schema.Column{Name: "int_col", Type: arrow.PrimitiveTypes.Int64},
			// Should not be unwrapped
			schema.Column{Name: "test_struct", Type: types.ExtensionTypes.JSON},
			// Should be unwrapped
			schema.Column{
				Name: "non_embedded_embedded_string",
				Type: arrow.BinaryTypes.String,
			},
			// should be PK
			schema.Column{
				Name:       "non_embedded_int_col",
				Type:       arrow.PrimitiveTypes.Int64,
				PrimaryKey: true,
			},
		},
	}
	expectedTestSliceStruct = schema.Table{
		Name: "test_struct",
		Columns: schema.ColumnList{
			{
				Name: "int_col",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}

	expectedTableWithPKs = schema.Table{
		Name: "test_pk_struct",
		Columns: schema.ColumnList{
			{
				Name:       "parent",
				Type:       arrow.BinaryTypes.String,
				PrimaryKey: true,
			},
			{
				Name:       "name",
				Type:       arrow.BinaryTypes.String,
				PrimaryKey: true,
			},
			{
				Name: "version",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}

	expectedFunnyTable = schema.Table{
		Name: "test_funny_struct",
		Columns: schema.ColumnList{
			{
				Name: "a_funny_looking_field",
				Type: arrow.BinaryTypes.String,
			},
			{
				Name: "camel_case_name",
				Type: arrow.BinaryTypes.String,
			},
		},
	}

	expectedTestTableStructWithAny = schema.Table{
		Name: "test_embedded_struct_with_any",
		Columns: schema.ColumnList{
			{
				Name: "int_col",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}

	expectedTestTableStructWithCustomAny = schema.Table{
		Name: "test_embedded_struct_with_custom_any",
		Columns: schema.ColumnList{
			{
				Name: "int_col",
				Type: arrow.PrimitiveTypes.Int64,
			},
			{
				Name: "properties",
				Type: types.ExtensionTypes.JSON,
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
			name: "should unwrap all embedded structs when option is set and use top-level field as PK",
			args: args{
				testStruct: testStructWithEmbeddedStruct{},
				options: []StructTransformerOption{
					WithUnwrapAllEmbeddedStructs(),
					WithPrimaryKeys("IntCol"),
				},
			},
			want: expectedTestTableEmbeddedStructWithTopLevelPK,
		},
		{
			name: "should unwrap all embedded structs when option is set and use its field as PK",
			args: args{
				testStruct: testStructWithEmbeddedStruct{},
				options: []StructTransformerOption{
					WithUnwrapAllEmbeddedStructs(),
					WithPrimaryKeys("EmbeddedString"),
				},
			},
			want: expectedTestTableEmbeddedStructWithUnwrappedPK,
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
			name: "should unwrap specific structs when option is set and use top level field as PK",
			args: args{
				testStruct: testStructWithNonEmbeddedStruct{},
				options: []StructTransformerOption{
					WithUnwrapStructFields("NonEmbedded"),
					WithPrimaryKeys("IntCol"),
				},
			},
			want: expectedTestTableNonEmbeddedStructWithTopLevelPK,
		},
		{
			name: "should unwrap specific structs when option is set and use its field as PK",
			args: args{
				testStruct: testStructWithNonEmbeddedStruct{},
				options: []StructTransformerOption{
					WithUnwrapStructFields("NonEmbedded"),
					WithPrimaryKeys("NonEmbedded.IntCol"),
				},
			},
			want: expectedTestTableNonEmbeddedStructWithUnwrappedPK,
		},
		{
			name: "should generate table from slice struct",
			args: args{
				testStruct: testSliceStruct{},
			},
			want: expectedTestSliceStruct,
		},
		{
			name: "Should configure primary keys when options are set",
			args: args{
				testStruct: testPKStruct{},
				options: []StructTransformerOption{
					WithPrimaryKeys("Parent", "Name"),
				},
			},
			want: expectedTableWithPKs,
		},
		{
			name: "Should return an error when a PK Field is not found",
			args: args{
				testStruct: testPKStruct{},
				options: []StructTransformerOption{
					WithPrimaryKeys("Parent", "Name", "InvalidColumns"),
				},
			},
			want:    expectedTableWithPKs,
			wantErr: true,
		},
		{
			name: "Should properly transform structs with funny looking fields",
			args: args{
				testStruct: testFunnyStruct{},
			},
			want: expectedFunnyTable,
		},
		{
			name: "Ignore any-type fields by default",
			args: args{
				testStruct: testStructWithAny{},
			},
			want: expectedTestTableStructWithAny,
		},
		{
			name: "Should be able to override any-type fields with a type",
			args: args{
				testStruct: testStructWithAny{},
				options: []StructTransformerOption{
					WithTypeTransformer(func(f reflect.StructField) (arrow.DataType, error) {
						if f.Type.Kind() == reflect.Interface {
							return types.ExtensionTypes.JSON, nil
						}
						return nil, nil
					}),
				},
			},
			want: expectedTestTableStructWithCustomAny,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			table := schema.Table{
				Name:    "test",
				Columns: schema.ColumnList{},
			}
			transformer := TransformWithStruct(tt.args.testStruct, tt.args.options...)
			err := transformer(&table)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}
			if tt.wantErr {
				t.Fatal("expected error, got none")
			}

			for i, col := range table.Columns {
				if !arrow.TypeEqual(col.Type, tt.want.Columns[i].Type) {
					t.Fatalf("column %q does not match expected type. got %v, want %v", col.Name, col.Type, tt.want.Columns[i].Type)
				}
			}
			for _, exc := range tt.want.Columns {
				if c := table.Column(exc.Name); c == nil {
					t.Fatalf("column %q not found. want: %v", exc.Name, exc.Type)
				}
			}

			if diff := cmp.Diff(table.PrimaryKeys(), tt.want.PrimaryKeys()); diff != "" {
				t.Fatalf("table does not match expected. diff (-got, +want): %v", diff)
			}
		})
	}
}
