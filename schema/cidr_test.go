package schema

import (
	"net"
	"testing"
)

func TestCIDR_LessThan(t *testing.T) {
	_, ip1, _ := net.ParseCIDR("192.0.2.1/24")
	_, ip2, _ := net.ParseCIDR("192.0.3.1/24")
	cases := []struct {
		a    CIDR
		b    CIDR
		want bool
	}{
		{
			a:    CIDR{IPNet: ip1, Status: Present},
			b:    CIDR{IPNet: ip2, Status: Present},
			want: true,
		},
		{
			a:    CIDR{IPNet: ip2, Status: Present},
			b:    CIDR{IPNet: ip1, Status: Present},
			want: false,
		},
		{
			a:    CIDR{IPNet: ip1, Status: Present},
			b:    CIDR{IPNet: ip1, Status: Present},
			want: false,
		},
		{
			a:    CIDR{IPNet: nil, Status: Undefined},
			b:    CIDR{IPNet: ip1, Status: Present},
			want: true,
		},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
