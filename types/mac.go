package types

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/goccy/go-json"
)

type MACBuilder struct {
	*array.ExtensionBuilder
}

func NewMACBuilder(builder *array.ExtensionBuilder) *MACBuilder {
	return &MACBuilder{ExtensionBuilder: builder}
}

func (b *MACBuilder) Append(v net.HardwareAddr) {
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).Append(v[:])
}

func (b *MACBuilder) UnsafeAppend(v net.HardwareAddr) {
	b.ExtensionBuilder.Builder.(*array.BinaryBuilder).UnsafeAppend(v[:])
}

func (b *MACBuilder) AppendValues(v []net.HardwareAddr, valid []bool) {
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

func (b *MACBuilder) AppendValueFromString(s string) error {
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

func (b *MACBuilder) UnmarshalOne(dec *json.Decoder) error {
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

func (b *MACBuilder) Unmarshal(dec *json.Decoder) error {
	for dec.More() {
		if err := b.UnmarshalOne(dec); err != nil {
			return err
		}
	}
	return nil
}

func (b *MACBuilder) UnmarshalJSON(data []byte) error {
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

func (b *MACBuilder) NewMACArray() *MACArray {
	return b.NewExtensionArray().(*MACArray)
}

// MACArray is a simple array which is a wrapper around a BinaryArray
type MACArray struct {
	array.ExtensionArrayBase
}

func (a *MACArray) String() string {
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

func (a *MACArray) Value(i int) net.HardwareAddr {
	if a.IsNull(i) {
		return nil
	}
	b := a.Storage().(*array.Binary).Value(i)
	if len(b) == 0 {
		const minMACLen = 6
		return make(net.HardwareAddr, minMACLen)
	}

	return net.HardwareAddr(b)
}

func (a *MACArray) ValueStr(i int) string {
	switch {
	case a.IsNull(i):
		return array.NullValueStr
	default:
		return a.Value(i).String()
	}
}

func (a *MACArray) MarshalJSON() ([]byte, error) {
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

func (a *MACArray) GetOneForMarshal(i int) any {
	arr := a.Storage().(*array.Binary)
	if a.IsValid(i) {
		return net.HardwareAddr(arr.Value(i)).String()
	}
	return nil
}

// MACType is a simple extension type that represents a BinaryType
// to be used for representing MAC addresses.
type MACType struct {
	arrow.ExtensionBase
}

// NewMACType is a convenience function to create an instance of MACType
// with the correct storage type
func NewMACType() *MACType {
	return &MACType{ExtensionBase: arrow.ExtensionBase{Storage: &arrow.BinaryType{}}}
}

// ArrayType returns TypeOf(MACArray{}) for constructing MAC arrays
func (*MACType) ArrayType() reflect.Type {
	return reflect.TypeOf(MACArray{})
}

func (*MACType) ExtensionName() string {
	return "mac"
}

func (*MACType) String() string {
	return "mac"
}

// Serialize returns "mac-serialized" for testing proper metadata passing
func (*MACType) Serialize() string {
	return "mac-serialized"
}

// Deserialize expects storageType to be FixedSizeBinaryType{ByteWidth: 16} and the data to be
// "MAC-serialized" in order to correctly create a MACType for testing deserialize.
func (*MACType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "mac-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.BinaryType{}) {
		return nil, fmt.Errorf("invalid storage type for MACType: %s", storageType.Name())
	}
	return NewMACType(), nil
}

// ExtensionEquals returns true if both extensions have the same name
func (u *MACType) ExtensionEquals(other arrow.ExtensionType) bool {
	return u.ExtensionName() == other.ExtensionName()
}

func (*MACType) NewBuilder(bldr *array.ExtensionBuilder) array.Builder {
	return NewMACBuilder(bldr)
}
