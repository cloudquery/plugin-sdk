package schema

import "encoding/json"

type CIDR Inet

type CIDRTransformer interface {
	TransformCIDR(*CIDR) interface{}
}

func (*CIDR) Type() ValueType {
	return TypeCIDR
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

func (dst *CIDR) Set(src interface{}) error {
	return (*Inet)(dst).Set(src)
}

func (dst CIDR) Get() interface{} {
	return (Inet)(dst).Get()
}

func (dst *CIDR) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*Inet)(dst))
}
