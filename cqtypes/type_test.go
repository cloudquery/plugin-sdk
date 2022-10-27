package cqtypes

import (
	"errors"
	"net"
	"testing"
)

// Test for renamed types
type _string string
type _bool bool
type _int8 int8
type _int16 int16
type _int16Slice []int16
type _int32Slice []int32
type _int64Slice []int64
type _float32Slice []float32
type _float64Slice []float64
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
