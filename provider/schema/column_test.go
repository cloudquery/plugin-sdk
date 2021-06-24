package schema

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/cloudquery/faker/v3"

	"github.com/thoas/go-funk"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type validateFixture struct {
	Column     Column
	TestValues []interface{}
	BadValues  []interface{}
}

type SomeString string

type SomeInt int

type SomeInt16 int16

func GenerateMac() net.HardwareAddr {
	mac, _ := net.ParseMAC(faker.MacAddress())
	return mac
}
func GenerateMacPtr() *net.HardwareAddr {
	mac, _ := net.ParseMAC(faker.MacAddress())
	return &mac
}

func GenerateIPv4Ptr() *net.IP {
	ip := net.ParseIP(faker.IPv4())
	return &ip
}

func GenerateIPv6Ptr() *net.IP {
	ip := net.ParseIP(faker.IPv6())
	return &ip
}

func GenerateCIDR() *net.IPNet {
	_, mask, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", faker.IPv4(), rand.Int31n(16)+16))
	return mask
}

var validateFixtures = []validateFixture{
	{
		Column:     Column{Type: TypeBigInt},
		TestValues: []interface{}{5, 300, funk.PtrOf(555), SomeInt(555)},
		BadValues:  []interface{}{"a", funk.PtrOf("abc"), SomeInt16(555)},
	},
	{
		Column:     Column{Type: TypeSmallInt},
		TestValues: []interface{}{SomeInt16(555)},
		BadValues:  []interface{}{"a", funk.PtrOf("abc")},
	},
	{
		Column:     Column{Type: TypeString},
		TestValues: []interface{}{"abcd", "aaaaaaa", funk.PtrOf("asdasd"), SomeString("Asda")},
		BadValues:  []interface{}{funk.PtrOf(555), 555, time.Now()},
	},
	{
		Column:     Column{Type: TypeTimestamp},
		TestValues: []interface{}{time.Now()},
	},
	{
		Column:     Column{Type: TypeUUID},
		TestValues: []interface{}{uuid.New(), uuid.New().String()},
		BadValues:  []interface{}{uuid.New().String()[:5], 5555555},
	},
	{
		Column:     Column{Type: TypeJSON},
		TestValues: []interface{}{make(map[string]interface{}), make(map[string]string), []byte{11, 11, 11, 11}},
	},
	{
		Column:     Column{Type: TypeBool},
		TestValues: []interface{}{true, false, funk.PtrOf(false), funk.PtrOf(true)},
	},
	{
		Column:     Column{Type: TypeIntArray},
		TestValues: []interface{}{[]int{1, 2, 3}, []SomeInt{SomeInt(3)}},
		BadValues:  []interface{}{[]interface{}{1, 2, 3}},
	},
	{
		Column:     Column{Type: TypeStringArray},
		TestValues: []interface{}{[]string{"a", "b", "c"}, []*string{funk.PtrOf("a").(*string)}, []SomeString{SomeString("lol")}},
		BadValues:  []interface{}{[]interface{}{1, 2, 3}},
	},
	{
		Column:     Column{Type: TypeMacAddr},
		TestValues: []interface{}{GenerateMac(), GenerateMac(), GenerateMacPtr()},
		BadValues:  []interface{}{"asdasdsadads", -55, 44, "00:33:44:55:77:55"},
	},
	{
		Column:     Column{Type: TypeMacAddrArray},
		TestValues: []interface{}{[]net.HardwareAddr{GenerateMac(), GenerateMac()}, []*net.HardwareAddr{GenerateMacPtr(), GenerateMacPtr()}},
		BadValues:  []interface{}{"asdasdsadads", -55, 44, "00:33:44:55:77:55"},
	},
	{
		Column:     Column{Type: TypeInet},
		TestValues: []interface{}{net.ParseIP("127.0.0.1"), net.ParseIP("2b15:800f:a66b:0:1278:b7ad:6052:f444"), GenerateIPv4Ptr(), GenerateIPv6Ptr()},
		BadValues:  []interface{}{"asdasdsadads", "127.0.0.1", "333"},
	},

	{
		Column:     Column{Type: TypeInetArray},
		TestValues: []interface{}{[]net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("2b15:800f:a66b:0:1278:b7ad:6052:f444")}, []*net.IP{GenerateIPv4Ptr(), GenerateIPv6Ptr()}},
		BadValues:  []interface{}{"asdasdsadads", "127.0.0.1", net.ParseIP("127.0.0.1"), []*net.HardwareAddr{GenerateMacPtr(), GenerateMacPtr()}},
	},
	{
		Column:     Column{Type: TypeCIDR},
		TestValues: []interface{}{GenerateCIDR(), GenerateCIDR()},
		BadValues:  []interface{}{"asdasdsadads", 555, "127.0.0.1/24", net.IP{}},
	},
	{
		Column:     Column{Type: TypeCIDRArray},
		TestValues: []interface{}{[]*net.IPNet{GenerateCIDR(), GenerateCIDR()}, []*net.IPNet{}, []net.IPNet{}},
		BadValues:  []interface{}{"asdasdsadads", 555, "127.0.0.1/24", net.IPNet{}, net.IP{}},
	},
}

func TestValidateType(t *testing.T) {
	for _, f := range validateFixtures {
		for _, v := range f.TestValues {
			assert.Nil(t, f.Column.ValidateType(v))
		}
		for _, v := range f.BadValues {
			assert.Error(t, f.Column.ValidateType(v))
		}
	}
}

func TestValueTypeFromString(t *testing.T) {
	assert.Equal(t, ValueTypeFromString("String"), TypeString)
	// case insensitive
	assert.Equal(t, ValueTypeFromString("Json"), TypeJSON)
	assert.Equal(t, ValueTypeFromString("JSON"), TypeJSON)
	assert.Equal(t, ValueTypeFromString("bigint"), TypeBigInt)
	assert.Equal(t, ValueTypeFromString("Blabla"), TypeInvalid)
}

func BenchmarkColumn_ValidateTypeInt(b *testing.B) {
	col := Column{Type: TypeInt}
	for n := 0; n < b.N; n++ {
		_ = col.ValidateType(555)
	}
}

func BenchmarkColumn_ValidateTypeString(b *testing.B) {
	col := Column{Type: TypeString}
	for n := 0; n < b.N; n++ {
		_ = col.ValidateType("Asdsad")
	}
}

func BenchmarkColumn_ValidateTypeBadValue(b *testing.B) {
	col := Column{Type: TypeInt}
	for n := 0; n < b.N; n++ {
		_ = col.ValidateType("Asdsad")
	}
}

func BenchmarkColumn_ValidateTypeMap(b *testing.B) {
	col := Column{Type: TypeInt}
	m := make(map[string]interface{})
	for n := 0; n < b.N; n++ {
		_ = col.ValidateType(m)
	}
}
