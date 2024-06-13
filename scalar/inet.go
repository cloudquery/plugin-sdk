package scalar

import (
	"encoding"
	"fmt"
	"net"
	"strings"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type Inet struct {
	Valid bool
	Value *net.IPNet
}

func (s *Inet) IsValid() bool {
	return s.Valid
}

func (*Inet) DataType() arrow.DataType {
	return types.ExtensionTypes.Inet
}

func (s *Inet) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Inet)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value.String() == r.Value.String()
}

func (s *Inet) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.Value.String()
}

func (s *Inet) Get() any {
	return s.Value
}

func (s *Inet) Set(val any) error {
	if val == nil {
		return nil
	}

	if sc, ok := val.(Scalar); ok {
		if !sc.IsValid() {
			s.Valid = false
			return nil
		}
		return s.Set(sc.Get())
	}

	switch value := val.(type) {
	case net.IPNet:
		s.Value = &value
	case net.IP:
		if len(value) == 0 {
			return nil
		}
		bitCount := len(value) * 8
		mask := net.CIDRMask(bitCount, bitCount)
		s.Value = &net.IPNet{Mask: mask, IP: value}
	case string:
		ip, ipnet, err := net.ParseCIDR(value)
		if err != nil {
			ip := net.ParseIP(value)
			if ip == nil {
				return &ValidationError{Type: types.ExtensionTypes.Inet, Msg: "cannot parse string as IP", Value: value}
			}

			if ipv4 := maybeGetIPv4(value, ip); ipv4 != nil {
				ipnet = &net.IPNet{IP: ipv4, Mask: net.CIDRMask(32, 32)}
			} else {
				ipnet = &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
			}
		} else {
			ipnet.IP = ip
			if ipv4 := maybeGetIPv4(value, ipnet.IP); ipv4 != nil {
				ipnet.IP = ipv4
				if len(ipnet.Mask) == 16 {
					ipnet.Mask = ipnet.Mask[12:] // Not sure this is ever needed.
				}
			}
		}
		s.Value = ipnet
	case *net.IPNet:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *net.IP:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *string:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	default:
		if tv, ok := value.(encoding.TextMarshaler); ok {
			text, err := tv.MarshalText()
			if err != nil {
				return &ValidationError{Type: types.ExtensionTypes.UUID, Msg: "cannot marshal text", Err: err, Value: val}
			}
			return s.Set(string(text))
		}
		if sv, ok := value.(fmt.Stringer); ok {
			return s.Set(sv.String())
		}
		if originalSrc, ok := underlyingPtrType(val); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: types.ExtensionTypes.UUID, Msg: noConversion, Value: value}
	}
	s.Valid = true
	return nil
}

// Convert the net.IP to IPv4, if appropriate.
//
// When parsing a string to a net.IP using net.ParseIP() and the like, we get a
// 16 byte slice for IPv4 addresses as well as IPv6 addresses. This function
// calls To4() to convert them to a 4 byte slice. This is useful as it allows
// users of the net.IP check for IPv4 addresses based on the length and makes
// it clear we are handling IPv4 as opposed to IPv6 or IPv4-mapped IPv6
// addresses.
func maybeGetIPv4(input string, ip net.IP) net.IP {
	// Do not do this if the provided input looks like IPv6. This is because
	// To4() on IPv4-mapped IPv6 addresses converts them to IPv4, which behave
	// different in some cases.
	if strings.Contains(input, ":") {
		return nil
	}

	return ip.To4()
}
