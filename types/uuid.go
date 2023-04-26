package types

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/goccy/go-json"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/google/uuid"
)

var UUID = NewUUIDType()

type UUIDBuilder struct {
	*array.ExtensionBuilder
}

func NewUUIDBuilder(builder *array.ExtensionBuilder) *UUIDBuilder {
	b := &UUIDBuilder{
		ExtensionBuilder: builder,
	}
	return b
}

func (b *UUIDBuilder) Append(v uuid.UUID) {
	b.ExtensionBuilder.Builder.(*array.FixedSizeBinaryBuilder).Append(v[:])
}

func (b *UUIDBuilder) UnsafeAppend(v uuid.UUID) {
	b.ExtensionBuilder.Builder.(*array.FixedSizeBinaryBuilder).UnsafeAppend(v[:])
}

func (b *UUIDBuilder) AppendValueFromString(s string) error {
	if s == array.NullValueStr {
		b.AppendNull()
		return nil
	}

	b.Append(uuid.MustParse(s))
	return nil
}

func (b *UUIDBuilder) AppendValues(v []uuid.UUID, valid []bool) {
	data := make([][]byte, len(v))
	for i := range v {
		if !valid[i] {
			continue
		}
		data[i] = v[i][:]
	}
	b.ExtensionBuilder.Builder.(*array.FixedSizeBinaryBuilder).AppendValues(data, valid)
}

func (b *UUIDBuilder) UnmarshalOne(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}

	var val uuid.UUID
	switch v := t.(type) {
	case string:
		val, err = uuid.Parse(v)
		if err != nil {
			return err
		}
	case []byte:
		val, err = uuid.ParseBytes(v)
		if err != nil {
			return err
		}
	case nil:
		b.AppendNull()
		return nil
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprint(t),
			Type:   reflect.TypeOf([]byte{}),
			Offset: dec.InputOffset(),
			Struct: fmt.Sprintf("FixedSizeBinary[%d]", 16),
		}
	}

	b.Append(val)
	return nil
}

func (b *UUIDBuilder) Unmarshal(dec *json.Decoder) error {
	for dec.More() {
		if err := b.UnmarshalOne(dec); err != nil {
			return err
		}
	}
	return nil
}

func (b *UUIDBuilder) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("uuid builder must unpack from json array, found %s", delim)
	}

	return b.Unmarshal(dec)
}

// UUIDArray is a simple array which is a FixedSizeBinary(16)
type UUIDArray struct {
	array.ExtensionArrayBase
}

func (a *UUIDArray) String() string {
	arr := a.Storage().(*array.FixedSizeBinary)
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
			fmt.Fprintf(o, "%q", a.Value(i))
		}
	}
	o.WriteString("]")
	return o.String()
}

func (a *UUIDArray) Value(i int) uuid.UUID {
	if a.IsNull(i) {
		return uuid.Nil
	}
	return uuid.Must(uuid.FromBytes(a.Storage().(*array.FixedSizeBinary).Value(i)))
}

func (a *UUIDArray) ValueStr(i int) string {
	switch {
	case a.IsNull(i):
		return array.NullValueStr
	default:
		return a.Value(i).String()
	}
}

func (a *UUIDArray) MarshalJSON() ([]byte, error) {
	arr := a.Storage().(*array.FixedSizeBinary)
	values := make([]any, a.Len())
	for i := 0; i < a.Len(); i++ {
		if a.IsValid(i) {
			values[i] = uuid.Must(uuid.FromBytes(arr.Value(i))).String()
		}
	}
	return json.Marshal(values)
}

func (a *UUIDArray) GetOneForMarshal(i int) any {
	if a.IsNull(i) {
		return nil
	}
	return uuid.Must(uuid.FromBytes(a.Storage().(*array.FixedSizeBinary).Value(i)))
}

// UUIDType is a simple extension type that represents a FixedSizeBinary(16)
// to be used for representing UUIDs
type UUIDType struct {
	arrow.ExtensionBase
}

// NewUUIDType is a convenience function to create an instance of UuidType
// with the correct storage type
func NewUUIDType() *UUIDType {
	return &UUIDType{
		ExtensionBase: arrow.ExtensionBase{
			Storage: &arrow.FixedSizeBinaryType{ByteWidth: 16}}}
}

// ArrayType returns TypeOf(UuidArray) for constructing uuid arrays
func (UUIDType) ArrayType() reflect.Type {
	return reflect.TypeOf(UUIDArray{})
}

func (UUIDType) ExtensionName() string {
	return "uuid"
}

func (e UUIDType) String() string {
	return fmt.Sprintf("extension_type<storage=%s>", e.Storage)
}

func (e UUIDType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"name":"%s","metadata":%s}`, e.ExtensionName(), e.Serialize())), nil
}

// Serialize returns "uuid-serialized" for testing proper metadata passing
func (UUIDType) Serialize() string {
	return "uuid-serialized"
}

// Deserialize expects storageType to be FixedSizeBinaryType{ByteWidth: 16} and the data to be
// "uuid-serialized" in order to correctly create a UuidType for testing deserialize.
func (UUIDType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "uuid-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.FixedSizeBinaryType{ByteWidth: 16}) {
		return nil, fmt.Errorf("invalid storage type for UuidType: %s", storageType.Name())
	}
	return NewUUIDType(), nil
}

// ExtensionEquals for UUIDType just checks that the names are equal
func (e UUIDType) ExtensionEquals(other arrow.ExtensionType) bool {
	return e.ExtensionName() == other.ExtensionName()
}

func (UUIDType) NewBuilder(builder *array.ExtensionBuilder) array.Builder {
	return NewUUIDBuilder(builder)
}
