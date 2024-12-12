package types

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/goccy/go-json"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

type JSONBuilder struct {
	*array.ExtensionBuilder
}

func NewJSONBuilder(mem memory.Allocator) *JSONBuilder {
	return &JSONBuilder{ExtensionBuilder: array.NewExtensionBuilder(mem, NewJSONType())}
}

func (b *JSONBuilder) AppendBytes(v []byte) {
	if v == nil {
		b.AppendNull()
		return
	}

	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).Append(v)
}

func (b *JSONBuilder) Append(v any) {
	if v == nil {
		b.AppendNull()
		return
	}

	// per https://github.com/cloudquery/plugin-sdk/issues/622
	data, err := json.MarshalWithOption(v, json.DisableHTMLEscape())
	if err != nil {
		panic(err)
	}

	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).Append(data)
}

func (b *JSONBuilder) UnsafeAppend(v any) {
	// per https://github.com/cloudquery/plugin-sdk/issues/622
	data, err := json.MarshalWithOption(v, json.DisableHTMLEscape())
	if err != nil {
		panic(err)
	}

	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).UnsafeAppend(data)
}

func (b *JSONBuilder) AppendValueFromString(s string) error {
	if s == array.NullValueStr {
		b.AppendNull()
		return nil
	}
	return b.UnmarshalOne(json.NewDecoder(strings.NewReader(s)))
}

func (b *JSONBuilder) AppendValues(v []any, valid []bool) {
	if len(v) != len(valid) && len(valid) != 0 {
		panic("len(v) != len(valid) && len(valid) != 0")
	}

	data := make([][]byte, len(v))
	var err error
	for i := range v {
		if len(valid) > 0 && !valid[i] {
			continue
		}
		// per https://github.com/cloudquery/plugin-sdk/issues/622
		data[i], err = json.MarshalWithOption(v[i], json.DisableHTMLEscape())
		if err != nil {
			panic(err)
		}
	}
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).AppendValues(data, valid)
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

func (b *JSONBuilder) Unmarshal(dec *json.Decoder) error {
	for dec.More() {
		if err := b.UnmarshalOne(dec); err != nil {
			return err
		}
	}
	return nil
}

func (b *JSONBuilder) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("JSON builder must unpack from JSON array, found %s", delim)
	}

	return b.Unmarshal(dec)
}

func (b *JSONBuilder) NewJSONArray() *JSONArray {
	return b.NewExtensionArray().(*JSONArray)
}

// JSONArray is a simple array which is a Binary
type JSONArray struct {
	array.ExtensionArrayBase
}

func (a *JSONArray) String() string {
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
			fmt.Fprintf(o, "%q", a.ValueStr(i))
		}
	}
	o.WriteString("]")
	return o.String()
}

func (a *JSONArray) Value(i int) any {
	if a.IsNull(i) {
		return nil
	}

	var data any
	err := json.Unmarshal(a.Storage().(*array.Binary).Value(i), &data)
	if err != nil {
		panic(fmt.Errorf("invalid json: %w", err))
	}
	return data
}

func (a *JSONArray) ValueStr(i int) string {
	switch {
	case a.IsNull(i):
		return array.NullValueStr
	default:
		return string(a.GetOneForMarshal(i).(json.RawMessage))
	}
}

func (a *JSONArray) MarshalJSON() ([]byte, error) {
	values := make([]json.RawMessage, a.Len())
	for i := 0; i < a.Len(); i++ {
		if a.IsNull(i) {
			continue
		}
		values[i] = a.GetOneForMarshal(i).(json.RawMessage)
	}
	// per https://github.com/cloudquery/plugin-sdk/issues/622
	return json.MarshalWithOption(values, json.DisableHTMLEscape())
}

func (a *JSONArray) GetOneForMarshal(i int) any {
	if a.IsNull(i) {
		return nil
	}
	return json.RawMessage(a.Storage().(*array.Binary).Value(i))
}

// JSONType is a simple extension type that represents a BinaryType
// to be used for representing JSONs
type JSONType struct {
	arrow.ExtensionBase
}

// NewJSONType is a convenience function to create an instance of JSONType
// with the correct storage type
func NewJSONType() *JSONType {
	return &JSONType{ExtensionBase: arrow.ExtensionBase{Storage: &arrow.BinaryType{}}}
}

// ArrayType returns TypeOf(JSONArray{}) for constructing JSON arrays
func (*JSONType) ArrayType() reflect.Type {
	return reflect.TypeOf(JSONArray{})
}

func (*JSONType) ExtensionName() string {
	return "json"
}

func (*JSONType) String() string {
	return "json"
}

func (e *JSONType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"name":"%s","metadata":%s}`, e.ExtensionName(), e.Serialize())), nil
}

// Serialize returns "json-serialized" for testing proper metadata passing
func (*JSONType) Serialize() string {
	return "json-serialized"
}

// Deserialize expects storageType to be BinaryBuilder and the data to be
// "json-serialized" in order to correctly create a JSONType for testing deserialize.
func (*JSONType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "json-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.BinaryType{}) {
		return nil, fmt.Errorf("invalid storage type for *JSONType: %s", storageType.Name())
	}
	return NewJSONType(), nil
}

// ExtensionEquals returns true if both extensions have the same name
func (e *JSONType) ExtensionEquals(other arrow.ExtensionType) bool {
	return e.ExtensionName() == other.ExtensionName()
}

func (*JSONType) NewBuilder(mem memory.Allocator) array.Builder {
	return NewJSONBuilder(mem)
}
