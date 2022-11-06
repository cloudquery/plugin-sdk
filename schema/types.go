package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Status byte

const (
	Undefined Status = iota
	Null
	Present
)

type InfinityModifier int8

type ValueType int

type deprecatedTypes []ValueType

const (
	Infinity         InfinityModifier = 1
	None             InfinityModifier = 0
	NegativeInfinity InfinityModifier = -Infinity
)

var deprecatedTypesValues = deprecatedTypes{
	TypeTimeIntervalDeprecated,
}

func (v deprecatedTypes) isDeprecated(t ValueType) bool {
	for _, dt := range v {
		if dt == t {
			return true
		}
	}
	return false
}

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeUUID
	TypeString
	TypeByteArray
	TypeStringArray
	TypeIntArray
	TypeTimestamp
	TypeJSON
	TypeUUIDArray
	TypeInet
	TypeInetArray
	TypeCIDR
	TypeCIDRArray
	TypeMacAddr
	TypeMacAddrArray
	TypeTimeIntervalDeprecated
	TypeEnd
)

func (r *ValueType) UnmarshalJSON(data []byte) (err error) {
	var valueType int
	if err := json.Unmarshal(data, &valueType); err != nil {
		return err
	}
	if ValueType(valueType) <= TypeInvalid || ValueType(valueType) >= TypeEnd {
		*r = TypeInvalid
	} else {
		*r = ValueType(valueType)
	}
	return nil
}

func (r ValueType) String() string {
	switch r {
	case TypeBool:
		return "TypeBool"
	case TypeInt:
		return "TypeInt"
	case TypeFloat:
		return "TypeFloat"
	case TypeUUID:
		return "TypeUUID"
	case TypeString:
		return "TypeString"
	case TypeJSON:
		return "TypeJSON"
	case TypeIntArray:
		return "TypeIntArray"
	case TypeStringArray:
		return "TypeStringArray"
	case TypeTimestamp:
		return "TypeTimestamp"
	case TypeByteArray:
		return "TypeByteArray"
	case TypeUUIDArray:
		return "TypeUUIDArray"
	case TypeInetArray:
		return "TypeInetArray"
	case TypeInet:
		return "TypeInet"
	case TypeMacAddrArray:
		return "TypeMacAddrArray"
	case TypeMacAddr:
		return "TypeMacAddr"
	case TypeCIDRArray:
		return "TypeCIDRArray"
	case TypeCIDR:
		return "TypeCIDR"
	case TypeInvalid:
		return "TypeInvalid"
	default:
		return fmt.Sprintf("Unknown(%d)", r)
	}
}

type CQType interface {
	Set(v interface{}) error
	Get() interface{}
	String() string
	Equal(CQType) bool
	Type() ValueType
}

type CQTypes []CQType

type cqTypeWrapper struct {
	Type  ValueType `json:"type"`
	Value CQType    `json:"value"`
}

func (c CQTypes) MarshalJSON() ([]byte, error) {
	res := make([]*cqTypeWrapper, len(c))
	for i, v := range c {
		res[i] = &cqTypeWrapper{
			Type:  v.Type(),
			Value: v,
		}
	}
	return json.Marshal(res)
}

func (c *CQTypes) UnmarshalJSON(b []byte) error {
	var res []map[string]json.RawMessage
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	cqTypes := make(CQTypes, 0, len(res))
	for i := range res {
		var t ValueType
		if err := json.Unmarshal(res[i]["type"], &t); err != nil {
			return err
		}
		if t <= TypeInvalid || t >= TypeEnd {
			// this means we dont support it yet on the destination side so we skip this CQType
			continue
		}
		v := NewCqTypeFromValueType(t)
		if err := json.Unmarshal(res[i]["value"], v); err != nil {
			return err
		}

		cqTypes = append(cqTypes, v)
	}
	*c = cqTypes
	return nil
}

func (c CQTypes) Len() int {
	return len(c)
}

func (c CQTypes) Equal(other CQTypes) bool {
	if other == nil {
		return false
	}
	if len(c) != len(other) {
		return false
	}
	for i := range c {
		if c[i] == nil && other[i] == nil {
			continue
		}
		if c[i] == nil || other[i] == nil {
			return false
		}
		if !c[i].Equal(other[i]) {
			return false
		}
	}
	return true
}

func (c CQTypes) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range c {
		if i > 0 {
			sb.WriteString(", ")
		}
		if v == nil {
			sb.WriteString("nil")
		} else {
			sb.WriteString(v.String())
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func NewCqTypeFromValueType(typ ValueType) CQType {
	switch typ {
	case TypeBool:
		return &Bool{}
	case TypeByteArray:
		return &Bytea{}
	case TypeCIDRArray:
		return &CIDRArray{}
	case TypeCIDR:
		return &CIDR{}
	case TypeFloat:
		return &Float8{}
	case TypeInetArray:
		return &InetArray{}
	case TypeInet:
		return &Inet{}
	case TypeIntArray:
		return &Int8Array{}
	case TypeInt:
		return &Int8{}
	case TypeJSON:
		return &JSON{}
	case TypeMacAddrArray:
		return &MacaddrArray{}
	case TypeMacAddr:
		return &Macaddr{}
	case TypeStringArray:
		return &TextArray{}
	case TypeString:
		return &Text{}
	case TypeTimestamp:
		return &Timestamptz{}
	case TypeUUIDArray:
		return &UUIDArray{}
	case TypeUUID:
		return &UUID{}
	default:
		panic("unknown type " + typ.String())
	}
}
