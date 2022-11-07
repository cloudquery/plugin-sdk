package schema

import (
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

func mustParseCIDR(t testing.TB, s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		t.Fatal(err)
	}

	return ipnet
}

func TestInetSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result Inet
	}{
		{source: mustParseCIDR(t, "127.0.0.1/32"), result: Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present}},
		{source: mustParseCIDR(t, "127.0.0.1/32").IP, result: Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present}},
		{source: "127.0.0.1/32", result: Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present}},
		{source: "1.2.3.4/24", result: Inet{IPNet: &net.IPNet{IP: net.ParseIP("1.2.3.4").To4(), Mask: net.CIDRMask(24, 32)}, Status: Present}},
		{source: "10.0.0.1", result: Inet{IPNet: mustParseInet(t, "10.0.0.1"), Status: Present}},
		{source: "2607:f8b0:4009:80b::200e", result: Inet{IPNet: mustParseInet(t, "2607:f8b0:4009:80b::200e"), Status: Present}},
		{source: net.ParseIP(""), result: Inet{Status: Null}},
		{source: "0.0.0.0/8", result: Inet{IPNet: mustParseInet(t, "0.0.0.0/8"), Status: Present}},
		{source: "::ffff:0.0.0.0/104", result: Inet{IPNet: &net.IPNet{IP: net.ParseIP("::ffff:0.0.0.0"), Mask: net.CIDRMask(104, 128)}, Status: Present}},
		{source: textMarshaler{"127.0.0.1"}, result: Inet{IPNet: mustParseInet(t, "127.0.0.1"), Status: Present}},
		{source: func(s string) fmt.Stringer {
			var b strings.Builder
			b.WriteString(s)
			return &b
		}("127.0.0.1"), result: Inet{IPNet: mustParseInet(t, "127.0.0.1"), Status: Present}},
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
