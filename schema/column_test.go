package schema

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

type validateFixture struct {
	Column     Column
	TestValues []interface{}
	BadValues  []interface{}
}

type SomeString string

type SomeInt int

type SomeInt16 int16

var validateFixtures = []validateFixture{
	{
		Column:     Column{Type: TypeInt},
		TestValues: []interface{}{5, 300, funk.PtrOf(555), SomeInt16(555), SomeInt(555)},
		BadValues:  []interface{}{"a", funk.PtrOf("abc")},
	},
	{
		Column:     Column{Type: TypeFloat},
		TestValues: []interface{}{555.5},
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
		Column:     Column{Type: TypeTimeInterval},
		TestValues: []interface{}{time.Minute},
		BadValues:  []interface{}{int64(2), -3, 5.0, time.Now()},
	},
	{
		Column:     Column{Type: TypeUUID},
		TestValues: []interface{}{uuid.New(), uuid.New().String()},
		BadValues:  []interface{}{uuid.New().String()[:5], 5555555},
	},
	{
		Column:     Column{Type: TypeJSON},
		TestValues: []interface{}{make(map[string]interface{}), make(map[string]string), []byte{11, 11, 11, 11}, []interface{}{struct{ Test int }{Test: 1}}, []Column{{Name: "test"}}},
	},
	{
		Column:     Column{Type: TypeBool},
		TestValues: []interface{}{true, false, funk.PtrOf(false), funk.PtrOf(true)},
	},
	{
		Column:     Column{Type: TypeIntArray},
		TestValues: []interface{}{[]int{1, 2, 3}, []SomeInt{SomeInt(3)}, []int16{1, 2, 3}},
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

func GenerateMac() net.HardwareAddr {
	hw, err := net.ParseMAC(`00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01`)
	if err != nil {
		panic(err)
	}
	return hw
}
func GenerateMacPtr() *net.HardwareAddr {
	r := GenerateMac()
	return &r
}

func GenerateIPv4Ptr() *net.IP {
	r := net.ParseIP("127.0.0.1")
	return &r
}

func GenerateIPv6Ptr() *net.IP {
	r := net.ParseIP("2001:db8::68")
	return &r
}

func GenerateCIDR() *net.IPNet {
	ip := GenerateIPv4Ptr()
	_, mask, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.String(), rand.Int31n(16)+16))
	return mask
}

// func TestValidateType(t *testing.T) {
// 	for _, f := range validateFixtures {
// 		t.Run(f.Column.Type.String(), func(t *testing.T) {
// 			for _, v := range f.TestValues {
// 				assert.Nil(t, f.Column.ValidateType(v))
// 			}
// 			for _, v := range f.BadValues {
// 				assert.Error(t, f.Column.ValidateType(v))
// 			}
// 		})
// 	}
// }


// func TestColumnJsonMarshal(t *testing.T) {
// 	// we are testing column json marshalling to make sure
// 	// this can be easily sent over the wire
// 	expected := Column{
// 		Name: "test",
// 		Type: TypeJSON,
// 	}
// 	b, err := json.Marshal(expected)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	got := Column{}
// 	if err := json.Unmarshal(b, &got); err != nil {
// 		t.Fatal(err)
// 	}
// 	if !reflect.DeepEqual(expected, got) {
// 		t.Fatalf("expected %v got %v", expected, got)
// 	}
// }
