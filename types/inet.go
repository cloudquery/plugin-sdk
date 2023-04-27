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

func NewInetBuilder(builder *array.ExtensionBuilder) *InetBuilder {
	return &InetBuilder{ExtensionBuilder: builder}
}

func (b *InetBuilder) Append(v *net.IPNet) {
	if v == nil {
		b.AppendNull()
		return
	}
	b.ExtensionBuilder.Builder.(*array.StringBuilder).Append(v.String())
}

func (b *InetBuilder) UnsafeAppend(v *net.IPNet) {
	b.ExtensionBuilder.Builder.(*array.StringBuilder).UnsafeAppend([]byte(v.String()))
}

func (b *InetBuilder) AppendValues(v []*net.IPNet, valid []bool) {
	data := make([]string, len(v))
	for i, v := range v {
		if !valid[i] {
			continue
		}
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
	b.Append(data)
	return nil
}

func (b *InetBuilder) UnmarshalOne(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}

	var val *net.IPNet
	switch v := t.(type) {
	case string:
		_, val, err = net.ParseCIDR(v)
		if err != nil {
			return err
		}
	case []byte:
		_, val, err = net.ParseCIDR(string(v))
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

func (a *InetArray) String() string {
	arr := a.Storage().(*array.String)
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

func (a *InetArray) Value(i int) *net.IPNet {
	if a.IsNull(i) {
		return nil
	}
	_, ipnet, err := net.ParseCIDR(a.Storage().(*array.String).Value(i))
	if err != nil {
		panic(fmt.Errorf("invalid ip+net: %w", err))
	}

	return ipnet
}

func (a *InetArray) ValueStr(i int) string {
	switch {
	case a.IsNull(i):
		return array.NullValueStr
	default:
		return a.Value(i).String()
	}
}

func (a *InetArray) GetOneForMarshal(i int) any {
	if val := a.Value(i); val != nil {
		return val.String()
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
	return &InetType{ExtensionBase: arrow.ExtensionBase{Storage: &arrow.StringType{}}}
}

// ArrayType returns TypeOf(InetArray{}) for constructing Inet arrays
func (*InetType) ArrayType() reflect.Type {
	return reflect.TypeOf(InetArray{})
}

func (*InetType) ExtensionName() string {
	return "inet"
}

// Serialize returns "inet-serialized" for testing proper metadata passing
func (*InetType) Serialize() string {
	return "inet-serialized"
}

// Deserialize expects storageType to be StringType and the data to be
// "inet-serialized" in order to correctly create a InetType for testing deserialize.
func (*InetType) Deserialize(storageType arrow.DataType, data string) (arrow.ExtensionType, error) {
	if data != "inet-serialized" {
		return nil, fmt.Errorf("type identifier did not match: '%s'", data)
	}
	if !arrow.TypeEqual(storageType, &arrow.StringType{}) {
		return nil, fmt.Errorf("invalid storage type for InetType: %s", storageType.Name())
	}
	return NewInetType(), nil
}

// ExtensionEquals returns true if both extensions have the same name
func (u *InetType) ExtensionEquals(other arrow.ExtensionType) bool {
	return u.ExtensionName() == other.ExtensionName()
}

func (*InetType) NewBuilder(bldr *array.ExtensionBuilder) array.Builder {
	return NewInetBuilder(bldr)
}
