package schema

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type innerStruct struct {
	Value string
}

type testStruct struct {
	Inner      innerStruct
	Value      int
	unexported bool
}

var pathTestTable = &Table{
	Columns: []Column{
		{
			Name: "test",
			Type: TypeString,
		},
		{
			Name: "int_value",
			Type: TypeInt,
		},
		{
			Name: "unexported",
			Type: TypeBool,
		},
	},
}

func TestPathResolver(t *testing.T) {
	r1 := PathResolver("Inner.Value")
	r2 := PathResolver("Value")
	r3 := PathResolver("unexported")
	resource := NewResourceData(pathTestTable, nil, testStruct{Inner: innerStruct{Value: "bla"}, Value: 5, unexported: false}, nil)
	err := r1(context.TODO(), nil, resource, Column{Name: "test"})

	assert.Nil(t, err)
	assert.Equal(t, resource.Get("test"), "bla")

	err = r2(context.TODO(), nil, resource, Column{Name: "int_value"})

	assert.Nil(t, err)
	assert.Equal(t, resource.Get("int_value"), 5)

	err = r3(context.TODO(), nil, resource, Column{Name: "unexported"})
	assert.Nil(t, err)
	assert.Nil(t, resource.Get("unexported"))
}

func TestInterfaceSlice(t *testing.T) {
	var sType []interface{}
	var names []string
	names = append(names, "first", "second")
	assert.IsTypef(t, sType, interfaceSlice(names), "")
	assert.IsTypef(t, sType, interfaceSlice(&names), "")
	assert.IsTypef(t, sType, interfaceSlice(1), "")
	assert.IsTypef(t, sType, interfaceSlice(innerStruct{"asdsa"}), "")
	assert.IsTypef(t, sType, interfaceSlice(&innerStruct{"asdsa"}), "")
	pSlice := []*innerStruct{{"asdsa"}, {"asdsa"}, {"asdsa"}}
	assert.IsTypef(t, sType, interfaceSlice(pSlice), "")
	assert.IsTypef(t, sType, interfaceSlice(&pSlice), "")

}

var dateTestTable = &Table{
	Columns: []Column{
		{
			Name: "date",
			Type: TypeTimestamp,
		},
	},
}

type testDateStruct struct {
	Date string
}

func TestDateTimeResolver(t *testing.T) {
	r1 := DateResolver("Date")
	resource := NewResourceData(dateTestTable, nil, testDateStruct{Date: "2011-10-05T14:48:00.000Z"}, nil)
	err := r1(context.TODO(), nil, resource, Column{Name: "date"})

	assert.Nil(t, err)
	t1 := time.Date(2011, 10, 5, 14, 48, 0, 0, time.UTC)
	assert.Equal(t, resource.Get("date"), &t1)

	r2 := DateResolver("Date", time.RFC822)
	resource = NewResourceData(dateTestTable, nil, testDateStruct{Date: "2011-10-05T14:48:00.000Z"}, nil)
	err = r2(context.TODO(), nil, resource, Column{Name: "date"})

	assert.Error(t, err)

	resource = NewResourceData(dateTestTable, nil, testDateStruct{Date: "03 Jan 06 15:04 EST"}, nil)
	err = r2(context.TODO(), nil, resource, Column{Name: "date"})
	assert.Nil(t, err)

	t2 := time.Date(2006, 1, 3, 15, 4, 0, 0, time.UTC)
	assert.Equal(t, t2.Unix(), resource.Get("date").(*time.Time).UTC().Unix())

	r3 := DateResolver("Date", time.RFC822, "2006-01-02")
	resource = NewResourceData(dateTestTable, nil, testDateStruct{Date: "2011-10-05"}, nil)
	err = r3(context.TODO(), nil, resource, Column{Name: "date"})
	assert.Nil(t, err)

	t3 := time.Date(2011, 10, 5, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, resource.Get("date"), &t3)
}

var networkTestTable = &Table{
	Columns: []Column{
		{
			Name: "ip",
			Type: TypeInet,
		},
		{
			Name: "mac",
			Type: TypeMacAddr,
		},
		{
			Name: "net",
			Type: TypeCIDR,
		},
	},
}

type testNetStruct struct {
	IP  string
	MAC string
	Net string
}

