package scalar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
)

type textMarshaler struct {
	Text string
}

func (t textMarshaler) MarshalText() (text []byte, err error) {
	return []byte(t.Text), err
}

// nolint:unparam
func mustParseCIDR(t testing.TB, s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		t.Fatal(err)
	}

	return ipnet
}

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

func TestInetSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result Inet
	}{
		{source: mustParseCIDR(t, "127.0.0.1/32"), result: Inet{Value: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
		{source: mustParseCIDR(t, "127.0.0.1/32").IP, result: Inet{Value: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
		{source: "127.0.0.1/32", result: Inet{Value: mustParseCIDR(t, "127.0.0.1/32"), Valid: true}},
		{source: "1.2.3.4/24", result: Inet{Value: &net.IPNet{IP: net.ParseIP("1.2.3.4").To4(), Mask: net.CIDRMask(24, 32)}, Valid: true}},
		{source: "10.0.0.1", result: Inet{Value: mustParseInet(t, "10.0.0.1"), Valid: true}},
		{source: "2607:f8b0:4009:80b::200e", result: Inet{Value: mustParseInet(t, "2607:f8b0:4009:80b::200e"), Valid: true}},
		{source: net.ParseIP(""), result: Inet{}},
		{source: "0.0.0.0/8", result: Inet{Value: mustParseInet(t, "0.0.0.0/8"), Valid: true}},
		{source: "::ffff:0.0.0.0/104", result: Inet{Value: &net.IPNet{IP: net.ParseIP("::ffff:0.0.0.0"), Mask: net.CIDRMask(104, 128)}, Valid: true}},
		{source: textMarshaler{"127.0.0.1"}, result: Inet{Value: mustParseInet(t, "127.0.0.1"), Valid: true}},
		{source: func(s string) fmt.Stringer {
			var b strings.Builder
			b.WriteString(s)
			return &b
		}("127.0.0.1"), result: Inet{Value: mustParseInet(t, "127.0.0.1"), Valid: true}},
		{source: &Inet{Value: &net.IPNet{IP: net.ParseIP("::ffff:0.0.0.0"), Mask: net.CIDRMask(104, 128)}, Valid: true}, result: Inet{Value: &net.IPNet{IP: net.ParseIP("::ffff:0.0.0.0"), Mask: net.CIDRMask(104, 128)}, Valid: true}},
		{source: (*net.IP)(nil), result: Inet{Value: nil, Valid: false}},
	}

	for i, tt := range successfulTests {
		var r Inet
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestInetMarshalUnmarshal(t *testing.T) {
	var r Inet
	err := r.Set("10.244.0.0/24")
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var r2 Inet
	err = json.Unmarshal(b, &r2)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Equal(&r2) {
		t.Errorf("%v != %v", r, r2)
	}

	// workaround this Golang bug: https://github.com/golang/go/issues/35727
	if !bytes.Equal(r.Value.Mask, r2.Value.Mask) {
		t.Errorf("%v != %v", r.Value.Mask, r2.Value.Mask)
	}
	if !net.IP.Equal(r.Value.IP, r2.Value.IP) {
		t.Errorf("%v != %v", r.Value.IP, r2.Value.IP)
	}
}
