package schema

import "encoding/json"

type CIDR Inet

type CIDRTransformer interface {
	TransformCIDR(*CIDR) any
}

func (dst *CIDR) GetStatus() Status {
	return dst.Status
}

func (*CIDR) Type() ValueType {
	return TypeCIDR
}

func (dst *CIDR) Size() int {
	return len(dst.IPNet.IP) + len(dst.IPNet.Mask)
}

func (dst *CIDR) String() string {
	if dst.Status == Present {
		return dst.IPNet.String()
	}

	return ""
}

func (dst *CIDR) Equal(src CQType) bool {
	if src == nil {
		return false
	}
	s, ok := src.(*CIDR)
	if !ok {
		return false
	}
	return dst.Status == s.Status && dst.IPNet.String() == s.IPNet.String()
}

func (dst *CIDR) Set(src any) error {
	return (*Inet)(dst).Set(src)
}

func (dst CIDR) Get() any {
	return (Inet)(dst).Get()
}

func (dst *CIDR) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*Inet)(dst))
}
