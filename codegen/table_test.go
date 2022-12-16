package codegen

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/caser"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
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
	testStructWithCustomType struct {
		TimeCol time.Time `json:"time_col,omitempty"`
	}

	customStruct                 struct{}
	testStructForCustomResolvers struct {
		TimeCol time.Time    `json:"time_col,omitempty"`
		Custom  customStruct `json:"custom"`
	}

	testStructCaseCheck struct {
		IPAddress   string
		CDNs        string
		MyCDN       string
		CIDR        int
		IPV6        int
		IPv6Test    int
		Ipv6Address int
		AccountID   string
		PostgreSQL  string
		IDs         string
	}

	testSliceStruct []struct {
		IntCol int
	}
)

var (
	expectedColumns = []ColumnDefinition{
		{
			Name:     "int_col",
			Type:     schema.TypeInt,
			Resolver: `schema.PathResolver("IntCol")`,
			Options:  schema.ColumnCreationOptions{PrimaryKey: true},
		},
		{
			Name:     "int64_col",
			Type:     schema.TypeInt,
			Resolver: `schema.PathResolver("Int64Col")`,
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
			Name:     "int_pointer_array_col",
			Type:     schema.TypeIntArray,
			Resolver: `schema.PathResolver("IntPointerArrayCol")`,
		},
		{
			Name:     "string_array_col",
			Type:     schema.TypeStringArray,
			Resolver: `schema.PathResolver("StringArrayCol")`,
		},
		{
			Name:     "string_pointer_array_col",
			Type:     schema.TypeStringArray,
			Resolver: `schema.PathResolver("StringPointerArrayCol")`,
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
			Resolver: `schema.PathResolver("JSONTag")`,
		},
		{
			Name:     "no_json_tag",
			Type:     schema.TypeString,
			Resolver: `schema.PathResolver("NoJSONTag")`,
		},
	}
	expectedTestTable = TableDefinition{
		Name:                "test_struct",
		Columns:             expectedColumns,
		nameTransformer:     DefaultNameTransformer,
		typeTransformer:     DefaultTypeTransformer,
		resolverTransformer: DefaultResolverTransformer,
	}
	expectedTestTableEmbeddedStruct = TableDefinition{
		Name:                "test_struct",
		Columns:             append(expectedColumns, ColumnDefinition{Name: "embedded_string", Type: schema.TypeString, Resolver: `schema.PathResolver("EmbeddedString")`}),
		nameTransformer:     DefaultNameTransformer,
		typeTransformer:     DefaultTypeTransformer,
		resolverTransformer: DefaultResolverTransformer,
	}
	expectedTestTableNonEmbeddedStruct = TableDefinition{
		Name: "test_struct",
		Columns: ColumnDefinitions{
			// Should not be unwrapped
			ColumnDefinition{Name: "test_struct", Type: schema.TypeJSON, Resolver: `schema.PathResolver("TestStruct")`},
			// Should be unwrapped
			ColumnDefinition{
				Name:     "non_embedded_embedded_string",
				Type:     schema.TypeString,
				Resolver: `schema.PathResolver("NonEmbedded.EmbeddedString")`,
				Options:  schema.ColumnCreationOptions{PrimaryKey: true},
			},
		},
		nameTransformer:     DefaultNameTransformer,
		typeTransformer:     DefaultTypeTransformer,
		resolverTransformer: DefaultResolverTransformer,
	}
	expectedTestTableStructForCustomResolvers = TableDefinition{
		Name: "test_struct",
		Columns: ColumnDefinitions{
			{
				Name:     "time_col",
				Type:     schema.TypeTimestamp,
				Resolver: `schema.PathResolver("TimeCol")`,
				Options:  schema.ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name:     "custom",
				Type:     schema.TypeTimestamp,
				Resolver: `customResolver("Custom")`,
			},
		},
		nameTransformer:     DefaultNameTransformer,
		typeTransformer:     customTypeTransformer,
		resolverTransformer: customResolverTransformer,
	}
	expectedTestSliceStruct = TableDefinition{
		Name: "test_struct",
		Columns: ColumnDefinitions{
			{
				Name:     "int_col",
				Type:     schema.TypeInt,
				Resolver: `schema.PathResolver("IntCol")`,
				Options:  schema.ColumnCreationOptions{PrimaryKey: true},
			},
		},
		nameTransformer:     DefaultNameTransformer,
		typeTransformer:     customTypeTransformer,
		resolverTransformer: customResolverTransformer,
	}
)

func customResolverTransformer(field reflect.StructField, path string) (string, error) {
	switch reflect.New(field.Type).Interface().(type) {
	case customStruct, *customStruct:
		return `customResolver("` + path + `")`, nil
	default:
		return "", nil
	}
}

func customTypeTransformer(field reflect.StructField) (schema.ValueType, error) {
	switch reflect.New(field.Type).Interface().(type) {
	case customStruct, *customStruct:
		return schema.TypeTimestamp, nil
	default:
		return schema.TypeInvalid, nil
	}
}

func customNameTransformer(field reflect.StructField) (string, error) {
	name := field.Name
	if jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]; len(jsonTag) > 0 {
		// return empty string if the field is not related api response
		if jsonTag == "-" {
			return "", nil
		}
		name = jsonTag
	}

	c := caser.New(caser.WithCustomInitialisms(map[string]bool{"CDN": true, "IP": true, "IPv6": true, "IPV6": true, "CIDR": true}))
	return c.ToSnake(name), nil
}

