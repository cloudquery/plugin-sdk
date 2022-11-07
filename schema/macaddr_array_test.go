package schema

import (
	"net"
	"testing"
)

func TestMacaddrArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result MacaddrArray
	}{
		{
			source: []net.HardwareAddr{mustParseMacaddr(t, "01:23:45:67:89:ab")},
			result: MacaddrArray{
				Elements:   []Macaddr{{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: (([]net.HardwareAddr)(nil)),
			result: MacaddrArray{Status: Null},
		},
		{
			source: [][]net.HardwareAddr{
				{mustParseMacaddr(t, "01:23:45:67:89:ab")},
				{mustParseMacaddr(t, "cd:ef:01:23:45:67")}},
			result: MacaddrArray{
				Elements: []Macaddr{
					{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
					{Addr: mustParseMacaddr(t, "cd:ef:01:23:45:67"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][][][]net.HardwareAddr{
				{{{
					mustParseMacaddr(t, "01:23:45:67:89:ab"),
					mustParseMacaddr(t, "cd:ef:01:23:45:67"),
					mustParseMacaddr(t, "89:ab:cd:ef:01:23")}}},
				{{{
					mustParseMacaddr(t, "45:67:89:ab:cd:ef"),
					mustParseMacaddr(t, "fe:dc:ba:98:76:54"),
					mustParseMacaddr(t, "32:10:fe:dc:ba:98")}}}},
			result: MacaddrArray{
				Elements: []Macaddr{
					{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
					{Addr: mustParseMacaddr(t, "cd:ef:01:23:45:67"), Status: Present},
					{Addr: mustParseMacaddr(t, "89:ab:cd:ef:01:23"), Status: Present},
					{Addr: mustParseMacaddr(t, "45:67:89:ab:cd:ef"), Status: Present},
					{Addr: mustParseMacaddr(t, "fe:dc:ba:98:76:54"), Status: Present},
					{Addr: mustParseMacaddr(t, "32:10:fe:dc:ba:98"), Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
		{
			source: [2][1]net.HardwareAddr{
				{mustParseMacaddr(t, "01:23:45:67:89:ab")},
				{mustParseMacaddr(t, "cd:ef:01:23:45:67")}},
			result: MacaddrArray{
				Elements: []Macaddr{
					{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
					{Addr: mustParseMacaddr(t, "cd:ef:01:23:45:67"), Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [2][1][1][3]net.HardwareAddr{
				{{{
					mustParseMacaddr(t, "01:23:45:67:89:ab"),
					mustParseMacaddr(t, "cd:ef:01:23:45:67"),
					mustParseMacaddr(t, "89:ab:cd:ef:01:23")}}},
				{{{
					mustParseMacaddr(t, "45:67:89:ab:cd:ef"),
					mustParseMacaddr(t, "fe:dc:ba:98:76:54"),
					mustParseMacaddr(t, "32:10:fe:dc:ba:98")}}}},
			result: MacaddrArray{
				Elements: []Macaddr{
					{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
					{Addr: mustParseMacaddr(t, "cd:ef:01:23:45:67"), Status: Present},
					{Addr: mustParseMacaddr(t, "89:ab:cd:ef:01:23"), Status: Present},
					{Addr: mustParseMacaddr(t, "45:67:89:ab:cd:ef"), Status: Present},
					{Addr: mustParseMacaddr(t, "fe:dc:ba:98:76:54"), Status: Present},
					{Addr: mustParseMacaddr(t, "32:10:fe:dc:ba:98"), Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r MacaddrArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
