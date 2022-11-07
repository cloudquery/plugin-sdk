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

// This is needed only for backward compatibility.
// can be removed in sdk v2
func valueTypeFromOverTheWireString(s string) ValueType {
	switch s {
	case "Bool":
		return TypeBool
	case "Int8":
		return TypeInt
	case "Float8":
		return TypeFloat
	case "UUID":
		return TypeUUID
	case "Text":
		return TypeString
	case "JSON":
		return TypeJSON
	case "Int8Array":
		return TypeIntArray
	case "TextArray":
		return TypeStringArray
	case "Timestamptz":
		return TypeTimestamp
	case "Bytea":
		return TypeByteArray
	case "UUIDArray":
		return TypeUUIDArray
	case "InetArray":
		return TypeInetArray
	case "Inet":
		return TypeInet
	case "MacaddrArray":
		return TypeMacAddrArray
	case "Macaddr":
		return TypeMacAddr
	case "CIDRArray":
		return TypeCIDRArray
	case "CIDR":
		return TypeCIDR
	default:
		return TypeInvalid
	}
}

// this is for backward compatibility
// can be removed in sdk v2
func (r ValueType) overTheWireString() string {
	switch r {
	case TypeBool:
		return "Bool"
	case TypeInt:
		return "Int8"
	case TypeFloat:
		return "Float8"
	case TypeUUID:
		return "UUID"
	case TypeString:
		return "Text"
	case TypeJSON:
		return "JSON"
	case TypeIntArray:
		return "Int8Array"
	case TypeStringArray:
		return "TextArray"
	case TypeTimestamp:
		return "Timestamptz"
	case TypeByteArray:
		return "Bytea"
	case TypeUUIDArray:
		return "UUIDArray"
	case TypeInetArray:
		return "InetArray"
	case TypeInet:
		return "Inet"
	case TypeMacAddrArray:
		return "MacaddrArray"
	case TypeMacAddr:
		return "Macaddr"
	case TypeCIDRArray:
		return "CIDRArray"
	case TypeCIDR:
		return "CIDR"
	case TypeInvalid:
		return "TypeInvalid"
	default:
		return fmt.Sprintf("Unknown(%d)", r)
	}
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

func (c CQTypes) MarshalJSON() ([]byte, error) {
	res := make([]map[string]interface{}, len(c))
	for i, v := range c {
		res[i] = map[string]interface{}{
			"type":  v.Type().overTheWireString(),
			"value": v,
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
		var s string
		if err := json.Unmarshal(res[i]["type"], &s); err != nil {
			return fmt.Errorf("failed to unmarshal CQType type: %w", err)
		}
		t := valueTypeFromOverTheWireString(s)
		if t <= TypeInvalid || t >= TypeEnd {
			// this means we dont support it yet on the destination side so we skip this CQType
			continue
		}
		v := NewCqTypeFromValueType(t)
		if err := json.Unmarshal(res[i]["value"], v); err != nil {
			return fmt.Errorf("failed to unmarshal CQType value: %w", err)
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