var netTests = []testNetStruct{
	{IP: "192.168.1.12", MAC: "2C:54:91:88:C9:E3", Net: "192.168.0.1/24"},
	{IP: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", MAC: "2C-54-91-88-C9-E3", Net: "2002::1234:abcd:ffff:c0a8:101/64"},
	{IP: "::1234:5678", MAC: "2C-54-91-88-C9-E3", Net: "::1234:5678/12"},
}
var netTestsFails = []testNetStruct{
	{IP: "192.168.1/", MAC: "2C:54:91:88:C9", Net: "192.168.0.1-24"},
	{IP: "::1234:5678:", MAC: "2C:54-91-88-C9-E3", Net: "2002::1234:abcd:ffff:c0a8:101-64"},
}

func TestNetResolvers(t *testing.T) {
	r1 := IPAddressResolver("IP")
	r2 := MACAddressResolver("MAC")
	r3 := IPNetResolver("Net")
	for _, r := range netTests {
		resource := NewResourceData(networkTestTable, nil, r, nil)
		err := r1(context.TODO(), nil, resource, Column{Name: "ip"})
		assert.Nil(t, err)
		err = r2(context.TODO(), nil, resource, Column{Name: "mac"})
		assert.Nil(t, err)
		err = r3(context.TODO(), nil, resource, Column{Name: "net"})
		assert.Nil(t, err)
	}
	for _, r := range netTestsFails {
		resource := NewResourceData(networkTestTable, nil, r, nil)
		err := r1(context.TODO(), nil, resource, Column{Name: "ip"})
		assert.Error(t, err)
		err = r2(context.TODO(), nil, resource, Column{Name: "mac"})
		assert.Error(t, err)
		err = r3(context.TODO(), nil, resource, Column{Name: "net"})
		assert.Error(t, err)
	}
}

var TransformersTestTable = &Table{
	Columns: []Column{
		{
			Name: "string_to_int",
			Type: TypeInt,
		},
		{
			Name: "float_to_int",
			Type: TypeInt,
		},
		{
			Name: "int_to_string",
			Type: TypeString,
		},
		{
			Name: "float_to_string",
			Type: TypeString,
		},
	},
}

type testTransformersStruct struct {
	Int      int
	String   string
	Float    float64
	BadFloat string
}

func TestTransformersResolvers(t *testing.T) {
	r1 := StringResolver("Int")
	r2 := StringResolver("Float")
	r3 := IntResolver("String")
	r4 := IntResolver("Float")
	r5 := IntResolver("BadFloat")
	resource := NewResourceData(TransformersTestTable, nil, testTransformersStruct{Int: 10, Float: 10.2, String: "123", BadFloat: "10,1"}, nil)
	err := r1(context.TODO(), nil, resource, Column{Name: "int_to_string"})
	assert.Nil(t, err)
	assert.Equal(t, resource.Get("int_to_string"), "10")

	err = r2(context.TODO(), nil, resource, Column{Name: "float_to_string"})
	assert.Nil(t, err)
	assert.Equal(t, resource.Get("float_to_string"), "10.2")

	err = r3(context.TODO(), nil, resource, Column{Name: "string_to_int"})
	assert.Nil(t, err)
	assert.Equal(t, resource.Get("string_to_int"), 123)

	err = r4(context.TODO(), nil, resource, Column{Name: "float_to_int"})
	assert.Nil(t, err)
	assert.Equal(t, resource.Get("float_to_int"), 10)

	err = r5(context.TODO(), nil, resource, Column{Name: "float_to_int"})
	assert.Error(t, err)
}

var UUIDTestTable = &Table{
	Columns: []Column{
		{
			Name: "uuid",
			Type: TypeUUID,
		},
	},
}

type testUUIDStruct struct {
	UUID    string
	BadUUID string
}

func TestUUIDResolver(t *testing.T) {
	r1 := UUIDResolver("UUID")
	r2 := UUIDResolver("BadUUID")
	resource := NewResourceData(UUIDTestTable, nil, testUUIDStruct{UUID: "123e4567-e89b-12d3-a456-426614174000", BadUUID: "123e4567-e89b-12d3-a456-4266141740001"}, nil)

	err := r1(context.TODO(), nil, resource, Column{Name: "uuid"})
	assert.Nil(t, err)
	uuid, err := uuid.FromString("123e4567-e89b-12d3-a456-426614174000")
	assert.Nil(t, err)
	assert.Equal(t, uuid, resource.Get("uuid"))

	err = r2(context.TODO(), nil, resource, Column{Name: "uuid"})
	assert.Error(t, err)
}
