// nolint:dupl
package schema

import (
	"net"
	"testing"
)

func TestInetArraySet(t *testing.T) {
	successfulTests := []struct {
		source any
		result InetArray
	}{
		{
			source: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
			result: InetArray{
				Elements:   []Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: (([]*net.IPNet)(nil)),
			result: InetArray{Status: Null},
		},
		{
			source: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
			result: InetArray{
				Elements:   []Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: (([]net.IP)(nil)),
			result: InetArray{Status: Null},
		},
		{
			source: [][]net.IP{{mustParseCIDR(t, "127.0.0.1/32").IP}, {mustParseCIDR(t, "10.0.0.1/32").IP}},
			result: InetArray{
				Elements: []Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present},
					{IPNet: mustParseCIDR(t, "10.0.0.1/32"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][][][]*net.IPNet{
				{{{
					mustParseCIDR(t, "127.0.0.1/24"),
					mustParseCIDR(t, "10.0.0.1/24"),
					mustParseCIDR(t, "172.16.0.1/16")}}},
				{{{
					mustParseCIDR(t, "192.168.0.1/16"),
					mustParseCIDR(t, "224.0.0.1/24"),
					mustParseCIDR(t, "169.168.0.1/16")}}}},
			result: InetArray{
				Elements: []Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "172.16.0.1/16"), Status: Present},
					{IPNet: mustParseCIDR(t, "192.168.0.1/16"), Status: Present},
					{IPNet: mustParseCIDR(t, "224.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "169.168.0.1/16"), Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
		{
			source: [2][1]net.IP{{mustParseCIDR(t, "127.0.0.1/32").IP}, {mustParseCIDR(t, "10.0.0.1/32").IP}},
			result: InetArray{
				Elements: []Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: Present},
					{IPNet: mustParseCIDR(t, "10.0.0.1/32"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [2][1][1][3]*net.IPNet{
				{{{
					mustParseCIDR(t, "127.0.0.1/24"),
					mustParseCIDR(t, "10.0.0.1/24"),
					mustParseCIDR(t, "172.16.0.1/16")}}},
				{{{
					mustParseCIDR(t, "192.168.0.1/16"),
					mustParseCIDR(t, "224.0.0.1/24"),
					mustParseCIDR(t, "169.168.0.1/16")}}}},
			result: InetArray{
				Elements: []Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "172.16.0.1/16"), Status: Present},
					{IPNet: mustParseCIDR(t, "192.168.0.1/16"), Status: Present},
					{IPNet: mustParseCIDR(t, "224.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "169.168.0.1/16"), Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r InetArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestInetArray_Size(t *testing.T) {
	tests := []struct {
		source InetArray
		want   int
	}{
		{
			source: InetArray{
				Elements: []Inet{
					{IPNet: mustParseCIDR(t, "127.0.0.1/24"), Status: Present},
					{IPNet: mustParseCIDR(t, "10.0.0.1/24"), Status: Present},
				},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1}},
				Status: Present,
			},
			want: 16,
		},
	}
	for _, tt := range tests {
		if tt.source.Size() != tt.want {
			t.Errorf("%v.Size() = %d, want %v", tt.source, tt.source.Size(), tt.want)
		}
	}
}
