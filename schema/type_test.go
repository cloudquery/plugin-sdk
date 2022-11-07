package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
)

// Test for renamed types
type _string string
type _bool bool
type _int8 int8

// type _int16 int16
// type _int16Slice []int16
// type _int32Slice []int32
// type _int64Slice []int64
// type _float32Slice []float32
// type _float64Slice []float64
type _byteSlice []byte

func mustParseInet(t testing.TB, s string) *net.IPNet {
	ip, ipnet, err := net.ParseCIDR(s)
	if err == nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			ipnet.IP = ipv4
		}
		return ipnet
	}

	// May be bare IP address.
	//
	ip = net.ParseIP(s)
	if ip == nil {
		t.Fatal(errors.New("unable to parse inet address"))
	}
	ipnet = &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
	if ipv4 := ip.To4(); ipv4 != nil {
		ipnet.IP = ipv4
		ipnet.Mask = net.CIDRMask(32, 32)
	}
	return ipnet
}

func mustParseMacaddr(t testing.TB, s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		t.Fatal(err)
	}

	return addr
}

func TestCQTypesMarshal(t *testing.T) {
	cqTypes := CQTypes{
		&Bool{Bool: true, Status: Present},
	}
	b, err := json.Marshal(cqTypes)
	fmt.Println(string(b))
	if err != nil {
		t.Fatal(err)
	}
	var res CQTypes
	if err := json.Unmarshal(b, &res); err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1, got %d", len(res))
	}

	if err := json.Unmarshal([]byte(`[{"type": "Bool", "value": {"Bool":true,"Status":1}}, {"type": "UnknownType"}]`), &res); err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1, got %d", len(res))
	}
}

func TestValueTypeFromOverTheWireString(t *testing.T) {
	for i := TypeInvalid + 1; i < TypeEnd; i++ {
		if deprecatedTypesValues.isDeprecated(i) {
			continue
		}
		v := valueTypeFromOverTheWireString(i.overTheWireString())
		if v != i {
			t.Fatalf("expected %d, got %d", i, v)
		}
	}
}

func TestValueTypeString(t *testing.T) {
	for i := TypeInvalid + 1; i < TypeEnd; i++ {
		if deprecatedTypesValues.isDeprecated(i) {
			continue
		}
		if i.String() == "TypeInvalid" || strings.HasPrefix(i.String(), "Unknown") {
			t.Fatalf("invalid string %s for type %d", i, i)
		}
	}
}

func TestAllTypesRegistered(t *testing.T) {
	for i := TypeInvalid + 1; i < TypeEnd; i++ {
		if deprecatedTypesValues.isDeprecated(i) {
			continue
		}
		v := NewCqTypeFromValueType(i)
		if v.Type() != i {
			t.Fatalf("expected %d, got %d", i, v.Type())
		}
	}
}
