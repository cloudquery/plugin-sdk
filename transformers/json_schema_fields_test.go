package transformers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/faker"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	gojson "github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

type SchemaTestSecret struct {
	SecretName string `json:"secretName"`
}

type SchemaTestObjectReference struct {
	Name string `json:"name"`
}

type SchemaTestConfigMapVolumeSource struct {
	SchemaTestObjectReference `json:",inline"` //nolint:revive
	Optional                  *bool            `json:"optional,omitempty"`
}

type SchemaTestVolumeSource struct {
	Secret    *SchemaTestSecret                `json:"secret,omitempty"`
	ConfigMap *SchemaTestConfigMapVolumeSource `json:"configMap,omitempty"`
}

type SchemaTestVolume struct {
	Name                   string           `json:"name"`
	SchemaTestVolumeSource `json:",inline"` //nolint:revive
}

type SchemaTestQuantity struct {
	milli  int64
	Format string
}

func (q SchemaTestQuantity) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%dm", q.milli))
}

type SchemaTestHTTPGetAction struct {
	Path string                `json:"path,omitempty"`
	Port SchemaTestIntOrString `json:"port"`
}

type SchemaTestProbeHandler struct {
	HTTPGet *SchemaTestHTTPGetAction `json:"httpGet,omitempty"`
}

type SchemaTestProbe struct {
	SchemaTestProbeHandler `json:",inline"` //nolint:revive
	PeriodSeconds          int32            `json:"periodSeconds,omitempty"`
}

type SchemaTestResourceRequirements struct {
	Limits map[string]SchemaTestQuantity `json:"limits,omitempty"`
}

type SchemaTestContainer struct {
	Name          string                         `json:"name"`
	LivenessProbe *SchemaTestProbe               `json:"livenessProbe,omitempty"`
	Resources     SchemaTestResourceRequirements `json:"resources,omitempty"`
}

type SchemaTestPodSpec struct {
	Volumes    []SchemaTestVolume    `json:"volumes,omitempty"`
	Containers []SchemaTestContainer `json:"containers"`
}

type SchemaTestNamed struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type SchemaTestUntaggedName struct {
	Name string
}

type SchemaTestTaggedName struct {
	Name int `json:"Name"`
}

type SchemaTestLeft struct {
	SchemaTestObjectReference
}

type SchemaTestRight struct {
	SchemaTestObjectReference
}

type SchemaTestIntOrString struct {
	Type   int
	IntVal int32
	StrVal string
}

func (v SchemaTestIntOrString) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.IntVal)
}

type SchemaTestPointerMarshaler struct {
	Value string
}

