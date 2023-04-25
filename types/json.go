package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/goccy/go-json"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
)

type JSONBuilder struct {
	*array.ExtensionBuilder
}

func NewJSONBuilder(bldr *array.ExtensionBuilder) *JSONBuilder {
	b := &JSONBuilder{
		ExtensionBuilder: bldr,
	}
	return b
}

func (b *JSONBuilder) Append(v any) {
	if v == nil {
		b.AppendNull()
		return
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).Append(bytes)
}

func (b *JSONBuilder) UnsafeAppend(v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).UnsafeAppend(bytes)
}

func (b *JSONBuilder) AppendValueFromString(s string) error {
	if s == array.NullValueStr {
		b.AppendNull()
		return nil
	}
	return b.UnmarshalOne(json.NewDecoder(strings.NewReader(s)))
}

func (b *JSONBuilder) AppendValues(v []any, valid []bool) {
	data := make([][]byte, len(v))
	for i := range v {
		bytes, err := json.Marshal(v[i])
		if err != nil {
			panic(err)
		}
		data[i] = bytes
	}
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).AppendValues(data, valid)
}

func (b *JSONBuilder) UnmarshalJSON(data []byte) error {
	var a []any
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	valid := make([]bool, len(a))
	for i := range a {
		valid[i] = a[i] != nil
	}
	b.AppendValues(a, valid)
	return nil
}

func (b *JSONBuilder) UnmarshalOne(dec *json.Decoder) error {
	var buf any
	err := dec.Decode(&buf)
	if err != nil {
		return err
	}
	if buf == nil {
		b.AppendNull()
	} else {
		b.Append(buf)
	}
	return nil
}

// JSONArray is a simple array which is a Binary
type JSONArray struct {
	array.ExtensionArrayBase
}

func (a JSONArray) String() string {
	arr := a.Storage().(*array.Binary)
	o := new(strings.Builder)
	o.WriteString("[")
	for i := 0; i < arr.Len(); i++ {
		if i > 0 {
			o.WriteString(" ")
		}
		switch {
		case a.IsNull(i):
			o.WriteString(array.NullValueStr)
		default:
			fmt.Fprintf(o, "\"%s\"", arr.Value(i))
		}
	}
	o.WriteString("]")
	return o.String()
}

func (a *JSONArray) Value(i int) []byte {
	if a.IsNull(i) {
		return nil
	}
	return a.Storage().(*array.Binary).Value(i)
}

func (a *JSONArray) ValueStr(i int) string {
	switch {
	case a.IsNull(i):
		return array.NullValueStr
	default:
		return string(a.Value(i))
	}
}

func (a *JSONArray) MarshalJSON() ([]byte, error) {
	arr := a.Storage().(*array.Binary)
	vals := make([]any, a.Len())
	for i := 0; i < a.Len(); i++ {
		if a.IsValid(i) {
			err := json.Unmarshal(arr.Value(i), &vals[i])
			if err != nil {
				panic(fmt.Errorf("invalid json: %w", err))
			}
		} else {
			vals[i] = nil
		}
	}
	return json.Marshal(vals)
}

func (a *JSONArray) GetOneForMarshal(i int) any {
	arr := a.Storage().(*array.Binary)
	if a.IsValid(i) {
		var data any
		err := json.Unmarshal(arr.Value(i), &data)
		if err != nil {
			panic(fmt.Errorf("invalid json: %w", err))
		}
		return data
	}
	return nil
}

// JSONType is a simple extension type that represents a BinaryType
// to be used for representing JSONs
type JSONType struct {
	arrow.ExtensionBase
}

// NewJSONType is a convenience function to create an instance of JSONType
// with the correct storage type
func NewJSONType() *JSONType {
	return &JSONType{
		ExtensionBase: arrow.ExtensionBase{
			Storage: &arrow.BinaryType{}}}
}

// ArrayType returns TypeOf(JSONType) for constructing JSON arrays
func (JSONType) ArrayType() reflect.Type {
	return reflect.TypeOf(JSONArray{})
}

func (JSONType) ExtensionName() string {
	return "json"
}

func (e JSONType) String() string {
	return fmt.Sprintf("extension_type<storage=%s>", e.Storage)
}

func (e JSONType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"name":"%s","metadata":%s}`, e.ExtensionName(), e.Serialize())), nil
}

// Serialize returns "json-serialized" for testing proper metadata passing
func (JSONType) Serialize() string {
	return "json-serialized"
}

// Deserialize expects storageType to be BinaryBuilder and the data to be
// "json-serialized" in order to correctly create a JSONType for testing deserialize.
func (JSONType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "json-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.BinaryType{}) {
		return nil, fmt.Errorf("invalid storage type for JSONType: %s", storageType.Name())
	}
	return NewJSONType(), nil
}

// ExtensionEquals returns true if both extensions have the same name
func (e JSONType) ExtensionEquals(other arrow.ExtensionType) bool {
	return e.ExtensionName() == other.ExtensionName()
}

func (JSONType) NewBuilder(bldr *array.ExtensionBuilder) array.Builder {
	return NewJSONBuilder(bldr)
}