func TestTableFromGoStruct(t *testing.T) {
	type args struct {
		testStruct any
		options    []TableOption
	}

	tests := []struct {
		name    string
		args    args
		want    TableDefinition
		wantErr bool
	}{
		{
			name: "should generate table from struct with default options",
			args: args{
				testStruct: testStruct{},
				options:    []TableOption{WithPKColumns("int_col")},
			},
			want: expectedTestTable,
		},
		{
			name: "should unwrap all embedded structs when option is set",
			args: args{
				testStruct: testStructWithEmbeddedStruct{},
				options: []TableOption{
					WithPKColumns("int_col"),
					WithUnwrapAllEmbeddedStructs(),
				},
			},
			want: expectedTestTableEmbeddedStruct,
		},
		{
			name: "should unwrap specific structs when option is set",
			args: args{
				testStruct: testStructWithNonEmbeddedStruct{},
				options: []TableOption{
					WithPKColumns("non_embedded_embedded_string"),
					WithUnwrapStructFields("NonEmbedded"),
				},
			},
			want: expectedTestTableNonEmbeddedStruct,
		},
		{
			name: "should override schema type when option is set",
			args: args{
				testStruct: testStructWithCustomType{},
				options: []TableOption{
					WithPKColumns("time_col"),
					WithTypeTransformer(func(t reflect.StructField) (schema.ValueType, error) {
						switch t.Type {
						case reflect.TypeOf(time.Time{}), reflect.TypeOf(&time.Time{}):
							return schema.TypeJSON, nil
						default:
							return schema.TypeInvalid, nil
						}
					}),
				},
			},
			want: TableDefinition{Name: "test_struct",
				// We expect the time column to be of type JSON, since we override the type of `time.Time` to be JSON
				Columns: ColumnDefinitions{{
					Name:     "time_col",
					Type:     schema.TypeJSON,
					Resolver: `schema.PathResolver("TimeCol")`,
					Options:  schema.ColumnCreationOptions{PrimaryKey: true},
				}},
				nameTransformer: DefaultNameTransformer},
		},
		{
			name: "should handle default and custom acronyms correctly",
			args: args{
				testStruct: testStructCaseCheck{},
				options: []TableOption{
					WithPKColumns("ip_address"),
					WithNameTransformer(customNameTransformer),
				},
			},
			want: TableDefinition{Name: "test_struct",
				// We expect the time column to be of type JSON, since we override the type of `time.Time` to be JSON
				Columns: ColumnDefinitions{
					{
						Name:     "ip_address",
						Type:     schema.TypeString,
						Resolver: `schema.PathResolver("IPAddress")`,
						Options:  schema.ColumnCreationOptions{PrimaryKey: true},
					},
					{Name: "cdns", Type: schema.TypeString, Resolver: `schema.PathResolver("CDNs")`},
					{Name: "my_cdn", Type: schema.TypeString, Resolver: `schema.PathResolver("MyCDN")`},
					{Name: "cidr", Type: schema.TypeInt, Resolver: `schema.PathResolver("CIDR")`},
					{Name: "ipv6", Type: schema.TypeInt, Resolver: `schema.PathResolver("IPV6")`},
					{Name: "ipv6_test", Type: schema.TypeInt, Resolver: `schema.PathResolver("IPv6Test")`},
					{Name: "ipv6_address", Type: schema.TypeInt, Resolver: `schema.PathResolver("Ipv6Address")`},
					{Name: "account_id", Type: schema.TypeString, Resolver: `schema.PathResolver("AccountID")`},
					{Name: "postgre_sql", Type: schema.TypeString, Resolver: `schema.PathResolver("PostgreSQL")`},
					{Name: "ids", Type: schema.TypeString, Resolver: `schema.PathResolver("IDs")`},
				},
				nameTransformer: DefaultNameTransformer},
		},
		{
			name: "should use custom resolvers",
			args: args{
				testStruct: testStructForCustomResolvers{},
				options: []TableOption{
					WithPKColumns("time_col"),
					WithTypeTransformer(customTypeTransformer),
					WithResolverTransformer(customResolverTransformer),
				},
			},
			want: expectedTestTableStructForCustomResolvers,
		},
		{
			name: "should generate table from slice struct",
			args: args{
				testStruct: testSliceStruct{},
				options:    []TableOption{WithPKColumns("int_col")},
			},
			want: expectedTestSliceStruct,
		},
		{
			name: "should error if there are undefined pk columns",
			args: args{
				testStruct: testSliceStruct{},
				options:    []TableOption{WithPKColumns("int_col", "some_col")},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := NewTableFromStruct("test_struct", tt.args.testStruct, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(&tt.want, table,
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

func TestTableDefinition_Check(t *testing.T) {
	tests := []struct {
		name       string
		definition *TableDefinition
		err        error
	}{
		{
			name: "nil",
			err:  fmt.Errorf("nil table definition"),
		},
		{
			name:       "no table name",
			definition: new(TableDefinition),
			err:        fmt.Errorf("empty table name"),
		},
		{
			name:       "no columns",
			definition: &TableDefinition{Name: "no_columns"},
			err:        fmt.Errorf("no_columns: no columns"),
		},
		{
			name: "no column name",
			definition: &TableDefinition{
				Name:    "no_column_name",
				Columns: ColumnDefinitions{{}},
			},
			err: fmt.Errorf("no_column_name: empty column name"),
		},
		{
			name: "invalid type",
			definition: &TableDefinition{
				Name:    "invalid_type",
				Columns: ColumnDefinitions{{Name: "col1"}},
			},
			err: fmt.Errorf("invalid_type->col1: invalid column type"),
		},
		{
			name: "duplicate name",
			definition: &TableDefinition{
				Name: "duplicate_name",
				Columns: ColumnDefinitions{
					{Name: "col1", Type: schema.TypeInt},
					{Name: "col1", Type: schema.TypeInt},
				},
			},
			err: fmt.Errorf("duplicate_name->col1: duplicate column name"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.err, tt.definition.Check())
		})
	}
}