func (v *SchemaTestPointerMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

type SchemaTestText struct {
	Raw string
}

func (v SchemaTestText) MarshalText() ([]byte, error) {
	return []byte(v.Raw), nil
}

type SchemaTestMarshalerList []SchemaTestSecret

func (v SchemaTestMarshalerList) MarshalJSON() ([]byte, error) {
	return json.Marshal(len(v))
}

type schemaTestUnexportedEmbed struct {
	Exported string `json:"exported"`
}

func TestJSONTypeSchemaMatchesEncoder(t *testing.T) {
	tests := []struct {
		name       string
		testStruct any
		want       string
	}{
		{
			name: "flattens inline embedded structs at every level",
			testStruct: struct {
				Volume SchemaTestVolume `json:"volume"`
			}{},
			want: `{"configMap":{"name":"utf8","optional":"bool"},"name":"utf8","secret":{"secretName":"utf8"}}`,
		},
		{
			name: "flattens an embedded struct without a json tag",
			testStruct: struct {
				Ref struct {
					SchemaTestObjectReference
				} `json:"ref"`
			}{},
			want: `{"name":"utf8"}`,
		},
		{
			name: "flattens an embedded struct pointer",
			testStruct: struct {
				Ref struct {
					*SchemaTestObjectReference
				} `json:"ref"`
			}{},
			want: `{"name":"utf8"}`,
		},
		{
			name: "flattens the exported fields of an unexported embedded struct",
			testStruct: struct {
				Item struct {
					schemaTestUnexportedEmbed
				} `json:"item"`
			}{},
			want: `{"exported":"utf8"}`,
		},
		{
			name: "keeps an embedded struct with a json name as a key",
			testStruct: struct {
				Ref struct {
					SchemaTestObjectReference `json:"reference"`
				} `json:"ref"`
			}{},
			want: `{"reference":{"name":"utf8"}}`,
		},
		{
			name: "prefers the shallower field when names collide",
			testStruct: struct {
				Item struct {
					Name string `json:"name"`
					SchemaTestObjectReference
				} `json:"item"`
			}{},
			want: `{"name":"utf8"}`,
		},
		{
			name: "drops tagged fields that collide at the same depth",
			testStruct: struct {
				Item struct {
					SchemaTestNamed
					SchemaTestObjectReference //nolint:govet
				} `json:"item"`
			}{},
			want: `{"id":"utf8"}`,
		},
		{
			name: "prefers the tagged field when names collide at the same depth",
			testStruct: struct {
				Item struct {
					SchemaTestUntaggedName
					SchemaTestTaggedName
				} `json:"item"`
			}{},
			want: `{"Name":"int64"}`,
		},
		{
			name: "drops fields of the same embedded type reached twice at the same depth",
			testStruct: struct {
				Item struct {
					SchemaTestLeft
					SchemaTestRight        //nolint:govet
					ID              string `json:"id"`
				} `json:"item"`
			}{},
			want: `{"id":"utf8"}`,
		},
		{
			name: "uses json for a type with MarshalJSON",
			testStruct: struct {
				Item struct {
					Port SchemaTestIntOrString `json:"port"`
				} `json:"item"`
			}{},
			want: `{"port":"json"}`,
		},
		{
			name: "uses json for a type with a pointer MarshalJSON",
			testStruct: struct {
				Item struct {
					Value *SchemaTestPointerMarshaler `json:"value"`
				} `json:"item"`
			}{},
			want: `{"value":"json"}`,
		},
		{
			name: "uses json for a column type with MarshalJSON",
			testStruct: struct {
				Port SchemaTestIntOrString `json:"port"`
			}{},
			want: `"json"`,
		},
		{
			name: "uses json for a slice type with MarshalJSON",
			testStruct: struct {
				Item struct {
					Secrets SchemaTestMarshalerList `json:"secrets"`
				} `json:"item"`
			}{},
			want: `{"secrets":"json"}`,
		},
		{
			name: "uses utf8 for a type with MarshalText",
			testStruct: struct {
				Item struct {
					Text SchemaTestText `json:"text"`
				} `json:"item"`
			}{},
			want: `{"text":"utf8"}`,
		},
		{
			name: "keeps the timestamp type for time.Time",
			testStruct: struct {
				Item struct {
					CreatedAt time.Time `json:"created_at"`
				} `json:"item"`
			}{},
			want: `{"created_at":"timestamp[us, tz=UTC]"}`,
		},
		{
			name: "skips a field tagged with a dash",
			testStruct: struct {
				Item struct {
					Name   string `json:"name"`
					Secret string `json:"-"`
				} `json:"item"`
			}{},
			want: `{"name":"utf8"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := schema.Table{Name: "test"}
			require.NoError(t, TransformWithStruct(tt.testStruct)(&table))
			require.Len(t, table.Columns, 1)
			require.Equal(t, tt.want, table.Columns[0].TypeSchema)
		})
	}
}

func TestJSONTypeSchemaKeysMatchMarshalledValue(t *testing.T) {
	tests := []struct {
		name   string
		column func() any
	}{
		{
			name:   "kubernetes-shaped volume",
			column: func() any { return &struct{ Volume SchemaTestVolume }{} },
		},
		{
			name:   "kubernetes-shaped pod spec",
			column: func() any { return &struct{ Spec SchemaTestPodSpec }{} },
		},
		{
			name: "colliding embedded structs",
			column: func() any {
				return &struct {
					Item struct {
						SchemaTestNamed
						SchemaTestObjectReference //nolint:govet
						SchemaTestUntaggedName
						SchemaTestTaggedName
						SchemaTestLeft
						SchemaTestRight //nolint:govet
					}
				}{}
			},
		},
		{
			name: "types with their own marshallers",
			column: func() any {
				return &struct {
					Item struct {
						Port    SchemaTestIntOrString
						Value   *SchemaTestPointerMarshaler
						Text    SchemaTestText
						Secrets SchemaTestMarshalerList
					}
				}{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := tt.column()
			require.NoError(t, faker.FakeObject(column, faker.WithMaxDepth(30)))

			table := schema.Table{Name: "test"}
			require.NoError(t, TransformWithStruct(column)(&table))
			require.Len(t, table.Columns, 1)

			var typeSchema any
			require.NoError(t, json.Unmarshal([]byte(table.Columns[0].TypeSchema), &typeSchema))

			marshalled, err := gojson.MarshalWithOption(firstFieldValue(column), gojson.DisableHTMLEscape())
			require.NoError(t, err)
			var value any
			require.NoError(t, json.Unmarshal(marshalled, &value))

			var mismatches []string
			compareJSONKeys(typeSchema, value, "$", &mismatches)
			slices.Sort(mismatches)
			require.Empty(t, mismatches, "type schema: %s\nvalue: %s", table.Columns[0].TypeSchema, marshalled)
		})
	}
}

func firstFieldValue(column any) any {
	return reflect.ValueOf(column).Elem().Field(0).Interface()
}

func compareJSONKeys(typeSchema, value any, path string, mismatches *[]string) {
	if value == nil {
		return
	}
	switch s := typeSchema.(type) {
	case []any:
		values, ok := value.([]any)
		if !ok {
			*mismatches = append(*mismatches, fmt.Sprintf("%s: schema has a list, value is %T", path, value))
			return
		}
		for _, v := range values {
			compareJSONKeys(s[0], v, path+"[]", mismatches)
		}
	case map[string]any:
		object, ok := value.(map[string]any)
		if !ok {
			*mismatches = append(*mismatches, fmt.Sprintf("%s: schema has an object, value is %T", path, value))
			return
		}
		if mapValueSchema, isMap := s["utf8"]; isMap && len(s) == 1 {
			for _, v := range object {
				compareJSONKeys(mapValueSchema, v, path+".*", mismatches)
			}
			return
		}
		for key, keySchema := range s {
			v, ok := object[key]
			if !ok {
				*mismatches = append(*mismatches, fmt.Sprintf("%s.%s: in schema, not in value", path, key))
				continue
			}
			compareJSONKeys(keySchema, v, path+"."+key, mismatches)
		}
		for key := range object {
			if _, ok := s[key]; !ok {
				*mismatches = append(*mismatches, fmt.Sprintf("%s.%s: in value, not in schema", path, key))
			}
		}
	}
}
