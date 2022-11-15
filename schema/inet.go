//nolint:revive
package schema

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type InetTransformer interface {
	TransformInet(*Inet) interface{}
}

// workaround this Golang bug: https://github.com/golang/go/issues/35727
type inetWrapper struct {
	IPNet  *net.IPNet
	Status Status
}

// Inet represents both inet and cidr PostgreSQL types.
type Inet struct {
	IPNet  *net.IPNet
	Status Status
}

func (*Inet) Type() ValueType {
	return TypeInet
}

func (dst *Inet) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Inet)
	if !ok {
		return false
	}
	return dst.Status == s.Status && dst.IPNet.String() == s.IPNet.String()
}

func (dst *Inet) String() string {
	if dst.Status == Present {
		return dst.IPNet.String()
	} else {
		return ""
	}
}

func (dst *Inet) Set(src interface{}) error {
	if src == nil {
		*dst = Inet{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case net.IPNet:
		*dst = Inet{IPNet: &value, Status: Present}
	case net.IP:
		if len(value) == 0 {
			*dst = Inet{Status: Null}
		} else {
			bitCount := len(value) * 8
			mask := net.CIDRMask(bitCount, bitCount)
			*dst = Inet{IPNet: &net.IPNet{Mask: mask, IP: value}, Status: Present}
		}
	case string:
		ip, ipnet, err := net.ParseCIDR(value)
		if err != nil {
			ip := net.ParseIP(value)
			if ip == nil {
				return fmt.Errorf("unable to parse inet address: %s", value)
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

		*dst = Inet{IPNet: ipnet, Status: Present}
	case *net.IPNet:
		if value == nil {
			*dst = Inet{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *net.IP:
		if value == nil {
			*dst = Inet{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Inet{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if tv, ok := src.(encoding.TextMarshaler); ok {
			text, err := tv.MarshalText()
			if err != nil {
				return fmt.Errorf("cannot marshal %v: %w", value, err)
			}
			return dst.Set(string(text))
		}
		if sv, ok := src.(fmt.Stringer); ok {
			return dst.Set(sv.String())
		}
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Inet", value)
	}

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

// workaround this Golang bug: https://github.com/golang/go/issues/35727
func (dst *Inet) UnmarshalJSON(b []byte) error {
	tmp := inetWrapper{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*dst = Inet{Status: tmp.Status}
	if dst.Status == Present {
		if err := dst.Set(tmp.IPNet.String()); err != nil {
			return err
		}
	}

	return nil
}

func (dst Inet) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.IPNet
	case Null:
		return nil
	default:
		return dst.Status
	}
}
