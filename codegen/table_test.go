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

	tests := []struct {
		name            string
		TableDefinition TableDefinition
		want            ColumnDefinitions
		wantErr         bool
	}{
		{
			name: "should generate table from struct with default options",
			TableDefinition: TableDefinition{
				NameOverride: "tables",
				Struct:       testStruct{},
				PKColumns:    []string{"int_col"},
			},
			want: expectedColumns,
		},
		{
			name: "should unwrap all embedded structs when option is set",
			TableDefinition: TableDefinition{
				NameOverride:                  "tables",
				Struct:                        testStructWithEmbeddedStruct{},
				PKColumns:                     []string{"int_col"},
				UnwrapAllEmbeddedStructFields: true,
			},
			want: append(expectedColumns, ColumnDefinition{Name: "embedded_string", Type: schema.TypeString, Resolver: `schema.PathResolver("EmbeddedString")`}),
		},
		{
			name: "should unwrap specific structs when option is set",
			TableDefinition: TableDefinition{
				NameOverride:         "tables",
				Struct:               testStructWithNonEmbeddedStruct{},
				PKColumns:            []string{"non_embedded_embedded_string"},
				StructFieldsToUnwrap: []string{"NonEmbedded"},
			},
			want: ColumnDefinitions{
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
		},
		{
			name: "should override schema type when option is set",
			TableDefinition: TableDefinition{
				NameOverride: "tables",
				Struct:       testStructWithCustomType{},
				PKColumns:    []string{"time_col"},
				TypeTransformer: func(t reflect.StructField) (schema.ValueType, error) {
					switch t.Type {
					case reflect.TypeOf(time.Time{}), reflect.TypeOf(&time.Time{}):
						return schema.TypeJSON, nil
					default:
						return schema.TypeInvalid, nil
					}
				},
			},
			want: ColumnDefinitions{{
				Name:     "time_col",
				Type:     schema.TypeJSON,
				Resolver: `schema.PathResolver("TimeCol")`,
				Options:  schema.ColumnCreationOptions{PrimaryKey: true},
			}},
		},
		{
			name: "should handle default and custom acronyms correctly",
			TableDefinition: TableDefinition{
				NameOverride:    "tables",
				Struct:          testStructCaseCheck{},
				PKColumns:       []string{"ip_address"},
				NameTransformer: customNameTransformer,
			},
			want: ColumnDefinitions{
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
		},
		{
			name: "should use custom resolvers",
			TableDefinition: TableDefinition{
				NameOverride:        "tables",
				Struct:              testStructForCustomResolvers{},
				PKColumns:           []string{"time_col"},
				TypeTransformer:     customTypeTransformer,
				ResolverTransformer: customResolverTransformer,
			},
			want: ColumnDefinitions{
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
		},
		{
			name: "should generate table from slice struct",
			TableDefinition: TableDefinition{
				NameOverride: "tables",
				Struct:       testSliceStruct{},
				PKColumns:    []string{"int_col"},
			},
			want: ColumnDefinitions{
				{
					Name:     "int_col",
					Type:     schema.TypeInt,
					Resolver: `schema.PathResolver("IntCol")`,
					Options:  schema.ColumnCreationOptions{PrimaryKey: true},
				},
			},
		},
		{
			name: "should error if there are undefined pk columns",
			TableDefinition: TableDefinition{
				NameOverride: "tables",
				Struct:       testStruct{},
				PKColumns:    []string{"int_col", "some_col"},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			columns, err := test.TableDefinition.Columns()
			if err != nil != test.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if diff := cmp.Diff(test.want, columns); diff != "" {
				t.Fatalf("columns does not match expected. diff (-got, +want): %s", diff)
			}
			buf := bytes.NewBufferString("")
			if err := test.TableDefinition.GenerateTemplate(buf); err != nil {
				t.Fatal(err)
			}
			fmt.Println(buf.String())
		})
	}

}

type EmptyStruct struct {
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
			err:        fmt.Errorf("unable to generate table name: missing fields (pluginName, service, subService)"),
		},
		{
			name: "no columns",
			definition: &TableDefinition{
				NameOverride: "no_columns",
				Struct:       EmptyStruct{},
			},
			err: fmt.Errorf("no_columns: no columns"),
		},
		{
			name: "no column name",
			definition: &TableDefinition{
				NameOverride: "no_column_names",
				Struct:       testStruct{},
				ExtraColumns: ColumnDefinitions{{}},
			},
			err: fmt.Errorf("no_column_names: empty column name"),
		},
		{
			name: "invalid type",
			definition: &TableDefinition{
				NameOverride: "invalid_types",
				Struct:       EmptyStruct{},
				ExtraColumns: ColumnDefinitions{{Name: "col1"}},
			},
			err: fmt.Errorf("invalid_types->col1: invalid column type"),
		},
		{
			name: "duplicate name",
			definition: &TableDefinition{
				NameOverride: "duplicate_names",
				Struct:       testStruct{},
				ExtraColumns: ColumnDefinitions{
					{Name: "col1", Type: schema.TypeInt},
					{Name: "col1", Type: schema.TypeInt},
				},
			},
			err: fmt.Errorf("duplicate_names->col1: duplicate column name"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.err, tt.definition.Check())
		})
	}
}
