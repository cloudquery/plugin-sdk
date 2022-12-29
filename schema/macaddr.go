//nolint:revive
package schema

import (
	"fmt"
	"net"
)

type MacaddrTransformer interface {
	TransformMacaddr(*Macaddr) any
}

type Macaddr struct {
	Addr   net.HardwareAddr
	Status Status
}

func (*Macaddr) Type() ValueType {
	return TypeMacAddr
}

func (dst *Macaddr) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Macaddr)
	if !ok {
		return false
	}
	return dst.Status == s.Status && dst.Addr.String() == s.Addr.String()
}

func (dst *Macaddr) LessThan(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*Macaddr)
	if !ok {
		return false
	}
	if dst.Status != s.Status {
		return dst.Status < s.Status
	}
	for i := range dst.Addr {
		if dst.Addr[i] != s.Addr[i] {
			return dst.Addr[i] < s.Addr[i]
		}
	}
	return false
}

func (dst *Macaddr) String() string {
	if dst.Status == Present {
		return dst.Addr.String()
	} else {
		return ""
	}
}

func (dst *Macaddr) Set(src any) error {
	if src == nil {
		*dst = Macaddr{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() any }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case net.HardwareAddr:
		addr := make(net.HardwareAddr, len(value))
		copy(addr, value)
		*dst = Macaddr{Addr: addr, Status: Present}
	case string:
		addr, err := net.ParseMAC(value)
		if err != nil {
			return err
		}
		*dst = Macaddr{Addr: addr, Status: Present}
	case *net.HardwareAddr:
		if value == nil {
			*dst = Macaddr{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Macaddr{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Macaddr", value)
	}

	return nil
}

func (dst Macaddr) Get() any {
	switch dst.Status {
	case Present:
		return dst.Addr
	case Null:
		return nil
	default:
		return dst.Status
	}
}
