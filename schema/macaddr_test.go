package schema

import (
	"testing"
)

func TestMacaddrSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result Macaddr
	}{
		{
			source: mustParseMacaddr(t, "01:23:45:67:89:ab"),
			result: Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
		},
		{
			source: "01:23:45:67:89:ab",
			result: Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r Macaddr
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestMacAddr_LessThan(t *testing.T) {
	macAddr1 := mustParseMacaddr(t, "01:23:45:67:89:ab")
	macAddr2 := mustParseMacaddr(t, "01:23:45:67:89:ac")
	cases := []struct {
		a    Macaddr
		b    Macaddr
		want bool
	}{
		{a: Macaddr{Addr: macAddr1, Status: Present}, b: Macaddr{Addr: macAddr2, Status: Present}, want: true},
		{a: Macaddr{Addr: macAddr2, Status: Present}, b: Macaddr{Addr: macAddr1, Status: Present}, want: false},
		{a: Macaddr{Addr: macAddr1, Status: Undefined}, b: Macaddr{Addr: macAddr1, Status: Present}, want: true},
		{a: Macaddr{Addr: macAddr1, Status: Present}, b: Macaddr{Addr: macAddr1, Status: Undefined}, want: false},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
