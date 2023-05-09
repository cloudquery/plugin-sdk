package types

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/goccy/go-json"
)

type MacBuilder struct {
	*array.ExtensionBuilder
}

func NewMacBuilder(builder *array.ExtensionBuilder) *MacBuilder {
	return &MacBuilder{ExtensionBuilder: builder}
}

func (b *MacBuilder) Append(v net.HardwareAddr) {
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).Append(v[:])
}

func (b *MacBuilder) UnsafeAppend(v net.HardwareAddr) {
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).UnsafeAppend(v[:])
}

func (b *MacBuilder) AppendValues(v []net.HardwareAddr, valid []bool) {
	if len(v) != len(valid) && len(valid) != 0 {
		panic("len(v) != len(valid) && len(valid) != 0")
	}

	data := make([][]byte, len(v))
	for i, v := range v {
		if len(valid) > 0 && !valid[i] {
			continue
		}
		data[i] = v
	}
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).AppendValues(data, valid)
}

func (b *MacBuilder) AppendValueFromString(s string) error {
	if s == array.NullValueStr {
		b.AppendNull()
		return nil
	}
	data, err := net.ParseMAC(s)
	if err != nil {
		return err
	}
	b.Append(data)
	return nil
}

func (b *MacBuilder) UnmarshalOne(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}

	var val net.HardwareAddr
	switch v := t.(type) {
	case string:
		data, err := net.ParseMAC(v)
		if err != nil {
			return err
		}
		val = data
	case []byte:
		val = net.HardwareAddr(v)
	case nil:
		b.AppendNull()
		return nil
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprint(t),
			Type:   reflect.TypeOf([]byte{}),
			Offset: dec.InputOffset(),
			Struct: "Binary",
		}
	}

	b.Append(val)
	return nil
}

func (b *MacBuilder) Unmarshal(dec *json.Decoder) error {
	for dec.More() {
		if err := b.UnmarshalOne(dec); err != nil {
			return err
		}
	}
	return nil
}

func (b *MacBuilder) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("fixed size binary builder must unpack from json array, found %s", delim)
	}

	return b.Unmarshal(dec)
}

func (b *MacBuilder) NewMacArray() *MacArray {
	return b.NewExtensionArray().(*MacArray)
}

// MacArray is a simple array which is a wrapper around a BinaryArray
type MacArray struct {
	array.ExtensionArrayBase
}

func (a *MacArray) String() string {
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
			fmt.Fprintf(o, "%q", a.Value(i))
		}
	}
	o.WriteString("]")
	return o.String()
}

func (a *MacArray) Value(i int) net.HardwareAddr {
	if a.IsNull(i) {
		return nil
	}
	return net.HardwareAddr(a.Storage().(*array.Binary).Value(i))
}

func (a *MacArray) ValueStr(i int) string {
	switch {
	case a.IsNull(i):
		return array.NullValueStr
	default:
		return a.Value(i).String()
	}
}

func (a *MacArray) MarshalJSON() ([]byte, error) {
	arr := a.Storage().(*array.Binary)
	values := make([]any, a.Len())
	for i := 0; i < a.Len(); i++ {
		if a.IsValid(i) {
			values[i] = net.HardwareAddr(arr.Value(i)).String()
		} else {
			values[i] = nil
		}
	}
	return json.Marshal(values)
}

func (a *MacArray) GetOneForMarshal(i int) any {
	arr := a.Storage().(*array.Binary)
	if a.IsValid(i) {
		return net.HardwareAddr(arr.Value(i)).String()
	}
	return nil
}

// MacType is a simple extension type that represents a BinaryType
// to be used for representing mac addresses.
type MacType struct {
	arrow.ExtensionBase
}

// NewMacType is a convenience function to create an instance of MacType
// with the correct storage type
func NewMacType() *MacType {
	return &MacType{ExtensionBase: arrow.ExtensionBase{Storage: &arrow.BinaryType{}}}
}

// ArrayType returns TypeOf(MacArray{}) for constructing MAC arrays
func (*MacType) ArrayType() reflect.Type {
	return reflect.TypeOf(MacArray{})
}

func (*MacType) ExtensionName() string {
	return "mac"
}

// Serialize returns "mac-serialized" for testing proper metadata passing
func (*MacType) Serialize() string {
	return "mac-serialized"
}

// Deserialize expects storageType to be FixedSizeBinaryType{ByteWidth: 16} and the data to be
// "mac-serialized" in order to correctly create a MacType for testing deserialize.
func (*MacType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "mac-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.BinaryType{}) {
		return nil, fmt.Errorf("invalid storage type for MacType: %s", storageType.Name())
	}
	return NewInetType(), nil
}

// ExtensionEquals returns true if both extensions have the same name
func (u *MacType) ExtensionEquals(other arrow.ExtensionType) bool {
	return u.ExtensionName() == other.ExtensionName()
}

func (*MacType) NewBuilder(bldr *array.ExtensionBuilder) array.Builder {
	return NewMacBuilder(bldr)
}
