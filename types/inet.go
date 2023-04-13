package types

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/goccy/go-json"
)

type InetBuilder struct {
	*array.ExtensionBuilder
}

func NewInetBuilder(bldr *array.ExtensionBuilder) *InetBuilder {
	b := &InetBuilder{
		ExtensionBuilder: bldr,
	}
	return b
}

func (b *InetBuilder) Append(v net.IPNet) {
	b.ExtensionBuilder.Builder.(*array.StringBuilder).Append(v.String())
}

func (b *InetBuilder) UnsafeAppend(v net.IPNet) {
	b.ExtensionBuilder.Builder.(*array.StringBuilder).UnsafeAppend([]byte(v.String()))
}

func (b *InetBuilder) AppendValues(v []net.IPNet, valid []bool) {
	data := make([]string, len(v))
	for i, v := range v {
		data[i] = v.String()
	}
	b.ExtensionBuilder.Builder.(*array.StringBuilder).AppendValues(data, valid)
}

func (b *InetBuilder) AppendValueFromString(s string) error {
	if s == array.NullValueStr {
		b.AppendNull()
		return nil
	}
	_, data, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}
	b.Append(*data)
	return nil
}

func (b *InetBuilder) UnmarshalOne(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}

	var val net.IPNet
	switch v := t.(type) {
	case string:
		_, data, err := net.ParseCIDR(v)
		if err != nil {
			return err
		}
		val = *data
	case []byte:
		_, data, err := net.ParseCIDR(string(v))
		if err != nil {
			return err
		}
		val = *data
	case nil:
		b.AppendNull()
		return nil
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprint(t),
			Type:   reflect.TypeOf([]byte{}),
			Offset: dec.InputOffset(),
			Struct: "String",
		}
	}

	b.Append(val)
	return nil
}

func (b *InetBuilder) Unmarshal(dec *json.Decoder) error {
	for dec.More() {
		if err := b.UnmarshalOne(dec); err != nil {
			return err
		}
	}
	return nil
}

func (b *InetBuilder) UnmarshalJSON(data []byte) error {
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

// InetArray is a simple array which is a FixedSizeBinary(16)
type InetArray struct {
	array.ExtensionArrayBase
}

func (a InetArray) String() string {
	arr := a.Storage().(*array.String)
	o := new(strings.Builder)
	o.WriteString("[")
	for i := 0; i < arr.Len(); i++ {
		if i > 0 {
			o.WriteString(" ")
		}
		switch {
		case a.IsNull(i):
			o.WriteString("(null)")
		default:
			fmt.Fprintf(o, "\"%s\"", arr.Value(i))
		}
	}
	o.WriteString("]")
	return o.String()
}

func (a *InetArray) ValueStr(i int) string {
	arr := a.Storage().(*array.String)
	switch {
	case a.IsNull(i):
		return "(null)"
	default:
		return arr.Value(i)
	}
}

func (a *InetArray) GetOneForMarshal(i int) any {
	arr := a.Storage().(*array.String)
	if a.IsValid(i) {
		_, ipnet, err := net.ParseCIDR(arr.Value(i))
		if err != nil {
			panic(fmt.Errorf("invalid ip+net: %w", err))
		}
		return ipnet.String()
	}
	return nil
}

// InetType is a simple extension type that represents a StringType
// to be used for representing IP Addresses and CIDRs
type InetType struct {
	arrow.ExtensionBase
}

// NewInetType is a convenience function to create an instance of InetType
// with the correct storage type
func NewInetType() *InetType {
	return &InetType{
		ExtensionBase: arrow.ExtensionBase{
			Storage: &arrow.StringType{}}}
}

func (InetType) ArrayType() reflect.Type {
	return reflect.TypeOf(InetArray{})
}

func (InetType) ExtensionName() string {
	return "inet"
}

// Serialize returns "inet-serialized" for testing proper metadata passing
func (InetType) Serialize() string {
	return "inet-serialized"
}

// Deserialize expects storageType to be StringType and the data to be
// "inet-serialized" in order to correctly create a InetType for testing deserialize.
func (InetType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "inet-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.StringType{}) {
		return nil, fmt.Errorf("invalid storage type for InetType: %s", storageType.Name())
	}
	return NewInetType(), nil
}

// InetType are equal if both are named "inet"
func (u InetType) ExtensionEquals(other arrow.ExtensionType) bool {
	return u.ExtensionName() == other.ExtensionName()
}

func (InetType) NewBuilder(bldr *array.ExtensionBuilder) array.Builder {
	return NewInetBuilder(bldr)
}
